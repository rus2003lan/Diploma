package provider

import (
	"context"

	"diploma-project/internal/config"
	"diploma-project/internal/repository/report"
	"diploma-project/internal/repository/sqlmap"
)

type RepositoryContainer struct {
	cfg *config.Config

	reportRepo *report.Repository
	sqlRepo    *sqlmap.Repository

	cc ClientContainer
}

func NewRepositoryContainer(
	cfg *config.Config,
	cc ClientContainer,
) *RepositoryContainer {
	return &RepositoryContainer{
		cfg:        cfg,
		sqlRepo:    nil,
		cc:         cc,
		reportRepo: nil,
	}
}

func (c *RepositoryContainer) ReportRepo(ctx context.Context) *report.Repository {
	if c.reportRepo != nil {
		return c.reportRepo
	}

	c.reportRepo = report.New(c.cc.Elastic(ctx))

	return c.reportRepo
}

func (c *RepositoryContainer) SQLMapRepo(ctx context.Context) *sqlmap.Repository {
	if c.sqlRepo != nil {
		return c.sqlRepo
	}

	c.sqlRepo = sqlmap.New(c.cc.Ceph(ctx), c.cfg.Ceph.Bucket)

	return c.sqlRepo
}
