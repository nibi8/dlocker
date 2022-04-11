package dlocker

import (
	"context"

	"github.com/nibi8/dlocker/models"
)

// Manages locks
type Locker interface {
	// Set a lock (and wait for previous) or return ErrNoLuck if no luck or an unexpected error occurs
	LockWithWait(
		ctx context.Context,
		lock models.Lock,
	) (lockCtx context.Context, cancel context.CancelFunc, err error)

	// Release lock or return ErrNoLuck if no luck or an unexpected error occurs
	Unlock(
		lockCtx context.Context,
	) (err error)
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
	// ErrNotSupported returns in case storage does not support updates and created lock record expires (for example Redis values)
	UpdateLockRecord(
		ctx context.Context,
		jobName string,
		version string,
		patch models.LockRecordPatch,
	) (err error)
}
