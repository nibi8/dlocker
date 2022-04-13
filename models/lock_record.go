package models

import (
	"time"

	"github.com/google/uuid"
)

type LockRecord struct {
	LockName    string
	Version     string
	DurationSec int
	State       LockRecordState
	Dt          time.Time
}

type LockRecordState string

const (
	LockRecordStateNone   LockRecordState = ""
	LockRecordStateLock   LockRecordState = "lock"
	LockRecordStateUnlock LockRecordState = "unlock"
)

func (state LockRecordState) IsLock() bool {
	return state == LockRecordStateLock || state == LockRecordStateNone
}

func NewLockRecord(
	lock Lock,
) LockRecord {
	lr := LockRecord{
		LockName:    lock.Name,
		Version:     uuid.New().String(),
		DurationSec: lock.GetDurationSec(),
		State:       LockRecordStateLock,
		Dt:          time.Now(),
	}
	return lr
}

type LockRecordPatch struct {
	Version     string
	DurationSec int
	State       LockRecordState
	Dt          time.Time
}

func NewLockRecordPatchForCapture(
	durationSec int,
) LockRecordPatch {
	lr := LockRecordPatch{
		Version:     uuid.New().String(),
		DurationSec: durationSec,
		State:       LockRecordStateLock,
		Dt:          time.Now(),
	}
	return lr
}

func NewLockRecordPatchForRelease(curVersion string) LockRecordPatch {
	patch := LockRecordPatch{
		Version:     curVersion, // todo: ? new version
		DurationSec: 1,          // todo: ? set to 0
		State:       LockRecordStateUnlock,
		Dt:          time.Now(),
	}
	return patch
}

func (lr *LockRecord) ApplyPatch(patch LockRecordPatch) {
	lr.Version = patch.Version
	lr.DurationSec = patch.DurationSec
	lr.State = patch.State
	lr.Dt = patch.Dt
}
