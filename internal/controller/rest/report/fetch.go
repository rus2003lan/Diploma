package report

import (
	"context"
	"fmt"

	api "diploma-project/api/web/gen"
)

func (h *Handler) ReportFetch(
	ctx context.Context,
	request api.ReportFetchRequestObject,
) (api.ReportFetchResponseObject, error) {
	report, err := h.s.Fetch(ctx, mapParamsToReportFetchQuery(request))
	if err != nil {
		return nil, fmt.Errorf("fetch report: %w", err)
	}

	return api.ReportFetch200ApplicationJSONCharsetUTF8Response(mapReportToWeb(*report)), nil
}
