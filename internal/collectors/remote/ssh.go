package remote

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"github.com/howieduhzit/bubblefetch/internal/collectors"
	"github.com/howieduhzit/bubblefetch/internal/config"
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
	hostParts := strings.Split(c.host, "@")
	if len(hostParts) == 2 {
		user = hostParts[0]
		c.host = hostParts[1]
	}

	// Load SSH key
	keyPath := c.config.SSH.KeyPath
	if keyPath == "" {
		home, _ := os.UserHomeDir()
		keyPath = filepath.Join(home, ".ssh", "id_rsa")
	}

	key, err := os.ReadFile(keyPath)
	if err != nil {
		// Try id_ed25519 as fallback
		keyPath = filepath.Join(filepath.Dir(keyPath), "id_ed25519")
		key, err = os.ReadFile(keyPath)
		if err != nil {
			return fmt.Errorf("failed to read SSH key: %v", err)
		}
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return fmt.Errorf("failed to parse SSH key: %v", err)
	}

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: Use known_hosts
		Timeout:         10 * time.Second,
	}

	addr := net.JoinHostPort(c.host, strconv.Itoa(port))
	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %v", addr, err)
	}

	c.client = client
	return nil
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

	commands := map[string]string{
		"os":       "cat /etc/os-release 2>/dev/null || uname -s",
		"kernel":   "uname -r",
		"hostname": "hostname",
		"uptime":   "cat /proc/uptime 2>/dev/null || uptime",
		"cpu":      "cat /proc/cpuinfo 2>/dev/null | grep 'model name' | head -1 | cut -d: -f2 | xargs",
		"memory":   "cat /proc/meminfo 2>/dev/null | grep -E 'MemTotal|MemAvailable'",
		"disk":     "df -B1 / 2>/dev/null | tail -1",
		"shell":    "echo $SHELL",
		"term":     "echo $TERM",
		"de":       "echo $XDG_CURRENT_DESKTOP",
		"wm":       "echo $XDG_SESSION_TYPE",
		"gpu":      "lspci 2>/dev/null | grep -iE 'vga|3d|display'",
		"network":  "ip -o addr show 2>/dev/null | grep -v 'lo' | grep 'inet '",
		"battery":  "cat /sys/class/power_supply/BAT*/capacity 2>/dev/null | head -1",
	}

	resultChan := make(chan cmdResult, len(commands))

	for name, cmd := range commands {
		go func(n, command string) {
			output, err := c.runCommand(command)
			resultChan <- cmdResult{name: n, output: output, err: err}
		}(name, cmd)
	}

	// Collect all results
	results := make(map[string]string)
	for i := 0; i < len(commands); i++ {
		result := <-resultChan
		if result.err == nil {
			results[result.name] = strings.TrimSpace(result.output)
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

	output, err := session.CombinedOutput(cmd)
	return string(output), err
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
