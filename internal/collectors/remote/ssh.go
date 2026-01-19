package remote

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/howieduhzit/bubblefetch/internal/collectors"
	"github.com/howieduhzit/bubblefetch/internal/config"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/knownhosts"
)

type SSHCollector struct {
	host   string
	config *config.Config
	client *ssh.Client
}

func New(host string, cfg *config.Config) *SSHCollector {
	return &SSHCollector{
		host:   host,
		config: cfg,
	}
}

func (c *SSHCollector) Connect() error {
	// Determine user
	user := c.config.SSH.User
	if user == "" {
		user = os.Getenv("USER")
	}

	// Determine port
	port := c.config.SSH.Port
	if port == 0 {
		port = 22
	}

	// Parse host (support user@host format)
	inputHost := c.host
	if strings.Contains(inputHost, "@") {
		parts := strings.SplitN(inputHost, "@", 2)
		user = parts[0]
		inputHost = parts[1]
	}

	host, explicitPort := splitHostPort(inputHost)

	sshCfg := loadSSHConfig("")
	hostCfg := sshCfg.match(host)

	if hostCfg.HostName != "" {
		host = hostCfg.HostName
	}
	if hostCfg.User != "" && c.config.SSH.User == "" {
		user = hostCfg.User
	}
	if !explicitPort && hostCfg.Port != "" {
		if parsed, err := strconv.Atoi(hostCfg.Port); err == nil {
			port = parsed
		}
	}

	authMethods, err := buildSSHAuthMethods(c.config.SSH.KeyPath, hostCfg.IdentityFiles, hostCfg.IdentityAgent)
	if err != nil {
		return err
	}

	hostKeyCallback, err := buildHostKeyCallback(c.config.SSH.KnownHostsPath)
	if err != nil {
		return err
	}

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: authMethods,
		HostKeyCallback: hostKeyCallback,
		Timeout:         10 * time.Second,
	}

	addr := net.JoinHostPort(host, strconv.Itoa(port))
	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %v", addr, err)
	}

	c.client = client
	return nil
}

type sshHostConfig struct {
	Patterns      []string
	HostName      string
	User          string
	Port          string
	IdentityFiles []string
	IdentityAgent string
}

type sshConfigEntries []sshHostConfig

func loadSSHConfig(configPath string) sshConfigEntries {
	if configPath == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			configPath = filepath.Join(home, ".ssh", "config")
		}
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil
	}
	defer file.Close()

	var entries sshConfigEntries
	var current *sshHostConfig

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		key := strings.ToLower(fields[0])
		value := strings.Join(fields[1:], " ")

		if key == "host" {
			if current != nil {
				entries = append(entries, *current)
			}
			current = &sshHostConfig{
				Patterns: fields[1:],
			}
			continue
		}

		if current == nil {
			continue
		}

		switch key {
		case "hostname":
			current.HostName = value
		case "user":
			current.User = value
		case "port":
			current.Port = value
		case "identityfile":
			current.IdentityFiles = append(current.IdentityFiles, expandHome(value))
		case "identityagent":
			current.IdentityAgent = expandHome(value)
		}
	}

	if current != nil {
		entries = append(entries, *current)
	}

	return entries
}

func (entries sshConfigEntries) match(host string) sshHostConfig {
	var matched sshHostConfig
	for _, entry := range entries {
		for _, pattern := range entry.Patterns {
			if matchHost(pattern, host) {
				matched = mergeHostConfig(matched, entry)
				break
			}
		}
	}
	return matched
}

func mergeHostConfig(base, override sshHostConfig) sshHostConfig {
	if override.HostName != "" {
		base.HostName = override.HostName
	}
	if override.User != "" {
		base.User = override.User
	}
	if override.Port != "" {
		base.Port = override.Port
	}
	if len(override.IdentityFiles) > 0 {
		base.IdentityFiles = append(base.IdentityFiles, override.IdentityFiles...)
	}
	if override.IdentityAgent != "" {
		base.IdentityAgent = override.IdentityAgent
	}
	return base
}

func matchHost(pattern, host string) bool {
	if pattern == "*" {
		return true
	}
	if strings.ContainsAny(pattern, "*?") {
		ok, err := path.Match(pattern, host)
		return err == nil && ok
	}
	return strings.EqualFold(pattern, host)
}

