package entity

import "time"

// Session mirrors the sessions table.
type Session struct {
	ID         string     `json:"id"`
	PlayerID   *string    `json:"player_id"`
	Device     *string    `json:"device"`
	IsFinished bool       `json:"is_finished"`
	StartedAt  time.Time  `json:"started_at"`
	EndedAt    *time.Time `json:"ended_at"`
}
