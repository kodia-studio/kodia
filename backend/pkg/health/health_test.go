package health

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockChecker is a test implementation of Checker
type MockChecker struct {
	name  string
	err   error
	delay time.Duration
}

func (m *MockChecker) Name() string {
	return m.name
}

func (m *MockChecker) Check(ctx context.Context) error {
	if m.delay > 0 {
		time.Sleep(m.delay)
	}
	return m.err
}

// TestGatherBasicStats tests gathering basic health statistics
func TestGatherBasicStats(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stats, err := Gather(ctx)

	require.NoError(t, err)
	require.NotNil(t, stats)

	// Verify timestamp
	assert.WithinDuration(t, time.Now(), stats.Timestamp, 1*time.Second)

	// Verify CPU usage
	assert.GreaterOrEqual(t, stats.CPUUsage, 0.0)
	assert.LessOrEqual(t, stats.CPUUsage, 100.0)

	// Verify memory stats
	assert.Greater(t, stats.MemoryTotal, uint64(0))
	assert.Greater(t, stats.MemoryUsage, 0.0)
	assert.LessOrEqual(t, stats.MemoryUsage, 100.0)

	// Verify goroutine count
	assert.Greater(t, stats.Goroutines, 0)

	// Verify initial status
	assert.Equal(t, "up", stats.Status)
}

// TestGatherWithPassingCheckers tests gathering stats with passing health checkers
func TestGatherWithPassingCheckers(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	checker1 := &MockChecker{name: "database", err: nil}
	checker2 := &MockChecker{name: "cache", err: nil}

	stats, err := Gather(ctx, checker1, checker2)

	require.NoError(t, err)
	assert.Equal(t, "up", stats.Status)
	assert.Len(t, stats.Checks, 2)

	// Verify both checkers passed
	assert.Equal(t, "up", stats.Checks["database"].Status)
	assert.Equal(t, "up", stats.Checks["cache"].Status)
	assert.Empty(t, stats.Checks["database"].Error)
	assert.Empty(t, stats.Checks["cache"].Error)
}

// TestGatherWithFailingChecker tests gathering stats with failing health checkers
func TestGatherWithFailingCheckers(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	checker1 := &MockChecker{name: "database", err: nil}
	checker2 := &MockChecker{name: "cache", err: errors.New("connection refused")}

	stats, err := Gather(ctx, checker1, checker2)

	require.NoError(t, err)

	// Overall status should be degraded
	assert.Equal(t, "degraded", stats.Status)

	// Database should be up
	assert.Equal(t, "up", stats.Checks["database"].Status)
	assert.Empty(t, stats.Checks["database"].Error)

	// Cache should be down with error
	assert.Equal(t, "down", stats.Checks["cache"].Status)
	assert.Equal(t, "connection refused", stats.Checks["cache"].Error)
}

// TestGatherWithMultipleFailures tests status with multiple checker failures
func TestGatherWithMultipleFailures(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	checkers := []Checker{
		&MockChecker{name: "db", err: errors.New("timeout")},
		&MockChecker{name: "cache", err: errors.New("unreachable")},
		&MockChecker{name: "queue", err: nil},
	}

	stats, err := Gather(ctx, checkers...)

	require.NoError(t, err)
	assert.Equal(t, "degraded", stats.Status)
	assert.Equal(t, "down", stats.Checks["db"].Status)
	assert.Equal(t, "down", stats.Checks["cache"].Status)
	assert.Equal(t, "up", stats.Checks["queue"].Status)
}

// TestGatherNoCheckers tests gathering with no checkers
func TestGatherNoCheckers(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stats, err := Gather(ctx)

	require.NoError(t, err)
	assert.Equal(t, "up", stats.Status)
	assert.Len(t, stats.Checks, 0)
}

// TestGatherContextTimeout tests behavior when context times out
func TestGatherContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	checker := &MockChecker{
		name:  "slow-checker",
		delay: 100 * time.Millisecond,
	}

	// This may timeout during Gather itself, or the checker may not complete
	// Either way, we should not panic
	stats, err := Gather(ctx, checker)

	// We may or may not get an error depending on timing
	if err == nil {
		assert.NotNil(t, stats)
	}
}

