package logger

// defaultLevel is the default logger level
const defaultLevel = "debug"

// Config is the logger configuration struct
type Config struct {
	Level string
}

// GetLevel returns logger level
func (c *Config) GetLevel() string {
	if c.Level == "" {
		return defaultLevel
	}

	return c.Level
}