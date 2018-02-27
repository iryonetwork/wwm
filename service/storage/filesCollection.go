package storage

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/iryonetwork/wwm/gen/storage/models"
)

const labelFilesCollection = "filesCollection"

type filesCollection map[string]models.FileDescriptor

type filesCollectionFile []models.FileDescriptor

func (c *filesCollection) Update(fd *models.FileDescriptor) {
	// files collection holds only latest version of the file so we can safely overwrite if passed file is newer
	old, ok := (*c)[fd.Name]
	if !ok {
		(*c)[fd.Name] = *fd
		return
	}

	if old.Created.String() < fd.Created.String() {
		(*c)[fd.Name] = *fd
	}

	return
}

func (c *filesCollection) Remove(fd *models.FileDescriptor) {
	delete(*c, fd.Name)
}

func (c *filesCollection) GetFile() (io.Reader, error) {
	// transform map to slice and encode
	files := filesCollectionFile{}
	for _, fd := range *c {
		files = append(files, fd)
	}

	b, err := json.Marshal(files)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}

func FilesCollection(r io.ReadCloser) (*filesCollection, error) {
	filesMap := filesCollection{}
	files := filesCollectionFile{}

	d := json.NewDecoder(r)
	err := d.Decode(&files)
	if err != nil {
		return nil, err
	}

	for _, fd := range files {
		filesMap[fd.Name] = fd
	}

	return &filesMap, nil
}
