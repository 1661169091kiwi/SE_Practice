package repo

import (
	"se_practice/backend/internal/db"
	"se_practice/backend/internal/model"
)

type RatingRepo struct{}

func NewRatingRepo() *RatingRepo {
	return &RatingRepo{}
}

func (r *RatingRepo) Upsert(matchID int64, studentID string, raterID string, score float64, comment string) error {
	_, err := db.Exec(`INSERT INTO player_ratings (match_id, student_id, rater_id, score, comment)
VALUES (?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE score=VALUES(score), comment=VALUES(comment)`,
		matchID, studentID, raterID, score, comment)
	return err
}

func (r *RatingRepo) ListByMatch(matchID int64) ([]model.PlayerRating, error) {
	rows, err := db.Query(`SELECT rating_id, match_id, student_id, rater_id, score, comment, created_at
FROM player_ratings WHERE match_id = ? ORDER BY created_at DESC`, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []model.PlayerRating
	for rows.Next() {
		var pr model.PlayerRating
		if err := rows.Scan(&pr.RatingID, &pr.MatchID, &pr.StudentID, &pr.RaterID, &pr.Score, &pr.Comment, &pr.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, pr)
	}
	return list, rows.Err()
}
