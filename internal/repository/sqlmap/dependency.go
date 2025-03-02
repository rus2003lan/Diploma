package sqlmap

import (
	"context"
	"diploma-project/pkg/ceph"
)

type Ceph interface {
	PutFile(ctx context.Context, bucket, filename, contentType string, data []byte) error
	FetchFile(ctx context.Context, bucket, filename string) (*ceph.File, error)
}
