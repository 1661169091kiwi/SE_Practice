package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAuthLogic_Flow(t *testing.T) {
	// Setup JWT Secret to match what the handlers expect
	os.Setenv("JWT_SECRET", "test-secret")

	// 1. Generate Tokens using the actual helper in auth.go (if accessible)
	// Since we are in package handler, we can access signToken.
	studentToken := signToken("student1", "student")
	collectorToken := signToken("collector1", "collector")

	// 2. Test Student Access to Collector Endpoint (Should be 403)
	t.Run("Student_Access_Collector_Data", func(t *testing.T) {
		// Method POST, and path matches what SubmitCollectorData expects
		req := httptest.NewRequest("POST", "/api/collector/football/events/1/data", nil)
		req.Header.Set("Authorization", "Bearer "+studentToken)
		w := httptest.NewRecorder()

		// Call the handler directly
		SubmitCollectorData(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected 403 Forbidden for student accessing collector data, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	// 3. Test Collector Access to Collector Endpoint (Should NOT be 403)
	t.Run("Collector_Access_Collector_Data", func(t *testing.T) {
		// We expect this to fail later at DB layer (500) or Body parsing (400)
		// But definitively NOT 403.

		req := httptest.NewRequest("POST", "/api/collector/football/events/1/data", nil)
		req.Header.Set("Authorization", "Bearer "+collectorToken)
		w := httptest.NewRecorder()

		SubmitCollectorData(w, req)

		if w.Code == http.StatusForbidden {
			t.Errorf("Expected non-403 for collector accessing collector data, got 403")
		} else {
			t.Logf("Collector access passed role check. Got code: %d (expected behavior without DB)", w.Code)
		}
	})

	// 4. Test Student Access to Ratings (Should NOT be 403 if rater_id matches)
	t.Run("Student_Access_Ratings", func(t *testing.T) {
		// Ratings handler checks AuthClaims then Body.
		// We need a valid JSON body to reach the logic.
		body := []byte(`{"student_id":"s1", "rater_id":"student1", "score":5, "comment":"good"}`)
		req := httptest.NewRequest("POST", "/api/matches/1/ratings", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+studentToken)
		w := httptest.NewRecorder()

		Ratings(w, req)

		// It checks claims.Sub ("student1") == body.RaterID ("student1").
		// Should pass auth check.
		// Then it calls service -> repo -> DB (nil) -> Error -> 500.

		if w.Code == http.StatusForbidden || w.Code == http.StatusUnauthorized {
			t.Errorf("Expected allowed access (likely 500/400) for student rating, got %d. Body: %s", w.Code, w.Body.String())
		} else {
			t.Logf("Student access to ratings passed auth check. Got code: %d (expected behavior without DB)", w.Code)
		}
	})

	// 5. Test Student Access to Ratings with Wrong RaterID (Should be 403)
	t.Run("Student_Access_Ratings_WrongID", func(t *testing.T) {
		body := []byte(`{"student_id":"s1", "rater_id":"other_student", "score":5, "comment":"good"}`)
		req := httptest.NewRequest("POST", "/api/matches/1/ratings", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+studentToken)
		w := httptest.NewRecorder()

		Ratings(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected 403 for student rating as another user, got %d. Body: %s", w.Code, w.Body.String())
		}
	})
}