// TestCheckResultFields tests CheckResult fields
func TestCheckResultFields(t *testing.T) {
	result := CheckResult{
		Status: "down",
		Error:  "connection timeout",
	}

	assert.Equal(t, "down", result.Status)
	assert.Equal(t, "connection timeout", result.Error)
}

// TestStatsFields tests all Stats fields are properly set
func TestStatsFields(t *testing.T) {
	ctx := context.Background()
	stats, err := Gather(ctx)

	require.NoError(t, err)

	// All fields should be non-nil or have valid values
	assert.False(t, stats.Timestamp.IsZero())
	assert.Greater(t, stats.MemoryTotal, uint64(0))
	assert.Greater(t, stats.Goroutines, 0)
	assert.NotEmpty(t, stats.Status)
	assert.NotNil(t, stats.Checks)
}

// TestStatsMemoryRatios tests memory ratios are consistent
func TestStatsMemoryRatios(t *testing.T) {
	ctx := context.Background()
	stats, err := Gather(ctx)

	require.NoError(t, err)

	// Used memory should not exceed total memory
	assert.LessOrEqual(t, stats.MemoryUsed, stats.MemoryTotal)

	// Memory usage percentage should be between 0 and 100
	assert.GreaterOrEqual(t, stats.MemoryUsage, 0.0)
	assert.LessOrEqual(t, stats.MemoryUsage, 100.0)
}

// TestStatsDiskRatios tests disk ratios are consistent
func TestStatsDiskRatios(t *testing.T) {
	ctx := context.Background()
	stats, err := Gather(ctx)

	require.NoError(t, err)

	if stats.DiskTotal > 0 {
		// Used disk should not exceed total disk
		assert.LessOrEqual(t, stats.DiskUsed, stats.DiskTotal)

		// Disk usage percentage should be valid
		assert.GreaterOrEqual(t, stats.DiskUsage, 0.0)
	}
}

// TestCheckerName tests Checker.Name() method
func TestCheckerName(t *testing.T) {
	checker := &MockChecker{name: "my-checker"}

	assert.Equal(t, "my-checker", checker.Name())
}

// TestCheckerCheck tests Checker.Check() method
func TestCheckerCheck(t *testing.T) {
	// Passing checker
	checker := &MockChecker{name: "test", err: nil}
	err := checker.Check(context.Background())
	assert.NoError(t, err)

	// Failing checker
	checker = &MockChecker{name: "test", err: errors.New("test error")}
	err = checker.Check(context.Background())
	assert.Error(t, err)
	assert.Equal(t, "test error", err.Error())
}

// TestGatherMultipleCallsConsistent tests that multiple calls produce consistent structure
func TestGatherMultipleCallsConsistent(t *testing.T) {
	ctx := context.Background()

	stats1, err1 := Gather(ctx)
	require.NoError(t, err1)

	stats2, err2 := Gather(ctx)
	require.NoError(t, err2)

	// Both should have same structure
	assert.Equal(t, stats1.Status, stats2.Status)
	assert.Greater(t, stats1.Goroutines, 0)
	assert.Greater(t, stats2.Goroutines, 0)
}

// TestGatherWithManyCheckers tests gathering with many checkers
func TestGatherWithManyCheckers(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	checkers := []Checker{}
	for i := 0; i < 10; i++ {
		checkers = append(checkers, &MockChecker{
			name: "checker-" + string(rune(i)),
			err:  nil,
		})
	}

	stats, err := Gather(ctx, checkers...)

	require.NoError(t, err)
	assert.Equal(t, "up", stats.Status)
	assert.Len(t, stats.Checks, 10)
}

// BenchmarkGather benchmarks the Gather function
func BenchmarkGather(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Gather(ctx)
	}
}

// BenchmarkGatherWithCheckers benchmarks Gather with multiple checkers
func BenchmarkGatherWithCheckers(b *testing.B) {
	ctx := context.Background()
	checkers := []Checker{
		&MockChecker{name: "db", err: nil},
		&MockChecker{name: "cache", err: nil},
		&MockChecker{name: "queue", err: nil},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Gather(ctx, checkers...)
	}
}
