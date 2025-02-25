package report

import (
	"context"
	"fmt"

	api "diploma-project/api/web/gen"
)

func (h *Handler) ReportList(
	ctx context.Context,
	request api.ReportListRequestObject,
) (api.ReportListResponseObject, error) {
	res, err := h.s.Search(ctx, mapSearchParamsToReportSearchQuery(request))
	if err != nil {
		return nil, fmt.Errorf("search reports: %w", err)
	}

	return api.ReportList200ApplicationJSONCharsetUTF8Response(mapReportsToWeb(res)), nil
}
