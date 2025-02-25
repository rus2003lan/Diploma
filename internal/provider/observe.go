package provider

import (
	"context"

	"diploma-project/internal/config"
	"diploma-project/pkg/logger"

	"go.uber.org/zap"
)

type ObserveContainer struct {
	logger *zap.SugaredLogger
}

func NewObserveContainer(
	cfg *config.Config,
) *ObserveContainer {
	return &ObserveContainer{
		logger: nil,
	}
}

func (c *ObserveContainer) Logger(_ context.Context) *zap.SugaredLogger {
	if c.logger != nil {
		return c.logger
	}

	c.logger = logger.NewZapLogger()

	return c.logger
}
