package filesDataExporter

import "github.com/iryonetwork/wwm/gen/storage/models"

// byCreated implements sort.Interface for []*model.FileDescriptor based on
// the Created field with reverse order (timestamp ascending)
type descByCreated []*models.FileDescriptor

func (c descByCreated) Len() int           { return len(c) }
func (c descByCreated) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c descByCreated) Less(i, j int) bool { return c[i].Created.String() > c[j].Created.String() }
