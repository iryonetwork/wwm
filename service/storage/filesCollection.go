package storage

import (
	"bytes"
	"encoding/json"
	"io"
	"sort"

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
	// first, keys (filenames) are extracted and sorted alphabetically to ensure
	// deterministic order in the resulting file as order ranging over maps
	// is not in golang.
	var filenames []string
	for filename := range *c {
		filenames = append(filenames, filename)
	}
	sort.Strings(filenames)

	// transform map of file descriptors to slice for encoding
	files := filesCollectionFile{}
	for _, filename := range filenames {
		files = append(files, (*c)[filename])
	}

	var b []byte
	buf := bytes.NewBuffer(b)
	enc := json.NewEncoder(buf)

	err := enc.Encode(files)
	if err != nil {
		return nil, err
	}

	return buf, nil
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
