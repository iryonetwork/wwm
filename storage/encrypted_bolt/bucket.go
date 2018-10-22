package bolt

import (
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/coreos/bbolt"
)

type b struct {
	*bolt.Bucket
}

type Bucket struct {
	*b
	aesgcm cipher.AEAD
}

func (eb *Bucket) Bucket(name []byte) *Bucket {
	bb := eb.b.Bucket.Bucket(name)
	if bb == nil {
		return nil
	}
	return &Bucket{b: &b{bb}, aesgcm: eb.aesgcm}
}

func (eb *Bucket) CreateBucket(name []byte) (*Bucket, error) {
	nb, err := eb.b.Bucket.CreateBucket(name)
	if err != nil {
		return nil, err
	}
	return &Bucket{b: &b{nb}, aesgcm: eb.aesgcm}, nil
}

func (eb *Bucket) CreateBucketIfNotExists(name []byte) (*Bucket, error) {
	nb, err := eb.b.Bucket.CreateBucketIfNotExists(name)
	if err != nil {
		return nil, err
	}
	return &Bucket{b: &b{nb}, aesgcm: eb.aesgcm}, nil
}

func (eb *Bucket) Cursor() *Cursor {
	return &Cursor{Cursor: eb.b.Bucket.Cursor(), aesgcm: eb.aesgcm}
}

const nonceLength = 12

func (eb *Bucket) Get(key []byte) []byte {
	data := eb.b.Bucket.Get(key)

	decrypted, err := decrypt(data, eb.aesgcm)
	if err != nil {
		return nil
	}

	return decrypted
}

func (eb *Bucket) Put(key, value []byte) error {
	nonce := make([]byte, nonceLength)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	return eb.b.Bucket.Put(key, append(nonce, eb.aesgcm.Seal(nil, nonce, value, nil)...))
}

func (eb *Bucket) ForEach(fn func(k, v []byte) error) error {
	return eb.b.Bucket.ForEach(func(k, v []byte) error {
		// value is bucket
		if v == nil {
			return fn(k, v)
		}

		decrypted, err := decrypt(v, eb.aesgcm)
		if err != nil {
			return ErrDecrypt
		}
		return fn(k, decrypted)
	})
}

func decrypt(data []byte, aesgcm cipher.AEAD) ([]byte, error) {
	if len(data) < nonceLength {
		return nil, ErrDecrypt
	}

	decrypted, err := aesgcm.Open(nil, data[:nonceLength], data[nonceLength:], nil)
	if err != nil {
		return nil, err
	}
	if decrypted == nil {
		return []byte{}, nil
	}

	return decrypted, nil
}
