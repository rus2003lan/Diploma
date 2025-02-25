package report

import (
	"context"
	"fmt"

	"diploma-project/internal/model"
)

func (s *Service) Search(ctx context.Context, q model.ReportSearchQuery) ([]model.Report, error) {
	reports, err := s.r.Search(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("get list of reports: %w", err)
	}

	return reports, nil
}
