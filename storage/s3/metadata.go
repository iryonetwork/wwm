package s3

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/iryonetwork/wwm/gen/storage/models"
	"github.com/iryonetwork/wwm/storage/s3/object"
)

type metadata struct {
	filename  string
	version   string
	operation Operation
	created   time.Time
	checksum  string
	archetype string
}

var utc, _ = time.LoadLocation("UTC")

func metadataFromKey(key string) (*metadata, error) {
	items := strings.SplitN(key, ".", 6)
	md := &metadata{
		filename:  items[0],
		version:   items[1],
		operation: Operation(items[2]),
		checksum:  items[4],
		archetype: items[5],
	}

	// validate operation
	if md.operation != Write && md.operation != Delete {
		return nil, fmt.Errorf("Invalid operation %s", md.operation)
	}

	// convert timestamp
	if len(items[3]) < 13 {
		return nil, fmt.Errorf("Invalid timestamp length (%s, %d)", items[3], len(items[3]))
	}
	s, err := strconv.ParseInt(items[3], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse timestamp (%s)", items[3])
	}

	md.created = time.Unix(s/1000, s%1000*1000000).In(utc)

	return md, nil
}

func metadataFromNewFile(newFile *object.NewObjectInfo) (*metadata, error) {
	md := &metadata{
		filename:  newFile.Name,
		version:   newFile.Version,
		operation: Operation(newFile.Operation),
		checksum:  newFile.Checksum,
		archetype: newFile.Archetype,
		created:   time.Time(newFile.Created),
	}

	// validate operation
	if md.operation != Write && md.operation != Delete {
		return nil, fmt.Errorf("Invalid operation %s", md.operation)
	}

	return md, nil
}

func metadataFromFileDescriptor(fd *models.FileDescriptor) (*metadata, error) {
	md := &metadata{
		filename:  fd.Name,
		version:   fd.Version,
		operation: Operation(fd.Operation),
		checksum:  fd.Checksum,
		archetype: fd.Archetype,
	}

	// parse created
	c, err := time.Parse(strfmt.RFC3339Millis, fd.Created.String())
	if err != nil {
		return nil, err
	}
	md.created = c

	return md, nil
}

func (m *metadata) String() string {
	return fmt.Sprintf("%s.%s.%s.%d.%s.%s",
		m.filename,
		m.version,
		m.operation,
		m.created.UnixNano()/1000000,
		m.checksum,
		m.archetype)
}

func (m *metadata) FileDescriptor() *models.FileDescriptor {
	return &models.FileDescriptor{
		Name:      m.filename,
		Version:   m.version,
		Operation: string(m.operation),
		Created:   strfmt.DateTime(m.created),
		Archetype: m.archetype,
		Checksum:  m.checksum,
	}
}
