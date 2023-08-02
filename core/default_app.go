package core

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"golang-monorepo-boilerplate/core/config"
	"golang-monorepo-boilerplate/core/log"
	"golang-monorepo-boilerplate/core/persistence"
	"golang-monorepo-boilerplate/internal/health"
	"golang-monorepo-boilerplate/ui"
)

type DefaultApp struct {
	l *log.Logger
	c *config.Config
	r *echo.Echo
	p *persistence.Persister

	healthHandler *health.Handler
}

func (a *DefaultApp) withPersistence(p *persistence.Persister) App {
	a.p = p
	return a
}

func (a *DefaultApp) Persister() *persistence.Persister {
	return a.p
}

func (a *DefaultApp) HealthHandler() *health.Handler {
	if a.healthHandler == nil {
		a.healthHandler = health.NewHandler(a)
	}
	return a.healthHandler
}

func (a *DefaultApp) registerUiApp(ctx context.Context) {
	a.Router().GET(
		"/*",
		echo.StaticDirectoryHandler(ui.DistDirFS, false),
		middleware.Gzip(),
	)
}

func (a *DefaultApp) Router() *echo.Echo {
	return a.r
}

func (a *DefaultApp) Serve(ctx context.Context) error {
	err := a.Router().Start(":" + a.Config().Port())
	if err != nil {
		panic(errors.WithStack(err))
	}
	return nil
}

func (a *DefaultApp) withRouter(e *echo.Echo) App {
	a.r = e
	return a
}

func (a *DefaultApp) registerApiRoutes(ctx context.Context) {
	apiGroup := a.Router().Group("/api")
	a.HealthHandler().RegisterRoutes(apiGroup)
}

func (a *DefaultApp) Init(ctx context.Context) error {
	a.registerApiRoutes(ctx)
	a.registerUiApp(ctx)

	runner, err := a.Persister().MigrateRunner()
	if err != nil {
		return err
	}

	err = runner.Run(false)
	if err != nil {
		return err
	}

	return nil
}

func (a *DefaultApp) withLogger(l *log.Logger) App {
	a.l = l
	return a
}

func (a *DefaultApp) withConfig(c *config.Config) App {
	a.c = c
	return a
}

func (a *DefaultApp) Logger() *log.Logger {
	return a.l
}

func (a *DefaultApp) Config() *config.Config {
	return a.c
}
