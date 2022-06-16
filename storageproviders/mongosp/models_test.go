package mongosp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nibi8/dlocker/models"
)

func TestFromToTest(t *testing.T) {
	lock, err := models.NewLock(
		"unique-lock-name",
		60,
		10,
	)
	require.NoError(t, err)

	lr := models.NewLockRecord(lock)

	lrdb := FromLockRecord(lr)

	lrRestored := ToLockRecord(lrdb)

	assert.Equal(t, lr, lrRestored)
}
