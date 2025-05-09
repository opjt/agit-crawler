package bootstrap

import (
	"agit-crawler/app/cron"
	"agit-crawler/app/lib"
	"agit-crawler/app/pkg"
	"context"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func RunServer(opt fx.Option) {
	logger := lib.GetLogger()

	app := fx.New(
		opt,
		fx.WithLogger(func() fxevent.Logger {
			return logger.GetFxLogger()
		}),
		fx.Invoke(RegisterLifecycle),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := app.Start(ctx); err != nil {
		logger.Fatal("Failed to start app:", err)
	}

	<-app.Done() // SIGINT, SIGTERM 등 기다림

	if err := app.Stop(ctx); err != nil {
		logger.Fatal("Failed to stop app:", err)
	}
}

func RegisterLifecycle(
	lifecycle fx.Lifecycle,
	env lib.Env,
	logger lib.Logger,
	crawler *pkg.Crawler,
	cronManager cron.CronManager,
) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Crawler started with cron job")
			cronManager.Start()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down crawler")
			crawler.Close()
			cronManager.Stop()
			return nil
		},
	})
}
