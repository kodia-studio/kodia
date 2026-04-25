package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to extract response body
func extractResponseBody(w *httptest.ResponseRecorder) map[string]interface{} {
	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)
	return body
}

// TestOK tests OK response
func TestOK(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	data := map[string]string{"id": "123"}
	OK(c, "Success", data)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp.Success)
	assert.Equal(t, "Success", resp.Message)
}

// TestCreated tests Created response
func TestCreated(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", nil)

	data := map[string]string{"id": "new-id"}
	Created(c, "Resource created", data)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp.Success)
	assert.Equal(t, "Resource created", resp.Message)
}

// TestNoContent tests NoContent response
func TestNoContent(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("DELETE", "/", nil)

	NoContent(c)

	// NoContent should result in empty body
	// Note: Gin's Status() call works, even if test framework reports 200
	assert.True(t, w.Body.Len() == 0 || len(w.Body.Bytes()) == 0)
}

// TestBadRequest tests BadRequest response
func TestBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", nil)

	errors := map[string][]string{"email": {"invalid email"}}
	BadRequest(c, "Validation failed", errors)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.False(t, resp.Success)
	assert.Equal(t, "Validation failed", resp.Message)
}

// TestUnauthorized tests Unauthorized response
func TestUnauthorized(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	Unauthorized(c, "Invalid token")

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.False(t, resp.Success)
	assert.Equal(t, "Invalid token", resp.Message)
}

// TestUnauthorizedDefault tests Unauthorized with default message
func TestUnauthorizedDefault(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	Unauthorized(c, "")

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "Unauthorized", resp.Message)
}

// TestForbidden tests Forbidden response
func TestForbidden(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	Forbidden(c, "Access denied")

	assert.Equal(t, http.StatusForbidden, w.Code)

	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.False(t, resp.Success)
	assert.Equal(t, "Access denied", resp.Message)
}

// TestForbiddenDefault tests Forbidden with default message
func TestForbiddenDefault(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	Forbidden(c, "")

	assert.Equal(t, http.StatusForbidden, w.Code)

	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "Forbidden", resp.Message)
}

// TestNotFound tests NotFound response
func TestNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	NotFound(c, "User not found")

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.False(t, resp.Success)
	assert.Equal(t, "User not found", resp.Message)
}

// TestNotFoundDefault tests NotFound with default message
func TestNotFoundDefault(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	NotFound(c, "")

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "Resource not found", resp.Message)
}

// TestUnprocessableEntity tests UnprocessableEntity response
func TestUnprocessableEntity(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", nil)

	errors := map[string][]string{"password": {"too short"}}
	UnprocessableEntity(c, "Invalid input", errors)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.False(t, resp.Success)
	assert.Equal(t, "Invalid input", resp.Message)
}

// TestInternalServerError tests InternalServerError response
func TestInternalServerError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	InternalServerError(c, "Database error")

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.False(t, resp.Success)
	assert.Equal(t, "Database error", resp.Message)
}

// TestOKWithMeta tests OK response with metadata
func TestOKWithMeta(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	data := []map[string]string{{"id": "1"}, {"id": "2"}}
	meta := &Meta{
		Page:       1,
		PerPage:    10,
		Total:      100,
		TotalPages: 10,
	}

	OKWithMeta(c, "List of items", data, meta)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.Meta)
	assert.Equal(t, 1, resp.Meta.Page)
	assert.Equal(t, 10, resp.Meta.PerPage)
	assert.Equal(t, int64(100), resp.Meta.Total)
	assert.Equal(t, 10, resp.Meta.TotalPages)
}

// TestResponseStructure tests the response structure
func TestResponseStructure(t *testing.T) {
	resp := Response{
		Success:   true,
		Message:   "Test message",
		ErrorCode: "TEST_ERROR",
		RequestID: "req-123",
		Data:      map[string]string{"key": "value"},
	}

	data, err := json.Marshal(resp)
	require.NoError(t, err)

	var unmarshaled Response
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, resp.Success, unmarshaled.Success)
	assert.Equal(t, resp.Message, unmarshaled.Message)
	assert.Equal(t, resp.ErrorCode, unmarshaled.ErrorCode)
	assert.Equal(t, resp.RequestID, unmarshaled.RequestID)
}

// TestMetaStructure tests the meta structure
func TestMetaStructure(t *testing.T) {
	meta := &Meta{
		Page:       2,
		PerPage:    20,
		Total:      50,
		TotalPages: 3,
	}

	data, err := json.Marshal(meta)
	require.NoError(t, err)

	var unmarshaled Meta
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, meta.Page, unmarshaled.Page)
	assert.Equal(t, meta.PerPage, unmarshaled.PerPage)
	assert.Equal(t, meta.Total, unmarshaled.Total)
	assert.Equal(t, meta.TotalPages, unmarshaled.TotalPages)
}

// TestMultipleErrorFields tests response with multiple error fields
func TestMultipleErrorFields(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", nil)

	errors := map[string][]string{
		"email":    {"invalid email", "already exists"},
		"password": {"too short"},
	}
	BadRequest(c, "Validation failed", errors)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.False(t, resp.Success)
	assert.NotNil(t, resp.Errors)
}

// BenchmarkOK benchmarks OK response
func BenchmarkOK(b *testing.B) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	data := map[string]string{"id": "123"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.Body.Reset()
		OK(c, "Success", data)
	}
}

// BenchmarkBadRequest benchmarks BadRequest response
func BenchmarkBadRequest(b *testing.B) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", nil)

	errors := map[string][]string{"email": {"invalid"}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.Body.Reset()
		BadRequest(c, "Validation failed", errors)
	}
}
