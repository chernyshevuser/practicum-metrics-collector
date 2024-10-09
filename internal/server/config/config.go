package config

import (
	"flag"

	configgetter "github.com/chernyshevuser/practicum-metrics-collector/tools/configgetter"
	"github.com/chernyshevuser/practicum-metrics-collector/tools/logger"
)

const (
	AddrEnv            = configgetter.ConfigKey("ADDRESS")
	StoreIntervalEnv   = configgetter.ConfigKey("STORE_INTERVAL")
	FileStoragePathEnv = configgetter.ConfigKey("FILE_STORAGE_PATH")
	RestoreEnv         = configgetter.ConfigKey("RESTORE")
	DatabaseDsnEnv     = configgetter.ConfigKey("DATABASE_DSN")
	HashKeyEnv         = configgetter.ConfigKey("KEY")
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

	addr, err := configgetter.GetConfigString(AddrEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		Addr = addr
	}

	storeInterval, err := configgetter.GetConfigInt64(StoreIntervalEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		StoreInterval = storeInterval
	}

	fileStoragePath, err := configgetter.GetConfigString(FileStoragePathEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		FileStoragePath = fileStoragePath
	}

	restore, err := configgetter.GetConfigBool(RestoreEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		Restore = restore
	}

	databaseDsn, err := configgetter.GetConfigString(DatabaseDsnEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		DatabaseDsn = databaseDsn
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
