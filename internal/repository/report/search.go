package report

import (
	"context"
	"diploma-project/internal/model"
	"fmt"
)

func (r *Repository) Search(ctx context.Context, q model.ReportSearchQuery) ([]model.Report, error) {
	resp, err := r.es.Search(ctx)
	if err != nil {
		return nil, fmt.Errorf("reports search: %w", mapElasticErrToModelErr(err))
	}

	res, err := mapElasticFetchRespToModel(resp)
	if err != nil {
		return nil, fmt.Errorf("map repo reports to model: %w", err)
	}

	return res, nil
}
