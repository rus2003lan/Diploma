package sqlmap

import (
	"context"
	"fmt"

	"diploma-project/internal/model"
)

func (s *Service) Fetch(ctx context.Context, cmd model.SQLMapCommand) (string, error) {
	report, err := s.r.Fetch(ctx, cmd)
	if err != nil {
		return "", fmt.Errorf("fetch sqlmap report: %w", err)
	}

	return report, nil
}
