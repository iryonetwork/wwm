package bolt

import (
	"crypto/cipher"

	"github.com/coreos/bbolt"
)

type Tx struct {
	*bolt.Tx
	aesgcm cipher.AEAD
}

func (tx *Tx) DB() *DB {
	return &DB{DB: tx.Tx.DB(), aesgcm: tx.aesgcm}
}

func (tx *Tx) Bucket(name []byte) *Bucket {
	bb := tx.Tx.Bucket(name)
	if bb == nil {
		return nil
	}
	return &Bucket{b: &b{Bucket: bb}, aesgcm: tx.aesgcm}
}

func (tx *Tx) CreateBucket(name []byte) (*Bucket, error) {
	bb, err := tx.Tx.CreateBucket(name)
	if err != nil {
		return nil, err
	}
	return &Bucket{b: &b{Bucket: bb}, aesgcm: tx.aesgcm}, nil
}

func (tx *Tx) CreateBucketIfNotExists(name []byte) (*Bucket, error) {
	bb, err := tx.Tx.CreateBucketIfNotExists(name)
	if err != nil {
		return nil, err
	}
	return &Bucket{b: &b{Bucket: bb}, aesgcm: tx.aesgcm}, nil
}

func (tx *Tx) ForEach(fn func(name []byte, b *Bucket) error) error {
	return tx.Tx.ForEach(func(name []byte, bb *bolt.Bucket) error {
		return fn(name, &Bucket{b: &b{Bucket: bb}, aesgcm: tx.aesgcm})
	})
}

func (tx *Tx) Cursor() *Cursor {
	return &Cursor{Cursor: tx.Tx.Cursor(), aesgcm: tx.aesgcm}
}
