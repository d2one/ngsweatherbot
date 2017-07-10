package main

import "github.com/boltdb/bolt"

func (this *BoltStorage) writer() {
	for data := range this.writerChan {
		bucket := data[0].(string)
		keyId := data[1].(string)
		dataBytes := data[2].([]byte)
		err := this.DB.Update(func(tx *bolt.Tx) error {
			sesionBucket, err := tx.CreateBucketIfNotExists([]byte(bucket))
			if err != nil {
				return err
			}
			sesionBucket.Delete([]byte(keyId))
			return sesionBucket.Put([]byte(keyId), dataBytes)
		})
		if err != nil {
			// TODO: Handle instead of panic
			panic(err)
		}
	}
}

func NewBoltStorage(dbPath string) *BoltStorage {
	db, err := bolt.Open(dbPath, 0666, nil)
	writerChan := make(chan [3]interface{})
	boltStorage := &BoltStorage{DB: db, writerChan: writerChan}
	boltStorage.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return err
		}
		return nil
	})
	go boltStorage.writer()
	if err != nil {
		panic(err)
	}
	return boltStorage
}