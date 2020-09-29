package config

var AppendMode bool

type Config struct {
	velocityLimitConfig
	Currency string
}
type velocityLimitConfig struct {
	DayLimit          float64 `json:"dailyLimit"`
	WeekLimit         float64 `json:"weeklyLimit"`
	MaxAttemptsPerDay uint    `json:"maxAttemptsPerDay"`
}

var Configuration Config
