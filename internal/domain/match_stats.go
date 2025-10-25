package domain

type MatchStats struct {
	MatchID      string  `json:"match_id"`
	TotalMoves   int     `json:"total_moves"`
	Duration     int64   `json:"duration"`
	PlayerScore  int     `json:"player_score"`
	OpponentScore int    `json:"opponent_score"`
	Accuracy     float64 `json:"accuracy"`
}
