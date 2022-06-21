package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLockRecord(t *testing.T) {

	lock, err := NewLock(
		"unique-lock-name",
		60,
		10,
	)
	require.NoError(t, err)

	lr := NewLockRecord(lock)

	assert.Equal(t, lr.LockName, lock.Name)
	assert.NotEmpty(t, lr.Version)
	assert.NotEmpty(t, lr.DurationSec)
	assert.Equal(t, lr.State, LockRecordStateLock)
	assert.NotEmpty(t, lr.Dt)
}

func TestLockRecordIsLock(t *testing.T) {

	lock, err := NewLock(
		"unique-lock-name",
		60,
		10,
	)
	require.NoError(t, err)

	lr := NewLockRecord(lock)

	assert.True(t, lr.State.IsLock())

	lr.State = LockRecordStateUnlock

	assert.False(t, lr.State.IsLock())
}

func TestNewLockRecordPatchForCapture(t *testing.T) {
	p := NewLockRecordPatchForCapture(60)

	assert.NotEmpty(t, p.Version)
	assert.NotEmpty(t, p.DurationSec)
	assert.Equal(t, p.State, LockRecordStateLock)
	assert.NotEmpty(t, p.Dt)
}

func TestNewLockRecordPatchForRelease(t *testing.T) {
	const version = "12345"
	p := NewLockRecordPatchForRelease(version)

	assert.Equal(t, p.Version, version)
	assert.Empty(t, p.DurationSec)
	assert.Equal(t, p.State, LockRecordStateUnlock)
	assert.NotEmpty(t, p.Dt)
}

func TestLockRecordApplyPatchForCapture(t *testing.T) {
	lock, err := NewLock(
		"unique-lock-name",
		60,
		10,
	)
	require.NoError(t, err)

	lr := NewLockRecord(lock)

	time.Sleep(10)

	const patchDurSec = 100
	p := NewLockRecordPatchForCapture(patchDurSec)

	lrOld := lr

	lr.ApplyPatch(p)

	assert.Equal(t, lr.LockName, lrOld.LockName)
	assert.NotEqual(t, lr.Version, lrOld.Version)
	assert.Equal(t, lr.DurationSec, patchDurSec)
	assert.Equal(t, lr.State, LockRecordStateLock)
	assert.NotEqual(t, lr.Dt, lrOld.Dt)
}

func TestLockRecordApplyPatchForRelease(t *testing.T) {
	lock, err := NewLock(
		"unique-lock-name",
		60,
		10,
	)
	require.NoError(t, err)

	lr := NewLockRecord(lock)

	time.Sleep(10)

	p := NewLockRecordPatchForRelease(lr.Version)

	lrOld := lr

	lr.ApplyPatch(p)

	assert.Equal(t, lr.LockName, lrOld.LockName)
	assert.Equal(t, lr.Version, lrOld.Version)
	assert.Empty(t, lr.DurationSec)
	assert.Equal(t, lr.State, LockRecordStateUnlock)
	assert.NotEqual(t, lr.Dt, lrOld.Dt)
}
