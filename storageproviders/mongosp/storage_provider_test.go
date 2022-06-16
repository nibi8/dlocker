package mongosp

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"github.com/nibi8/dlocker/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testDBName = "test"
const testCollectionName = "lockerTest"

func TestInit(t *testing.T) {
	ctx := context.Background()

	db := initDB(ctx, t)

	_ = db.Collection(testCollectionName).Drop(ctx)

	sp, err := NewStorageProvider(ctx, db, testCollectionName)
	require.NoError(t, err)

	assert.NotEmpty(t, sp.db)
	assert.NotEmpty(t, sp.collectionName)

	cursor, err := db.Collection(testCollectionName).Indexes().List(ctx)
	require.NoError(t, err)

	indexFound := false
	index := mongo.IndexSpecification{}
	for cursor.Next(ctx) {
		err = cursor.Decode(&index)
		require.NoError(t, err)
		v := index.KeysDocument.Lookup("jobname")
		if v.Value != nil {
			indexFound = true
			assert.True(t, *index.Unique)
		}
	}

	assert.True(t, indexFound)
}

func TestCreate(t *testing.T) {
	ctx := context.Background()

	db := initDB(ctx, t)

	_ = db.Collection(testCollectionName).Drop(ctx)

	sp, err := NewStorageProvider(ctx, db, testCollectionName)
	require.NoError(t, err)

	lock, err := models.NewLock("unique-lock-name", 60, 10)
	require.NoError(t, err)

	lr := models.NewLockRecord(lock)

	err = sp.CreateLockRecord(ctx, lr)
	require.NoError(t, err)

	err = sp.CreateLockRecord(ctx, lr)
	assert.True(t, errors.Is(err, models.ErrDuplicate))

}

func TestGet(t *testing.T) {
	ctx := context.Background()

	db := initDB(ctx, t)

	_ = db.Collection(testCollectionName).Drop(ctx)

	sp, err := NewStorageProvider(ctx, db, testCollectionName)
	require.NoError(t, err)

	lock, err := models.NewLock("unique-lock-name", 60, 10)
	require.NoError(t, err)

	_, err = sp.GetLockRecord(ctx, lock.Name)
	assert.True(t, errors.Is(err, models.ErrNotFound))

	lr := models.NewLockRecord(lock)
	// fix for db representation
	lr.Dt = time.Unix(time.Now().Unix(), 0).UTC()

	err = sp.CreateLockRecord(ctx, lr)
	require.NoError(t, err)

	lrRed, err := sp.GetLockRecord(ctx, lock.Name)
	require.NoError(t, err)

	assert.Equal(t, lr, lrRed)
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()

	db := initDB(ctx, t)

	_ = db.Collection(testCollectionName).Drop(ctx)

	sp, err := NewStorageProvider(ctx, db, testCollectionName)
	require.NoError(t, err)

	lock, err := models.NewLock("unique-lock-name", 60, 10)
	require.NoError(t, err)

	lr := models.NewLockRecord(lock)

	patch := models.NewLockRecordPatchForCapture(lr.DurationSec)
	err = sp.UpdateLockRecord(ctx, lr.LockName, lr.Version, patch)
	require.True(t, errors.Is(err, models.ErrNotFound))

	err = sp.CreateLockRecord(ctx, lr)
	require.NoError(t, err)

	err = sp.UpdateLockRecord(ctx, lr.LockName, lr.Version, patch)
	require.NoError(t, err)
}

func initDB(ctx context.Context, t *testing.T) *mongo.Database {

	constr := "mongodb://localhost:27017"

	constrEnv, envFound := os.LookupEnv("MONGO_CON_STR")
	if envFound {
		constr = constrEnv
	}

	opts := options.Client().ApplyURI(constr)

	// recommended option to prevent collisions
	opts = opts.SetWriteConcern(writeconcern.New(writeconcern.WMajority()))

	client, err := mongo.Connect(ctx, opts)
	require.NoError(t, err)

	err = client.Ping(ctx, nil)
	require.NoError(t, err)

	db := client.Database(testDBName)

	return db
}
