package report

import (
	"context"
	"diploma-project/internal/model"
	"encoding/json"
	"fmt"
)

func (r *Repository) Create(ctx context.Context, cmd model.ReportCreateCommand) error {
	data, err := json.Marshal(mapModelReportToRepo(cmd.Report))
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	body := fmt.Sprintf(`{"doc":%s, "doc_as_upsert": true}`, data)

	err = r.es.Update(ctx, cmd.Report.Id, []byte(body))
	if err != nil {
		return fmt.Errorf("upsert report: %w", mapElasticErrToModelErr(err))
	}

	return nil
}
