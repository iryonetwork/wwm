package bolt

import "github.com/coreos/bbolt"

type Cursor struct {
	*bolt.Cursor
	encryptionKey *[]byte
}

func (c *Cursor) Bucket() *Bucket {
	return &Bucket{b: &b{Bucket: c.Cursor.Bucket()}, encryptionKey: c.encryptionKey}
}

func (c *Cursor) First() (key, value []byte) {
	k, v := c.Cursor.First()
	decrypted, _ := decrypt(v, *c.encryptionKey)
	return k, decrypted
}

func (c *Cursor) Last() (key, value []byte) {
	k, v := c.Cursor.Last()
	decrypted, _ := decrypt(v, *c.encryptionKey)
	return k, decrypted
}

func (c *Cursor) Next() (key, value []byte) {
	k, v := c.Cursor.Next()
	decrypted, _ := decrypt(v, *c.encryptionKey)
	return k, decrypted
}

func (c *Cursor) Prev() (key, value []byte) {
	k, v := c.Cursor.Prev()
	decrypted, _ := decrypt(v, *c.encryptionKey)
	return k, decrypted
}

func (c *Cursor) Seek(seek []byte) (key, value []byte) {
	k, v := c.Cursor.Seek(seek)
	decrypted, _ := decrypt(v, *c.encryptionKey)
	return k, decrypted
}