func expandHome(value string) string {
	if strings.HasPrefix(value, "~") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, strings.TrimPrefix(value, "~"))
		}
	}
	return value
}

func splitHostPort(host string) (string, bool) {
	if strings.Contains(host, ":") {
		if h, _, err := net.SplitHostPort(host); err == nil {
			return h, true
		}
	}
	return host, false
}

func buildSSHAuthMethods(keyPath string, identityFiles []string, identityAgent string) ([]ssh.AuthMethod, error) {
	var methods []ssh.AuthMethod

	sock := identityAgent
	if sock == "" {
		sock = os.Getenv("SSH_AUTH_SOCK")
	}

	if sock != "" {
		conn, err := net.Dial("unix", sock)
		if err == nil {
			agentClient := agent.NewClient(conn)
			methods = append(methods, ssh.PublicKeysCallback(agentClient.Signers))
		}
	}

	keysToTry := []string{}
	if keyPath != "" {
		keysToTry = append(keysToTry, expandHome(keyPath))
	}
	keysToTry = append(keysToTry, identityFiles...)

	if len(keysToTry) == 0 {
		home, _ := os.UserHomeDir()
		keysToTry = []string{
			filepath.Join(home, ".ssh", "id_ed25519"),
			filepath.Join(home, ".ssh", "id_rsa"),
		}
	}

	for _, keyFile := range keysToTry {
		key, err := os.ReadFile(keyFile)
		if err != nil {
			continue
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			continue
		}

		methods = append(methods, ssh.PublicKeys(signer))
	}

	if len(methods) == 0 {
		return nil, fmt.Errorf("no usable SSH authentication methods found")
	}

	return methods, nil
}

func buildHostKeyCallback(knownHostsPath string) (ssh.HostKeyCallback, error) {
	path := knownHostsPath
	if path == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(home, ".ssh", "known_hosts")
		}
	}

	callback, err := knownhosts.New(path)
	if err == nil {
		return callback, nil
	}

	return ssh.InsecureIgnoreHostKey(), nil
}

func (c *SSHCollector) Collect() (*collectors.SystemInfo, error) {
	if c.client == nil {
		if err := c.Connect(); err != nil {
			return nil, err
		}
		defer c.client.Close()
	}

	info := &collectors.SystemInfo{}

	// Run commands in parallel
	type cmdResult struct {
		name   string
		output string
		err    error
	}

	commands := []struct {
		name string
		cmd  string
	}{
		{"os", "cat /etc/os-release 2>/dev/null || cat /usr/lib/os-release 2>/dev/null || uname -s"},
		{"kernel", "uname -r 2>/dev/null"},
		{"hostname", "hostname 2>/dev/null || cat /etc/hostname 2>/dev/null"},
		{"uptime", "cat /proc/uptime 2>/dev/null || uptime -p 2>/dev/null || uptime"},
		{"cpu", "awk -F: '/model name/ {print $2; exit}' /proc/cpuinfo 2>/dev/null | xargs || lscpu 2>/dev/null | awk -F: '/Model name/ {print $2}'"},
		{"memory", "awk '/MemTotal|MemAvailable/ {print}' /proc/meminfo 2>/dev/null"},
		{"disk", "df -B1 / 2>/dev/null | awk 'NR==2 {print}'"},
		{"shell", "echo $SHELL"},
		{"term", "echo $TERM"},
		{"de", "echo $XDG_CURRENT_DESKTOP"},
		{"wm", "echo $XDG_SESSION_TYPE"},
		{"gpu", "lspci 2>/dev/null | grep -iE 'vga|3d|display'"},
		{"network", "ip -o addr show 2>/dev/null | grep -v 'lo' | grep 'inet '"},
		{"battery", "cat /sys/class/power_supply/BAT*/capacity 2>/dev/null | head -1"},
	}

	results := make(map[string]string)
	for _, cmd := range commands {
		output, err := c.runCommand(cmd.cmd)
		clean := strings.TrimSpace(output)
		if clean != "" || err == nil {
			results[cmd.name] = clean
		}
	}

	// Parse results
	c.parseOS(info, results["os"])
	info.Kernel = results["kernel"]
	info.Hostname = results["hostname"]
	c.parseUptime(info, results["uptime"])
	info.CPU = results["cpu"]
	c.parseMemory(info, results["memory"])
	c.parseDisk(info, results["disk"])
	info.Shell = results["shell"]
	info.Terminal = results["term"]
	info.DE = results["de"]
	info.WM = results["wm"]
	c.parseGPU(info, results["gpu"])
	c.parseNetwork(info, results["network"])
	c.parseBattery(info, results["battery"])

	return info, nil
}

