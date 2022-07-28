package models

import (
	"fmt"
)

type Lock struct {
	// Unique lock name
	Name string

	// Duration of lock in seconds (TTL)
	ExecutionDurationSec int

	// Interval between locks execution in seconds
	SpanDurationSec int

	// Additional lock check period
	CheckPeriodSec int
}

func (l *Lock) SetCheckPeriod(checkPeriodSec int) (err error) {
	if checkPeriodSec < 0 {
		return ErrBadParam
	}
	l.CheckPeriodSec = checkPeriodSec
	return nil
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

func NewLockPnc(
	name string,
	executionDurationSec int,
	spanDurationSec int,
) (lock Lock) {
	lock, err := NewLock(
		name,
		executionDurationSec,
		spanDurationSec,
	)
	if err != nil {
		panic(err)
	}
	return lock
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
