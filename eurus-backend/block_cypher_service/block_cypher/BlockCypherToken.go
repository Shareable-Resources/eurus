package block_cypher

import "eurus-backend/foundation/database"

type BlockCypherToken struct {
	database.DbModel
	Email                  string
	Token                  string
	Coin                   string
	Chain                  string
	IsEnabled              bool
	Score                  int
	UsedCount              int
	HitsApiPerHour         int
	HitsApiPerDay          int
	HitsConfidencePerHour  int
	LimitApiPerHour        int
	LimitApiPerDay         int
	LimitConfidencePerHour int
}
