// Provides dummy implementation of KeyProvider interface from "github.com/iryonetwork/wwm/storage/s3" that always returns the same key
package keyProvider

type keyProvider struct {
	key string
}

func (k *keyProvider) Get(id string) (string, error) {
	return k.key, nil
}

func New(key string) *keyProvider {
	return &keyProvider{key}
}
