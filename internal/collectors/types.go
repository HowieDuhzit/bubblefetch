package collectors

// SystemInfo holds all collected system information
type SystemInfo struct {
	OS         string
	Kernel     string
	Hostname   string
	Uptime     string
	CPU        string
	Memory     MemoryInfo
	Disk       DiskInfo
	Shell      string
	Terminal   string
	Resolution string
	DE         string
	WM         string
	Theme      string
	Icons      string
	GPU        []string
	Network    []NetworkInfo
	Battery    BatteryInfo
	LocalIP    string
	PublicIP   string
}

type MemoryInfo struct {
	Used  uint64
	Total uint64
}

type DiskInfo struct {
	Used  uint64
	Total uint64
}

type NetworkInfo struct {
	Interface string
	IPv4      string
	IPv6      string
	MAC       string
}

type BatteryInfo struct {
	Present     bool
	Percentage  float64
	IsCharging  bool
	TimeRemain  string
}

// Collector interface for gathering system information
type Collector interface {
	Collect() (*SystemInfo, error)
}
