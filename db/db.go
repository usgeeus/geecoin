package db

import (
	"fmt"
	"os"

	"github.com/usgeeus/geecoin.git/utils"
	bolt "go.etcd.io/bbolt"
)

const (
	dbName       = "blockchain"
	dataBucket   = "data"
	blocksBucket = "blocks"
	checkpoint   = "checkpoint"
)

var db *bolt.DB

func getDBName() string {
	port := os.Args[2][6:]
	return fmt.Sprintf("%s_%s.db", dbName, port)
}

func DB() *bolt.DB {
	if db == nil {
		dbPointer, err := bolt.Open(getDBName(), 0600, nil)
		utils.HandlerErr(err)
		db = dbPointer
		err = db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(dataBucket))
			utils.HandlerErr(err)
			_, err = tx.CreateBucketIfNotExists([]byte(blocksBucket))
			return err
		})
		utils.HandlerErr(err)
	}
	return db
}

func Close() {
	DB().Close()
}

func SaveBlock(hash string, data []byte) {

	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		err := bucket.Put([]byte(hash), data)
		return err
	})
	utils.HandlerErr(err)
}

func SaveCheckpoint(data []byte) {
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket))
		err := bucket.Put([]byte(checkpoint), data)
		return err
	})
	utils.HandlerErr(err)
}

func Checkpoint() []byte {
	var data []byte
	DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(checkpoint))
		return nil
	})
	return data
}

func Block(hash string) []byte {
	var data []byte
	DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		data = bucket.Get([]byte(hash))
		return nil
	})
	return data
}

func EmptyBlocks() {
	DB().Update(func(tx *bolt.Tx) error {
		utils.HandlerErr(tx.DeleteBucket([]byte(blocksBucket)))
		_, err := tx.CreateBucket([]byte(blocksBucket))
		utils.HandlerErr(err)
		return nil
	})
}
