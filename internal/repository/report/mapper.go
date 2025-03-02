package report

import (
	"diploma-project/internal/model"
	"diploma-project/pkg/elastic"
	"encoding/json"
	"errors"
	"fmt"
)

func mapModelReportToRepo(r model.Report) report {
    return report{
    	Id:   r.Id,
    	URLs: mapModelURLsToRepo(r.URLs),
    }
}

func mapModelURLsToRepo(urls []model.URL) []url {
    res := make([]url, 0, len(urls))

	for _, u := range urls {
		res = append(res, url{
			URL:    u.URL,
			Method: u.Method,
			Params: mapModelParamsToRepo(u.Params),
		})
	}

	return res
}

func mapModelParamsToRepo(params []model.Param) []param {
    res := make([]param, len(params))

	for _, u := range params {
		res = append(res, param(u))
	}

	return res
}

func mapElasticErrToModelErr(err error) error {
	if errors.Is(err, elastic.ErrNotFound) {
		return errors.Join(model.ErrNotFound, err)
	}

	if errors.Is(err, elastic.ErrInvalid) {
		return errors.Join(model.ErrNotValid, err)
	}

	return err
}

func mapElasticFetchRespToModel(resp *elastic.FetchResponse) ([]model.Report, error) {
	reports := make([]model.Report, 0, len(resp.Hits))

	for _, v := range resp.Hits {
		r, err := mapElasticHitToModel(&v)
		if err != nil {
			return nil, fmt.Errorf("map report to model: %w", err)
		}

		reports = append(reports, r)
	}

	return reports, nil
}

func mapElasticHitToModel(hit *elastic.Hit) (model.Report, error) {
	d := new(report)
	if err := json.Unmarshal(hit.Source, d); err != nil {
		return model.Report{}, fmt.Errorf("unmarshal source: %w", err)
	}

	return mapRepoReportToModel(*d), nil
}

func mapRepoReportToModel(in report) model.Report {
	return model.Report{
		Id:   in.Id,
		URLs: mapRepoURLsToModel(in.URLs),
	}
}

func mapRepoURLsToModel(urls []url) []model.URL {
	res := make([]model.URL, 0, len(urls))

	for _, u := range urls {
		res = append(res, model.URL{
			URL:    u.URL,
			Method: u.Method,
			Params: mapRepoParamsToModel(u.Params),
		})
	}

	return res
}

func mapRepoParamsToModel(params []param) []model.Param {
	res := make([]model.Param, 0, len(params))

	for _, u := range params {
		res = append(res, model.Param(u))
	}

	return res
}
