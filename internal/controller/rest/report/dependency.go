package report

import (
	"context"

	"diploma-project/internal/model"
)

type (
	Service interface {
		Search(ctx context.Context, q model.ReportSearchQuery) ([]model.Report, error)
		Fetch(ctx context.Context, q model.ReportFetchQuery) (*model.Report, error)
		Create(ctx context.Context, q model.ReportCreateCommand) error
	}
)
