package provider

import (
	"context"

	"diploma-project/internal/config"
	"diploma-project/internal/service/report"
	"diploma-project/internal/service/sqlmap"
)

type ServiceContainer struct {
	cfg *config.Config

	rc RepositoryContainer

	rs *report.Service
	ss *sqlmap.Service
}

func NewServiceContainer(
	cfg *config.Config,
	rp *RepositoryContainer,
) *ServiceContainer {
	return &ServiceContainer{
		cfg: cfg,
		rc:  *rp,
		ss:  nil,
		rs:  nil,
	}
}

func (c *ServiceContainer) Report(ctx context.Context) *report.Service {
	if c.rs != nil {
		return c.rs
	}

	c.rs = report.New(c.rc.ReportRepo(ctx), c.SQLMap(ctx))

	return c.rs
}

func (c *ServiceContainer) SQLMap(ctx context.Context) *sqlmap.Service {
	if c.ss != nil {
		return c.ss
	}

	c.ss = sqlmap.New(c.rc.SQLMapRepo(ctx))

	return c.ss
}
