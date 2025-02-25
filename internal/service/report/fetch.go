package report

import (
	"context"
	"fmt"

	"diploma-project/internal/model"
)

func (s *Service) Fetch(ctx context.Context, q model.ReportFetchQuery) (*model.Report, error) {
	report, err := s.r.Fetch(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("fetch report: %w", err)
	}

	return report, nil
}