func (c *SSHCollector) runCommand(cmd string) (string, error) {
	session, err := c.client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.CombinedOutput(wrapRemoteCommand(cmd))
	return string(output), err
}

func wrapRemoteCommand(cmd string) string {
	prefix := "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin; LC_ALL=C; "
	return "/bin/sh -lc '" + escapeSingleQuotes(prefix+cmd) + "'"
}

func escapeSingleQuotes(input string) string {
	return strings.ReplaceAll(input, "'", "'\"'\"'")
}

func (c *SSHCollector) parseOS(info *collectors.SystemInfo, output string) {
	if strings.Contains(output, "PRETTY_NAME") {
		for _, line := range strings.Split(output, "\n") {
			if strings.HasPrefix(line, "PRETTY_NAME=") {
				info.OS = strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), "\"")
				return
			}
		}
	}
	info.OS = output
}

func (c *SSHCollector) parseUptime(info *collectors.SystemInfo, output string) {
	if strings.Contains(output, " ") {
		// Parse /proc/uptime format: "12345.67 12345.67"
		parts := strings.Fields(output)
		if len(parts) > 0 {
			if seconds, err := strconv.ParseFloat(parts[0], 64); err == nil {
				info.Uptime = formatUptime(uint64(seconds))
				return
			}
		}
	}
	// Fallback to raw uptime output
	info.Uptime = output
}

func (c *SSHCollector) parseMemory(info *collectors.SystemInfo, output string) {
	var total, available uint64
	for _, line := range strings.Split(output, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		value, _ := strconv.ParseUint(fields[1], 10, 64)
		value *= 1024 // Convert from KB to bytes

		if strings.HasPrefix(line, "MemTotal:") {
			total = value
		} else if strings.HasPrefix(line, "MemAvailable:") {
			available = value
		}
	}

	if total > 0 {
		info.Memory = collectors.MemoryInfo{
			Total: total,
			Used:  total - available,
		}
	}
}

func (c *SSHCollector) parseDisk(info *collectors.SystemInfo, output string) {
	fields := strings.Fields(output)
	if len(fields) >= 3 {
		total, _ := strconv.ParseUint(fields[1], 10, 64)
		used, _ := strconv.ParseUint(fields[2], 10, 64)
		info.Disk = collectors.DiskInfo{
			Total: total,
			Used:  used,
		}
	}
}

func (c *SSHCollector) parseGPU(info *collectors.SystemInfo, output string) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line != "" {
			// Extract GPU name after the PCI ID
			if idx := strings.Index(line, ":"); idx != -1 {
				gpu := strings.TrimSpace(line[idx+1:])
				gpu = strings.TrimPrefix(gpu, "VGA compatible controller: ")
				gpu = strings.TrimPrefix(gpu, "3D controller: ")
				gpu = strings.TrimPrefix(gpu, "Display controller: ")
				info.GPU = append(info.GPU, gpu)
			}
		}
	}
}

func (c *SSHCollector) parseNetwork(info *collectors.SystemInfo, output string) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 4 {
			iface := strings.TrimSuffix(fields[1], ":")
			ip := strings.Split(fields[3], "/")[0]

			netInfo := collectors.NetworkInfo{
				Interface: iface,
				IPv4:      ip,
			}
			info.Network = append(info.Network, netInfo)
			if info.LocalIP == "" {
				info.LocalIP = ip
			}
		}
	}
}

func (c *SSHCollector) parseBattery(info *collectors.SystemInfo, output string) {
	if output != "" {
		if percentage, err := strconv.ParseFloat(output, 64); err == nil {
			info.Battery = collectors.BatteryInfo{
				Present:    true,
				Percentage: percentage,
			}
		}
	}
}

func formatUptime(seconds uint64) string {
	duration := time.Duration(seconds) * time.Second
	days := duration / (24 * time.Hour)
	duration -= days * 24 * time.Hour
	hours := duration / time.Hour
	duration -= hours * time.Hour
	minutes := duration / time.Minute

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}
