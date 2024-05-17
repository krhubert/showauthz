package syncer

import (
	"context"

	"rift/memdb"
)

// Database interface provides methods to fetch data to be synced.
type Database interface {
	ShouldSync(ctx context.Context) (bool, error)
	SyncCompleted(ctx context.Context) error
	Members(ctx context.Context) ([]*memdb.Member, error)
	OffDays(ctx context.Context) ([]*memdb.OffDay, error)
}

type DB struct {
	db *memdb.DB
}

func (pd *DB) ShouldSync(ctx context.Context) (bool, error) {
	return true, nil // check if sync is required
}

func (pd *DB) SyncCompleted(ctx context.Context) error {
	return nil // update sync status
}

func (pd *DB) Members(ctx context.Context) ([]*memdb.Member, error) {
	return pd.db.AllMembers(), nil
}

func (pd *DB) OffDays(ctx context.Context) ([]*memdb.OffDay, error) {
	return pd.db.AllOffDays(), nil
}
