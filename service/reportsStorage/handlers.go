package reportsStorage

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/reportsStorage/models"
	"github.com/iryonetwork/wwm/gen/reportsStorage/restapi/operations"
	"github.com/iryonetwork/wwm/utils"
)

// Handlers describes the actions supported by the storage handlers
type Handlers interface {
	ReportList() operations.ReportListHandler
	ReportGet() operations.ReportGetHandler
	ReportGetVersion() operations.ReportGetVersionHandler
	ReportListVersions() operations.ReportListVersionsHandler
	ReportNew() operations.ReportNewHandler
	ReportUpdate() operations.ReportUpdateHandler
	ReportDelete() operations.ReportDeleteHandler
}

type handlers struct {
	service Service
	logger  zerolog.Logger
}

func (h *handlers) ReportList() operations.ReportListHandler {
	return operations.ReportListHandlerFunc(func(params operations.ReportListParams, principal *string) middleware.Responder {
		list, err := h.service.ReportList(params.HTTPRequest.Context(), params.ReportType)

		if err != nil {
			return operations.NewReportListInternalServerError().WithPayload(&models.Error{
				Code:    "server_error",
				Message: err.Error(),
			})
		}
		if len(list) == 0 {
			return operations.NewReportListNotFound()
		}

		return operations.NewReportListOK().WithPayload(list)
	})
}

func (h *handlers) ReportGet() operations.ReportGetHandler {
	return operations.ReportGetHandlerFunc(func(params operations.ReportGetParams, principal *string) middleware.Responder {
		r, fd, err := h.service.ReportGet(params.HTTPRequest.Context(), params.ReportType, params.FileName)

		if err != nil {
			switch err {
			case ErrNotFound:
				return operations.NewReportGetNotFound()
			default:
				return operations.NewReportGetInternalServerError().WithPayload(&models.Error{
					Code:    "server_error",
					Message: err.Error(),
				})
			}
		}

		return utils.UseProducer(operations.NewReportGetOK().
			WithPayload(r).
			WithXContentType(fd.ContentType).
			WithXCreated(fd.Created).
			WithXVersion(fd.Version).
			WithXChecksum(fd.Checksum).
			WithXName(fd.Name).
			WithXPath(fd.Path).
			WithXDataSince(fd.DataSince).
			WithXDataUntil(fd.DataUntil).
			WithXReportType(params.ReportType), utils.FileProducer)
	})
}

func (h *handlers) ReportGetVersion() operations.ReportGetVersionHandler {
	return operations.ReportGetVersionHandlerFunc(func(params operations.ReportGetVersionParams, principal *string) middleware.Responder {
		r, fd, err := h.service.ReportGetVersion(params.HTTPRequest.Context(), params.ReportType, params.FileName, params.Version)

		if err != nil {
			switch err {
			case ErrNotFound:
				return operations.NewReportGetVersionNotFound()
			default:
				return operations.NewReportGetVersionInternalServerError().WithPayload(&models.Error{
					Code:    "server_error",
					Message: err.Error(),
				})
			}
		}

		return utils.UseProducer(operations.NewReportGetVersionOK().
			WithPayload(r).
			WithXContentType(fd.ContentType).
			WithXCreated(fd.Created).
			WithXVersion(fd.Version).
			WithXChecksum(fd.Checksum).
			WithXName(fd.Name).
			WithXPath(fd.Path).
			WithXDataSince(fd.DataSince).
			WithXDataUntil(fd.DataUntil).
			WithXReportType(params.ReportType), utils.FileProducer)
	})
}

func (h *handlers) ReportListVersions() operations.ReportListVersionsHandler {
	return operations.ReportListVersionsHandlerFunc(func(params operations.ReportListVersionsParams, principal *string) middleware.Responder {
		list, err := h.service.ReportListVersions(params.HTTPRequest.Context(), params.ReportType, params.FileName, nil, nil)

		if err != nil {
			return operations.NewReportListVersionsInternalServerError().WithPayload(&models.Error{
				Code:    "server_error",
				Message: err.Error(),
			})
		}
		if len(list) == 0 {
			return operations.NewReportListVersionsNotFound()
		}

		return operations.NewReportListVersionsOK().WithPayload(list)
	})
}

func (h *handlers) ReportNew() operations.ReportNewHandler {
	return operations.ReportNewHandlerFunc(func(params operations.ReportNewParams, principal *string) middleware.Responder {
		defer params.File.Close()

		fd, err := h.service.ReportNew(params.HTTPRequest.Context(), params.ReportType, params.File, params.ContentType, params.DataSince, params.DataUntil)

		if err != nil {
			return operations.NewReportNewInternalServerError().WithPayload(&models.Error{
				Code:    "server_error",
				Message: err.Error(),
			})
		}

		return operations.NewReportNewCreated().WithPayload(fd)
	})
}

func (h *handlers) ReportUpdate() operations.ReportUpdateHandler {
	return operations.ReportUpdateHandlerFunc(func(params operations.ReportUpdateParams, principal *string) middleware.Responder {
		defer params.File.Close()

		fd, err := h.service.ReportUpdate(params.HTTPRequest.Context(), params.ReportType, params.FileName, params.File, params.ContentType)

		if err != nil {
			switch err {
			case ErrNotFound:
				return operations.NewReportUpdateNotFound()
			default:
				return operations.NewReportUpdateInternalServerError().WithPayload(&models.Error{
					Code:    "server_error",
					Message: err.Error(),
				})
			}
		}

		return operations.NewReportUpdateCreated().WithPayload(fd)
	})
}

func (h *handlers) ReportDelete() operations.ReportDeleteHandler {
	return operations.ReportDeleteHandlerFunc(func(params operations.ReportDeleteParams, principal *string) middleware.Responder {
		err := h.service.ReportDelete(params.HTTPRequest.Context(), params.ReportType, params.FileName)

		if err != nil {
			switch err {
			case ErrNotFound:
				return operations.NewReportDeleteNotFound()
			default:
				return operations.NewReportDeleteInternalServerError().WithPayload(&models.Error{
					Code:    "server_error",
					Message: err.Error(),
				})
			}
		}

		return operations.NewReportDeleteNoContent()
	})
}

// NewHandlers returns a new instance of authenticator handlers
func NewHandlers(service Service, logger zerolog.Logger) Handlers {
	logger = logger.With().Str("component", "service/storage/handlers").Logger()

	return &handlers{service: service, logger: logger}
}
