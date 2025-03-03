package sqlmap

import (
	"context"
	"fmt"
	"net/http"

	"diploma-project/internal/model"
)

const bufLen = 512

func (r *Repository) Create(ctx context.Context, cmd model.SQLMapCommand) error {
	data := cmd.Report

	// detect content type
	bufBound := len(data)
	if bufBound > bufLen {
		bufBound = bufLen
	}

	filetype := http.DetectContentType(data[:bufBound])

	err := r.ceph.PutFile(ctx, r.bucket, cmd.ID, filetype, data)
	if err != nil {
		return fmt.Errorf("put file to ceph: %w", mapCephErrToModelErr(err))
	}

	return nil
}
