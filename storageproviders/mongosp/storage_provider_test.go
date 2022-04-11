package mongosp

import (
	"context"
	"log"
	"testing"
	"time"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/nibi8/dlocker/models"
)

// todo: add tests

func TestStorageProvider(t *testing.T) {

	ctx := context.Background()

	// connect to mongodb
	constr := "mongodb://localhost:27017"
	
	constrEnv, envFound := os.LookupEnv("MONGO_CON_STR")
	if envFound {
		constr = constrEnv
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(constr))
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
		log.Fatal("client.Ping")
	}

	db := client.Database("test")

	collectionName := "lockerTest"

	_ = db.Collection(collectionName).Drop(ctx)

	// create storage provider
	sp, err := NewStorageProvider(ctx, db, collectionName)
	if err != nil {
		t.Error(err)
		return
	}

	lockName := "lock1"

	lock, err := models.NewLock(lockName, 20, 10)
	if err != nil {
		t.Error(err)
		return
	}

	lr := models.NewLockRecord(lock)
	lr.Dt = time.Unix(time.Now().Unix(), 0).UTC() // fix golang and db format for later compare in tests

	err = sp.CreateLockRecord(ctx, lr)
	if err != nil {
		t.Error(err)
		return
	}

	lrResp, err := sp.GetLockRecord(ctx, lock.Name)
	if err != nil {
		t.Error(err)
		return
	}

	if lr != lrResp {
		t.Errorf("GetLockRecord result differs")
		return
	}

	lrPatch := models.NewLockRecordPatchForCapture(lr.DurationSec)
	lrPatch.Dt = time.Unix(time.Now().Unix(), 0).UTC() // fix golang and db format for later compare in tests

	err = sp.UpdateLockRecord(ctx, lock.Name, lr.Version, lrPatch)
	if err != nil {
		t.Error(err)
		return
	}

	lrUpdated, err := sp.GetLockRecord(ctx, lock.Name)
	if err != nil {
		t.Error(err)
		return
	}

	lr.ApplyPatch(lrPatch)

	if lr != lrUpdated {
		t.Errorf("GetLockRecord (after UpdateLockRecord) result differs")
		return
	}

}
