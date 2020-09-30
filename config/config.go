package config

var AppendMode bool

type Config struct {
	VelocityLimitConfig
	Currency string
}
type VelocityLimitConfig struct {
	DayLimit          float64 `json:"dailyLimit"`
	WeekLimit         float64 `json:"weeklyLimit"`
	MaxAttemptsPerDay uint    `json:"maxAttemptsPerDay"`
}

var Configuration Config
