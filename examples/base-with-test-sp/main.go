package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nibi8/dlocker"
	"github.com/nibi8/dlocker/models"
	"github.com/nibi8/dlocker/storageproviders/testsp"
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

	lock2, err := models.NewLock("unique_lock_name_2", 10, 5)
	if err != nil {
		log.Fatal("NewLock")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(4)

	// lock1:
	go func() {
		defer wg.Done()
		captureLock(ctx, locker, lock1, "instace_1", false)

	}()
	go func() {
		defer wg.Done()
		time.Sleep(1 * time.Second)
		captureLock(ctx, locker, lock1, "instace_2", false)
	}()

	// lock2 with manual unlock:
	go func() {
		defer wg.Done()
		captureLock(ctx, locker, lock2, "instace_1", true)
	}()
	go func() {
		defer wg.Done()
		time.Sleep(1 * time.Second)
		captureLock(ctx, locker, lock2, "instace_2", true)
	}()

	wg.Wait()

	//<-ctx.Done()

}

func captureLock(
	ctx context.Context,
	locker dlocker.Locker,
	lock models.Lock,
	instanceName string,
	unlock bool,
) {

	lockPrintName := lock.Name
	if instanceName != "" {
		lockPrintName += " " + instanceName
	}

	msg := "lock success"
	lockCtx, _, err := locker.LockWithWait(ctx, lock)
	if err != nil {
		msg = fmt.Sprintf("lock failed with error = %v", err)
	}
	if unlock {
		defer locker.Unlock(ctx, lockCtx)
	}

	fmt.Printf(
		"[%v] %v: %v \n",
		time.Now().Format("15:04:05"),
		lockPrintName,
		msg,
	)

}
