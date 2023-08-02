package persistence

import (
	"golang-monorepo-boilerplate/core/config"
	"golang-monorepo-boilerplate/core/log"
	"golang-monorepo-boilerplate/core/persistence/migrate"
	"golang-monorepo-boilerplate/migrations"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type (
	persisterDependencies interface {
		config.Provider
		log.Provider
	}
	Persister struct {
		c *gorm.DB
		d persisterDependencies
	}
	Provider interface {
		Persister() *Persister
	}
)

func New(d persisterDependencies) (*Persister, error) {
	db, err := gorm.Open(sqlite.Open(d.Config().DB()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &Persister{
		c: db,
		d: d,
	}, nil
}

func (p *Persister) Persister() *gorm.DB {
	return p.c
}

func (p *Persister) MigrateRunner() (*migrate.Runner, error) {
	runner, err := migrate.NewRunner(p.c, migrations.Migrations, p.d)
	if err != nil {
		return nil, err
	}

	return runner, nil
}
