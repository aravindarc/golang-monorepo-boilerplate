package migrations

import (
	"golang-monorepo-boilerplate/core/persistence/migrate"
	"golang-monorepo-boilerplate/models"
	"gorm.io/gorm"
)

var Migrations migrate.MigrationsList

func init() {
	Migrations.Register(func(db *gorm.DB) error {
		err := db.Migrator().CreateTable(&models.Test{})
		if err != nil {
			return err
		}
		return nil
	}, func(db *gorm.DB) error {
		err := db.Migrator().DropTable(&models.Test{})
		if err != nil {
			return err
		}
		return nil
	})
}
