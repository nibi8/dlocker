package dlocker

import (
	"context"

	"github.com/nibi8/dlocker/models"
)

type LockContextImp struct {
	context.Context

	lock models.Lock
	lr   models.LockRecord
}

func (ctx LockContextImp) GetLock() models.Lock {
	return ctx.lock
}

func (ctx LockContextImp) GetLockRecord() models.LockRecord {
	return ctx.lr
}

func NewLockContext(
	ctx context.Context,
	lock models.Lock,
	lr models.LockRecord,
) LockContextImp {
	return LockContextImp{
		Context: ctx,
		lock:    lock,
		lr:      lr,
	}
}
