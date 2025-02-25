package report

import (
	"context"
	api "diploma-project/api/web/gen"
	"fmt"
)

func (h *Handler) ReportSave(
	ctx context.Context,
	request api.ReportSaveRequestObject,
) (api.ReportSaveResponseObject, error) {
	err := h.s.Create(ctx, mapParamsToReportCreateCommand(request))
	if err != nil {
		return nil, fmt.Errorf("create report: %w", err)
	}

	return api.ReportSave200ApplicationJSONCharsetUTF8Response("Success"), nil
}
