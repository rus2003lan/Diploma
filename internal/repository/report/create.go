package report

import (
	"context"
	"diploma-project/internal/model"
)

func (r *Repository) Create(ctx context.Context, cmd model.ReportCreateCommand) error {
	r.storage.Store(cmd.Report.Id, cmd.Report)

	return nil
}
