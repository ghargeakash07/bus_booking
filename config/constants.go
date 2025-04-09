// constants.go
package config

const (
	DefaultMarketStartTime = "09:30"
	DefaultMarketEndTime   = "23:59"
)

var (
	MarketStartTime = getEnv("MARKET_START_TIME", DefaultMarketStartTime)
	MarketEndTime   = getEnv("MARKET_END_TIME", DefaultMarketEndTime)
)


