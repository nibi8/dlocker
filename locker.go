package dlocker

import (
	"context"
	"errors"
	"time"

	"github.com/nibi8/dlocker/models"
)

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
) (lockCtx LockContext, cancel context.CancelFunc, err error) {

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
		dur := time.Duration(lr.DurationSec) * time.Second
	loop:
		for {
			var after <-chan time.Time
			if lock.CheckPeriodSec > 0 {
				after = time.After(time.Duration(lock.CheckPeriodSec) * time.Second)
			}
			select {
			case <-ctx.Done():
				return lockCtx, cancel, ctx.Err()
			case <-time.After(dur):
				break loop
			case <-after:
				if lrCheck, err := s.sp.GetLockRecord(ctx, lock.Name); err != nil {
					if !errors.Is(err, models.ErrNotFound) {
						// for storages with ttl keys
						lrFound = false
						break loop
					}
					// todo: ? process err
				} else {
					if lr.Version != lrCheck.Version {
						return lockCtx, cancel, models.ErrNoLuck
					}
					if lr.State == models.LockRecordStateUnlock {
						break loop
					}
				}
			}
		}
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

		lockCtx, cancel = createLockCtx(ctx, lock, lr)

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

	lockCtx, cancel = createLockCtx(ctx, lock, lr)

	return lockCtx, cancel, nil
}

func (s *LockerImp) ExtendLock(
	ctx context.Context,
	lockCtx LockContext,
) (newLockCtx LockContext, cancel context.CancelFunc, err error) {

	lr := lockCtx.GetLockRecord()
	lock := lockCtx.GetLock()

	lrPatch := models.NewLockRecordPatchForCapture(lr.DurationSec)
	err = s.sp.UpdateLockRecord(ctx, lr.LockName, lr.Version, lrPatch)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) ||
			errors.Is(err, models.ErrNotSupported) {
			return newLockCtx, cancel, models.ErrNoLuck
		}
		return newLockCtx, cancel, err
	}

	lr.ApplyPatch(lrPatch)

	newLockCtx, cancel = createLockCtx(ctx, lock, lr)

	return newLockCtx, cancel, nil
}

func (s *LockerImp) Unlock(
	ctx context.Context,
	lockCtx LockContext,
) (err error) {

	lrPatch := models.NewLockRecordPatchForRelease(lockCtx.GetLockRecord().Version)
	err = s.sp.UpdateLockRecord(ctx, lockCtx.GetLockRecord().LockName, lockCtx.GetLockRecord().Version, lrPatch)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) ||
			errors.Is(err, models.ErrNotSupported) {
			return models.ErrNoLuck
		}
		return err
	}
	return nil
}

func createLockCtx(
	ctx context.Context,
	lock models.Lock,
	lr models.LockRecord,
) (lockCtx LockContext, cancel context.CancelFunc) {

	ctx, cancel = context.WithTimeout(ctx, time.Duration(lock.ExecutionDurationSec)*time.Second)
	lockCtx = NewLockContext(ctx, lock, lr)

	return lockCtx, cancel
}
