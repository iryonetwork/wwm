package bolt

import (
	"crypto/aes"
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
	encryptionKey *[]byte
}

func (eb *Bucket) Bucket(name []byte) *Bucket {
	bb := eb.b.Bucket.Bucket(name)
	if bb == nil {
		return nil
	}
	return &Bucket{b: &b{bb}, encryptionKey: eb.encryptionKey}
}

func (eb *Bucket) CreateBucket(name []byte) (*Bucket, error) {
	nb, err := eb.b.Bucket.CreateBucket(name)
	if err != nil {
		return nil, err
	}
	return &Bucket{b: &b{nb}, encryptionKey: eb.encryptionKey}, nil
}

func (eb *Bucket) CreateBucketIfNotExists(name []byte) (*Bucket, error) {
	nb, err := eb.b.Bucket.CreateBucketIfNotExists(name)
	if err != nil {
		return nil, err
	}
	return &Bucket{b: &b{nb}, encryptionKey: eb.encryptionKey}, nil
}

func (eb *Bucket) Cursor() *Cursor {
	return &Cursor{Cursor: eb.b.Bucket.Cursor(), encryptionKey: eb.encryptionKey}
}

const nonceLength = 12

func (eb *Bucket) Get(key []byte) []byte {
	data := eb.b.Bucket.Get(key)

	decrypted, err := decrypt(data, *eb.encryptionKey)
	if err != nil {
		return nil
	}

	return decrypted
}

func (eb *Bucket) Put(key, value []byte) error {
	if eb.encryptionKey == nil {
		return ErrTxClosed
	}

	block, err := aes.NewCipher(*eb.encryptionKey)
	if err != nil {
		return err
	}

	nonce := make([]byte, nonceLength)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	return eb.b.Bucket.Put(key, append(nonce, aesgcm.Seal(nil, nonce, value, nil)...))
}

func (eb *Bucket) ForEach(fn func(k, v []byte) error) error {
	return eb.b.Bucket.ForEach(func(k, v []byte) error {
		// value is bucket
		if v == nil {
			return fn(k, v)
		}

		decrypted, err := decrypt(v, *eb.encryptionKey)
		if err != nil {
			return ErrDecrypt
		}
		return fn(k, decrypted)
	})
}

func decrypt(data, key []byte) ([]byte, error) {
	if len(data) < nonceLength {
		return nil, ErrDecrypt
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
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
