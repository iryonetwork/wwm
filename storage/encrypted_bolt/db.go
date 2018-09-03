package bolt

import (
	"errors"
	"fmt"
	"os"

	"github.com/coreos/bbolt"
)

type DB struct {
	*bolt.DB
	encryptionKey *[]byte
}

func Open(key []byte, path string, mode os.FileMode, options *Options) (*DB, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("Encryption key must be 32 bytes long")
	}

	db, err := bolt.Open(path, mode, options)
	if err != nil {
		return nil, err
	}

	return &DB{DB: db, encryptionKey: &key}, nil
}

func (db *DB) View(fn func(*Tx) error) error {
	if db.DB == nil {
		return ErrDatabaseNotOpen
	}

	return db.DB.View(func(tx *bolt.Tx) error {
		return fn(&Tx{Tx: tx, encryptionKey: db.encryptionKey})
	})
}

func (db *DB) Update(fn func(*Tx) error) error {
	if db.DB == nil {
		return ErrDatabaseNotOpen
	}

	return db.DB.Update(func(tx *bolt.Tx) error {
		return fn(&Tx{Tx: tx, encryptionKey: db.encryptionKey})
	})
}

func (db *DB) Begin(writable bool) (*Tx, error) {
	if db.DB == nil {
		return nil, ErrDatabaseNotOpen
	}

	tx, err := db.DB.Begin(writable)
	if err != nil {
		return nil, err
	}
	return &Tx{Tx: tx, encryptionKey: db.encryptionKey}, nil
}

func (db *DB) Batch(fn func(*Tx) error) error {
	return db.DB.Batch(func(tx *bolt.Tx) error {
		return fn(&Tx{Tx: tx, encryptionKey: db.encryptionKey})
	})
}

type Info = bolt.Info

type Options = bolt.Options

type PageInfo = bolt.PageInfo

type Stats = bolt.Stats

var (
	ErrDatabaseNotOpen = bolt.ErrDatabaseNotOpen
	ErrDatabaseOpen    = bolt.ErrDatabaseOpen
	ErrInvalid         = bolt.ErrInvalid
	ErrVersionMismatch = bolt.ErrVersionMismatch
	ErrChecksum        = bolt.ErrChecksum
	ErrTimeout         = bolt.ErrTimeout

	ErrTxNotWritable    = bolt.ErrTxNotWritable
	ErrTxClosed         = bolt.ErrTxClosed
	ErrDatabaseReadOnly = bolt.ErrDatabaseReadOnly

	ErrBucketNotFound     = bolt.ErrBucketNotFound
	ErrBucketExists       = bolt.ErrBucketExists
	ErrBucketNameRequired = bolt.ErrBucketNameRequired
	ErrKeyRequired        = bolt.ErrKeyRequired
	ErrKeyTooLarge        = bolt.ErrKeyTooLarge
	ErrValueTooLarge      = bolt.ErrValueTooLarge
	ErrIncompatibleValue  = bolt.ErrIncompatibleValue

	ErrDecrypt = errors.New("could not decrypt data")
)

const (
	MaxKeySize   = bolt.MaxKeySize
	MaxValueSize = bolt.MaxValueSize

	DefaultMaxBatchSize  = bolt.DefaultMaxBatchSize
	DefaultMaxBatchDelay = bolt.DefaultMaxBatchDelay
	DefaultAllocSize     = bolt.DefaultAllocSize

	DefaultFillPercent = bolt.DefaultFillPercent
	IgnoreNoSync       = bolt.IgnoreNoSync
)
