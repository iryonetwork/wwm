package storage

import (
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/storage/models"
	"github.com/iryonetwork/wwm/gen/storage/restapi/operations"
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
	SyncBucketList() operations.SyncBucketListHandler
	SyncFileList() operations.SyncFileListHandler
	SyncFileListVersions() operations.SyncFileListVersionsHandler
	SyncFileMetadata() operations.SyncFileMetadataHandler
	SyncFile() operations.SyncFileHandler
	SyncFileDelete() operations.SyncFileDeleteHandler
}

type handlers struct {
	service Service
	logger  zerolog.Logger
}

func (h *handlers) FileList() operations.FileListHandler {
	return operations.FileListHandlerFunc(func(params operations.FileListParams, principal *string) middleware.Responder {
		list, err := h.service.FileList(params.HTTPRequest.Context(), params.Bucket)

		if err != nil {
			return operations.NewFileListInternalServerError().WithPayload(&models.Error{
				Code:    "server_error",
				Message: err.Error(),
			})
		}
		if len(list) == 0 {
			return operations.NewFileListNotFound()
		}

		return operations.NewFileListOK().WithPayload(list)
	})
}

func (h *handlers) FileGet() operations.FileGetHandler {
	return operations.FileGetHandlerFunc(func(params operations.FileGetParams, principal *string) middleware.Responder {
		r, fd, err := h.service.FileGet(params.HTTPRequest.Context(), params.Bucket, params.FileID)

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
			WithXPath(fd.Path).
			WithXLabels(formatLabelsHeader(fd.Labels))
	})
}

func (h *handlers) FileGetVersion() operations.FileGetVersionHandler {
	return operations.FileGetVersionHandlerFunc(func(params operations.FileGetVersionParams, principal *string) middleware.Responder {
		r, fd, err := h.service.FileGetVersion(params.HTTPRequest.Context(), params.Bucket, params.FileID, params.Version)

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
			WithXPath(fd.Path).
			WithXLabels(formatLabelsHeader(fd.Labels))
	})
}

func (h *handlers) FileListVersions() operations.FileListVersionsHandler {
	return operations.FileListVersionsHandlerFunc(func(params operations.FileListVersionsParams, principal *string) middleware.Responder {
		list, err := h.service.FileListVersions(params.HTTPRequest.Context(), params.Bucket, params.FileID)

		if err != nil {
			return operations.NewFileListVersionsInternalServerError().WithPayload(&models.Error{
				Code:    "server_error",
				Message: err.Error(),
			})
		}
		if len(list) == 0 {
			return operations.NewFileListVersionsNotFound()
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

		fd, err := h.service.FileNew(params.HTTPRequest.Context(), params.Bucket, params.File, params.ContentType, archetype, params.Labels)

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

		fd, err := h.service.FileUpdate(params.HTTPRequest.Context(), params.Bucket, params.FileID, params.File, params.ContentType, archetype, params.Labels)

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
		err := h.service.FileDelete(params.HTTPRequest.Context(), params.Bucket, params.FileID)

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

func (h *handlers) SyncBucketList() operations.SyncBucketListHandler {
	return operations.SyncBucketListHandlerFunc(func(params operations.SyncBucketListParams, principal *string) middleware.Responder {
		list, err := h.service.BucketList(params.HTTPRequest.Context())

		if err != nil {
			return operations.NewSyncBucketListInternalServerError().WithPayload(&models.Error{
				Code:    "server_error",
				Message: err.Error(),
			})
		}
		if len(list) == 0 {
			return operations.NewSyncBucketListNotFound()
		}

		return operations.NewSyncBucketListOK().WithPayload(list)
	})
}

func (h *handlers) SyncFileList() operations.SyncFileListHandler {
	return operations.SyncFileListHandlerFunc(func(params operations.SyncFileListParams, principal *string) middleware.Responder {
		list, err := h.service.SyncFileList(params.HTTPRequest.Context(), params.Bucket)

		if err != nil {
			return operations.NewSyncFileListInternalServerError().WithPayload(&models.Error{
				Code:    "server_error",
				Message: err.Error(),
			})
		}
		if len(list) == 0 {
			return operations.NewSyncFileListNotFound()
		}

		return operations.NewSyncFileListOK().WithPayload(list)
	})
}

func (h *handlers) SyncFileListVersions() operations.SyncFileListVersionsHandler {
	return operations.SyncFileListVersionsHandlerFunc(func(params operations.SyncFileListVersionsParams, principal *string) middleware.Responder {
		list, err := h.service.FileListVersions(params.HTTPRequest.Context(), params.Bucket, params.FileID)

		if err != nil {
			return operations.NewSyncFileListVersionsInternalServerError().WithPayload(&models.Error{
				Code:    "server_error",
				Message: err.Error(),
			})
		}
		if len(list) == 0 {
			return operations.NewSyncFileListVersionsNotFound()
		}

		return operations.NewSyncFileListVersionsOK().WithPayload(list)
	})
}

func (h *handlers) SyncFileMetadata() operations.SyncFileMetadataHandler {
	return operations.SyncFileMetadataHandlerFunc(func(params operations.SyncFileMetadataParams, principal *string) middleware.Responder {
		_, fd, err := h.service.FileGetVersion(params.HTTPRequest.Context(), params.Bucket, params.FileID, params.Version)

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
			WithXPath(fd.Path).
			WithXLabels(formatLabelsHeader(fd.Labels))
	})
}

func (h *handlers) SyncFile() operations.SyncFileHandler {
	return operations.SyncFileHandlerFunc(func(params operations.SyncFileParams, principal *string) middleware.Responder {
		defer params.File.Close()
		archetype := swag.StringValue(params.Archetype)

		fd, err := h.service.SyncFile(
			params.HTTPRequest.Context(),
			params.Bucket,
			params.FileID,
			params.Version,
			params.File,
			params.ContentType,
			params.Created,
			archetype,
			params.Labels,
		)

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
		err := h.service.SyncFileDelete(params.HTTPRequest.Context(), params.Bucket, params.FileID, params.Version, params.Created)
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

// NewHandlers returns a new instance of authenticator handlers
func NewHandlers(service Service, logger zerolog.Logger) Handlers {
	logger = logger.With().Str("component", "service/storage/handlers").Logger()

	return &handlers{service: service, logger: logger}
}

func formatLabelsHeader(l []string) string {
	return strings.Join(l, "|")
}
