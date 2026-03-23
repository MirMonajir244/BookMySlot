package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBookingFlow_Integration(t *testing.T) {
	r, _, _ := SetupTestRouter()

	// 1. Register Coach & Get Token
	coachToken := registerAndGetToken(t, r, "Coach X", "coach_x@test.com", "coach")
	
	// 2. Set Availability
	availBody := map[string]interface{}{
		"coach_id":    1,
		"day_of_week": "Monday",
		"start_time":  "09:00",
		"end_time":    "11:00",
	}
	body, _ := json.Marshal(availBody)
	req, _ := http.NewRequest("POST", "/api/v1/coaches/availability", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+coachToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 for availability, got %d. Body: %s", w.Code, w.Body.String())
	}

	// 3. Register User & Get Token
	userToken := registerAndGetToken(t, r, "User Y", "user_y@test.com", "user")

	// 4. Get Available Slots
	req, _ = http.NewRequest("GET", "/api/v1/users/slots?coach_id=1&date=2025-10-27", nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for slots, got %d", w.Code)
	}

	var slotsResp struct {
		Count int `json:"count"`
	}
	json.Unmarshal(w.Body.Bytes(), &slotsResp)
	if slotsResp.Count == 0 {
		t.Error("expected slots to be generated")
	}

	// 5. Book a Slot
	bookingBody := map[string]interface{}{
		"user_id":  1,
		"coach_id": 1,
		"datetime": "2025-10-27T09:00:00Z",
	}
	body, _ = json.Marshal(bookingBody)
	req, _ = http.NewRequest("POST", "/api/v1/users/bookings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+userToken)
	req.Header.Set("Idempotency-Key", "unique-id-123")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 for booking, got %d. Body: %s", w.Code, w.Body.String())
	}

	// 6. Double Book (Should Fail)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/users/bookings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+userToken)
	req.Header.Set("Idempotency-Key", "different-id-456")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("expected 409 for double booking, got %d", w.Code)
	}
}

// Helper to dry up tests
func registerAndGetToken(t *testing.T, r http.Handler, name, email, role string) string {
	regBody := map[string]string{
		"name":     name,
		"email":    email,
		"password": "password123",
		"role":     role,
	}
	body, _ := json.Marshal(regBody)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("failed to register %s: %d", role, w.Code)
	}

	var resp struct {
		Token string `json:"token"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	return resp.Token
}
