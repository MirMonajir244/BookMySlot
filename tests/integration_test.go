package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MirMonajir244/BookMySlot/internal/router"
	"github.com/gin-gonic/gin"
)

// TestHealthCheck is a simple integration test for the health endpoint
func TestHealthCheck(t *testing.T) {
	// 1. Setup Gin in test mode
	gin.SetMode(gin.TestMode)
	
	// 2. Initialize router with nil handlers for health check test
	r := router.Setup(nil, nil, nil, nil)

	// 3. Create a request
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// 4. Perform request
	r.ServeHTTP(w, req)

	// 5. Assertions
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("expected status 'ok', got '%s'", response["status"])
	}
}

// NOTE: A full integration test would involve:
// 1. Starting a real PostgreSQL container (or using an in-memory DB like SQLite)
// 2. Wiring up the real Repository/Service/Handler chain
// 3. Using httptest to simulate high-level API calls (Register -> Login -> Book)
// 4. Asserting on the database state after the API call.
