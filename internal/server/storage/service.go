// Package storage provides an interface for managing metric data storage
// and operations like setting, retrieving, and actualizing metrics.
package storage

import (
	"context"
)

// Metric represents a structure for storing a metric.
// Each metric contains an ID, Type, Val for storing a float value, and Delta for storing an integer delta.
type Metric struct {
	ID    string  // ID is the identifier for the metric.
	Type  string  // Type indicates the type of the metric (e.g., gauge, counter).
	Val   float64 // Val represents the value of the metric for float-based metrics.
	Delta int64   // Delta represents the delta value for integer-based metrics (e.g., counters).
}

// BuildKey creates a unique hash key based on the metric's name and type.
// It uses a rolling hash algorithm with a prime number to avoid collisions.
// The resulting hash can be used as a unique identifier for the metric.
func BuildKey(metricName, metricType string) uint64 {
	var hash uint64
	const (
		prime uint64 = 31
		mod   uint64 = 1e9 + 7
	)

	for i := 0; i < len(metricName); i++ {
		hash = (hash*prime + uint64(metricName[i])) % mod
	}
	for i := 0; i < len(metricType); i++ {
		hash = (hash*prime + uint64(metricType[i])) % mod
	}

	return hash
}

// Storage is an interface that defines methods for managing metrics in a storage system.
type Storage interface {
	// Set stores a metric in the storage.
	Set(ctx context.Context, metric Metric) (err error)

	// Get retrieves a metric from the storage by its unique key.
	Get(ctx context.Context, key uint64) (*Metric, error)

	// GetAll retrieves all metrics stored in the storage.
	GetAll(ctx context.Context) (*[]Metric, error)

	// Lock ensures exclusive access to the storage, preventing concurrent modifications.
	Lock()

	// Unlock releases the exclusive access lock, allowing other goroutines to modify the storage.
	Unlock()

	// Actualize processes and updates the storage data, making sure that it's up-to-date.
	// It might be used to refresh or recalculate metrics.
	Actualize(ctx context.Context) error

	// Dump writes the current state of the storage to an external source.
	// This can be used to persist the metrics to a file or database.
	Dump(ctx context.Context) error

	// Ping checks if the storage is accessible and responsive.
	// This can be used to verify the health of the storage.
	Ping(ctx context.Context) error

	// Close gracefully shuts down the storage and releases any resources it holds.
	// It should be called when the storage is no longer needed.
	Close() error
}
