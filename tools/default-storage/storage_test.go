package defaultstorage_test

import (
	"sync"
	"testing"

	defaultstorage "github.com/chernyshevuser/practicum-metrics-collector/tools/default-storage"
)

func TestStorage_Set(t *testing.T) {
	store := defaultstorage.New[string]()

	key := uint64(1)
	val := "val"
	store.Set(key, val)

	received, exists := store.Get(key)
	if !exists {
		t.Errorf("expected key '%v' to exist", key)
	}
	if received != val {
		t.Errorf("expected value '%s', got '%v'", val, received)
	}

	_, exists = store.Get(uint64(2))
	if exists {
		t.Error("expected non-existing key to return false for exists")
	}
}

func TestStorage_Get(t *testing.T) {
	store := defaultstorage.New[string]()

	key := uint64(1)
	val := "val"
	store.Set(key, val)

	received, exists := store.Get(key)
	if !exists {
		t.Errorf("expected key '%v' to exist", key)
	}
	if received != val {
		t.Errorf("expected value '%s', got '%v'", val, received)
	}

	_, exists = store.Get(uint64(2))
	if exists {
		t.Error("expected non-existing key to return false for exists")
	}
}

func TestStorage_GetAll(t *testing.T) {
	store := defaultstorage.New[string]()

	all := store.GetAll()
	if len(all) != 0 {
		t.Errorf("expected empty slice, got length %d", len(all))
	}

	store.Set(1, "value1")
	store.Set(2, "value2")
	store.Set(3, "value3")

	all = store.GetAll()
	if len(all) != 3 {
		t.Errorf("expected 3 values, got %d", len(all))
	}

	expectedValues := map[string]bool{"value1": true, "value2": true, "value3": true}
	for _, v := range all {
		if !expectedValues[v.(string)] {
			t.Errorf("unexpected value found: %v", v)
		}
		delete(expectedValues, v.(string))
	}

	if len(expectedValues) != 0 {
		t.Error("not all expected values were returned by GetAll")
	}
}

func TestStorage_ConcurrentAccess(t *testing.T) {
	store := defaultstorage.New[int]()

	key1 := uint64(1)
	key2 := uint64(2)
	key3 := uint64(3)
	val1 := 100
	val2 := 200
	val3 := 300

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		store.Set(key1, val1)
	}()
	go func() {
		defer wg.Done()
		store.Set(key2, val2)
	}()
	go func() {
		defer wg.Done()
		store.Set(key3, val3)
	}()

	wg.Wait()

	_, exists := store.Get(key1)
	if !exists {
		t.Errorf("expected key '%v' to exist", key1)
	}
	_, exists = store.Get(key2)
	if !exists {
		t.Errorf("expected key '%v' to exist", key2)
	}
	_, exists = store.Get(key3)
	if !exists {
		t.Errorf("expected key '%v' to exist", key3)
	}
}
