package system

import (
	"os"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

// SystemInfo contains system monitoring data.
type SystemInfo struct {
	CPUPercent float64 `json:"cpuPercent"`
	MemUsed    uint64  `json:"memUsed"`
	MemTotal   uint64  `json:"memTotal"`
	MemPercent float64 `json:"memPercent"`
	GoVersion  string  `json:"goVersion"`
	OS         string  `json:"os"`
	Arch       string  `json:"arch"`
}

// ProcessInfo contains process monitoring data.
type ProcessInfo struct {
	PID        int32   `json:"pid"`
	CPUPercent float64 `json:"cpuPercent"`
	MemRSS     uint64  `json:"memRSS"`
	MemPercent float32 `json:"memPercent"`
	Uptime     int64   `json:"uptime"` // seconds
	Running    bool    `json:"running"`
}

// Monitor provides system monitoring functions.
type Monitor struct {
	startTime time.Time
}

// NewMonitor creates a new system monitor.
func NewMonitor() *Monitor {
	return &Monitor{
		startTime: time.Now(),
	}
}

// GetSystemInfo returns current system resource usage.
func (m *Monitor) GetSystemInfo() SystemInfo {
	info := SystemInfo{
		GoVersion: runtime.Version(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}

	if cpuPercents, err := cpu.Percent(0, false); err == nil && len(cpuPercents) > 0 {
		info.CPUPercent = cpuPercents[0]
	}

	if vmem, err := mem.VirtualMemory(); err == nil {
		info.MemUsed = vmem.Used
		info.MemTotal = vmem.Total
		info.MemPercent = vmem.UsedPercent
	}

	return info
}

// GetProcessInfo returns info about a specific process by PID.
func (m *Monitor) GetProcessInfo(pid int32) ProcessInfo {
	info := ProcessInfo{
		PID:     pid,
		Running: false,
	}

	proc, err := process.NewProcess(pid)
	if err != nil {
		return info
	}

	info.Running = true

	if cpuPct, err := proc.CPUPercent(); err == nil {
		info.CPUPercent = cpuPct
	}

	if memInfo, err := proc.MemoryInfo(); err == nil && memInfo != nil {
		info.MemRSS = memInfo.RSS
	}

	if memPct, err := proc.MemoryPercent(); err == nil {
		info.MemPercent = memPct
	}

	if createTime, err := proc.CreateTime(); err == nil {
		created := time.UnixMilli(createTime)
		info.Uptime = int64(time.Since(created).Seconds())
	}

	return info
}

// GetAppUptime returns the application uptime in seconds.
func (m *Monitor) GetAppUptime() int64 {
	return int64(time.Since(m.startTime).Seconds())
}

// GetPID returns the current process PID.
func GetPID() int {
	return os.Getpid()
}
