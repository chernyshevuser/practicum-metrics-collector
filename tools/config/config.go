package config

import (
	"flag"

	"github.com/chernyshevuser/practicum-metrics-collector/tools/logger"
)

type configKey string

const (
	RunAddrEnv           = configKey("RUN_ADDRESS")
	DatabaseURIEnv       = configKey("DATABASE_URI")
	AccrualSystemAddrEnv = configKey("ACCRUAL_SYSTEM_ADDRESS")
	CryptoKeyEnv         = configKey("CRYPTO_KEY")
	JwtSecretKeyEnv      = configKey("JWT_KEY")
)

var (
	RunAddr           string
	DatabaseURI       string
	AccrualSystemAddr string
	CryptoKey         string
	JwtSecretKey      string
)

func Setup(logger logger.Logger) {
	flag.StringVar(&RunAddr, "a", "localhost:8080", "runAddr")
	flag.StringVar(&DatabaseURI, "d", "", "dbUri")
	flag.StringVar(&AccrualSystemAddr, "r", "", "accrual system addr")
	flag.StringVar(&CryptoKey, "k", "examplekey123456", "crypto key")
	flag.StringVar(&JwtSecretKey, "j", "examplekey123456", "jwt crypto key")

	flag.Parse()

	runAddr, err := GetConfigString(RunAddrEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		RunAddr = runAddr
	}

	databaseURI, err := GetConfigString(DatabaseURIEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		DatabaseURI = databaseURI
	}

	accrualSystemAddr, err := GetConfigString(AccrualSystemAddrEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		AccrualSystemAddr = accrualSystemAddr
	}

	cryptoKey, err := GetConfigString(CryptoKeyEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		CryptoKey = cryptoKey
	}

	jwtSecretKey, err := GetConfigString(JwtSecretKeyEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		JwtSecretKey = jwtSecretKey
	}

	//TODO printing envs lol
	logger.Infow(
		"config",
		"runAddr", RunAddr,
		"dbUri", DatabaseURI,
		"accrual system addr", AccrualSystemAddr,
		"crypto key", CryptoKey,
		"jwt crypto key", JwtSecretKey,
	)
}
