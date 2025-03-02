package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"diploma-project/internal/config"
	webapi "diploma-project/internal/controller/rest"
	"diploma-project/internal/provider"

	cli "github.com/urfave/cli/v2"
)

type Server interface {
	Start() error
}

func runServer(c *cli.Context) error {
	configArg := flag.String("config", "./config.yaml", "path to config")
	flag.Parse()

	cfg, err := config.LoadConfig(*configArg)
	if err != nil {
		panic("can't load application config: " + err.Error())
	}

	obsContainer := provider.NewObserveContainer(cfg)
	defer obsContainer.Logger(c.Context).Sync()

	clientContainer := provider.NewClientContainer(
		*cfg,
		*obsContainer.Logger(c.Context),
	)

	repoContainer := provider.NewRepositoryContainer(
		cfg,
		*clientContainer,
	)

	serviceContainer := provider.NewServiceContainer(
		cfg,
		repoContainer,
	)

	var api Server

	di := provider.NewWebAPIContainer(
		cfg,
		serviceContainer.Report(c.Context),
		serviceContainer.SQLMap(c.Context),
	)

	api = webapi.New(
		webapi.Opts{
			Port: int(cfg.Server.Port),
			Handlers: webapi.Handler{
				ReportHandler: di.ReportHandler(c.Context),
				DocHandler:    di.DocHandler(c.Context),
				SQLMapHandler: di.SQLMapHandler(c.Context),
			},
			ValidateErrorFunc:    di.ValidateErrorFunc(c.Context),
			ErrorRespFunc:        di.ErrorFunc(c.Context),
			ValidateErrorHandler: di.ValidateErrorHandler(c.Context),
		},
		obsContainer.Logger(c.Context),
	)

	errCh := make(chan error)

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := api.Start(); err != nil {
			errCh <- fmt.Errorf("run webserver: %w", err)
		}
	}()

	select {
	case err := <-errCh:
		obsContainer.Logger(c.Context).Error("can't run application: " + err.Error())

		return err

	case sig := <-osSignals:
		obsContainer.Logger(c.Context).Info("graceful shutdown initiated", sig)

		return nil
	}
}

func main() {
	app := &cli.App{ //nolint:exhaustruct
		Name:  "reports",
		Usage: "reports",
		Flags: []cli.Flag{
			&cli.StringFlag{ //nolint:exhaustruct
				Name:  "config",
				Value: "./config.yml",
				Usage: "path to the config file",
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "web-api",
				Usage:  "run web api",
				Action: runServer,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal("can't run application: " + err.Error())
	}
}
