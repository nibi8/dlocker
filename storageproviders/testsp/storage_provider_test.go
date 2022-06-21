package testsp

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nibi8/dlocker/models"
)

func TestInit(t *testing.T) {
	_ = NewStorageProvider()
}

func TestCreate(t *testing.T) {
	ctx := context.Background()

	sp := NewStorageProvider()

	lock, err := models.NewLock("unique-lock-name", 60, 10)
	require.NoError(t, err)

	lr := models.NewLockRecord(lock)

	err = sp.CreateLockRecord(ctx, lr)
	require.NoError(t, err)

	err = sp.CreateLockRecord(ctx, lr)
	assert.True(t, errors.Is(err, models.ErrDuplicate))

}

func TestGet(t *testing.T) {
	ctx := context.Background()

	sp := NewStorageProvider()

	lock, err := models.NewLock("unique-lock-name", 60, 10)
	require.NoError(t, err)

	_, err = sp.GetLockRecord(ctx, lock.Name)
	assert.True(t, errors.Is(err, models.ErrNotFound))

	lr := models.NewLockRecord(lock)

	err = sp.CreateLockRecord(ctx, lr)
	require.NoError(t, err)

	lrRed, err := sp.GetLockRecord(ctx, lock.Name)
	require.NoError(t, err)

	assert.Equal(t, lr, lrRed)
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()

	sp := NewStorageProvider()

	lock, err := models.NewLock("unique-lock-name", 60, 10)
	require.NoError(t, err)

	lr := models.NewLockRecord(lock)

	patch := models.NewLockRecordPatchForCapture(lr.DurationSec)
	err = sp.UpdateLockRecord(ctx, lr.LockName, lr.Version, patch)
	require.True(t, errors.Is(err, models.ErrNotFound))

	err = sp.CreateLockRecord(ctx, lr)
	require.NoError(t, err)

	err = sp.UpdateLockRecord(ctx, lr.LockName, lr.Version, patch)
	require.NoError(t, err)
}
