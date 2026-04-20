package performance

import (
	"context"
	"sync"

	"gorm.io/gorm"
)

type contextKey string

const (
	queryCountKey contextKey = "kodia_query_count"
	// Threshold for warning about N+1 queries in a single request
	DefaultNPlusOneThreshold = 10
)

// QueryStats tracks query patterns within a context
type QueryStats struct {
	mu           sync.RWMutex
	TableCounts  map[string]int
	TotalQueries int
}

// NPlusOnePlugin tracks and warns about potential N+1 query patterns
type NPlusOnePlugin struct {
	Threshold int
	OnWarning func(tableName string, count int)
}

func NewNPlusOnePlugin(threshold int) *NPlusOnePlugin {
	if threshold <= 0 {
		threshold = DefaultNPlusOneThreshold
	}
	return &NPlusOnePlugin{
		Threshold: threshold,
	}
}

func (p *NPlusOnePlugin) Name() string {
	return "kodia:nplusone"
}

func (p *NPlusOnePlugin) Initialize(db *gorm.DB) error {
	return db.Callback().Query().After("gorm:query").Register("kodia:nplusone_after", p.afterQuery)
}

func (p *NPlusOnePlugin) afterQuery(db *gorm.DB) {
	if db.Statement.Context == nil {
		return
	}

	stats, ok := db.Statement.Context.Value(queryCountKey).(*QueryStats)
	if !ok || stats == nil {
		return
	}

	tableName := db.Statement.Table
	if tableName == "" && db.Statement.Schema != nil {
		tableName = db.Statement.Schema.Table
	}

	if tableName == "" {
		return
	}

	stats.mu.Lock()
	stats.TotalQueries++
	stats.TableCounts[tableName]++
	count := stats.TableCounts[tableName]
	stats.mu.Unlock()

	if count >= p.Threshold && p.OnWarning != nil {
		p.OnWarning(tableName, count)
	}
}

// InitContext initializes the query tracking stats in the context
func InitContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, queryCountKey, &QueryStats{
		TableCounts: make(map[string]int),
	})
}

// GetQueryStats retrieves query statistics from the context
func GetQueryStats(ctx context.Context) *QueryStats {
	stats, _ := ctx.Value(queryCountKey).(*QueryStats)
	return stats
}

// Middleware returns a Gin-compatible middleware (or adapter) that initializes the N+1 context
// Since this pkg should be independent of Gin, we just provide the context logic.
// The Gin middleware in internal/adapters will use this.
