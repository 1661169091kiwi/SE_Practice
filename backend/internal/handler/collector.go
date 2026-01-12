package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"se_practice/backend/internal/db"
	"se_practice/backend/internal/middleware"
	"se_practice/backend/internal/model"
	"se_practice/backend/internal/repo"
	"se_practice/backend/internal/service"
	"se_practice/backend/internal/util"
)

func GetCollectorMatchDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/collector/"), "/")
	if len(parts) < 3 || parts[1] != "events" {
		util.Error(w, http.StatusBadRequest, "invalid url format")
		return
	}

	sportType := parts[0]
	matchIDStr := parts[2]
	matchID, err := strconv.ParseInt(matchIDStr, 10, 64)
	if err != nil {
		util.Error(w, http.StatusBadRequest, "invalid match id")
		return
	}

	ms := repo.NewMatchRepo()
	match, err := ms.GetMatchByID(matchID)
	if err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if match == nil {
		util.Error(w, http.StatusNotFound, "match not found")
		return
	}

	ur := repo.NewUserRepo()
	teamA, _ := ur.GetTeamByID(match.TeamAID)
	teamB, _ := ur.GetTeamByID(match.TeamBID)

	teamAName := "TBD"
	if teamA != nil {
		teamAName = teamA.TeamName
	}
	teamBName := "TBD"
	if teamB != nil {
		teamBName = teamB.TeamName
	}

	tr := repo.NewTeamRepo()
	membersA, _ := tr.GetTeamMembers(match.TeamAID)
	membersB, _ := tr.GetTeamMembers(match.TeamBID)
	var candidatesA []model.PlayerBrief
	for _, m := range membersA {
		candidatesA = append(candidatesA, model.PlayerBrief{
			ID:       m.StudentID,
			Name:     m.Name,
			Number:   m.JerseyNumber,
			Position: "",
		})
	}
	var candidatesB []model.PlayerBrief
	for _, m := range membersB {
		candidatesB = append(candidatesB, model.PlayerBrief{
			ID:       m.StudentID,
			Name:     m.Name,
			Number:   m.JerseyNumber,
			Position: "",
		})
	}
	autoFinishAt := match.Time.Add(72 * time.Hour).Format(time.RFC3339)

	// Fetch available teams for the event
	er := repo.NewEventRepo()
	eventTeams, _ := er.ListEventTeams(match.EventID)
	var availableTeams []model.CollectorTeamBrief
	for _, et := range eventTeams {
		availableTeams = append(availableTeams, model.CollectorTeamBrief{
			ID:   et.TeamID,
			Name: et.TeamName,
		})
	}

	resp := model.CollectorMatchDetailResponse{
		ID:              match.ID,
		EventName:       match.EventName,
		Name:            match.Name,
		TeamA:           teamAName,
		TeamB:           teamBName,
		TeamAID:         match.TeamAID,
		TeamBID:         match.TeamBID,
		Time:            match.Time.Format("2006-01-02 15:04"),
		Venue:           "主体育场",
		SportType:       sportType,
		Status:          match.Status,
		AutoFinishAt:    autoFinishAt,
		CandidatesTeamA: candidatesA,
		CandidatesTeamB: candidatesB,
		AvailableTeams:  availableTeams,
		CurrentData: model.CollectorMatchData{
			Scores: model.MatchScores{
				TeamA: match.ScoreA,
				TeamB: match.ScoreB,
			},
			Lineups: func() model.MatchLineups {
				var teamAL []model.PlayerBrief
				var teamBL []model.PlayerBrief
				seenA := map[string]struct{}{}
				seenB := map[string]struct{}{}
				lus, err := ms.GetMatchLineups(match.ID)
				if err == nil && len(lus) > 0 {
					for _, lu := range lus {
						p := model.PlayerBrief{
							ID:         strings.TrimSpace(lu.StudentID),
							Name:       strings.TrimSpace(lu.PlayerName),
							Number:     strings.TrimSpace(lu.JerseyNumber),
							Position:   strings.TrimSpace(lu.Position),
							IsStarting: lu.IsStarting,
						}
						var key string
						if strings.TrimSpace(lu.StudentID) != "" {
							key = "rid:" + strings.TrimSpace(lu.StudentID)
						} else {
							key = "ext:" + strings.ToLower(strings.TrimSpace(lu.PlayerName)) + "|" + strings.TrimSpace(lu.JerseyNumber) + "|" + strings.ToLower(strings.TrimSpace(lu.Position))
						}
						if lu.TeamID == match.TeamAID {
							if _, ok := seenA[key]; !ok {
								seenA[key] = struct{}{}
								teamAL = append(teamAL, p)
							}
						} else if lu.TeamID == match.TeamBID {
							if _, ok := seenB[key]; !ok {
								seenB[key] = struct{}{}
								teamBL = append(teamBL, p)
							}
						}
					}
				}
				return model.MatchLineups{
					TeamA: teamAL,
					TeamB: teamBL,
				}
			}(),
		},
	}
	util.OK(w, resp)
}

