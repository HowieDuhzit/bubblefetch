package local

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/howieduhzit/bubblefetch/internal/collectors"
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
	info.Shell = os.Getenv("SHELL")
	if info.Shell == "" {
		info.Shell = "unknown"
	}
	info.Terminal = os.Getenv("TERM")
	if info.Terminal == "" {
		info.Terminal = "unknown"
	}
	info.DE = os.Getenv("XDG_CURRENT_DESKTOP")
	info.WM = os.Getenv("XDG_SESSION_TYPE")

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
	}

	resultChan := make(chan result, 1)

	go func() {
		var r result

		// Run all collection in parallel using goroutines
		done := make(chan bool, 9)

		go func() {
			r.hostInfo, _ = host.Info()
			done <- true
		}()

		go func() {
			r.cpuInfo, _ = cpu.Info()
			done <- true
		}()

		go func() {
			r.memInfo, _ = mem.VirtualMemory()
			done <- true
		}()

		go func() {
			r.diskInfo, _ = disk.Usage("/")
			done <- true
		}()

		go func() {
			r.gpuInfo = detectGPU()
			done <- true
		}()

		go func() {
			r.netInfo = detectNetwork()
			done <- true
		}()

		go func() {
			r.batInfo = detectBattery()
			done <- true
		}()

		go func() {
			r.localIP = getLocalIP()
			done <- true
		}()

		go func() {
			if c.enablePublicIP {
				r.publicIP = detectPublicIP()
			}
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
	}

	if len(r.cpuInfo) > 0 {
		info.CPU = r.cpuInfo[0].ModelName
	} else {
		info.CPU = runtime.GOARCH
	}

	if r.memInfo != nil {
		info.Memory = collectors.MemoryInfo{
			Used:  r.memInfo.Used,
			Total: r.memInfo.Total,
		}
	}

	if r.diskInfo != nil {
		info.Disk = collectors.DiskInfo{
			Used:  r.diskInfo.Used,
			Total: r.diskInfo.Total,
		}
	}

	info.GPU = r.gpuInfo
	info.Network = r.netInfo
	info.Battery = r.batInfo
	info.LocalIP = r.localIP
	info.PublicIP = r.publicIP

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
