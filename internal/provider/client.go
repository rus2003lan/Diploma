package provider

import (
	"context"

	"diploma-project/internal/config"

	"diploma-project/pkg/ceph"
	"diploma-project/pkg/elastic"

	"go.uber.org/zap"
)

type ClientContainer struct {
	cfg config.Config

	elastic *elastic.ElClient
	ceph    *ceph.Client

	lg zap.SugaredLogger
}

func NewClientContainer(
	cfg config.Config,
	lg zap.SugaredLogger,
) *ClientContainer {
	return &ClientContainer{
		cfg:     cfg,
		elastic: nil,
		ceph:    nil,
		lg:      lg,
	}
}

func (c *ClientContainer) Elastic(ctx context.Context) *elastic.ElClient {
	if c.elastic != nil {
		return c.elastic
	}

	var err error

	c.elastic, err = elastic.NewElastic(*c.cfg.Elastic)
	if err != nil {
		panic(err)
	}

	err = c.elastic.CreateIndex(ctx, c.cfg.StartupIndexConfig)
	if err != nil {
		panic(err)
	}

	return c.elastic
}

func (c *ClientContainer) Ceph(ctx context.Context) *ceph.Client {
	if c.ceph != nil {
		return c.ceph
	}

	var err error

	c.ceph, err = ceph.New(c.cfg.Ceph)
	if err != nil {
		panic(err)
	}

	err = c.ceph.CreateBucket(ctx, c.cfg.Ceph.Bucket)
	if err != nil {
		panic(err)
	}

	return c.ceph
}
