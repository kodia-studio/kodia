// Package pagination provides standardized pagination helpers for Kodia Framework.
package pagination

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPage    = 1
	DefaultPerPage = 15
	MaxPerPage     = 100
)

// Params holds the parsed pagination parameters.
type Params struct {
	Page    int
	PerPage int
}

// Offset returns the SQL OFFSET value for the current page.
func (p *Params) Offset() int {
	return (p.Page - 1) * p.PerPage
}

// Limit returns the SQL LIMIT value.
func (p *Params) Limit() int {
	return p.PerPage
}

// TotalPages calculates the total number of pages.
func (p *Params) TotalPages(total int64) int {
	if p.PerPage == 0 {
		return 0
	}
	return int(math.Ceil(float64(total) / float64(p.PerPage)))
}

// FromContext parses pagination query parameters from a Gin context.
// Supports ?page=1&per_page=15
func FromContext(c *gin.Context) *Params {
	page := parseIntWithDefault(c.Query("page"), DefaultPage)
	perPage := parseIntWithDefault(c.Query("per_page"), DefaultPerPage)

	if page < 1 {
		page = DefaultPage
	}
	if perPage < 1 {
		perPage = DefaultPerPage
	}
	if perPage > MaxPerPage {
		perPage = MaxPerPage
	}

	return &Params{
		Page:    page,
		PerPage: perPage,
	}
}

func parseIntWithDefault(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return v
}

