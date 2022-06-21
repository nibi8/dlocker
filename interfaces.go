package dlocker

import (
	"context"

	"github.com/p8bin/dlocker/models"
)

// Manages locks
type Locker interface {
	// Set a lock (and wait for previous) or return ErrNoLuck if no luck or an unexpected error occurs
	LockWithWait(
		ctx context.Context,
		lock models.Lock,
	) (lockCtx LockContext, cancel context.CancelFunc, err error)

	// Extend existing lock
	ExtendLock(
		ctx context.Context,
		lockCtx LockContext,
	) (newLockCtx LockContext, cancel context.CancelFunc, err error)

	// Release lock or return ErrNoLuck if no luck or an unexpected error occurs
	Unlock(
		ctx context.Context,
		lockCtx LockContext,
	) (err error)
}

// Extends existing lock
type LockExtender interface {
	LockWithWait(ctx context.Context) (exCtx context.Context, err error)
	Unlock(ctx context.Context, wait bool)
}

// Creates lock records in a persistent storage.
type StorageProvider interface {
	// Returns LockRecord or error ErrNotFound if not found or unexpected error.
	GetLockRecord(
		ctx context.Context,
		jobName string,
	) (lr models.LockRecord, err error)

	// Creates LockRecord or returns error ErrDuplicate if already exists or unexpected error.
	CreateLockRecord(
		ctx context.Context,
		lr models.LockRecord,
	) (err error)

	// Updates LockRecord with new values or return error ErrNotFound if not found or unexpected error.
	// Returns ErrNotSupported in case storage does not support updates (in this case storage must handle ttl)
	UpdateLockRecord(
		ctx context.Context,
		jobName string,
		version string,
		patch models.LockRecordPatch,
	) (err error)
}

// Lock context
type LockContext interface {
	context.Context

	GetLock() models.Lock
	GetLockRecord() models.LockRecord
}
