package sqlmap

import (
	"context"
	"fmt"

	api "diploma-project/api/web/gen"
)

func (h *Handler) ReportSQLMap(
	ctx context.Context,
	request api.ReportSQLMapRequestObject,
) (api.ReportSQLMapResponseObject, error) {
	cmd := mapParamsToReportSQLMapCommand(request)
	resp, err := h.s.Fetch(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("create report: %w", err)
	}

	return api.ReportSQLMap200TextResponse(resp), nil
}
