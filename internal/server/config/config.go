package config

import (
	"flag"

	getter "github.com/chernyshevuser/practicum-metrics-collector/tools/config-getter"
	"github.com/chernyshevuser/practicum-metrics-collector/tools/crypto"
	"github.com/chernyshevuser/practicum-metrics-collector/tools/logger"
)

type configKey string

const (
	AddrEnv            = getter.ConfigKey("ADDRESS")
	StoreIntervalEnv   = getter.ConfigKey("STORE_INTERVAL")
	FileStoragePathEnv = getter.ConfigKey("FILE_STORAGE_PATH")
	RestoreEnv         = getter.ConfigKey("RESTORE")
	DatabaseDsnEnv     = getter.ConfigKey("DATABASE_DSN")
	HashKeyEnv         = getter.ConfigKey("KEY")
	CryptoKeyPathEnv   = getter.ConfigKey("CRYPTO_KEY")
)

var (
	Addr            string
	StoreInterval   int64
	FileStoragePath string
	Restore         bool
	DatabaseDsn     string
	HashKey         string
	CryptoKeyPath   string
	CryptoKey       string
)

func Setup(logger logger.Logger) {
	flag.StringVar(&Addr, "a", "localhost:8080", "server addr")
	flag.Int64Var(&StoreInterval, "i", 300, "store interval")
	flag.StringVar(&FileStoragePath, "f", "", "file storage path")
	flag.BoolVar(&Restore, "r", true, "restore flag")
	flag.StringVar(&DatabaseDsn, "d", "", "database data source name")
	flag.StringVar(&HashKey, "k", "", "hash key")
	flag.StringVar(&CryptoKeyPath, "crypto-key", "", "fpath with crypto key")

	flag.Parse()

	addr, err := getter.GetConfigString(AddrEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		Addr = addr
	}

	storeInterval, err := getter.GetConfigInt64(StoreIntervalEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		StoreInterval = storeInterval
	}

	fileStoragePath, err := getter.GetConfigString(FileStoragePathEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		FileStoragePath = fileStoragePath
	}

	restore, err := getter.GetConfigBool(RestoreEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		Restore = restore
	}

	databaseDsn, err := getter.GetConfigString(DatabaseDsnEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		DatabaseDsn = databaseDsn
	}

	hashKey, err := getter.GetConfigString(HashKeyEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		HashKey = hashKey
	}

	cryptoKeyPath, err := getter.GetConfigString(CryptoKeyPathEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		CryptoKeyPath = cryptoKeyPath
	}

	if CryptoKeyPath != "" {
		CryptoKey, err = crypto.LoadFromFile(CryptoKeyPath)
		if err != nil {
			logger.Errorw("can't parse file with crypto key", "msg", err)
		}
	}
}
