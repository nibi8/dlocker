package dlocker

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/p8bin/dlocker/models"
	"github.com/p8bin/dlocker/storageproviders/testsp"
)

// todo: add detailed tests with mocks

func TestFiniteExtends(t *testing.T) {
	ctx := context.Background()

	sp := testsp.NewStorageProvider()

	locker := NewLocker(sp)

	lock, err := models.NewLock(
		"unique-lock-name",
		4,
		1,
	)
	require.NoError(t, err)

	lockEx := LockExtender(NewLockExtender(locker, lock, 4))

	rCtx, err := lockEx.LockWithWait(ctx)
	require.NoError(t, err)

	for i := 0; ; i++ {
		if rCtx.Err() != nil {
			time.Sleep(1 * time.Second)
			return
		}
		fmt.Println(i)
		time.Sleep(1 * time.Second)
	}

}

func TestInfiniteExtends(t *testing.T) {
	ctx := context.Background()

	sp := testsp.NewStorageProvider()

	locker := NewLocker(sp)

	lock, err := models.NewLock(
		"unique-lock-name",
		4,
		1,
	)
	require.NoError(t, err)

	lockEx := NewLockExtender(locker, lock, 0)

	rCtx, err := lockEx.LockWithWait(ctx)
	require.NoError(t, err)

	go func() {
		time.Sleep(25 * time.Second)
		lockEx.Unlock(ctx, true)
	}()

	for i := 0; ; i++ {
		if rCtx.Err() != nil {
			time.Sleep(1 * time.Second)
			return
		}
		fmt.Println(i)
		time.Sleep(1 * time.Second)
	}

}

func TestInit(t *testing.T) {
	ctx := context.Background()

	sp := testsp.NewStorageProvider()

	locker := NewLocker(sp)

	lock, err := models.NewLock(
		"unique-lock-name",
		4,
		1,
	)
	require.NoError(t, err)

	lockEx := NewLockExtender(locker, lock, 4)

	assert.NotEmpty(t, lockEx.locker)
	assert.NotEmpty(t, lockEx.lock)
	assert.NotEmpty(t, lockEx.count)
	assert.Empty(t, lockEx.runCtx)
	assert.Empty(t, lockEx.runCancel)
	assert.Empty(t, lockEx.mu)
	assert.NotNil(t, lockEx.quit)

	_, err = lockEx.LockWithWait(ctx)
	require.NoError(t, err)

	assert.NoError(t, lockEx.runCtx.Err())
	assert.NotEmpty(t, lockEx.runCancel)
	assert.NotEmpty(t, lockEx.mu)

	lockEx.Unlock(ctx, true)

	//time.Sleep(5 * time.Second)

	assert.Error(t, lockEx.runCtx.Err())

}
