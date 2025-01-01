package db

import (
	"fmt"

	"github.com/boltdb/bolt"
)

var defaultBucket = []byte("Default")

// Bolt Database
type Database struct {
	db *bolt.DB
}

// Db instance
func NewDatabase(dbPath string) (db *Database, closeFunc func() error, err error) {
	boltDb, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, nil, err
	}
	closeFunc = boltDb.Close
	db = &Database{db: boltDb}

	if err := db.createDefaultBucket(); err != nil {
		closeFunc()
		return nil, nil, fmt.Errorf("error While creating Default bucket: %w", err)
	}

	return &Database{db: boltDb}, closeFunc, nil

}

// Default Bucket Creation
func (d *Database) createDefaultBucket() error {
	return d.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(defaultBucket)
		return err
	})
}

// Seting Key
func (d *Database) SetKey(key string, value []byte) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)
		return b.Put([]byte(key), value)
	})
}

// Getting value
func (d *Database) GetKey(key string) ([]byte, error) {
	var ans []byte
	err := d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)
		ans = b.Get([]byte(key))
		return nil
	})

	if err == nil {
		return ans, nil
	}

	return nil, err

}
