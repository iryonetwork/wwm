package storage

import (
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/storage/models"
	operations "github.com/iryonetwork/wwm/gen/storage/restapi/operations/storage"
)

// Handlers describes the actions supported by the storage handlers
type Handlers interface {
	FileList() operations.FileListHandler
	FileGet() operations.FileGetHandler
	FileGetVersion() operations.FileGetVersionHandler
	FileListVersions() operations.FileListVersionsHandler
	FileNew() operations.FileNewHandler
	FileUpdate() operations.FileUpdateHandler
	FileDelete() operations.FileDeleteHandler
	SyncFileMetadata() operations.SyncFileMetadataHandler
	SyncFile() operations.SyncFileHandler
	SyncFileDelete() operations.SyncFileDeleteHandler
	Authorizer() runtime.Authorizer
	GetUserIDFromToken(token string) (*string, error)
}

type handlers struct {
	service Service
	logger  zerolog.Logger
}

func (h *handlers) FileList() operations.FileListHandler {
	return operations.FileListHandlerFunc(func(params operations.FileListParams, principal *string) middleware.Responder {
		list, err := h.service.FileList(params.Bucket)
		if err != nil {
			return operations.NewFileListInternalServerError().WithPayload(&models.Error{
				Code:    "server_error",
				Message: err.Error(),
			})
		}

		return operations.NewFileListOK().WithPayload(list)
	})
}

func (h *handlers) FileGet() operations.FileGetHandler {
	return operations.FileGetHandlerFunc(func(params operations.FileGetParams, principal *string) middleware.Responder {
		r, fd, err := h.service.FileGet(params.Bucket, params.FileID)
		if err != nil {
			switch err {
			case ErrNotFound:
				return operations.NewFileGetNotFound()
			default:
				return operations.NewFileGetInternalServerError().WithPayload(&models.Error{
					Code:    "server_error",
					Message: err.Error(),
				})
			}
		}

		return operations.NewFileGetOK().
			WithPayload(r).
			WithContentType(fd.ContentType).
			WithXCreated(fd.Created).
			WithXVersion(fd.Version).
			WithXArchetype(fd.Archetype).
			WithXChecksum(fd.Checksum).
			WithXName(fd.Name).
			WithXPath(fd.Path)
	})
}

func (h *handlers) FileGetVersion() operations.FileGetVersionHandler {
	return operations.FileGetVersionHandlerFunc(func(params operations.FileGetVersionParams, principal *string) middleware.Responder {
		r, fd, err := h.service.FileGetVersion(params.Bucket, params.FileID, params.Version)
		if err != nil {
			switch err {
			case ErrNotFound:
				return operations.NewFileGetVersionNotFound()
			default:
				return operations.NewFileGetVersionInternalServerError().WithPayload(&models.Error{
					Code:    "server_error",
					Message: err.Error(),
				})
			}
		}

		return operations.NewFileGetVersionOK().
			WithPayload(r).
			WithContentType(fd.ContentType).
			WithXCreated(fd.Created).
			WithXVersion(fd.Version).
			WithXArchetype(fd.Archetype).
			WithXChecksum(fd.Checksum).
			WithXName(fd.Name).
			WithXPath(fd.Path)
	})
}

func (h *handlers) FileListVersions() operations.FileListVersionsHandler {
	return operations.FileListVersionsHandlerFunc(func(params operations.FileListVersionsParams, principal *string) middleware.Responder {
		list, err := h.service.FileListVersions(params.Bucket, params.FileID)
		if err != nil {
			return operations.NewFileListVersionsInternalServerError().WithPayload(&models.Error{
				Code:    "server_error",
				Message: err.Error(),
			})
		}

		return operations.NewFileListVersionsOK().WithPayload(list)
	})
}

func (h *handlers) FileNew() operations.FileNewHandler {
	return operations.FileNewHandlerFunc(func(params operations.FileNewParams, principal *string) middleware.Responder {
		var archetype string
		if params.Archetype != nil {
			archetype = *params.Archetype
		}
		defer params.File.Close()

		// call service
		fd, err := h.service.FileNew(params.Bucket, params.File, params.ContentType, archetype)
		if err != nil {
			return operations.NewFileNewInternalServerError().WithPayload(&models.Error{
				Code:    "server_error",
				Message: err.Error(),
			})
		}

		return operations.NewFileNewCreated().WithPayload(fd)
	})
}

