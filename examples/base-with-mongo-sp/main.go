package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"github.com/nibi8/dlocker"
	"github.com/nibi8/dlocker/models"
	"github.com/nibi8/dlocker/storageproviders/mongosp"
)

func main() {

	ctx := context.Background()

	// connect to mongodb
	constr := "mongodb://localhost:27017"

	constrEnv, envFound := os.LookupEnv("MONGO_CON_STR")
	if envFound {
		constr = constrEnv
	}

	opts := options.Client().ApplyURI(constr)
	// recommended option to prevent collisions
	opts = opts.SetWriteConcern(writeconcern.New(writeconcern.WMajority()))
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal("mongo.Connect")
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println(err)
		log.Fatal("client.Ping")
	}

	db := client.Database("test")

	// create storage provider
	sp, err := mongosp.NewStorageProvider(ctx, db, "lockerTest")
	if err != nil {
		log.Fatal("mongosp.NewStorageProvider")
	}

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
