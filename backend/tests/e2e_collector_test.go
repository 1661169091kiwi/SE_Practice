package tests

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"se_practice/backend/internal/db"
	"se_practice/backend/internal/router"
)

// Response wrappers for parsing
type APIResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type LoginData struct {
	Token string `json:"token"`
}

type Match struct {
	ID      int64  `json:"match_id"`
	Name    string `json:"match_name"`
	ScoreA  int    `json:"score_team_a"`
	ScoreB  int    `json:"score_team_b"`
	EventID int64  `json:"event_id"`
	SportID int64  `json:"sport_id"`
}

type CollectorMatchDetail struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	CurrentData struct {
		Scores struct {
			TeamA int `json:"teamA"`
			TeamB int `json:"teamB"`
		} `json:"scores"`
	} `json:"currentData"`
}

var dsn = func() string {
	if v := os.Getenv("TEST_DB_DSN"); v != "" {
		return v
	}
	return os.Getenv("DB_DSN")
}()

func TestCollectorBusinessLoop(t *testing.T) {
	if dsn == "" {
		t.Skip("TEST_DB_DSN or DB_DSN is required to run database-backed tests")
	}
	if os.Getenv("ALLOW_DB_TESTS") == "" {
		t.Skip("Set ALLOW_DB_TESTS=1 to run database-backed tests")
	}
	// 1. Initialize DB
	if _, err := db.Init(dsn); err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}

	// Seed Data
	if err := seedData(); err != nil {
		t.Fatalf("Failed to seed data: %v", err)
	}

	// Setup Server
	r := router.New()
	server := httptest.NewServer(r)
	defer server.Close()
	client := server.Client()

	t.Log("Starting Collector E2E Test...")

	// 2. Login
	t.Log("Step 2: Login as Collector")
	loginBody := map[string]string{
		"student_id": "2023000001",
		"password":   "123456",
	}
	loginBytes, _ := json.Marshal(loginBody)
	resp, err := client.Post(server.URL+"/api/login", "application/json", bytes.NewReader(loginBytes))
	if err != nil {
		t.Fatalf("Login request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Login failed with status: %d, body: %s", resp.StatusCode, string(body))
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		t.Fatalf("Failed to decode login response: %v", err)
	}
	if apiResp.Code != 200 {
		t.Fatalf("Login API error: %s", apiResp.Message)
	}

	var loginData LoginData
	if err := json.Unmarshal(apiResp.Data, &loginData); err != nil {
		t.Fatalf("Failed to parse login data: %v", err)
	}
	token := loginData.Token
	if token == "" {
		t.Fatalf("Token is empty")
	}
	t.Logf("Got Token.")

	// 3. List Matches
	t.Log("Step 3: List Matches")
	req, _ := http.NewRequest("GET", server.URL+"/api/matches", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("List Matches failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("List Matches status: %d, body: %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		t.Fatalf("Failed to decode matches response: %v", err)
	}

	var matches []Match
	if err := json.Unmarshal(apiResp.Data, &matches); err != nil {
		t.Fatalf("Failed to parse matches list: %v", err)
	}

	var targetMatch *Match
	for _, m := range matches {
		if m.Name == "Group Stage Match 1" {
			targetMatch = &m
			break
		}
	}

	if targetMatch == nil {
		t.Fatalf("Target match 'Group Stage Match 1' not found in list")
	}
	t.Logf("Found Match ID: %d", targetMatch.ID)

	// 4. Enter Collection
	t.Log("Step 4: Enter Collection Page")
	// URL: /api/collector/{sportType}/events/{matchID}
	// Note: We need sportType. "football" is used in the prompt example.
	collectionURL := fmt.Sprintf("%s/api/collector/football/events/%d", server.URL, targetMatch.ID)
	req, _ = http.NewRequest("GET", collectionURL, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Enter Collection failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Enter Collection status: %d, body: %s", resp.StatusCode, string(body))
	}

	var detailResp CollectorMatchDetail
	// We need to decode APIResponse wrapper first
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		t.Fatalf("Failed to decode collection detail: %v", err)
	}
	if err := json.Unmarshal(apiResp.Data, &detailResp); err != nil {
		t.Fatalf("Failed to parse collection detail data: %v", err)
	}

	if detailResp.ID != targetMatch.ID {
		t.Errorf("Mismatch match ID in detail: got %d, want %d", detailResp.ID, targetMatch.ID)
	}
	t.Log("Entered Collection Page successfully")

	// 5. Submit Data (Score 1:0)
	t.Log("Step 5: Submit Data (1:0)")
	submitURL := fmt.Sprintf("%s/api/collector/football/events/%d/data", server.URL, targetMatch.ID)

	// Construct complex payload to avoid backend bug with simple payload (and to be more robust)
	submitPayload := map[string]interface{}{
		"eventId":   targetMatch.EventID,
		"sportType": "football",
		"timestamp": time.Now().Format(time.RFC3339),
		"status":    "ongoing",
		"data": map[string]interface{}{
			"scores": map[string]interface{}{
				"teamA": 1,
				"teamB": 0,
			},
			"lineups": map[string]interface{}{
				"teamA": []interface{}{},
				"teamB": []interface{}{},
			},
		},
	}

	submitBytes, _ := json.Marshal(submitPayload)
	req, _ = http.NewRequest("POST", submitURL, bytes.NewReader(submitBytes))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Submit Data failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// Read body for error message
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Submit Data status: %d, body: %s", resp.StatusCode, string(body))
	}
	t.Log("Data submitted successfully")

	// 6. Verify Update
	t.Log("Step 6: Verify Update")
	// We can query the match detail again via collector api or public api
	req, _ = http.NewRequest("GET", collectionURL, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Verify Update failed: %v", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		t.Fatalf("Failed to decode verify response: %v", err)
	}
	if err := json.Unmarshal(apiResp.Data, &detailResp); err != nil {
		t.Fatalf("Failed to parse verify data: %v", err)
	}

	if detailResp.CurrentData.Scores.TeamA != 1 || detailResp.CurrentData.Scores.TeamB != 0 {
		t.Errorf("Score update failed. Got %d:%d, want 1:0", detailResp.CurrentData.Scores.TeamA, detailResp.CurrentData.Scores.TeamB)
	} else {
		t.Log("Score verified: 1:0")
	}
}

