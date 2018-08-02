package reports

//go:generate ../bin/mockgen.sh reports Storage $GOFILE

import (
	"github.com/go-openapi/strfmt"
)

type (
	// Storage describes reports storage public API
	Storage interface {
		// Insert inserts new file or updated file data
		Insert(fileID, version, patientID string, timestamp strfmt.DateTime, data string) error

		// Remove deletes file data
		Remove(fileID string) error

		// Get looks up file by ID and version (optional)
		Get(fileID, version string) (*File, error)

		// Exists check if file with given ID and version (optional) already exists in DB
		Exists(fileID, version string) (bool, error)

		// Find looks up files by patientID, data key values, createdStart, createdEnd
		Find(patientID string, keyValues map[string]string, createdStart *strfmt.DateTime, createdAt *strfmt.DateTime) (*[]File, error)

		// Close closes the storage connection
		Close() error
	}

	File struct {
		FileID    string
		Version   string
		PatientID string
		CreatedAt strfmt.DateTime
		UpdatedAt strfmt.DateTime
		Data      string
	}
)
