package report

import (
	"context"
	"diploma-project/internal/model"
)

func (r *Repository) Fetch(ctx context.Context, q model.ReportFetchQuery) (*model.Report, error) {
	report, ok := r.storage.Load(q.ID)
	if !ok {
		return nil, model.ErrNotFound
	}

	res := report.(model.Report)

	return &res, nil
}
