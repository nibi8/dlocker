//
// !!! For test use only
//
package testsp

import (
	"context"
	"fmt"
	"sync"

	"github.com/p8bin/dlocker/models"
)

//
// !!! For test use only
//

type StorageProvider struct {
	locks sync.Map
}

func NewStorageProvider() *StorageProvider {
	return &StorageProvider{}
}

type MemLock struct {
	Lock *sync.Mutex
	Lr   models.LockRecord
}

func NewMemLock(lr models.LockRecord) MemLock {
	return MemLock{
		Lock: &sync.Mutex{},
		Lr:   lr,
	}
}

func (sp *StorageProvider) GetLockRecord(
	ctx context.Context,
	lockName string,
) (lr models.LockRecord, err error) {

	value, ok := sp.locks.Load(lockName)
	if !ok {
		return lr, models.ErrNotFound
	}

	jl, ok := value.(MemLock)
	if !ok {
		return lr, fmt.Errorf("cast error")
	}

	return jl.Lr, nil
}

func (sp *StorageProvider) CreateLockRecord(
	ctx context.Context,
	lr models.LockRecord,
) (err error) {

	jl := NewMemLock(lr)
	_, loaded := sp.locks.LoadOrStore(lr.LockName, jl)
	if loaded {
		return models.ErrDuplicate
	}

	return nil
}

func (sp *StorageProvider) UpdateLockRecord(
	ctx context.Context,
	lockName string,
	version string,
	patch models.LockRecordPatch,
) (err error) {

	value, ok := sp.locks.Load(lockName)
	if !ok {
		return models.ErrNotFound
	}
	jl, ok := value.(MemLock)
	if !ok {
		return fmt.Errorf("cast error")
	}

	jl.Lock.Lock()
	defer jl.Lock.Unlock()

	// reread
	value, ok = sp.locks.Load(lockName)
	if !ok {
		return models.ErrNotFound
	}
	jl, ok = value.(MemLock)
	if !ok {
		return fmt.Errorf("cast error")
	}

	if jl.Lr.Version != version {
		return models.ErrNoLuck
	}

	jl.Lr.ApplyPatch(patch)

	sp.locks.Store(lockName, jl)

	return nil
}
