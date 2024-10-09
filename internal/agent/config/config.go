package config

import (
	"flag"

	configgetter "github.com/chernyshevuser/practicum-metrics-collector/tools/configgetter"

	"github.com/chernyshevuser/practicum-metrics-collector/tools/logger"
)

const (
	AddrEnv           = configgetter.ConfigKey("ADDRESS")
	ReportIntervalEnv = configgetter.ConfigKey("REPORT_INTERVAL")
	PollIntervalEnv   = configgetter.ConfigKey("POLL_INTERVAL")
	HashKeyEnv        = configgetter.ConfigKey("KEY")
	RateLimitEnv      = configgetter.ConfigKey("RATE_LIMIT")
	FixedIVStrEnv     = configgetter.ConfigKey("SYPHER")
)

var (
	Addr           string
	ReportInterval int64
	PollInterval   int64
	HashKey        string
	RateLimit      int64
	FixedIVStr     string
)

func Setup(logger logger.Logger) {
	flag.StringVar(&Addr, "a", "localhost:8080", "addr")
	flag.Int64Var(&ReportInterval, "r", 10, "report")
	flag.Int64Var(&PollInterval, "p", 2, "poll")
	flag.StringVar(&HashKey, "k", "", "hash key")
	flag.Int64Var(&RateLimit, "l", 2, "rate limit")
	flag.StringVar(&FixedIVStr, "S", "1234567890123456", "sypher")

	flag.Parse()

	addr, err := configgetter.GetConfigString(AddrEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		Addr = addr
	}

	reportInterval, err := configgetter.GetConfigInt64(ReportIntervalEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		ReportInterval = reportInterval
	}

	pollInterval, err := configgetter.GetConfigInt64(PollIntervalEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		PollInterval = pollInterval
	}

	hashKey, err := configgetter.GetConfigString(HashKeyEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		HashKey = hashKey
	}

	rateLimit, err := configgetter.GetConfigInt64(RateLimitEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		RateLimit = rateLimit
	}

	fixedIVStr, err := configgetter.GetConfigString(FixedIVStrEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		FixedIVStr = fixedIVStr
	}
}
