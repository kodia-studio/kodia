// Package health provides system health check and statistics gathering for Kodia.
package health

import (
	"context"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

// Stats represents the system health statistics.
type Stats struct {
	Timestamp   time.Time `json:"timestamp"`
	CPUUsage    float64   `json:"cpu_usage_percent"`
	MemoryTotal uint64    `json:"memory_total_bytes"`
	MemoryUsed  uint64    `json:"memory_used_bytes"`
	MemoryUsage float64   `json:"memory_usage_percent"`
	DiskTotal   uint64    `json:"disk_total_bytes"`
	DiskUsed    uint64    `json:"disk_used_bytes"`
	DiskUsage   float64   `json:"disk_usage_percent"`
	Goroutines  int       `json:"goroutines"`
	Status      string    `json:"status"` // "up", "degraded", "down"
}

// Gather collects current system statistics.
func Gather(ctx context.Context) (*Stats, error) {
	// CPU
	cpuPercs, err := cpu.PercentWithContext(ctx, 0, false)
	cpuUsage := 0.0
	if err == nil && len(cpuPercs) > 0 {
		cpuUsage = cpuPercs[0]
	}

	// Memory
	vm, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return nil, err
	}

	// Disk
	du, err := disk.UsageWithContext(ctx, "/")
	if err != nil {
		// Fallback for environments where / might not be the root
		du = &disk.UsageStat{}
	}

	stats := &Stats{
		Timestamp:   time.Now(),
		CPUUsage:    cpuUsage,
		MemoryTotal: vm.Total,
		MemoryUsed:  vm.Used,
		MemoryUsage: vm.UsedPercent,
		DiskTotal:   du.Total,
		DiskUsed:    du.Used,
		DiskUsage:   du.UsedPercent,
		Goroutines:  runtime.NumGoroutine(),
		Status:      "up",
	}

	return stats, nil
}
