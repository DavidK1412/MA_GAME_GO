package entity

import (
	"encoding/json"
	"time"
)

// Match mirrors the matches table.
type Match struct {
	ID           string          `json:"id"`
	SessionID    string          `json:"session_id"`
	DifficultyID int             `json:"difficulty_id"`
	LevelN       int             `json:"level_n"`
	IsActive     bool            `json:"is_active"`
	StartedAt    time.Time       `json:"started_at"`
	EndedAt      *time.Time      `json:"ended_at"`
	Outcome      *string         `json:"outcome"`
	Meta         json.RawMessage `json:"meta"`
}
