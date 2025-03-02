package report

import (
	"context"
	"diploma-project/pkg/elastic"
)

type Elastic interface {
	GetDoc(ctx context.Context, id string) (*elastic.Hit, error)
	Search(ctx context.Context) (*elastic.FetchResponse, error)
	Update(ctx context.Context, id string, body []byte) error
}
