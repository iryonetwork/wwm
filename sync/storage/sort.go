package storage

import "github.com/iryonetwork/wwm/gen/storage/models"

// byCreated implements sort.Interface for []*model.FileDescriptor based on
// the Created field with reverse order (timestamp ascending)
type ascByCreated []*models.FileDescriptor

func (c ascByCreated) Len() int           { return len(c) }
func (c ascByCreated) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c ascByCreated) Less(i, j int) bool { return c[i].Created.String() < c[j].Created.String() }
