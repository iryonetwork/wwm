package bolt

import (
	"github.com/coreos/bbolt"
)

type Tx struct {
	*bolt.Tx
	encryptionKey *[]byte
}

func (tx *Tx) DB() *DB {
	return &DB{tx.Tx.DB(), tx.encryptionKey}
}

func (tx *Tx) Bucket(name []byte) *Bucket {
	bb := tx.Tx.Bucket(name)
	if bb == nil {
		return nil
	}
	return &Bucket{&b{bb}, tx.encryptionKey}
}

func (tx *Tx) CreateBucket(name []byte) (*Bucket, error) {
	bb, err := tx.Tx.CreateBucket(name)
	if err != nil {
		return nil, err
	}
	return &Bucket{&b{bb}, tx.encryptionKey}, nil
}

func (tx *Tx) CreateBucketIfNotExists(name []byte) (*Bucket, error) {
	bb, err := tx.Tx.CreateBucketIfNotExists(name)
	if err != nil {
		return nil, err
	}
	return &Bucket{&b{bb}, tx.encryptionKey}, nil
}

func (tx *Tx) ForEach(fn func(name []byte, b *Bucket) error) error {
	return tx.Tx.ForEach(func(name []byte, bb *bolt.Bucket) error {
		return fn(name, &Bucket{&b{bb}, tx.encryptionKey})
	})
}

func (tx *Tx) Cursor() *Cursor {
	return &Cursor{tx.Tx.Cursor(), tx.encryptionKey}
}
