package model

import "time"

type PlayerRating struct {
	RatingID  int64     `json:"rating_id"`
	MatchID   int64     `json:"match_id"`
	StudentID string    `json:"student_id"`
	RaterID   string    `json:"rater_id"`
	Score     float64   `json:"score"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}
