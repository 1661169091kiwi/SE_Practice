package repo

import (
	"database/sql"
	"se_practice/backend/internal/db"
	"se_practice/backend/internal/model"
)

type EventRepo struct{}

func NewEventRepo() *EventRepo {
	return &EventRepo{}
}

// CreateEvent 创建赛事
func (r *EventRepo) CreateEvent(event *model.Event) (int64, error) {
	query := `INSERT INTO events 
	          (event_name, sport_id, season, round, format_type, start_date, end_date, status) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	return db.Insert(query,
		event.Name,
		event.SportID,
		event.Season,
		event.Round,
		event.Format,
		event.StartDate,
		event.EndDate,
		event.Status,
	)
}

// GetEventByID 根据ID获取赛事
func (r *EventRepo) GetEventByID(id int64) (*model.Event, error) {
	query := `SELECT event_id, event_name, sport_id, season, round, format_type, 
	                 start_date, end_date, status, created_at 
	          FROM events WHERE event_id = ?`

	var event model.Event
	err := db.QueryRow(query, id).Scan(
		&event.ID,
		&event.Name,
		&event.SportID,
		&event.Season,
		&event.Round,
		&event.Format,
		&event.StartDate,
		&event.EndDate,
		&event.Status,
		&event.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &event, err
}

// ListEvents 列出赛事
func (r *EventRepo) ListEvents(status string) ([]model.Event, error) {
	query := `SELECT event_id, event_name, sport_id, season, round, format_type, 
	                 start_date, end_date, status, created_at 
	          FROM events`

	var args []interface{}
	if status != "" {
		query += " WHERE status = ?"
		args = append(args, status)
	}
	query += " ORDER BY start_date DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []model.Event
	for rows.Next() {
		var e model.Event
		if err := rows.Scan(
			&e.ID, &e.Name, &e.SportID, &e.Season, &e.Round,
			&e.Format, &e.StartDate, &e.EndDate, &e.Status, &e.CreatedAt,
		); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

// UpdateEventStatus 更新赛事状态
func (r *EventRepo) UpdateEventStatus(id int64, status string) error {
	query := "UPDATE events SET status = ? WHERE event_id = ?"
	_, err := db.Exec(query, status, id)
	return err
}
