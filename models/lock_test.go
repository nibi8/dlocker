package models

import (
	"fmt"
	"testing"
)

// todo: add tests

func TestNewLock(t *testing.T) {
	lock, err := NewLock(
		"unique-lock-name",
		30,
		10,
	)

	if err != nil {
		t.Error(err)
	}

	fmt.Println("Total lock period:", lock.GetDurationSec())

}
