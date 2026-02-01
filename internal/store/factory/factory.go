// Package factory provides store creation functions.
package factory

import (
	"fmt"

	"github.com/ossydotpy/veil/internal/store"
	"github.com/ossydotpy/veil/internal/store/sqlite"
)

// NewStore creates a new store instance based on the store type.
func NewStore(storeType, dbPath string) (store.Store, error) {
	switch storeType {
	case "sqlite":
		return sqlite.NewSqliteStore(dbPath)
	default:
		return nil, fmt.Errorf("unsupported store type: %s (supported: sqlite)", storeType)
	}
}
