package entity

import (
	"encoding/json"
	"time"
)

// Move mirrors the moves table.
type Move struct {
	ID              string          `json:"id"`
	MatchID         string          `json:"match_id"`
	Seq             int             `json:"seq"`
	OccurredAt      time.Time       `json:"occurred_at"`
	ElapsedMs       int             `json:"elapsed_ms"`
	FromIdx         int             `json:"from_idx"`
	ToIdx           int             `json:"to_idx"`
	MoveKind        int16           `json:"move_kind"`
	FrogSide        int16           `json:"frog_side"`
	IsCorrect       bool            `json:"is_correct"`
	Interruption    bool            `json:"interruption"`
	BoardBefore     json.RawMessage `json:"board_before"`
	BoardAfter      json.RawMessage `json:"board_after"`
	BranchingFactor *int            `json:"branching_factor"`
	Buclicidad      *float64        `json:"buclicidad"`
}
