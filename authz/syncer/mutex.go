package syncer

import (
	"context"
	"errors"

	"github.com/go-redsync/redsync/v4"
)

type Mutex interface {
	Lock(ctx context.Context) error
	Unlock(ctx context.Context)
}

type RedisMutes struct {
	lock *redsync.Mutex
}

func NewRedisMutex(redc *redsync.Redsync) *RedisMutes {
	const mutexLockKey = "authz_syncer"
	return &RedisMutes{
		lock: redc.NewMutex(mutexLockKey),
	}
}

func (r *RedisMutes) Lock(ctx context.Context) error {
	return r.lock.LockContext(ctx)
}

func (r *RedisMutes) Unlock(ctx context.Context) {
	_, _ = r.lock.UnlockContext(ctx)
}

func isMutexLocked(err error) bool {
	return errors.Is(err, redsync.ErrFailed)
}
