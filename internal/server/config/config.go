package config

import (
	"flag"

	"github.com/chernyshevuser/practicum-metrics-collector/tools/logger"
)

type configKey string

const (
	AddrEnv            = configKey("ADDRESS")
	StoreIntervalEnv   = configKey("STORE_INTERVAL")
	FileStoragePathEnv = configKey("FILE_STORAGE_PATH")
	RestoreEnv         = configKey("RESTORE")
	DatabaseDsnEnv     = configKey("DATABASE_DSN")
	HashKeyEnv         = configKey("KEY")
)

var (
	Addr            string
	StoreInterval   int64
	FileStoragePath string
	Restore         bool
	DatabaseDsn     string
	HashKey         string
)

func Setup(logger logger.Logger) {
	flag.StringVar(&Addr, "a", "localhost:8080", "server addr")
	flag.Int64Var(&StoreInterval, "i", 300, "store interval")
	flag.StringVar(&FileStoragePath, "f", "", "file storage path")
	flag.BoolVar(&Restore, "r", true, "restore flag")
	flag.StringVar(&DatabaseDsn, "d", "", "database data source name")
	flag.StringVar(&HashKey, "k", "", "hash key")

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

	storeInterval, err := GetConfigInt64(StoreIntervalEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		StoreInterval = storeInterval
	}

	fileStoragePath, err := GetConfigString(FileStoragePathEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		FileStoragePath = fileStoragePath
	}

	restore, err := GetConfigBool(RestoreEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		Restore = restore
	}

	databaseDsn, err := GetConfigString(DatabaseDsnEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		DatabaseDsn = databaseDsn
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

	logger.Infow(
		"envs",
		"addr", Addr,
		"storeInterval", StoreInterval,
		"fileStoragePath", FileStoragePath,
		"restore", Restore,
		"databaseDsn", DatabaseDsn,
		"hashKey", HashKey,
	)
}
