package local

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/howieduhzit/bubblefetch/internal/collectors"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type LocalCollector struct {
	enablePublicIP bool
}

func New(enablePublicIP bool) *LocalCollector {
	return &LocalCollector{
		enablePublicIP: enablePublicIP,
	}
}

func (c *LocalCollector) Collect() (*collectors.SystemInfo, error) {
	info := &collectors.SystemInfo{}

	// Collect fast env vars immediately
	start := time.Now()
	info.Shell = os.Getenv("SHELL")
	if info.Shell == "" {
		info.Shell = "unknown"
	}
	collectors.AddModuleCost(info, "Shell", time.Since(start))

	start = time.Now()
	info.Terminal = os.Getenv("TERM")
	if info.Terminal == "" {
		info.Terminal = "unknown"
	}
	collectors.AddModuleCost(info, "Terminal", time.Since(start))

	start = time.Now()
	info.DE = os.Getenv("XDG_CURRENT_DESKTOP")
	collectors.AddModuleCost(info, "DE", time.Since(start))

	start = time.Now()
	info.WM = os.Getenv("XDG_SESSION_TYPE")
	collectors.AddModuleCost(info, "WM", time.Since(start))

	// Collect slower metrics in parallel
	type result struct {
		hostInfo *host.InfoStat
		cpuInfo  []cpu.InfoStat
		memInfo  *mem.VirtualMemoryStat
		diskInfo *disk.UsageStat
		gpuInfo  []string
		netInfo  []collectors.NetworkInfo
		batInfo  collectors.BatteryInfo
		localIP  string
		publicIP string
		hostDur  time.Duration
		cpuDur   time.Duration
		memDur   time.Duration
		diskDur  time.Duration
		gpuDur   time.Duration
		netDur   time.Duration
		batDur   time.Duration
		localDur time.Duration
		pubDur   time.Duration
	}

	resultChan := make(chan result, 1)

	go func() {
		var r result

		// Run all collection in parallel using goroutines
		done := make(chan bool, 9)

		go func() {
			start := time.Now()
			r.hostInfo, _ = host.Info()
			r.hostDur = time.Since(start)
			done <- true
		}()

		go func() {
			start := time.Now()
			r.cpuInfo, _ = cpu.Info()
			r.cpuDur = time.Since(start)
			done <- true
		}()

		go func() {
			start := time.Now()
			r.memInfo, _ = mem.VirtualMemory()
			r.memDur = time.Since(start)
			done <- true
		}()

		go func() {
			start := time.Now()
			r.diskInfo, _ = disk.Usage("/")
			r.diskDur = time.Since(start)
			done <- true
		}()

		go func() {
			start := time.Now()
			r.gpuInfo = detectGPU()
			r.gpuDur = time.Since(start)
			done <- true
		}()

		go func() {
			start := time.Now()
			r.netInfo = detectNetwork()
			r.netDur = time.Since(start)
			done <- true
		}()

		go func() {
			start := time.Now()
			r.batInfo = detectBattery()
			r.batDur = time.Since(start)
			done <- true
		}()

		go func() {
			start := time.Now()
			r.localIP = getLocalIP()
			r.localDur = time.Since(start)
			done <- true
		}()

		go func() {
			start := time.Now()
			if c.enablePublicIP {
				r.publicIP = detectPublicIP()
			}
			r.pubDur = time.Since(start)
			done <- true
		}()

		// Wait for all to complete
		for i := 0; i < 9; i++ {
			<-done
		}

		resultChan <- r
	}()

	// Wait for parallel collection
	r := <-resultChan

	// Process results
	if r.hostInfo != nil {
		info.OS = fmt.Sprintf("%s %s", r.hostInfo.Platform, r.hostInfo.PlatformVersion)
		info.Kernel = r.hostInfo.KernelVersion
		info.Hostname = r.hostInfo.Hostname
		info.Uptime = formatUptime(r.hostInfo.Uptime)
		collectors.AddModuleCost(info, "OS", r.hostDur)
		collectors.AddModuleCost(info, "Kernel", r.hostDur)
		collectors.AddModuleCost(info, "Host", r.hostDur)
		collectors.AddModuleCost(info, "Uptime", r.hostDur)
	}

	if len(r.cpuInfo) > 0 {
		info.CPU = r.cpuInfo[0].ModelName
	} else {
		info.CPU = runtime.GOARCH
	}
	collectors.AddModuleCost(info, "CPU", r.cpuDur)

	if r.memInfo != nil {
		info.Memory = collectors.MemoryInfo{
			Used:  r.memInfo.Used,
			Total: r.memInfo.Total,
		}
	}
	collectors.AddModuleCost(info, "Memory", r.memDur)

	if r.diskInfo != nil {
		info.Disk = collectors.DiskInfo{
			Used:  r.diskInfo.Used,
			Total: r.diskInfo.Total,
		}
	}
	collectors.AddModuleCost(info, "Disk", r.diskDur)

	info.GPU = r.gpuInfo
	info.Network = r.netInfo
	info.Battery = r.batInfo
	info.LocalIP = r.localIP
	info.PublicIP = r.publicIP
	collectors.AddModuleCost(info, "GPU", r.gpuDur)
	collectors.AddModuleCost(info, "Network", r.netDur)
	collectors.AddModuleCost(info, "Local IP", r.localDur)
	collectors.AddModuleCost(info, "Battery", r.batDur)
	if c.enablePublicIP {
		collectors.AddModuleCost(info, "Public IP", r.pubDur)
	}

	return info, nil
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
