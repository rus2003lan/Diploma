package webapi

import (
	"context"
	"net/http"

	api "diploma-project/api/web/gen"
)

type (
	Logger interface {
		Info(msg string, args ...any)
		Warn(msg string, args ...any)
		Error(msg string, args ...any)
	}

	ReportHandler interface {
		ReportList(
			ctx context.Context,
			request api.ReportListRequestObject,
		) (api.ReportListResponseObject, error)
		ReportFetch(
			ctx context.Context,
			request api.ReportFetchRequestObject,
		) (api.ReportFetchResponseObject, error)
		ReportSave(
			ctx context.Context,
			request api.ReportSaveRequestObject,
		) (api.ReportSaveResponseObject, error)
	}

	SQLMapHandler interface {
		ReportSQLMap(
			ctx context.Context,
			request api.ReportSQLMapRequestObject,
		) (api.ReportSQLMapResponseObject, error)
	}

	DocHandler interface {
		Doc(
			ctx context.Context,
			request api.DocRequestObject,
		) (api.DocResponseObject, error)
	}

	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
)
