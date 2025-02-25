package sqlmap

import (
	api "diploma-project/api/web/gen"
	"diploma-project/internal/model"
)

func mapParamsToReportSQLMapCommand(
	request api.ReportSQLMapRequestObject,
) model.SQLMapCommand {
	return model.SQLMapCommand{
		ID: request.Id,
	}
}