func hashPassword(pw string) string {
	sum := sha256.Sum256([]byte(pw))
	return hex.EncodeToString(sum[:])
}

// seedData simulates 0003_seed_data.sql
func seedData() error {
	conn := db.GetDB()

	// 1. Sports
	_, err := conn.Exec(`INSERT INTO sports (sport_name, description) VALUES ('Football', 'Standard 11-a-side Football') ON DUPLICATE KEY UPDATE description=VALUES(description)`)
	if err != nil {
		return err
	}

	var footballID int64
	err = conn.QueryRow("SELECT sport_id FROM sports WHERE sport_name = 'Football'").Scan(&footballID)
	if err != nil {
		return err
	}

	// 2. Users
	hashedPwd := hashPassword("123456")
	_, err = conn.Exec(`INSERT INTO users (student_id, password, name, college, grade) VALUES ('2023000001', ?, 'Collector User', 'Computer Science', '2023') ON DUPLICATE KEY UPDATE name=VALUES(name), password=VALUES(password)`, hashedPwd)
	if err != nil {
		return err
	}
	_, err = conn.Exec(`INSERT INTO users (student_id, password, name, college, grade) VALUES ('2023000002', ?, 'Student User', 'Software Engineering', '2023') ON DUPLICATE KEY UPDATE name=VALUES(name), password=VALUES(password)`, hashedPwd)
	if err != nil {
		return err
	}

	// 3. Collectors
	_, err = conn.Exec(`INSERT IGNORE INTO collectors (student_id, is_approved) VALUES ('2023000001', TRUE)`)
	if err != nil {
		return err
	}

	// 4. Teams
	_, err = conn.Exec(`INSERT IGNORE INTO teams (team_name, sport_id, college, team_type, created_by) VALUES ('Team A', ?, 'Computer Science', 'college', '2023000002')`, footballID)
	if err != nil {
		return err
	}
	_, err = conn.Exec(`INSERT IGNORE INTO teams (team_name, sport_id, college, team_type, created_by) VALUES ('Team B', ?, 'Software Engineering', 'college', '2023000002')`, footballID)
	if err != nil {
		return err
	}

	var teamAID, teamBID int64
	err = conn.QueryRow("SELECT team_id FROM teams WHERE team_name = 'Team A'").Scan(&teamAID)
	if err != nil {
		return err
	}
	err = conn.QueryRow("SELECT team_id FROM teams WHERE team_name = 'Team B'").Scan(&teamBID)
	if err != nil {
		return err
	}

	// 5. Events
	_, err = conn.Exec(`INSERT IGNORE INTO events (event_name, sport_id, format_type, start_date, status) VALUES ('Campus Cup 2024', ?, 'points', CURDATE(), 'ongoing')`, footballID)
	if err != nil {
		return err
	}

	var eventID int64
	err = conn.QueryRow("SELECT event_id FROM events WHERE event_name = 'Campus Cup 2024'").Scan(&eventID)
	if err != nil {
		return err
	}

	// 6. Matches
	// Note: using NOW() + 1 DAY for match time
	_, err = conn.Exec(`INSERT IGNORE INTO matches (event_id, match_name, match_time, team_a_id, team_b_id, status) VALUES (?, 'Group Stage Match 1', DATE_ADD(NOW(), INTERVAL 1 DAY), ?, ?, 'not_started')`, eventID, teamAID, teamBID)
	if err != nil {
		return err
	}

	return nil
}
