package dlocker

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/nibi8/dlocker/models"
	"github.com/nibi8/dlocker/storageproviders/testsp"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// todo: add detailed tests with mocks

func TestNewLocker(t *testing.T) {
	ctx := context.Background()

	sp := testsp.NewStorageProvider()
	locker := NewLocker(sp)

	lock, err := models.NewLock(
		"unique-lock-name",
		60,
		10,
	)
	require.NoError(t, err)
	if err != nil {
		t.Error(err)
	}

	_, _, err = locker.LockWithWait(ctx, lock)
	require.NoError(t, err)
}

func TestNewLockerWithOpt(t *testing.T) {
	ctx := context.Background()

	sp := testsp.NewStorageProvider()
	locker := NewLocker(sp)

	lock, err := models.NewLock(
		"unique-lock-name",
		60,
		10,
	)
	require.NoError(t, err)

	err = lock.SetCheckPeriod(5)
	require.NoError(t, err)

	_, _, err = locker.LockWithWait(ctx, lock)
	require.NoError(t, err)
}

func TestLockWithWait(t *testing.T) {
	ctx := context.Background()

	sp := testsp.NewStorageProvider()
	locker := NewLocker(sp)

	lock, err := models.NewLock(
		"unique-lock-name",
		2,
		1,
	)
	require.NoError(t, err)

	dt := time.Now()
	_, _, err = locker.LockWithWait(ctx, lock)
	require.NoError(t, err)
	assert.True(t, time.Since(dt).Seconds() < 1)

	dt = time.Now()
	lockCtx, _, err := locker.LockWithWait(ctx, lock)
	require.NoError(t, err)
	assert.True(t, time.Since(dt).Seconds() > 1)

	<-lockCtx.Done()

	dt = time.Now()
	_, _, err = locker.LockWithWait(ctx, lock)
	require.NoError(t, err)
	// wait anyway
	assert.True(t, time.Since(dt).Seconds() > 1)
}

func TestExtendLock(t *testing.T) {

	ctx := context.Background()

	sp := testsp.NewStorageProvider()
	locker := NewLocker(sp)

	lock, err := models.NewLock(
		"unique-lock-name",
		2,
		1,
	)
	require.NoError(t, err)

	lockCtx, _, err := locker.LockWithWait(ctx, lock)
	require.NoError(t, err)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, _, err := locker.LockWithWait(ctx, lock)
		assert.True(t, errors.Is(err, models.ErrNoLuck), err)
	}()

	time.Sleep(1 * time.Second)

	newLockCtx, _, err := locker.ExtendLock(ctx, lockCtx)
	require.NoError(t, err)

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, _, err := locker.LockWithWait(ctx, lock)
		assert.True(t, errors.Is(err, models.ErrNoLuck), err)
	}()

	time.Sleep(1 * time.Second)

	_, _, err = locker.ExtendLock(ctx, newLockCtx)
	require.NoError(t, err)

	wg.Wait()
}

func TestUnlock(t *testing.T) {
	ctx := context.Background()

	sp := testsp.NewStorageProvider()
	locker := NewLocker(sp)

	lock, err := models.NewLock(
		"unique-lock-name",
		2,
		1,
	)
	require.NoError(t, err)

	lockCtx, _, err := locker.LockWithWait(ctx, lock)
	require.NoError(t, err)

	err = locker.Unlock(ctx, lockCtx)
	require.NoError(t, err)

	dt := time.Now()
	_, _, err = locker.LockWithWait(ctx, lock)
	require.NoError(t, err)
	assert.True(t, time.Since(dt).Seconds() < 1)
}
