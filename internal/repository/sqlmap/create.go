package sqlmap

import (
	"context"
	"crypto/sha256"
	"diploma-project/internal/model"
	"fmt"
	"net/http"
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

	hash := fmt.Sprintf("%x", sha256.Sum256(data))

	err := r.ceph.PutFile(ctx, r.bucket, hash, filetype, data)
	if err != nil {
		return fmt.Errorf("put file to ceph: %w", mapCephErrToModelErr(err))
	}

	return nil
}
