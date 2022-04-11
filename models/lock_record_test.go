package models

import (
	"fmt"
	"testing"
)

// todo: add tests

func TestNewLockRecord(t *testing.T) {

	lock, err := NewLock(
		"unique-lock-name",
		30,
		10,
	)

	if err != nil {
		t.Error(err)
	}

	lr := NewLockRecord(lock)

	fmt.Println("Lock record new unique version:", lr.Version)

}
