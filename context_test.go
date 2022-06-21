package dlocker

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/p8bin/dlocker/models"
)

func TestNewLockContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), "key", "val")

	lock, err := models.NewLock(
		"unique-lock-name",
		60,
		10,
	)
	require.NoError(t, err)

	lr := models.NewLockRecord(lock)

	lCtx := NewLockContext(
		ctx,
		lock,
		lr,
	)

	assert.NotEmpty(t, lCtx.Context)
	assert.NotEmpty(t, lCtx.lock)
	assert.NotEmpty(t, lCtx.lr)

	assert.NotEmpty(t, lCtx.GetLock())
	assert.NotEmpty(t, lCtx.GetLockRecord())
}
