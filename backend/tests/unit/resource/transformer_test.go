package resource

import (
	"testing"

	"github.com/kodia-studio/kodia/pkg/resource"
)

func TestTransformContextSince(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		compare  string
		expected bool
	}{
		{
			name:     "v1 since v1",
			version:  "v1",
			compare:  "v1",
			expected: true,
		},
		{
			name:     "v2 since v1",
			version:  "v2",
			compare:  "v1",
			expected: true,
		},
		{
			name:     "v1 since v2",
			version:  "v1",
			compare:  "v2",
			expected: false,
		},
		{
			name:     "v3 since v2",
			version:  "v3",
			compare:  "v2",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := resource.TransformContext{Version: tt.version}
			result := ctx.Since(tt.compare)
			if result != tt.expected {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTransformContextUntil(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		compare  string
		expected bool
	}{
		{
			name:     "v1 until v1",
			version:  "v1",
			compare:  "v1",
			expected: true,
		},
		{
			name:     "v1 until v2",
			version:  "v1",
			compare:  "v2",
			expected: true,
		},
		{
			name:     "v2 until v1",
			version:  "v2",
			compare:  "v1",
			expected: false,
		},
		{
			name:     "v2 until v3",
			version:  "v2",
			compare:  "v3",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := resource.TransformContext{Version: tt.version}
			result := ctx.Until(tt.compare)
			if result != tt.expected {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTransformContextOnly(t *testing.T) {
	tests := []struct {
		name        string
		fields      []string
		input       map[string]any
		expected    map[string]any
		expectedLen int
	}{
		{
			name:   "no fields filter returns all",
			fields: nil,
			input: map[string]any{
				"id":    "1",
				"name":  "John",
				"email": "john@example.com",
			},
			expectedLen: 3,
		},
		{
			name:   "filter to id and name",
			fields: []string{"id", "name"},
			input: map[string]any{
				"id":    "1",
				"name":  "John",
				"email": "john@example.com",
			},
			expected: map[string]any{
				"id":   "1",
				"name": "John",
			},
			expectedLen: 2,
		},
		{
			name:   "filter with non-existent field",
			fields: []string{"id", "nonexistent"},
			input: map[string]any{
				"id":   "1",
				"name": "John",
			},
			expected: map[string]any{
				"id": "1",
			},
			expectedLen: 1,
		},
		{
			name:        "empty fields returns all",
			fields:      []string{},
			input:       map[string]any{"id": "1", "name": "John"},
			expectedLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := resource.TransformContext{Fields: tt.fields}
			result := ctx.Only(tt.input)

			if len(result) != tt.expectedLen {
				t.Errorf("length: got %d, want %d", len(result), tt.expectedLen)
			}

			if tt.expected != nil {
				for k, v := range tt.expected {
					if result[k] != v {
						t.Errorf("field %s: got %v, want %v", k, result[k], v)
					}
				}
			}
		})
	}
}

func TestTransformFuncItem(t *testing.T) {
	type User struct {
		ID   string
		Name string
	}

	transformer := resource.TransformFunc[*User](func(u *User, ctx resource.TransformContext) map[string]any {
		if u == nil {
			return map[string]any{}
		}
		data := map[string]any{
			"id":   u.ID,
			"name": u.Name,
		}
		return ctx.Only(data)
	})

	ctx := resource.TransformContext{Version: "v1"}
	user := &User{ID: "123", Name: "John"}
	result := resource.Item(user, transformer, ctx)

	if result["id"] != "123" {
		t.Errorf("id: got %v, want %v", result["id"], "123")
	}
	if result["name"] != "John" {
		t.Errorf("name: got %v, want %v", result["name"], "John")
	}
}

func TestTransformFuncCollection(t *testing.T) {
	type User struct {
		ID   string
		Name string
	}

	transformer := resource.TransformFunc[*User](func(u *User, ctx resource.TransformContext) map[string]any {
		if u == nil {
			return map[string]any{}
		}
		return map[string]any{
			"id":   u.ID,
			"name": u.Name,
		}
	})

	ctx := resource.TransformContext{Version: "v1"}
	users := []*User{
		{ID: "1", Name: "Alice"},
		{ID: "2", Name: "Bob"},
		{ID: "3", Name: "Charlie"},
	}

	result := resource.Collection(users, transformer, ctx)

	if len(result) != 3 {
		t.Errorf("length: got %d, want %d", len(result), 3)
	}

	for i, user := range users {
		if result[i]["id"] != user.ID {
			t.Errorf("item %d id: got %v, want %v", i, result[i]["id"], user.ID)
		}
	}
}

func TestTransformFuncNilHandling(t *testing.T) {
	type User struct {
		ID string
	}

	transformer := resource.TransformFunc[*User](func(u *User, ctx resource.TransformContext) map[string]any {
		if u == nil {
			return map[string]any{}
		}
		return map[string]any{"id": u.ID}
	})

	ctx := resource.TransformContext{Version: "v1"}
	result := resource.Item(nil, transformer, ctx)

	if len(result) != 0 {
		t.Errorf("nil item should return empty map, got %d fields", len(result))
	}
}

func TestTransformContextVersionComparison(t *testing.T) {
	tests := []struct {
		name           string
		version        string
		sinceExpected  map[string]bool // version: expected result
		untilExpected  map[string]bool
	}{
		{
			name:    "v1 comparisons",
			version: "v1",
			sinceExpected: map[string]bool{
				"v1": true,
				"v2": false,
				"v3": false,
			},
			untilExpected: map[string]bool{
				"v1": true,
				"v2": true,
				"v3": true,
			},
		},
		{
			name:    "v2 comparisons",
			version: "v2",
			sinceExpected: map[string]bool{
				"v1": true,
				"v2": true,
				"v3": false,
			},
			untilExpected: map[string]bool{
				"v1": false,
				"v2": true,
				"v3": true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := resource.TransformContext{Version: tt.version}

			for compare, expected := range tt.sinceExpected {
				result := ctx.Since(compare)
				if result != expected {
					t.Errorf("Since(%s): got %v, want %v", compare, result, expected)
				}
			}

			for compare, expected := range tt.untilExpected {
				result := ctx.Until(compare)
				if result != expected {
					t.Errorf("Until(%s): got %v, want %v", compare, result, expected)
				}
			}
		})
	}
}
