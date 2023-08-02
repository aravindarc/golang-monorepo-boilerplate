package core

import (
	"context"
	"github.com/cenkalti/backoff"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang-monorepo-boilerplate/core/config"
	"golang-monorepo-boilerplate/core/log"
	"golang-monorepo-boilerplate/core/persistence"
	"golang-monorepo-boilerplate/internal/health"
	"strings"
	"time"
)

type App interface {
	Init(ctx context.Context) error
	withLogger(l *log.Logger) App
	withConfig(c *config.Config) App
	withRouter(e *echo.Echo) App
	withPersistence(p *persistence.Persister) App

	registerApiRoutes(ctx context.Context)
	registerUiApp(ctx context.Context)

	Serve(ctx context.Context) error

	Router() *echo.Echo

	config.Provider
	log.Provider
	persistence.Provider

	health.HandlerProvider
}

func New(ctx context.Context, cmd *cobra.Command) (App, error) {
	var app App = &DefaultApp{}
	app.withConfig(config.NewConfig(app, cmd))
	app.withLogger(log.NewLogger("gmb", "main", strings.ToLower(app.Config().LogLevel())))
	app.withRouter(newRouter())
	bc := backoff.NewExponentialBackOff()
	bc.MaxElapsedTime = time.Minute * 5
	bc.Reset()
	err := errors.WithStack(
		backoff.Retry(func() error {
			persister, err := persistence.New(app)
			if err != nil {
				return errors.WithStack(err)
			}
			app.withPersistence(persister)
			return nil
		}, bc),
	)
	if err != nil {
		return nil, err
	}

	return app, nil
}