func SubmitCollectorData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	// 仅 POST 采集数据需要采集员角色
	if !middleware.CheckRole(w, r, "collector") {
		return
	}

	p := strings.TrimPrefix(r.URL.Path, "/api/collector/")
	parts := strings.Split(p, "/")
	if len(parts) != 4 || parts[1] != "events" || parts[3] != "data" || parts[0] == "" || parts[2] == "" {
		util.Error(w, http.StatusBadRequest, "invalid path")
		return
	}

	sportType := parts[0]
	matchIDStr := parts[2]
	matchID, err := strconv.ParseInt(matchIDStr, 10, 64)
	if err != nil {
		util.Error(w, http.StatusBadRequest, "invalid matchId")
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		util.Error(w, http.StatusBadRequest, "failed to read request body")
		return
	}
	// We don't close r.Body here because the server does it, but we can if we want.

	var envScoresOnly struct {
		EventID       interface{} `json:"eventId"`
		SportType     string      `json:"sportType"`
		Timestamp     string      `json:"timestamp"`
		CollectorName string      `json:"collectorName"`
		Status        string      `json:"status"`
		Data          struct {
			Scores model.MatchScores `json:"scores"`
		} `json:"data"`
		EventInfo struct {
			TeamAId interface{} `json:"teamAId"`
			TeamBId interface{} `json:"teamBId"`
		} `json:"eventInfo"`
	}
	var envLineupsLoose struct {
		Data struct {
			Lineups struct {
				TeamA []struct {
					ID         interface{} `json:"id"`
					Name       string      `json:"name"`
					Position   string      `json:"position"`
					Number     string      `json:"number"`
					IsStarting bool        `json:"isStarting"`
				} `json:"teamA"`
				TeamB []struct {
					ID         interface{} `json:"id"`
					Name       string      `json:"name"`
					Position   string      `json:"position"`
					Number     string      `json:"number"`
					IsStarting bool        `json:"isStarting"`
				} `json:"teamB"`
			} `json:"lineups"`
		} `json:"data"`
	}
	var simple model.CollectorSubmitDataRequest

	log.Printf("[Collector] Received submit data: %s", string(bodyBytes))
	errScores := json.Unmarshal(bodyBytes, &envScoresOnly)
	useEnv := errScores == nil && envScoresOnly.SportType != ""
	log.Printf("[Collector] useEnv: %v", useEnv)
	if !useEnv {
		if err := json.Unmarshal(bodyBytes, &simple); err != nil {
			util.Error(w, http.StatusBadRequest, "invalid request body")
			return
		}
		if simple.ScoreTeamA < 0 || simple.ScoreTeamB < 0 {
			util.Error(w, http.StatusBadRequest, "scores cannot be negative")
			return
		}
	} else if matchID < 0 {
		// 检查是否包含队伍信息，如果是 TBD 比赛且有队伍信息，则实例化比赛
		taIDStr := toString(envScoresOnly.EventInfo.TeamAId)
		tbIDStr := toString(envScoresOnly.EventInfo.TeamBId)
		if taIDStr != "" && tbIDStr != "" {
			taID, _ := strconv.ParseInt(taIDStr, 10, 64)
			tbID, _ := strconv.ParseInt(tbIDStr, 10, 64)
			if taID > 0 && tbID > 0 {
				log.Printf("[Collector] Instantiating knockout match %d with teams %d vs %d", matchID, taID, tbID)
				mrepo := repo.NewMatchRepo()
				realMatchID, err := mrepo.InstantiateKnockoutMatch(-matchID, taID, tbID)
				if err != nil {
					log.Printf("[Collector] Failed to instantiate match: %v", err)
					// 不报错，尝试继续（虽然可能会失败）
				} else {
					log.Printf("[Collector] Match instantiated as ID: %d", realMatchID)
					matchID = realMatchID
				}
			}
		}
	}

	mrepo := repo.NewMatchRepo()
	match, err := mrepo.GetMatchByID(matchID)
	if err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if match == nil {
		util.Error(w, http.StatusNotFound, "match not found")
		return
	}

	tx, err := db.BeginTransaction()
	if err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	scoreA := 0
	scoreB := 0
	if useEnv {
		scoreA = envScoresOnly.Data.Scores.TeamA
		scoreB = envScoresOnly.Data.Scores.TeamB
	} else {
		scoreA = int(simple.ScoreTeamA)
		scoreB = int(simple.ScoreTeamB)
	}
	if _, err = tx.Exec("UPDATE matches SET score_team_a = ?, score_team_b = ? WHERE match_id = ?", scoreA, scoreB, matchID); err != nil {
		_ = tx.Rollback()
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 更新队伍ID（如果有）
	if useEnv {
		taIDStr := toString(envScoresOnly.EventInfo.TeamAId)
		tbIDStr := toString(envScoresOnly.EventInfo.TeamBId)
		if taIDStr != "" && tbIDStr != "" {
			taID, _ := strconv.ParseInt(taIDStr, 10, 64)
			tbID, _ := strconv.ParseInt(tbIDStr, 10, 64)
			if taID > 0 && tbID > 0 {
				if _, err = tx.Exec("UPDATE matches SET team_a_id = ?, team_b_id = ? WHERE match_id = ?", taID, tbID, matchID); err != nil {
					_ = tx.Rollback()
					util.Error(w, http.StatusInternalServerError, err.Error())
					return
				}
				match.TeamAID = taID
				match.TeamBID = tbID
			}
		}
	}

	// 如果传入了状态且不为空，更新状态
	if useEnv && envScoresOnly.Status != "" {
		statusVal := envScoresOnly.Status
		if statusVal == "in_progress" {
			statusVal = "ongoing"
		}
		shouldUpdate := true
		if match.Status == "finished" && statusVal != "finished" {
			shouldUpdate = false
		}
		if shouldUpdate {
			if _, err = tx.Exec("UPDATE matches SET status = ? WHERE match_id = ?", statusVal, matchID); err != nil {
				_ = tx.Rollback()
				util.Error(w, http.StatusInternalServerError, err.Error())
				return
			}
			match.Status = statusVal
		}
	}

	lrepo := repo.NewLineupRepo()
	ur := repo.NewUserRepo()
	upsert := func(items []model.LineupItem, teamID int64) error {
		for _, it := range items {
			sid := strings.TrimSpace(it.StudentID)
			if sid == "" {
				log.Printf("[Collector] lineup external (empty sid) name=%s teamID=%d", strings.TrimSpace(it.Name), teamID)
				if strings.TrimSpace(it.Name) == "" {
					return fmt.Errorf("bad_request: player name required")
				}
				if err := lrepo.UpsertExternalTx(tx, matchID, teamID, strings.TrimSpace(it.Name), it.Position, it.JerseyNumber, it.IsStarting); err != nil {
					return err
				}
				continue
			}
			exists, _ := ur.UserExists(sid)
			if !exists {
				log.Printf("[Collector] lineup external (user not found) sid=%s name=%s teamID=%d", sid, strings.TrimSpace(it.Name), teamID)
				if strings.TrimSpace(it.Name) == "" {
					return fmt.Errorf("bad_request: player name required")
				}
				if err := lrepo.UpsertExternalTx(tx, matchID, teamID, strings.TrimSpace(it.Name), it.Position, it.JerseyNumber, it.IsStarting); err != nil {
					return err
				}
				continue
			}
			log.Printf("[Collector] lineup update match_lineups sid=%s teamID=%d", sid, teamID)
			res, err := tx.Exec("UPDATE match_lineups SET position = ?, jersey_number = ?, is_starting = ? WHERE match_id = ? AND team_id = ? AND student_id = ?", it.Position, it.JerseyNumber, it.IsStarting, matchID, teamID, sid)
			if err != nil {
				if db.IsForeignKeyError(err) {
					log.Printf("[Collector] lineup update FK error sid=%s -> external", sid)
					if err := lrepo.UpsertExternalTx(tx, matchID, teamID, strings.TrimSpace(it.Name), it.Position, it.JerseyNumber, it.IsStarting); err != nil {
						return err
					}
					continue
				}
				return err
			}
			rows, _ := res.RowsAffected()
			if rows == 0 {
				log.Printf("[Collector] lineup insert match_lineups sid=%s teamID=%d", sid, teamID)
				if _, err = tx.Exec("INSERT INTO match_lineups (match_id, team_id, student_id, position, is_starting, jersey_number) VALUES (?, ?, ?, ?, ?, ?)", matchID, teamID, sid, it.Position, it.IsStarting, it.JerseyNumber); err != nil {
					if db.IsForeignKeyError(err) {
						log.Printf("[Collector] lineup insert FK error sid=%s -> external", sid)
						if err := lrepo.UpsertExternalTx(tx, matchID, teamID, strings.TrimSpace(it.Name), it.Position, it.JerseyNumber, it.IsStarting); err != nil {
							return err
						}
						continue
					}
					return err
				}
			}
		}
		return nil
	}

	var teamAItems []model.LineupItem
	var teamBItems []model.LineupItem
	if useEnv {
		_ = json.Unmarshal(bodyBytes, &envLineupsLoose)
		for _, p := range envLineupsLoose.Data.Lineups.TeamA {
			teamAItems = append(teamAItems, model.LineupItem{
				StudentID:    toString(p.ID),
				Name:         p.Name,
				Position:     p.Position,
				JerseyNumber: p.Number,
				IsStarting:   p.IsStarting,
			})
		}
		for _, p := range envLineupsLoose.Data.Lineups.TeamB {
			teamBItems = append(teamBItems, model.LineupItem{
				StudentID:    toString(p.ID),
				Name:         p.Name,
				Position:     p.Position,
				JerseyNumber: p.Number,
				IsStarting:   p.IsStarting,
			})
		}
	} else {
		teamAItems = simple.TeamA
		teamBItems = simple.TeamB
	}
	normKey := func(it model.LineupItem) string {
		sid := strings.TrimSpace(it.StudentID)
		if sid != "" {
			return "rid:" + sid
		}
		return "ext:" + strings.ToLower(strings.TrimSpace(it.Name)) + "|" + strings.TrimSpace(it.JerseyNumber) + "|" + strings.ToLower(strings.TrimSpace(it.Position))
	}
	dedup := func(items []model.LineupItem) []model.LineupItem {
		seen := map[string]struct{}{}
		var out []model.LineupItem
		for _, it := range items {
			k := normKey(it)
			if _, ok := seen[k]; ok {
				continue
			}
			seen[k] = struct{}{}
			out = append(out, it)
		}
		return out
	}
	teamAItems = dedup(teamAItems)
	teamBItems = dedup(teamBItems)

	if err = upsert(teamAItems, match.TeamAID); err != nil {
		_ = tx.Rollback()
		msg := err.Error()
		if strings.HasPrefix(msg, "bad_request:") {
			util.Error(w, http.StatusBadRequest, strings.TrimPrefix(msg, "bad_request: "))
		} else {
			util.Error(w, http.StatusInternalServerError, msg)
		}
		return
	}
	if err = upsert(teamBItems, match.TeamBID); err != nil {
		_ = tx.Rollback()
		msg := err.Error()
		if strings.HasPrefix(msg, "bad_request:") {
			util.Error(w, http.StatusBadRequest, strings.TrimPrefix(msg, "bad_request: "))
		} else {
			util.Error(w, http.StatusInternalServerError, msg)
		}
		return
	}

	// 删除未在本次提交中的旧阵容（支持删除）
	ms := repo.NewMatchRepo()
	existing, err := ms.GetMatchLineups(matchID)
	if err != nil {
		_ = tx.Rollback()
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	presentA := map[string]struct{}{}
	presentB := map[string]struct{}{}
	for _, it := range teamAItems {
		sid := strings.TrimSpace(it.StudentID)
		if sid == "" {
			k := "ext:" + strings.ToLower(strings.TrimSpace(it.Name)) + "|" + strings.TrimSpace(it.JerseyNumber) + "|" + strings.ToLower(strings.TrimSpace(it.Position))
			presentA[k] = struct{}{}
			continue
		}
		exists, _ := ur.UserExists(sid)
		if exists {
			presentA["rid:"+sid] = struct{}{}
		} else {
			k := "ext:" + strings.ToLower(strings.TrimSpace(it.Name)) + "|" + strings.TrimSpace(it.JerseyNumber) + "|" + strings.ToLower(strings.TrimSpace(it.Position))
			presentA[k] = struct{}{}
		}
	}
	for _, it := range teamBItems {
		sid := strings.TrimSpace(it.StudentID)
		if sid == "" {
			k := "ext:" + strings.ToLower(strings.TrimSpace(it.Name)) + "|" + strings.TrimSpace(it.JerseyNumber) + "|" + strings.ToLower(strings.TrimSpace(it.Position))
			presentB[k] = struct{}{}
			continue
		}
		exists, _ := ur.UserExists(sid)
		if exists {
			presentB["rid:"+sid] = struct{}{}
		} else {
			k := "ext:" + strings.ToLower(strings.TrimSpace(it.Name)) + "|" + strings.TrimSpace(it.JerseyNumber) + "|" + strings.ToLower(strings.TrimSpace(it.Position))
			presentB[k] = struct{}{}
		}
	}
	for _, lu := range existing {
		var key string
		if strings.TrimSpace(lu.StudentID) != "" {
			key = "rid:" + strings.TrimSpace(lu.StudentID)
		} else {
			key = "ext:" + strings.ToLower(strings.TrimSpace(lu.PlayerName)) + "|" + strings.TrimSpace(lu.JerseyNumber) + "|" + strings.ToLower(strings.TrimSpace(lu.Position))
		}
		if lu.TeamID == match.TeamAID {
			if _, ok := presentA[key]; !ok {
				if strings.TrimSpace(lu.StudentID) != "" {
					if err := lrepo.DeleteTx(tx, matchID, lu.TeamID, strings.TrimSpace(lu.StudentID)); err != nil {
						_ = tx.Rollback()
						util.Error(w, http.StatusInternalServerError, err.Error())
						return
					}
				} else {
					if err := lrepo.DeleteExternalTx(tx, matchID, lu.TeamID, strings.TrimSpace(lu.PlayerName), strings.TrimSpace(lu.Position), strings.TrimSpace(lu.JerseyNumber)); err != nil {
						_ = tx.Rollback()
						util.Error(w, http.StatusInternalServerError, err.Error())
						return
					}
				}
			}
		} else if lu.TeamID == match.TeamBID {
			if _, ok := presentB[key]; !ok {
				if strings.TrimSpace(lu.StudentID) != "" {
					if err := lrepo.DeleteTx(tx, matchID, lu.TeamID, strings.TrimSpace(lu.StudentID)); err != nil {
						_ = tx.Rollback()
						util.Error(w, http.StatusInternalServerError, err.Error())
						return
					}
				} else {
					if err := lrepo.DeleteExternalTx(tx, matchID, lu.TeamID, strings.TrimSpace(lu.PlayerName), strings.TrimSpace(lu.Position), strings.TrimSpace(lu.JerseyNumber)); err != nil {
						_ = tx.Rollback()
						util.Error(w, http.StatusInternalServerError, err.Error())
						return
					}
				}
			}
		}
	}

	if err = tx.Commit(); err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 触发积分榜重算 (如果是已结束比赛)
	if match.Status == "finished" {
		go func() {
			ss := service.NewStandingsService()
			_ = ss.RecalculateStandings(match.EventID)
		}()
	}

	util.JSON(w, http.StatusOK, util.APIResponse{Code: 200, Message: "ok", Data: map[string]interface{}{
		"sportType":       sportType,
		"matchId":         matchID,
		"updatedScore":    true,
		"lineupsUpserted": len(teamAItems) + len(teamBItems),
	}})
}

func CollectorEntry(w http.ResponseWriter, r *http.Request) {
	p := strings.TrimPrefix(r.URL.Path, "/api/collector/")
	ps := strings.Split(p, "/")
	if len(ps) >= 3 && ps[1] == "events" {
		if r.Method == http.MethodGet && len(ps) == 3 {
			GetCollectorMatchDetail(w, r)
			return
		}
		if r.Method == http.MethodPost && len(ps) == 4 && ps[3] == "data" {
			SubmitCollectorData(w, r)
			return
		}
	}
	util.Error(w, http.StatusNotFound, "not found")
}

func toString(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		if t == float64(int64(t)) {
			return strconv.FormatInt(int64(t), 10)
		}
		return fmt.Sprintf("%v", t)
	case int:
		return strconv.Itoa(t)
	case int64:
		return strconv.FormatInt(t, 10)
	default:
		return fmt.Sprintf("%v", v)
	}
}
