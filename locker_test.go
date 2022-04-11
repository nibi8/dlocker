package dlocker

import (
	"context"
	"testing"

	"github.com/nibi8/dlocker/models"
	"github.com/nibi8/dlocker/storageproviders/testsp"
)

// todo: add tests

func TestNewLocker(t *testing.T) {
	sp := testsp.NewStorageProvider()
	locker := NewLocker(sp)

	lock, err := models.NewLock(
		"unique-lock-name",
		30,
		10,
	)
	if err != nil {
		t.Error(err)
	}

	_, _, err = locker.LockWithWait(context.Background(), lock)
	if err != nil {
		t.Error(err)
	}

}
