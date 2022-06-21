package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/p8bin/dlocker"
	"github.com/p8bin/dlocker/models"
	"github.com/p8bin/dlocker/storageproviders/testsp"
)

func main() {

	// create storage provider
	sp := testsp.NewStorageProvider()

	// create locker
	locker := dlocker.NewLocker(sp)

	// create locks

	lock1, err := models.NewLock("unique_lock_name_1", 10, 5)
	if err != nil {
		log.Fatal("NewLock")
	}
	err = lock1.SetCheckPeriod(2)
	if err != nil {
		log.Fatal("SetCheckPeriod")
	}
	lockEx1 := dlocker.NewLockExtender(locker, lock1, 0)

	lock2, err := models.NewLock("unique_lock_name_2", 10, 5)
	if err != nil {
		log.Fatal("NewLock")
	}
	lockEx2 := dlocker.NewLockExtender(locker, lock2, 4)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(4)

	// lock1:
	go func() {
		defer wg.Done()
		exCtx, err := captureLock(ctx, locker, lock1, "instace_1", false, lockEx1)
		if err != nil {
			return
		}
		<-exCtx.Done()

	}()
	go func() {
		defer wg.Done()
		time.Sleep(1 * time.Second)
		exCtx, err := captureLock(ctx, locker, lock1, "instace_2", false, lockEx1)
		if err != nil {
			return
		}
		<-exCtx.Done()
	}()

	// lock2 with manual unlock:
	go func() {
		defer wg.Done()
		exCtx, err := captureLock(ctx, locker, lock2, "instace_1", true, lockEx2)
		if err != nil {
			return
		}
		<-exCtx.Done()
	}()
	go func() {
		defer wg.Done()
		time.Sleep(1 * time.Second)
		exCtx, err := captureLock(ctx, locker, lock2, "instace_2", true, lockEx2)
		if err != nil {
			return
		}
		<-exCtx.Done()
	}()

	wg.Wait()
}

func captureLock(
	ctx context.Context,
	locker dlocker.Locker,
	lock models.Lock,
	instanceName string,
	unlock bool,
	ex dlocker.LockExtender,
) (exCtx context.Context, err error) {

	lockPrintName := lock.Name
	if instanceName != "" {
		lockPrintName += " " + instanceName
	}

	msg := "lock success"
	exCtx, err = ex.LockWithWait(ctx)
	if err != nil {
		msg = fmt.Sprintf("lock failed with error = %v", err)
	}
	if unlock {
		go func() {
			time.Sleep(10 * time.Second)
			ex.Unlock(ctx, true)
		}()
	}

	fmt.Printf(
		"[%v] %v: %v \n",
		time.Now().Format("15:04:05"),
		lockPrintName,
		msg,
	)

	return exCtx, err
}
