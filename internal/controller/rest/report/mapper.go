package report

import (
	api "diploma-project/api/web/gen"
	"diploma-project/internal/model"
)

const (
	defaultLimit = 10
)

func mapParamsToReportCreateCommand(r api.ReportSaveRequestObject) model.ReportCreateCommand {
	return model.ReportCreateCommand{
		URL: r.Params.Url,
	}
}

func mapSearchParamsToReportSearchQuery(api.ReportListRequestObject) model.ReportSearchQuery {
	return model.ReportSearchQuery{
		Limit:  defaultLimit,
		Offset: 0,
	}
}

func mapParamsToReportFetchQuery(r api.ReportFetchRequestObject) model.ReportFetchQuery {
	return model.ReportFetchQuery{
		ID: r.Id,
	}
}

func mapReportsToWeb(r []model.Report) []api.Report {
	reports := make([]api.Report, 0, len(r))
	for _, report := range r {
		reports = append(reports, mapReportToWeb(report))
	}

	return reports
}

func mapReportToWeb(r model.Report) api.Report {
	return api.Report{
		Id:   r.Id,
		Urls: mapURLsToWeb(r.URLs),
	}
}

func mapURLsToWeb(p []model.URL) []api.URL {
	params := make([]api.URL, 0, len(p))
	for _, param := range p {
		params = append(params, api.URL{
			Method: param.Method,
			Url:    param.URL,
			Params: mapParamsToWeb(param.Params),
		})
	}

	return params
}

func mapParamsToWeb(v []model.Param) []api.Param {
	vulns := make([]api.Param, 0, len(v))
	for _, vuln := range v {
		vulns = append(vulns, api.Param{
			Name:     vuln.Name,
			Values:   vuln.Values,
			Patterns: vuln.Patterns,
		})
	}

	return vulns
}
