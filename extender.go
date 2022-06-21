package dlocker

import (
	"context"
	"sync"
	"time"

	"github.com/nibi8/dlocker/models"
)

type LockExtenderImp struct {
	locker Locker
	lock   models.Lock
	count  int

	runCtx    context.Context
	runCancel context.CancelFunc
	mu        sync.Mutex
	quit      chan bool
}

func NewLockExtender(
	locker Locker,
	lock models.Lock,
	count int,
) *LockExtenderImp {
	return &LockExtenderImp{
		locker: locker,
		lock:   lock,
		count:  count,

		quit: make(chan bool, 1),
	}
}

func (ex *LockExtenderImp) LockWithWait(
	ctx context.Context,
) (exCtx context.Context, err error) {

	if !ex.mu.TryLock() {
		return exCtx, models.ErrNoLuck
	}
	running := false
	defer func() {
		if !running {
			ex.mu.Unlock()
		}
	}()

	ex.runCtx, ex.runCancel = context.WithCancel(ctx)
	lockCtx, _, err := ex.locker.LockWithWait(ex.runCtx, ex.lock)
	if err != nil {
		return exCtx, err
	}

	running = true
	go ex.run(lockCtx)

	return ex.runCtx, nil
}

func (ex *LockExtenderImp) Unlock(ctx context.Context, wait bool) {
	select {
	case ex.quit <- true:
	default:
	}

	if wait {
		ex.mu.Lock()
		defer ex.mu.Unlock()
	}
}

func (ex *LockExtenderImp) run(
	lockCtx LockContext,
) {
	defer func() {
		if ex.runCtx.Err() == nil {
			ex.runCancel()
		}
		ex.mu.Unlock()
	}()

	dur := time.Duration(ex.lock.ExecutionDurationSec/2) * time.Second
	count := 0
	checkCount := false
	if ex.count > 0 {
		count = ex.count * 2
		checkCount = true
	}

	needUnlock := false

loop:
	for {

		if checkCount {
			if count < 1 {
				break loop
			}
			count--
		}

		select {
		case <-time.After(dur):
			// continue
		case <-lockCtx.Done():
			// case <-ex.runCtx.Done():
			break loop
		case <-ex.quit:
			needUnlock = true
			break loop
		}

		var err error
		lockCtx, _, err = ex.locker.ExtendLock(ex.runCtx, lockCtx)
		if err != nil {
			break loop
		}

	}

	if needUnlock && ex.runCtx.Err() == nil && lockCtx.Err() == nil {
		_ = ex.locker.Unlock(ex.runCtx, lockCtx)
		// todo: ? process err
	}

}
