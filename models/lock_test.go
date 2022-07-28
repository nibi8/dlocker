package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLock(t *testing.T) {
	lock, err := NewLock(
		"unique-lock-name",
		60,
		10,
	)
	require.NoError(t, err)

	assert.NotEmpty(t, lock.Name)
	assert.NotEmpty(t, lock.ExecutionDurationSec)
	assert.NotEmpty(t, lock.SpanDurationSec)
	assert.Empty(t, lock.CheckPeriodSec)
}

func TestNewLockPnc(t *testing.T) {
	lock := NewLockPnc(
		"unique-lock-name",
		60,
		10,
	)

	assert.NotEmpty(t, lock.Name)
	assert.NotEmpty(t, lock.ExecutionDurationSec)
	assert.NotEmpty(t, lock.SpanDurationSec)
	assert.Empty(t, lock.CheckPeriodSec)
}

func TestLockSetCheckPeriod(t *testing.T) {
	lock, err := NewLock(
		"unique-lock-name",
		60,
		10,
	)
	require.NoError(t, err)

	err = lock.SetCheckPeriod(20)
	require.NoError(t, err)

	assert.NotEmpty(t, lock.CheckPeriodSec)
}

func TestLockValidate(t *testing.T) {
	_, err := NewLock(
		"",
		60,
		10,
	)
	assert.Error(t, err)

	_, err = NewLock(
		"unique-lock-name",
		0,
		10,
	)
	require.Error(t, err)

	_, err = NewLock(
		"unique-lock-name",
		60,
		0,
	)
	require.Error(t, err)
}

func TestGetDurationSec(t *testing.T) {
	lock, err := NewLock(
		"unique-lock-name",
		60,
		10,
	)
	require.NoError(t, err)

	assert.Greater(t, lock.GetDurationSec(), lock.ExecutionDurationSec)
}
