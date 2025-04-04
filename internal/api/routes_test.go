package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/api/internal/setup"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var mockContinents = []struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}{
	{ID: 1, Name: "Europe"},
	{ID: 2, Name: "Asia"},
}

// ctx context.Context, sql string, args ...any
func mockQuery(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
	return &mockRows{data: mockContinents}, nil
}

// ctx context.Context, sql string, args ...any
func mockExec(_ context.Context, _ string, _ ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

// ctx context.Context, sql string, args ...any
func mockQueryRow(_ context.Context, _ string, _ ...any) pgx.Row {
	return &mockRow{}
}

type mockRow struct{}

func (m *mockRow) Scan(dest ...interface{}) error {
	if len(dest) >= 1 {
		if name, ok := dest[0].(*string); ok {
			*name = "Europe"
		}
	}
	return nil
}

type mockRows struct {
	data []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	index  int
	closed bool
}

func (m *mockRows) Close() {
	m.closed = true
}

func (m *mockRows) Err() error {
	return nil
}

func (m *mockRows) CommandTag() pgconn.CommandTag {
	return pgconn.CommandTag{}
}

func (m *mockRows) FieldDescriptions() []pgconn.FieldDescription {
	return nil
}

func (m *mockRows) Next() bool {
	if m.closed {
		return false
	}
	m.index++
	return m.index <= len(m.data)
}

func (m *mockRows) Scan(dest ...interface{}) error {
	if m.index > 0 && m.index <= len(m.data) {
		continent := m.data[m.index-1]
		if len(dest) >= 2 {
			if id, ok := dest[0].(*int); ok {
				*id = continent.ID
			}
			if name, ok := dest[1].(*string); ok {
				*name = continent.Name
			}
		}
		return nil
	}
	return nil
}

func (m *mockRows) Values() ([]interface{}, error) {
	return nil, nil
}

func (m *mockRows) RawValues() [][]byte {
	return nil
}

func (m *mockRows) Conn() *pgx.Conn {
	return nil
}

type mockPgxPool struct{}

func (m *mockPgxPool) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return mockExec(ctx, sql, args...)
}

func (m *mockPgxPool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return mockQueryRow(ctx, sql, args...)
}

func (m *mockPgxPool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return mockQuery(ctx, sql, args...)
}

///////////////////////////////////////////////////////////////////////////////
// Continents - OK
///////////////////////////////////////////////////////////////////////////////

func TestCreateContinent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockDBPool := &mockPgxPool{}
	cfg := setup.GetConfig()
	cfg.GinEngine = router
	cfg.PgPool = mockDBPool

	router.POST("/api/v1/continent", createContinent(mockDBPool.QueryRow))

	body := `{"name":"Europe"}`
	req, err := http.NewRequest(http.MethodPost, "/api/v1/continent", strings.NewReader(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, but got %d", http.StatusCreated, w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
}

func TestGetContinent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockDBPool := &mockPgxPool{}
	cfg := setup.GetConfig()
	cfg.GinEngine = router
	cfg.PgPool = mockDBPool

	router.GET("/api/v1/continent/:id", getContinent(mockDBPool.QueryRow))

	req, err := http.NewRequest(http.MethodGet, "/api/v1/continent/1", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["name"] != "Europe" {
		t.Errorf("Expected name 'Europe', but got '%s'", response["name"])
	}
}

func TestGetAllContinents(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockDBPool := &mockPgxPool{}
	cfg := setup.GetConfig()
	cfg.GinEngine = router
	cfg.PgPool = mockDBPool

	router.GET("/api/v1/continents", getAllContinents(mockDBPool.Query))

	req, err := http.NewRequest(http.MethodGet, "/api/v1/continents", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	var response []map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response) != 2 {
		t.Errorf("Expected 2 continents, but got %d", len(response))
	}
	if response[0]["name"] != "Europe" {
		t.Errorf("Expected name 'Europe', but got '%s'", response[0]["name"])
	}
	if response[1]["name"] != "Asia" {
		t.Errorf("Expected name 'Asia', but got '%s'", response[1]["name"])
	}
}

func TestUpdateContinent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockDBPool := &mockPgxPool{}
	cfg := setup.GetConfig()
	cfg.GinEngine = router
	cfg.PgPool = mockDBPool

	router.PUT("/api/v1/continent/:id", updateContinent(mockDBPool.Exec))

	body := `{"name":"Updated Europe"}`
	req, err := http.NewRequest(http.MethodPut, "/api/v1/continent/1", strings.NewReader(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["status"] != "updated" {
		t.Errorf("Expected status 'updated', but got '%s'", response["status"])
	}
}

func TestDeleteContinent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockDBPool := &mockPgxPool{}
	cfg := setup.GetConfig()
	cfg.GinEngine = router
	cfg.PgPool = mockDBPool

	router.DELETE("/api/v1/continent/:id", deleteContinent(mockDBPool.Exec))

	req, err := http.NewRequest(http.MethodDelete, "/api/v1/continent/1", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["status"] != "deleted" {
		t.Errorf("Expected status 'deleted', but got '%s'", response["status"])
	}
}
