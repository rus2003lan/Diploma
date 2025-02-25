package sqlmap

import (
	"context"

	"diploma-project/internal/model"
)

func (r *Repository) Fetch(ctx context.Context, q model.SQLMapCommand) (string, error) {
	report, ok := r.storage.Load(q.ID)
	if !ok {
		return "", model.ErrNotFound
	}

	res := report.(string)

	return res, nil
}
