package mongosp

import (
	"time"

	"github.com/nibi8/dlocker/models"
)

type LockRecordDB struct {
	LockName    string                 `bson:"jobname"`
	Version     string                 `bson:"version"`
	DurationSec int                    `bson:"durationsec"`
	State       models.LockRecordState `bson:"state"`
	Dt          time.Time              `bson:"dt"`
}

func FromLockRecord(in models.LockRecord) LockRecordDB {
	out := LockRecordDB{
		LockName:    in.LockName,
		Version:     in.Version,
		DurationSec: in.DurationSec,
		State:       in.State,
		Dt:          in.Dt,
	}
	return out
}

func ToLockRecord(in LockRecordDB) models.LockRecord {
	out := models.LockRecord{
		LockName:    in.LockName,
		Version:     in.Version,
		DurationSec: in.DurationSec,
		State:       in.State,
		Dt:          in.Dt,
	}
	return out
}
