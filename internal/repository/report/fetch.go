package report

import (
	"context"
	"diploma-project/internal/model"
	"fmt"
)

func (r *Repository) Fetch(ctx context.Context, q model.ReportFetchQuery) (*model.Report, error) {
	hit, err := r.es.GetDoc(ctx, q.ID)
	if err != nil {
		return nil, fmt.Errorf("get doc from elastic: %w", mapElasticErrToModelErr(err))
	}

	report, err := mapElasticHitToModel(hit)
	if err != nil {
		return nil, fmt.Errorf("map report to model: %w", err)
	}

	return &report, nil
}
