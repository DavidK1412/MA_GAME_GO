package domain

type Belief struct {
	ID          string             `json:"id"`
	AgentID     string             `json:"agent_id"`
	MatchID     string             `json:"match_id"`
	BeliefState map[string]float64 `json:"belief_state"`
	Confidence  float64            `json:"confidence"`
}
