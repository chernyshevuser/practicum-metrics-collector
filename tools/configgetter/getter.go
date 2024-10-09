package configgetter

import (
	"fmt"
	"os"
	"strconv"
)

type ConfigKey string

func GetConfigString(key ConfigKey) (string, error) {
	val := os.Getenv(string(key))
	if len(val) == 0 {
		return "", fmt.Errorf("error in getting %s value", key)
	}

	return val, nil
}

func GetConfigInt64(key ConfigKey) (int64, error) {
	valStr := os.Getenv(string(key))
	if len(valStr) == 0 {
		return 0, fmt.Errorf("error in getting %s value", key)
	}

	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("can't convert %s to int64, reason: %v", valStr, err)
	}

	return val, nil
}

func GetConfigFloat64(key ConfigKey) (float64, error) {
	valStr := os.Getenv(string(key))
	if len(valStr) == 0 {
		return 0, fmt.Errorf("error in getting %s value", key)
	}

	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return 0, fmt.Errorf("can't convert %s to int64, reason: %v", valStr, err)
	}

	return val, nil
}

func GetConfigBool(key ConfigKey) (bool, error) {
	valStr := os.Getenv(string(key))
	if len(valStr) == 0 {
		return false, fmt.Errorf("error in getting %s value", key)
	}

	val, err := strconv.ParseBool(valStr)
	if err != nil {
		return false, fmt.Errorf("can't convert %s to bool, reason: %v", valStr, err)
	}

	return val, nil
}
