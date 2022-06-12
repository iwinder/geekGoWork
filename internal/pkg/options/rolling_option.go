package configs

type RollingOption struct {
	TimeInMilliseconds int     `yaml:"time-in-milliseconds" mapstructure:"time-in-milliseconds"`
	LimitBucket        int     `yaml:"limit-bucket" mapstructure:"limit-bucket"`
	LimitCount         int64   `yaml:"limit-count" mapstructure:"limit-count"`
	ErrorPercentage    float64 `yaml:"error-percentage" mapstructure:"error-percentage"`
	BrokenTimePeriod   int64   `yaml:"broken-time-period" mapstructure:"broken-time-period"`
}

func NewRollingOption(timeInMilliseconds, limitBucket int, limitCount, brokenTimePeriod int64, errorPercentage float64) *RollingOption {
	if timeInMilliseconds <= 0 {
		timeInMilliseconds = 1000
	}
	if limitBucket <= 0 {
		limitBucket = 10
	}
	if limitCount <= 0 {
		limitCount = 100
	}
	if errorPercentage <= 0 {
		errorPercentage = 75
	}
	if brokenTimePeriod <= 0 {
		brokenTimePeriod = 100
	}
	return &RollingOption{
		TimeInMilliseconds: timeInMilliseconds,
		LimitBucket:        limitBucket,
		LimitCount:         limitCount,
		ErrorPercentage:    errorPercentage,
		BrokenTimePeriod:   brokenTimePeriod,
	}
}
