package domain

import "time"

type Move struct {
	ID        string    `json:"id"`
	MatchID   string    `json:"match_id"`
	PlayerID  string    `json:"player_id"`
	Position  string    `json:"position"`
	MoveType  string    `json:"move_type"`
	Timestamp time.Time `json:"timestamp"`
}
