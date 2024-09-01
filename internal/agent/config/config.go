package config

import (
	"flag"

	"github.com/chernyshevuser/practicum-metrics-collector/tools/logger"
)

type configKey string

const (
	AddrEnv           = configKey("ADDRESS")
	ReportIntervalEnv = configKey("REPORT_INTERVAL")
	PollIntervalEnv   = configKey("POLL_INTERVAL")
	HashKeyEnv        = configKey("KEY")
	RateLimitEnv      = configKey("RATE_LIMIT")
)

var (
	Addr           string
	ReportInterval int64
	PollInterval   int64
	HashKey        string
	RateLimit      int64
)

func Setup(logger logger.Logger) {
	flag.StringVar(&Addr, "a", "localhost:8080", "addr")
	flag.Int64Var(&ReportInterval, "r", 10, "report")
	flag.Int64Var(&PollInterval, "p", 2, "poll")
	flag.StringVar(&HashKey, "k", "", "hash key")
	flag.Int64Var(&RateLimit, "l", 2, "rate limit")

	flag.Parse()

	addr, err := GetConfigString(AddrEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		Addr = addr
	}

	reportInterval, err := GetConfigInt64(ReportIntervalEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		ReportInterval = reportInterval
	}

	pollInterval, err := GetConfigInt64(PollIntervalEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		PollInterval = pollInterval
	}

	hashKey, err := GetConfigString(HashKeyEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		HashKey = hashKey
	}

	rateLimit, err := GetConfigInt64(RateLimitEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		RateLimit = rateLimit
	}
}
