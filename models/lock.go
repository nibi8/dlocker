package models

import (
	"fmt"
)

type Lock struct {
	// Unique lock name
	Name string

	// Duration of lock in seconds
	ExecutionDurationSec int

	// Interval between locks execution in seconds
	SpanDurationSec int
}

func (j Lock) Validate() (err error) {
	if j.Name == "" {
		return fmt.Errorf(`Name == ""`)
	}

	if j.ExecutionDurationSec < 1 {
		return fmt.Errorf("ExecutionDurationSec < 1")
	}

	if j.SpanDurationSec < 1 {
		return fmt.Errorf("SpanDurationSec < 1")
	}

	return nil
}

// Total lock period
func (j Lock) GetDurationSec() int {
	return j.ExecutionDurationSec + j.SpanDurationSec
}

func NewLock(
	name string,
	executionDurationSec int,
	spanDurationSec int,
) (lock Lock, err error) {

	lock = Lock{
		Name:                 name,
		ExecutionDurationSec: executionDurationSec,
		SpanDurationSec:      spanDurationSec,
	}

	err = lock.Validate()
	if err != nil {
		return lock, err
	}

	return lock, nil
}
