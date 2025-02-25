package sqlmap

import (
	"context"
	"diploma-project/internal/model"
)

func (r *Repository) Create(ctx context.Context, cmd model.SQLMapCommand) error {
	r.storage.Store(cmd.ID, cmd.Report)

	return nil
}
