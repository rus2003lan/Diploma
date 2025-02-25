package sqlmap

import (
	"context"

	"diploma-project/internal/model"
)

type (
	Repo interface {
		Fetch(ctx context.Context, cmd model.SQLMapCommand) (string, error)
		Create(ctx context.Context, cmd model.SQLMapCommand) error
	}
)
