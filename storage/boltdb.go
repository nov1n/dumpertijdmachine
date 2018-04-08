package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"github.com/nov1n/dumpertijdmachine/types"
)

type BoltDB struct {
	db     *bolt.DB
	path   string
	bucket string
}

func NewBoltDB(path string, bucket string) (*BoltDB, error) {
	// create storage if not exist
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return nil, err
	}

	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	return &BoltDB{
		path:   path,
		bucket: bucket,
		db:     db,
	}, nil
}

func (s *BoltDB) PutDay(date time.Time, day *types.Day) error {
	data, err := proto.Marshal(day)
	if err != nil {
		return fmt.Errorf("unmarshal day: %s", err)
	}

	err = s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(s.bucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		b.Put([]byte(KeyForDate(date)), data)
		return nil
	})
	return err
}

func (s *BoltDB) GetDay(date []byte) (*types.Day, error) {
	day := &types.Day{}
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(s.bucket))
		if bucket == nil {
			return fmt.Errorf("bucket for date %s did not exist", date)
		}

		data := bucket.Get(date)
		if data == nil {
			return fmt.Errorf("day not found for date %s", date)
		}

		return proto.Unmarshal(data, day)
	})
	return day, err
}

func (s *BoltDB) Close() error {
	return s.db.Close()
}
