package entity

// Difficulty represents the difficulty table described in seed.sql.
type Difficulty struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	NumberOfBlocks int    `json:"number_of_blocks"`
}
