package sqlmap

import (
	"context"

	"diploma-project/internal/model"
)

type (
	Service interface {
		Fetch(ctx context.Context, cmd model.SQLMapCommand) (string, error)
	}
)
