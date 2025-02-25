package report

import (
	"context"
	"diploma-project/internal/model"
)

func (r *Repository) Search(ctx context.Context, q model.ReportSearchQuery) ([]model.Report, error) {
	var reports []model.Report
	r.storage.Range(func(key, value any) bool {
		reports = append(reports, value.(model.Report))
		return true
	})

	return reports, nil
}
