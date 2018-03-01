package s3

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"

	"github.com/iryonetwork/wwm/gen/storage/models"
	"github.com/iryonetwork/wwm/storage/s3/object"
)

type metadata struct {
	filename    string
	version     string
	operation   Operation
	created     time.Time
	checksum    string
	contentType string
	archetype   string
	labels      []string
}

var utc, _ = time.LoadLocation("UTC")

func metadataFromKey(key string) (*metadata, error) {
	items := strings.SplitN(key, ".", 8)

	// decode contentType
	ct, err := decode(items[5])
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode contentType from key")
	}
	// decode archetype
	arch, err := decode(items[6])
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode archetype from key")
	}
	labelsString, err := decode(items[7])
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode archetype from key")
	}
	labels := labelsStringToSlice(labelsString)

	md := &metadata{
		filename:    items[0],
		version:     items[1],
		operation:   Operation(items[2]),
		checksum:    items[4],
		contentType: ct,
		archetype:   arch,
		labels:      labels,
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
		filename:    newFile.Name,
		version:     newFile.Version,
		operation:   Operation(newFile.Operation),
		checksum:    newFile.Checksum,
		created:     time.Time(newFile.Created),
		contentType: newFile.ContentType,
		archetype:   newFile.Archetype,
		labels:      newFile.Labels,
	}

	// validate operation
	if md.operation != Write && md.operation != Delete {
		return nil, fmt.Errorf("Invalid operation %s", md.operation)
	}

	return md, nil
}

func metadataFromFileDescriptor(fd *models.FileDescriptor) (*metadata, error) {
	md := &metadata{
		filename:    fd.Name,
		version:     fd.Version,
		operation:   Operation(fd.Operation),
		checksum:    fd.Checksum,
		contentType: fd.ContentType,
		archetype:   fd.Archetype,
		labels:      fd.Labels,
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
	return fmt.Sprintf("%s.%s.%s.%d.%s.%s.%s.%s",
		m.filename,
		m.version,
		m.operation,
		m.created.UnixNano()/1000000,
		m.checksum,
		encode(m.contentType),
		encode(m.archetype),
		encode(sliceToLabelsString(m.labels)),
	)
}

func encode(src string) string {
	return base64.URLEncoding.EncodeToString([]byte(src))
}

func decode(base64URL string) (string, error) {
	b, err := base64.URLEncoding.DecodeString(base64URL)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func labelsStringToSlice(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, "|")
}

func sliceToLabelsString(s []string) string {
	return strings.Join(s, "|")
}
