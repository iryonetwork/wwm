package bolt

import (
	"github.com/coreos/bbolt"
)

type Tx struct {
	*bolt.Tx
	encryptionKey *[]byte
}

func (tx *Tx) DB() *DB {
	return &DB{DB: tx.Tx.DB(), encryptionKey: tx.encryptionKey}
}

func (tx *Tx) Bucket(name []byte) *Bucket {
	bb := tx.Tx.Bucket(name)
	if bb == nil {
		return nil
	}
	return &Bucket{b: &b{Bucket: bb}, encryptionKey: tx.encryptionKey}
}

func (tx *Tx) CreateBucket(name []byte) (*Bucket, error) {
	bb, err := tx.Tx.CreateBucket(name)
	if err != nil {
		return nil, err
	}
	return &Bucket{b: &b{Bucket: bb}, encryptionKey: tx.encryptionKey}, nil
}

func (tx *Tx) CreateBucketIfNotExists(name []byte) (*Bucket, error) {
	bb, err := tx.Tx.CreateBucketIfNotExists(name)
	if err != nil {
		return nil, err
	}
	return &Bucket{b: &b{Bucket: bb}, encryptionKey: tx.encryptionKey}, nil
}

func (tx *Tx) ForEach(fn func(name []byte, b *Bucket) error) error {
	return tx.Tx.ForEach(func(name []byte, bb *bolt.Bucket) error {
		return fn(name, &Bucket{b: &b{Bucket: bb}, encryptionKey: tx.encryptionKey})
	})
}

func (tx *Tx) Cursor() *Cursor {
	return &Cursor{Cursor: tx.Tx.Cursor(), encryptionKey: tx.encryptionKey}
}
