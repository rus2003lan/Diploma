package provider

import (
	"context"
	"errors"
	"net/http"

	api "diploma-project/api/web/gen"

	json "github.com/json-iterator/go"

	"diploma-project/internal/config"
	webapi "diploma-project/internal/controller/rest"
	"diploma-project/internal/model"

	"diploma-project/internal/controller/rest/doc"
	reportcontroller "diploma-project/internal/controller/rest/report"
	"diploma-project/internal/controller/rest/sqlmap"

	validator "github.com/oapi-codegen/nethttp-middleware"
)

type (
	captureErrorFunc func(context.Context, error)
	errorRespFunc    func(w http.ResponseWriter, r *http.Request, err error)
)

type ReportService interface {
	Search(ctx context.Context, q model.ReportSearchQuery) ([]model.Report, error)
	Fetch(ctx context.Context, q model.ReportFetchQuery) (*model.Report, error)
	Create(ctx context.Context, q model.ReportCreateCommand) error
}

type SQLService interface {
	Fetch(ctx context.Context, cmd model.SQLMapCommand) (string, error)
}

type WebAPIContainer struct {
	cfg *config.Config

	reportService ReportService
	sqlService    SQLService

	docHandler    webapi.DocHandler
	reportHandler webapi.ReportHandler
	sqlHandler    webapi.SQLMapHandler

	captureErrorFunc     captureErrorFunc
	errorRespFunc        errorRespFunc
	validateErrorFunc    errorRespFunc
	validateErrorHandler validator.ErrorHandler
}

//nolint:revive
func NewWebAPIContainer(
	cfg *config.Config,
	rs ReportService,
	ss SQLService,
) *WebAPIContainer {
	return &WebAPIContainer{
		cfg:                  cfg,
		reportService:        rs,
		sqlService:           ss,
		docHandler:           nil,
		reportHandler:        nil,
		sqlHandler:           nil,
		captureErrorFunc:     nil,
		errorRespFunc:        nil,
		validateErrorFunc:    nil,
		validateErrorHandler: nil,
	}
}

func (c *WebAPIContainer) ReportHandler(_ context.Context) webapi.ReportHandler {
	if c.reportHandler != nil {
		return c.reportHandler
	}

	c.reportHandler = reportcontroller.NewHandler(c.reportService)

	return c.reportHandler
}

func (c *WebAPIContainer) SQLMapHandler(_ context.Context) webapi.SQLMapHandler {
	if c.sqlHandler != nil {
		return c.sqlHandler
	}

	c.sqlHandler = sqlmap.NewHandler(c.sqlService)

	return c.sqlHandler
}

func (c *WebAPIContainer) DocHandler(_ context.Context) webapi.DocHandler {
	if c.docHandler != nil {
		return c.docHandler
	}

	c.docHandler = doc.NewHandler()

	return c.docHandler
}

func (c *WebAPIContainer) ErrorFunc(ctx context.Context) func(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	if c.errorRespFunc != nil {
		return c.errorRespFunc
	}

	c.errorRespFunc = func(w http.ResponseWriter, r *http.Request, err error) {
		var code int

		switch {
		case errors.Is(err, model.ErrNotFound):
			code = http.StatusNotFound

		case errors.Is(err, model.ErrNotValid):
			code = http.StatusBadRequest

		default:
			code = http.StatusInternalServerError
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(code)

		_ = json.NewEncoder(w).Encode(api.Error{
			Code:    code,
			Message: err.Error(),
		})
	}

	return c.errorRespFunc
}

func (c *WebAPIContainer) ValidateErrorFunc(_ context.Context) func(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	if c.validateErrorFunc != nil {
		return c.validateErrorFunc
	}

	c.validateErrorFunc = func(w http.ResponseWriter, _ *http.Request, err error) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		w.WriteHeader(http.StatusBadRequest)

		_ = json.NewEncoder(w).Encode(api.Error{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return c.validateErrorFunc
}

func (c *WebAPIContainer) ValidateErrorHandler(_ context.Context) validator.ErrorHandler {
	if c.validateErrorHandler != nil {
		return c.validateErrorHandler
	}

	c.validateErrorHandler = func(w http.ResponseWriter, message string, statusCode int) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		w.WriteHeader(statusCode)

		_ = json.NewEncoder(w).Encode(api.Error{
			Code:    statusCode,
			Message: message,
		})
	}

	return c.validateErrorHandler
}
