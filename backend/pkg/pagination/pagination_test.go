package pagination

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestParamsOffset tests offset calculation
func TestParamsOffset(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		perPage  int
		expected int
	}{
		{"First page", 1, 15, 0},
		{"Second page", 2, 15, 15},
		{"Third page", 3, 10, 20},
		{"Page 10", 10, 50, 450},
		{"Page 0 treated as 1", 0, 15, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Params{Page: tt.page, PerPage: tt.perPage}

			result := p.Offset()

			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestParamsLimit tests limit return
func TestParamsLimit(t *testing.T) {
	p := &Params{PerPage: 25}

	result := p.Limit()

	assert.Equal(t, 25, result)
}

// TestParamsTotalPages tests total pages calculation
func TestParamsTotalPages(t *testing.T) {
	tests := []struct {
		name      string
		perPage   int
		total     int64
		expected  int
	}{
		{"Exact pages", 15, 60, 4},
		{"With remainder", 15, 61, 5},
		{"Single page", 100, 50, 1},
		{"Large dataset", 10, 1000, 100},
		{"One item", 15, 1, 1},
		{"No items", 15, 0, 0},
		{"Zero per page", 0, 100, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Params{PerPage: tt.perPage}

			result := p.TotalPages(tt.total)

			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestFromContextDefaults tests default values from context
func TestFromContextDefaults(t *testing.T) {
	c, _ := gin.CreateTestContext(nil)

	p := FromContext(c)

	assert.Equal(t, DefaultPage, p.Page)
	assert.Equal(t, DefaultPerPage, p.PerPage)
	assert.Equal(t, "asc", p.SortDir)
	assert.Empty(t, p.Sort)
	assert.Empty(t, p.Search)
}

// TestFromContextWithValidValues tests parsing valid query parameters
func TestFromContextWithValidValues(t *testing.T) {
	c, _ := gin.CreateTestContext(nil)
	c.Request.URL.RawQuery = "page=3&per_page=25&sort=name&sort_dir=desc&search=test"

	p := FromContext(c)

	assert.Equal(t, 3, p.Page)
	assert.Equal(t, 25, p.PerPage)
	assert.Equal(t, "name", p.Sort)
	assert.Equal(t, "desc", p.SortDir)
	assert.Equal(t, "test", p.Search)
}

// TestFromContextPageBounds tests page parameter bounds
func TestFromContextPageBounds(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected int
	}{
		{"Negative page defaults to 1", "page=-5", DefaultPage},
		{"Zero page defaults to 1", "page=0", DefaultPage},
		{"Large page", "page=999999", 999999},
		{"Non-numeric page", "page=abc", DefaultPage},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(nil)
			c.Request.URL.RawQuery = tt.query

			p := FromContext(c)

			assert.Equal(t, tt.expected, p.Page)
		})
	}
}

// TestFromContextPerPageBounds tests per_page parameter bounds
func TestFromContextPerPageBounds(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected int
	}{
		{"Negative per_page defaults", "per_page=-10", DefaultPerPage},
		{"Zero per_page defaults", "per_page=0", DefaultPerPage},
		{"Valid per_page", "per_page=50", 50},
		{"Exceeds max", "per_page=999", MaxPerPage},
		{"Non-numeric per_page", "per_page=xyz", DefaultPerPage},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(nil)
			c.Request.URL.RawQuery = tt.query

			p := FromContext(c)

			assert.Equal(t, tt.expected, p.PerPage)
		})
	}
}

// TestFromContextSortDirection tests sort direction validation
func TestFromContextSortDirection(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected string
	}{
		{"Valid asc", "sort_dir=asc", "asc"},
		{"Valid desc", "sort_dir=desc", "desc"},
		{"Invalid defaults to asc", "sort_dir=invalid", "asc"},
		{"Empty defaults to asc", "sort_dir=", "asc"},
		{"Case sensitive - defaults to asc", "sort_dir=ASC", "asc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(nil)
			c.Request.URL.RawQuery = tt.query

			p := FromContext(c)

			assert.Equal(t, tt.expected, p.SortDir)
		})
	}
}

// TestFromContextSearch tests search parameter parsing
func TestFromContextSearch(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected string
	}{
		{"Simple search", "search=hello", "hello"},
		{"Search with spaces", "search=hello%20world", "hello world"},
		{"Empty search", "search=", ""},
		{"No search param", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(nil)
			c.Request.URL.RawQuery = tt.query

			p := FromContext(c)

			assert.Equal(t, tt.expected, p.Search)
		})
	}
}

// TestFromContextMultipleParams tests all parameters together
func TestFromContextMultipleParams(t *testing.T) {
	c, _ := gin.CreateTestContext(nil)
	c.Request.URL.RawQuery = "page=5&per_page=30&sort=email&sort_dir=desc&search=admin"

	p := FromContext(c)

	assert.Equal(t, 5, p.Page)
	assert.Equal(t, 30, p.PerPage)
	assert.Equal(t, "email", p.Sort)
	assert.Equal(t, "desc", p.SortDir)
	assert.Equal(t, "admin", p.Search)
	assert.Equal(t, 120, p.Offset())
	assert.Equal(t, 30, p.Limit())
}

// TestOffsetCalculations tests various offset calculations
func TestOffsetCalculations(t *testing.T) {
	tests := []struct {
		page    int
		perPage int
		offset  int
	}{
		{1, 10, 0},
		{2, 10, 10},
		{3, 10, 20},
		{1, 25, 0},
		{2, 25, 25},
		{10, 15, 135},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			p := &Params{Page: tt.page, PerPage: tt.perPage}
			assert.Equal(t, tt.offset, p.Offset())
		})
	}
}

// TestTotalPagesRounding tests proper rounding in total pages
func TestTotalPagesRounding(t *testing.T) {
	tests := []struct {
		perPage int
		total   int64
		pages   int
	}{
		{10, 0, 0},      // 0 / 10 = 0
		{10, 1, 1},      // 1 / 10 = 0.1 -> ceil to 1
		{10, 9, 1},      // 9 / 10 = 0.9 -> ceil to 1
		{10, 10, 1},     // 10 / 10 = 1
		{10, 11, 2},     // 11 / 10 = 1.1 -> ceil to 2
		{10, 100, 10},   // 100 / 10 = 10
		{10, 101, 11},   // 101 / 10 = 10.1 -> ceil to 11
	}

	for _, tt := range tests {
		p := &Params{PerPage: tt.perPage}
		assert.Equal(t, tt.pages, p.TotalPages(tt.total))
	}
}

// BenchmarkFromContext benchmarks context parsing
func BenchmarkFromContext(b *testing.B) {
	c, _ := gin.CreateTestContext(nil)
	c.Request.URL.RawQuery = "page=5&per_page=30&sort=name&sort_dir=desc&search=test"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FromContext(c)
	}
}

// BenchmarkOffset benchmarks offset calculation
func BenchmarkOffset(b *testing.B) {
	p := &Params{Page: 100, PerPage: 25}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Offset()
	}
}

// BenchmarkTotalPages benchmarks total pages calculation
func BenchmarkTotalPages(b *testing.B) {
	p := &Params{PerPage: 15}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.TotalPages(1000000)
	}
}
