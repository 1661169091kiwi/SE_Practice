package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"se_practice/backend/internal/db"
	"se_practice/backend/internal/router"
)

// Reuse APIResponse and LoginData from e2e_collector_test.go
// but for clarity I'll redefine them here to keep the test self-contained if needed,
// or I can rely on them if they were exported in the same package.
// Since they are in the same package `tests`, I can reuse them if they are exported.
// However, in e2e_collector_test.go they are defined as types, so they are accessible.

func TestAdminBusinessLoop(t *testing.T) {
	if dsn == "" {
		t.Skip("TEST_DB_DSN or DB_DSN is required to run database-backed tests")
	}
	// 1. Initialize DB
	if _, err := db.Init(dsn); err != nil {
		t.Skipf("Skip DB-backed test: %v", err)
	}

	// Seed Data for Admin
	if err := seedAdminData(); err != nil {
		t.Fatalf("Failed to seed admin data: %v", err)
	}

	// Setup Server
	r := router.New()
	server := httptest.NewServer(r)
	defer server.Close()
	client := server.Client()

	t.Log("Starting Admin E2E Test...")

	// 2. Login Admin
	t.Log("Step 2: Login Admin")
	loginBody := map[string]string{
		"student_id": "admin001",
		"password":   "password",
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
	adminToken := loginData.Token
	if adminToken == "" {
		t.Fatalf("Admin Token is empty")
	}
	t.Logf("Got Admin Token.")

	// 3. Create Team
	t.Log("Step 3: Create Team")
	// We need two teams for a match
	teamAId := createTeam(t, client, server.URL, adminToken, "Admin Team A")
	teamBId := createTeam(t, client, server.URL, adminToken, "Admin Team B")
	t.Logf("Created Teams: %d, %d", teamAId, teamBId)

	// 4. Create Event
	t.Log("Step 4: Create Event")
	eventPayload := map[string]interface{}{
		"event_name":  "Admin Test Event",
		"sport_id":    1, // Football from seed
		"format_type": "points",
		"start_date":  "2024-06-01T00:00:00Z",
		"end_date":    "2024-06-30T00:00:00Z",
		"status":      "ongoing",
	}
	eventBytes, _ := json.Marshal(eventPayload)
	req, _ := http.NewRequest("POST", server.URL+"/api/events/create", bytes.NewReader(eventBytes))
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Create Event failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Create Event status: %d, body: %s", resp.StatusCode, string(body))
	}

	// Assuming Create Event returns the created event or ID in data.
	// Let's decode to get ID if possible, or we query it.
	// Based on API conventions, it might return just success or the ID.
	// Let's assume standard response wrapper.
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		t.Fatalf("Failed to decode create event response: %v", err)
	}

	// Try to get Event ID. If not returned, query by name.
	// It's possible the API returns just "id" or the whole object.
	// Let's check DB directly to be safe and simple if parsing is complex.
	var eventID int64
	conn := db.GetDB()
	err = conn.QueryRow("SELECT event_id FROM events WHERE event_name = 'Admin Test Event' ORDER BY event_id DESC LIMIT 1").Scan(&eventID)
	if err != nil {
		t.Fatalf("Failed to find created event in DB: %v", err)
	}
	t.Logf("Created Event ID: %d", eventID)

	// 5. Create Match
	t.Log("Step 5: Create Match")
	matchPayload := map[string]interface{}{
		"event_id":   eventID,
		"match_name": "Admin Test Match",
		"match_time": "2024-06-02T14:00:00Z",
		"team_a_id":  teamAId,
		"team_b_id":  teamBId,
	}
	matchBytes, _ := json.Marshal(matchPayload)
	req, _ = http.NewRequest("POST", server.URL+"/api/matches/create", bytes.NewReader(matchBytes))
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Create Match failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Create Match status: %d, body: %s", resp.StatusCode, string(body))
	}
	t.Log("Match Created.")

	// 6. Verify Match
	t.Log("Step 6: Verify Match")
	verifyURL := fmt.Sprintf("%s/api/matches?event_id=%d", server.URL, eventID)
	req, _ = http.NewRequest("GET", verifyURL, nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Verify Match request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("Verify Match status: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		t.Fatalf("Failed to decode verify response: %v", err)
	}

	var matches []Match
	if err := json.Unmarshal(apiResp.Data, &matches); err != nil {
		t.Fatalf("Failed to parse matches list: %v", err)
	}

	found := false
	for _, m := range matches {
		if m.Name == "Admin Test Match" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("Created match not found in the list")
	}
	t.Log("Match verified in list.")

	// 6.5 Verify ListTeams (For Admin Dashboard Dropdown)
	t.Log("Step 6.5: Verify ListTeams")
	req, _ = http.NewRequest("GET", server.URL+"/api/teams/list", nil)
	// ListTeams is public or authenticated? Code says:
	// func ListTeams(w http.ResponseWriter, r *http.Request) { ... }
	// Router says: mux.HandleFunc("/api/teams/list", handler.ListTeams) -> It's public based on router definition!
	// But let's send token anyway to be safe or if middleware changes.
	req.Header.Set("Authorization", "Bearer "+adminToken)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("ListTeams request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("ListTeams status: %d, body: %s", resp.StatusCode, string(body))
	}
	t.Log("ListTeams endpoint works.")

	// 6.6 Verify ListEvents (For Admin Dashboard Dropdown)
	t.Log("Step 6.6: Verify ListEvents")
	req, _ = http.NewRequest("GET", server.URL+"/api/events/list", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("ListEvents request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("ListEvents status: %d, body: %s", resp.StatusCode, string(body))
	}
	t.Log("ListEvents endpoint works.")

	// 7. Permission Check (Student)
	t.Log("Step 7: Permission Check")
	// Login as student (created in seed)
	studentLoginBody := map[string]string{
		"student_id": "student001",
		"password":   "password",
	}
	studentLoginBytes, _ := json.Marshal(studentLoginBody)
	resp, err = client.Post(server.URL+"/api/login", "application/json", bytes.NewReader(studentLoginBytes))
	if err != nil {
		t.Fatalf("Student login failed: %v", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		t.Fatalf("Failed to decode student login: %v", err)
	}
	var studentData LoginData
	json.Unmarshal(apiResp.Data, &studentData)
	studentToken := studentData.Token

	// Try to create team with student token
	req, _ = http.NewRequest("POST", server.URL+"/api/teams/create", bytes.NewReader(matchBytes)) // Payload doesn't matter much if auth fails
	req.Header.Set("Authorization", "Bearer "+studentToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Permission check request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 403 {
		t.Fatalf("Expected 403 Forbidden for student creating team, got %d", resp.StatusCode)
	}
	t.Log("Permission check passed: Student cannot create team.")
}

func createTeam(t *testing.T, client *http.Client, serverURL, token, name string) int64 {
	payload := map[string]interface{}{
		"team_name":   name,
		"sport_id":    1, // Football
		"college":     "Test College",
		"team_type":   "college",
		"description": "Test Team",
		"created_by":  "admin001",
	}
	data, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", serverURL+"/api/teams/create", bytes.NewReader(data))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Create Team request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Create Team status: %d, body: %s", resp.StatusCode, string(body))
	}

	// Query DB for ID
	var teamID int64
	conn := db.GetDB()
	err = conn.QueryRow("SELECT team_id FROM teams WHERE team_name = ? ORDER BY team_id DESC LIMIT 1", name).Scan(&teamID)
	if err != nil {
		t.Fatalf("Failed to retrieve created team ID: %v", err)
	}
	return teamID
}

func seedAdminData() error {
	conn := db.GetDB()

	// 1. Ensure Sport (Football) exists
	_, err := conn.Exec(`INSERT INTO sports (sport_id, sport_name, description) VALUES (1, 'Football', 'Standard 11-a-side Football') ON DUPLICATE KEY UPDATE description=VALUES(description)`)
	if err != nil {
		return err
	}

	hashedPwd := hashPassword("password")

	// 2. Create Admin User
	_, err = conn.Exec(`INSERT INTO users (student_id, password, name, college, grade) VALUES ('admin001', ?, 'Admin User', 'Admin Dept', '2020') ON DUPLICATE KEY UPDATE password=VALUES(password)`, hashedPwd)
	if err != nil {
		return err
	}

	// 3. Grant Admin Rights
	_, err = conn.Exec(`INSERT INTO admins (student_id, role, permissions) VALUES ('admin001', 'super_admin', 'all') ON DUPLICATE KEY UPDATE role=VALUES(role)`)
	if err != nil {
		return err
	}

	// 4. Create Student User
	_, err = conn.Exec(`INSERT INTO users (student_id, password, name, college, grade) VALUES ('student001', ?, 'Student User', 'Student Dept', '2023') ON DUPLICATE KEY UPDATE password=VALUES(password)`, hashedPwd)
	if err != nil {
		return err
	}

	return nil
}
