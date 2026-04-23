// Package health provides system health check and statistics gathering for Kodia.
package health

import (
	"context"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	"gorm.io/gorm"
	"github.com/redis/go-redis/v9"
)

// Checker defines the interface for health checks.
type Checker interface {
	Name() string
	Check(ctx context.Context) error
}

// Stats represents the system health statistics.
type Stats struct {
	Timestamp   time.Time              `json:"timestamp"`
	CPUUsage    float64                `json:"cpu_usage_percent"`
	MemoryTotal uint64                 `json:"memory_total_bytes"`
	MemoryUsed  uint64                 `json:"memory_used_bytes"`
	MemoryUsage float64                `json:"memory_usage_percent"`
	DiskTotal   uint64                 `json:"disk_total_bytes"`
	DiskUsed    uint64                 `json:"disk_used_bytes"`
	DiskUsage   float64                `json:"disk_usage_percent"`
	Goroutines  int                    `json:"goroutines"`
	Status      string                 `json:"status"` // "up", "degraded", "down"
	Checks      map[string]CheckResult `json:"checks,omitempty"`
}

// CheckResult represents the result of an individual health check.
type CheckResult struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// Gather collects current system statistics and runs optional checkers.
func Gather(ctx context.Context, checkers ...Checker) (*Stats, error) {
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
		Checks:      make(map[string]CheckResult),
	}

	// Run checkers
	for _, c := range checkers {
		res := CheckResult{Status: "up"}
		if err := c.Check(ctx); err != nil {
			res.Status = "down"
			res.Error = err.Error()
			stats.Status = "degraded"
		}
		stats.Checks[c.Name()] = res
	}

	return stats, nil
}

// DBChecker checks the health of a GORM database connection.
type DBChecker struct {
	DB *gorm.DB
}

func (c *DBChecker) Name() string { return "database" }
func (c *DBChecker) Check(ctx context.Context) error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

// RedisChecker checks the health of a Redis connection.
type RedisChecker struct {
	Client *redis.Client
}

func (c *RedisChecker) Name() string { return "redis" }
func (c *RedisChecker) Check(ctx context.Context) error {
	return c.Client.Ping(ctx).Err()
}
