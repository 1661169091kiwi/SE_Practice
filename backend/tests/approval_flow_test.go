package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"se_practice/backend/internal/db"
	"se_practice/backend/internal/model"
	"se_practice/backend/internal/router"
)

func TestApprovalFlow(t *testing.T) {
	if dsn == "" {
		t.Skip("TEST_DB_DSN or DB_DSN is required to run database-backed tests")
	}
	// 1. Initialize DB
	if _, err := db.Init(dsn); err != nil {
		t.Skipf("Skip DB-backed test: %v", err)
	}
	// Initialize Schema to ensure new columns exist
	if err := db.InitSchema(); err != nil {
		t.Fatalf("Failed to init schema: %v", err)
	}

	// 2. Setup Test Data
	adminID := "admin_test_001"
	userID := "user_test_001"
	sportID := int64(0)
	approvedTeamID := int64(0)
	pendingTeamID := int64(0)

	// Clean up previous test data if any (optional, but good for idempotency)
	db.Exec("DELETE FROM team_members WHERE athlete_id IN (SELECT athlete_id FROM athletes WHERE student_id = ?)", userID)
	db.Exec("DELETE FROM athletes WHERE student_id = ?", userID)
	db.Exec("DELETE FROM teams WHERE created_by = ?", userID)
	// We don't delete users to avoid FK issues if they are used elsewhere, but we can ensure they exist

	// SHA256 of "123456"
	hashedPwd := "8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92"

	// Ensure Admin User Exists
	_, err := db.Exec("INSERT INTO users (student_id, password, name, college, grade) VALUES (?, ?, 'Test Admin', '', '') ON DUPLICATE KEY UPDATE name=VALUES(name), password=VALUES(password), college=VALUES(college), grade=VALUES(grade)", adminID, hashedPwd)
	if err != nil {
		t.Fatalf("Failed to create admin user: %v", err)
	}
	// Add to admins table
	_, err = db.Exec("INSERT INTO admins (student_id, role, permissions) VALUES (?, 'admin', '') ON DUPLICATE KEY UPDATE role=VALUES(role), permissions=VALUES(permissions)", adminID)
	if err != nil {
		t.Fatalf("Failed to add admin role: %v", err)
	}

	// Ensure Normal User Exists
	_, err = db.Exec("INSERT INTO users (student_id, password, name, college, grade) VALUES (?, ?, 'Test User', '', '') ON DUPLICATE KEY UPDATE name=VALUES(name), password=VALUES(password), college=VALUES(college), grade=VALUES(grade)", userID, hashedPwd)
	if err != nil {
		t.Fatalf("Failed to create normal user: %v", err)
	}

	// Ensure Sport Exists
	res, err := db.Exec("INSERT INTO sports (sport_name, description) VALUES ('Test Sport', 'For Testing') ON DUPLICATE KEY UPDATE description=VALUES(description)")
	if err != nil {
		t.Fatalf("Failed to seed sport: %v", err)
	}
	// Get Sport ID (simplest way if we don't know if it was inserted or updated)
	var sid int64
	err = db.QueryRow("SELECT sport_id FROM sports WHERE sport_name = 'Test Sport' LIMIT 1").Scan(&sid)
	if err != nil {
		t.Fatalf("Failed to get sport id: %v", err)
	}
	sportID = sid

	// Create an Approved Team (for athlete application)
	res, err = db.Exec("INSERT INTO teams (team_name, sport_id, college, team_type, created_by, is_approved) VALUES (?, ?, 'Test College', 'college', ?, TRUE)", "Approved Team "+time.Now().Format("150405"), sportID, adminID)
	if err != nil {
		t.Fatalf("Failed to create approved team: %v", err)
	}
	approvedTeamID, _ = res.LastInsertId()

	// Create a Pending Team (manually or via API later, let's do manually to setup state for approval test)
	res, err = db.Exec("INSERT INTO teams (team_name, sport_id, college, team_type, created_by, is_approved) VALUES (?, ?, 'Test College', 'college', ?, FALSE)", "Pending Team "+time.Now().Format("150405"), sportID, userID)
	if err != nil {
		t.Fatalf("Failed to create pending team: %v", err)
	}
	pendingTeamID, _ = res.LastInsertId()

	// Create a Pending Athlete Application
	// First create athlete record
	res, err = db.Exec("INSERT INTO athletes (student_id, sport_type, team_id) VALUES (?, 'Test Sport', ?)", userID, approvedTeamID)
	if err != nil {
		t.Fatalf("Failed to create athlete: %v", err)
	}
	athleteID, _ := res.LastInsertId()
	// Then create pending team member record
	res, err = db.Exec("INSERT INTO team_members (team_id, athlete_id, join_date, is_active, is_approved) VALUES (?, ?, CURDATE(), FALSE, FALSE)", approvedTeamID, athleteID)
	if err != nil {
		t.Fatalf("Failed to create pending team member: %v", err)
	}
	teamMemberID, _ := res.LastInsertId()

	// 3. Setup Server
	r := router.New()
	server := httptest.NewServer(r)
	defer server.Close()
	client := server.Client()

	// 4. Login as Admin
	loginBody, _ := json.Marshal(map[string]string{
		"student_id": adminID,
		"password":   "123456",
	})
	resp, err := client.Post(server.URL+"/api/login", "application/json", bytes.NewReader(loginBody))
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	var loginResp struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	// Read body for debugging
	bodyBytes, _ := io.ReadAll(resp.Body)
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	t.Logf("Login response: %s", string(bodyBytes))

	json.NewDecoder(resp.Body).Decode(&loginResp)
	token := loginResp.Data.Token
	if token == "" {
		t.Fatalf("Failed to get admin token")
	}

	// 5. Test: Get Pending Athletes
	t.Log("Testing Get Pending Athletes...")
	req, _ := http.NewRequest("GET", server.URL+"/api/admin/athletes/pending", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Get pending athletes failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
	var pendingAthletes struct {
		Code int                               `json:"code"`
		Data []model.PendingAthleteApplication `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&pendingAthletes)
	foundAthlete := false
	for _, app := range pendingAthletes.Data {
		if app.TeamMemberID == teamMemberID {
			foundAthlete = true
			break
		}
	}
	if !foundAthlete {
		t.Errorf("Created pending athlete application not found in list")
	}

	// 6. Test: Approve Athlete
	t.Log("Testing Approve Athlete...")
	approveAthleteBody, _ := json.Marshal(map[string]int64{
		"team_member_id": teamMemberID,
	})
	req, _ = http.NewRequest("POST", server.URL+"/api/admin/athletes/approve", bytes.NewReader(approveAthleteBody))
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Approve athlete failed: %v", err)
	}
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		t.Errorf("Expected 200, got %d. Body: %s", resp.StatusCode, string(body))
	}

	// Verify DB state
	var isApproved bool
	err = db.QueryRow("SELECT is_approved FROM team_members WHERE team_member_id = ?", teamMemberID).Scan(&isApproved)
	if err != nil {
		t.Fatalf("Failed to query team member status: %v", err)
	}
	if !isApproved {
		t.Errorf("Athlete should be approved, but is_approved is false")
	}

	// 7. Test: Get Pending Teams
	t.Log("Testing Get Pending Teams...")
	req, _ = http.NewRequest("GET", server.URL+"/api/admin/teams/pending", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Get pending teams failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
	var pendingTeams struct {
		Code int          `json:"code"`
		Data []model.Team `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&pendingTeams)
	foundTeam := false
	for _, team := range pendingTeams.Data {
		if team.ID == pendingTeamID {
			foundTeam = true
			break
		}
	}
	if !foundTeam {
		t.Errorf("Created pending team not found in list")
	}

	// 8. Test: Approve Team
	t.Log("Testing Approve Team...")
	approveTeamBody, _ := json.Marshal(map[string]int64{
		"team_id": pendingTeamID,
	})
	req, _ = http.NewRequest("POST", server.URL+"/api/admin/teams/approve", bytes.NewReader(approveTeamBody))
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Approve team failed: %v", err)
	}
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		t.Errorf("Expected 200, got %d. Body: %s", resp.StatusCode, string(body))
	}

	// Verify DB state
	err = db.QueryRow("SELECT is_approved FROM teams WHERE team_id = ?", pendingTeamID).Scan(&isApproved)
	if err != nil {
		t.Fatalf("Failed to query team status: %v", err)
	}
	if !isApproved {
		t.Errorf("Team should be approved, but is_approved is false")
	}

	t.Log("Approval Flow Test Completed Successfully")
}
