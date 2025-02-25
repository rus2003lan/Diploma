package report

import (
	"context"

	"diploma-project/internal/model"
)

type (
	ReportRepo interface {
		Search(ctx context.Context, q model.ReportSearchQuery) ([]model.Report, error)
		Fetch(ctx context.Context, q model.ReportFetchQuery) (*model.Report, error)
		Create(ctx context.Context, cmd model.ReportCreateCommand) error
	}

	SQLMapService interface {
		Create(ctx context.Context, cmd model.SQLMapCommand) error
	}
)
