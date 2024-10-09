package configgetter_test

import (
	"os"
	"strconv"
	"testing"

	"github.com/chernyshevuser/practicum-metrics-collector/tools/configgetter"
)

func TestGetConfigString(t *testing.T) {
	key := configgetter.ConfigKey("TEST_STRING")

	t.Run("successful string retrieval", func(t *testing.T) {
		expected := "test_value"
		os.Setenv(string(key), expected)
		defer os.Unsetenv(string(key))

		val, err := configgetter.GetConfigString(key)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if val != expected {
			t.Errorf("expected %s, got %s", expected, val)
		}
	})

	t.Run("empty string should return error", func(t *testing.T) {
		os.Unsetenv(string(key))

		_, err := configgetter.GetConfigString(key)
		if err == nil {
			t.Fatal("expected error, got none")
		}
	})
}

func TestGetConfigInt64(t *testing.T) {
	key := configgetter.ConfigKey("TEST_INT")

	t.Run("successful int64 retrieval", func(t *testing.T) {
		expected := int64(123456789)
		os.Setenv(string(key), strconv.FormatInt(expected, 10))
		defer os.Unsetenv(string(key))

		val, err := configgetter.GetConfigInt64(key)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if val != expected {
			t.Errorf("expected %d, got %d", expected, val)
		}
	})

	t.Run("invalid int64 format should return error", func(t *testing.T) {
		os.Setenv(string(key), "invalid_int")
		defer os.Unsetenv(string(key))

		_, err := configgetter.GetConfigInt64(key)
		if err == nil {
			t.Fatal("expected error, got none")
		}
	})

	t.Run("empty int64 should return error", func(t *testing.T) {
		os.Unsetenv(string(key))

		_, err := configgetter.GetConfigInt64(key)
		if err == nil {
			t.Fatal("expected error, got none")
		}
	})
}

func TestGetConfigFloat64(t *testing.T) {
	key := configgetter.ConfigKey("TEST_FLOAT")

	t.Run("successful float64 retrieval", func(t *testing.T) {
		expected := 123.456
		os.Setenv(string(key), strconv.FormatFloat(expected, 'f', -1, 64))
		defer os.Unsetenv(string(key))

		val, err := configgetter.GetConfigFloat64(key)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if val != expected {
			t.Errorf("expected %f, got %f", expected, val)
		}
	})

	t.Run("invalid float64 format should return error", func(t *testing.T) {
		os.Setenv(string(key), "invalid_float")
		defer os.Unsetenv(string(key))

		_, err := configgetter.GetConfigFloat64(key)
		if err == nil {
			t.Fatal("expected error, got none")
		}
	})

	t.Run("empty float64 should return error", func(t *testing.T) {
		os.Unsetenv(string(key))

		_, err := configgetter.GetConfigFloat64(key)
		if err == nil {
			t.Fatal("expected error, got none")
		}
	})
}

func TestGetConfigBool(t *testing.T) {
	key := configgetter.ConfigKey("TEST_BOOL")

	t.Run("successful bool retrieval", func(t *testing.T) {
		expected := true
		os.Setenv(string(key), strconv.FormatBool(expected))
		defer os.Unsetenv(string(key))

		val, err := configgetter.GetConfigBool(key)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if val != expected {
			t.Errorf("expected %t, got %t", expected, val)
		}
	})

	t.Run("invalid bool format should return error", func(t *testing.T) {
		os.Setenv(string(key), "invalid_bool")
		defer os.Unsetenv(string(key))

		_, err := configgetter.GetConfigBool(key)
		if err == nil {
			t.Fatal("expected error, got none")
		}
	})

	t.Run("empty bool should return error", func(t *testing.T) {
		os.Unsetenv(string(key))

		_, err := configgetter.GetConfigBool(key)
		if err == nil {
			t.Fatal("expected error, got none")
		}
	})
}
