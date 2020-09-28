package config

var RecoverMode bool

type VelocityLimitConfig struct {
	DayLimit          float64
	WeekLimit         float64
	MaxAttemptsPerDay uint
}
