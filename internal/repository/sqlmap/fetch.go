package sqlmap

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"diploma-project/internal/model"
)

func (r *Repository) Fetch(ctx context.Context, q model.SQLMapCommand) (string, error) {
	res, err := r.ceph.FetchFile(ctx, r.bucket, q.ID)
	if err != nil {
		return "", fmt.Errorf("fetch file: %w", mapCephErrToModelErr(err))
	}

	var buf bytes.Buffer

	_, err = io.Copy(&buf, res.Data)
	if err != nil {
		return "", fmt.Errorf("copy data: %w", err)
	}

	return buf.String(), nil
}
