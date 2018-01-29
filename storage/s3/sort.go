package s3

import "github.com/iryonetwork/wwm/gen/storage/models"

// byCreated implements sort.Interface for []*model.FileDescriptor based on
// the Created field.
type byCreated []*models.FileDescriptor

func (c byCreated) Len() int           { return len(c) }
func (c byCreated) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c byCreated) Less(i, j int) bool { return c[i].Created.String() > c[j].Created.String() }