func (h *handlers) FileUpdate() operations.FileUpdateHandler {
	return operations.FileUpdateHandlerFunc(func(params operations.FileUpdateParams, principal *string) middleware.Responder {
		var archetype string
		if params.Archetype != nil {
			archetype = *params.Archetype
		}
		defer params.File.Close()

		// call service
		fd, err := h.service.FileUpdate(params.Bucket, params.FileID, params.File, params.ContentType, archetype)
		if err != nil {
			switch err {
			case ErrNotFound:
				return operations.NewFileUpdateNotFound()
			default:
				return operations.NewFileUpdateInternalServerError().WithPayload(&models.Error{
					Code:    "server_error",
					Message: err.Error(),
				})
			}
		}

		return operations.NewFileUpdateCreated().WithPayload(fd)
	})
}

func (h *handlers) FileDelete() operations.FileDeleteHandler {
	return operations.FileDeleteHandlerFunc(func(params operations.FileDeleteParams, principal *string) middleware.Responder {
		err := h.service.FileDelete(params.Bucket, params.FileID)
		if err != nil {
			switch err {
			case ErrNotFound:
				return operations.NewFileDeleteNotFound()
			default:
				return operations.NewFileDeleteInternalServerError().WithPayload(&models.Error{
					Code:    "server_error",
					Message: err.Error(),
				})
			}
		}

		return operations.NewFileDeleteNoContent()
	})
}

func (h *handlers) SyncFileMetadata() operations.SyncFileMetadataHandler {
	return operations.SyncFileMetadataHandlerFunc(func(params operations.SyncFileMetadataParams, principal *string) middleware.Responder {
		_, fd, err := h.service.FileGetVersion(params.Bucket, params.FileID, params.Version)
		if err != nil {
			switch err {
			case ErrNotFound:
				return operations.NewSyncFileMetadataNotFound()
			default:
				h.logger.Error().Err(err).Msg("Failed to fetch the file to return metadata")
				return operations.NewSyncFileMetadataInternalServerError()
			}
		}

		return operations.NewSyncFileMetadataOK().
			WithContentType(fd.ContentType).
			WithXCreated(fd.Created).
			WithXVersion(fd.Version).
			WithXArchetype(fd.Archetype).
			WithXChecksum(fd.Checksum).
			WithXName(fd.Name).
			WithXPath(fd.Path)
	})
}

func (h *handlers) SyncFile() operations.SyncFileHandler {
	return operations.SyncFileHandlerFunc(func(params operations.SyncFileParams, principal *string) middleware.Responder {
		var archetype string
		if params.Archetype != nil {
			archetype = *params.Archetype
		}
		defer params.File.Close()

		// call service
		fd, err := h.service.SyncFile(params.Bucket, params.FileID, params.Version, params.File, params.ContentType, params.Created, archetype)
		if err != nil {
			switch err {
			case ErrAlreadyExists:
				return operations.NewSyncFileOK().WithPayload(fd)
			case ErrAlreadyExistsConflict:
				return operations.NewSyncFileConflict()
			default:
				return operations.NewSyncFileInternalServerError().WithPayload(&models.Error{
					Code:    "server_error",
					Message: err.Error(),
				})
			}
		}

		return operations.NewSyncFileCreated().WithPayload(fd)
	})
}

func (h *handlers) SyncFileDelete() operations.SyncFileDeleteHandler {
	return operations.SyncFileDeleteHandlerFunc(func(params operations.SyncFileDeleteParams, principal *string) middleware.Responder {
		err := h.service.SyncFileDelete(params.Bucket, params.FileID, params.Version, params.Created)
		if err != nil {
			switch err {
			case ErrNotFound:
				return operations.NewSyncFileDeleteNotFound()
			case ErrDeleted:
				return operations.NewSyncFileDeleteConflict()
			default:
				return operations.NewSyncFileDeleteInternalServerError().WithPayload(&models.Error{
					Code:    "server_error",
					Message: err.Error(),
				})
			}
		}

		return operations.NewSyncFileDeleteNoContent()
	})
}

func (h *handlers) Authorizer() runtime.Authorizer {
	return runtime.AuthorizerFunc(func(r *http.Request, principal interface{}) error {
		// @TODO
		return nil
	})
}

func (h *handlers) GetUserIDFromToken(token string) (*string, error) {
	userID := "USER"
	return &userID, nil
}

// NewHandlers returns a new instance of authenticator handlers
func NewHandlers(service Service, logger zerolog.Logger) Handlers {
	return &handlers{service: service, logger: logger}
}
