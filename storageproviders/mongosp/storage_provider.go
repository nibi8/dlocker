package mongosp

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/p8bin/dlocker/models"
)

type StorageProvider struct {
	db             *mongo.Database
	collectionName string
}

func NewStorageProvider(
	ctx context.Context,
	db *mongo.Database,
	collectionName string,
) (*StorageProvider, error) {

	sp := StorageProvider{
		db:             db,
		collectionName: collectionName,
	}

	err := sp.createIndexes(ctx)
	if err != nil {
		return nil, err
	}

	return &sp, nil
}

func (sp *StorageProvider) createIndexes(
	ctx context.Context,
) (err error) {

	_, err = sp.db.Collection(sp.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{
			"jobname": 1,
		},
		Options: options.Index().
			SetUnique(true).
			SetName("unique: jobname"),
	})
	if err != nil {
		return err
	}

	return nil
}

func (sp *StorageProvider) GetLockRecord(
	ctx context.Context,
	lockName string,
) (lr models.LockRecord, err error) {

	res := sp.db.Collection(sp.collectionName).FindOne(ctx, bson.M{
		"jobname": lockName,
	})
	lrDB := LockRecordDB{}
	err = res.Decode(&lrDB)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return lr, models.ErrNotFound
		}

		return lr, err
	}

	lr = ToLockRecord(lrDB)

	return lr, nil
}

func (sp *StorageProvider) CreateLockRecord(
	ctx context.Context,
	lr models.LockRecord,
) (err error) {

	lrDB := FromLockRecord(lr)

	_, err = sp.db.Collection(sp.collectionName).InsertOne(ctx, lrDB)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return models.ErrDuplicate
		}
		return err
	}

	return nil
}

func (sp *StorageProvider) UpdateLockRecord(
	ctx context.Context,
	lockName string,
	version string,
	patch models.LockRecordPatch,
) (err error) {
	ures := sp.db.Collection(sp.collectionName).FindOneAndUpdate(ctx,
		bson.M{
			"jobname": lockName,
			"version": version,
		},
		bson.M{"$set": bson.M{
			"version":     patch.Version,
			"durationsec": patch.DurationSec,
			"state":       patch.State,
			"dt":          patch.Dt,
		}},
	)
	err = ures.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.ErrNotFound
		}
		return err
	}

	return nil
}
