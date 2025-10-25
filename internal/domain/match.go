package domain

import "time"

type Match struct {
	ID         string     `json:"id"`
	PlayerID   string     `json:"player_id"`
	OpponentID string     `json:"opponent_id"`
	Status     string     `json:"status"`
	Difficulty Difficulty `json:"difficulty"`
	StartedAt  time.Time  `json:"started_at"`
	EndedAt    *time.Time `json:"ended_at,omitempty"`
	Winner     *string    `json:"winner,omitempty"`
}
