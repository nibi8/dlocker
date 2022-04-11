package dlocker

import (
	"context"
	"errors"
	"time"

	"github.com/nibi8/dlocker/models"
)

const lockCtxLrKey = "github.com/nibi8/dlocker/lock_ctx_lr_key"

type LockerImp struct {
	sp StorageProvider
}

func NewLocker(
	sp StorageProvider,
) *LockerImp {
	s := LockerImp{
		sp: sp,
	}
	return &s
}

// Set a lock (and wait for previous) or return ErrNoLuck if no luck or an unexpected error occurs
func (s *LockerImp) LockWithWait(
	ctx context.Context,
	lock models.Lock,
) (lockCtx context.Context, cancel context.CancelFunc, err error) {

	lrFound := false
	lr, err := s.sp.GetLockRecord(ctx, lock.Name)
	if err != nil {
		if !errors.Is(err, models.ErrNotFound) {
			return lockCtx, cancel, err
		}
		// not found
		// continue
		err = nil
		lrFound = false
	} else {
		lrFound = true
	}

	if lrFound && lr.State.IsLock() && lr.DurationSec > 0 {
		time.Sleep(time.Duration(lr.DurationSec) * time.Second)
	}

	if !lrFound {
		lr = models.NewLockRecord(lock)
		err = s.sp.CreateLockRecord(ctx, lr)
		if err != nil {
			if errors.Is(err, models.ErrDuplicate) {
				return lockCtx, cancel, models.ErrNoLuck
			}
			return lockCtx, cancel, err
		}

		// todo: ? init of lockCtx create or update
		lockCtx, cancel = newLockCtx(ctx, lock, lr)

		return lockCtx, cancel, nil
	}

	lrPatch := models.NewLockRecordPatchForCapture(lock.GetDurationSec())
	err = s.sp.UpdateLockRecord(ctx, lock.Name, lr.Version, lrPatch)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) ||
			errors.Is(err, models.ErrNotSupported) {
			return lockCtx, cancel, models.ErrNoLuck
		}
		return lockCtx, cancel, err
	}

	lr.ApplyPatch(lrPatch)

	// todo: ? init of lockCtx create or update
	lockCtx, cancel = newLockCtx(ctx, lock, lr)

	return lockCtx, cancel, nil
}

func (s *LockerImp) Unlock(
	lockCtx context.Context,
) (err error) {

	val := lockCtx.Value(lockCtxLrKey)
	if val == nil {
		return nil
	}

	lr, ok := val.(models.LockRecord)
	if !ok {
		return nil
	}

	lrPatch := models.NewLockRecordPatchForRelease(lr.Version)
	err = s.sp.UpdateLockRecord(lockCtx, lr.LockName, lr.Version, lrPatch)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) ||
			errors.Is(err, models.ErrNotSupported) {
			return models.ErrNoLuck
		}
		return err
	}
	return nil
}

func newLockCtx(
	ctx context.Context,
	lock models.Lock,
	lr models.LockRecord,
) (lockCtx context.Context, cancel context.CancelFunc) {

	lockCtx, cancel = context.WithTimeout(ctx, time.Duration(lock.ExecutionDurationSec)*time.Second)
	lockCtx = context.WithValue(lockCtx, lockCtxLrKey, lr)

	return lockCtx, cancel
}
