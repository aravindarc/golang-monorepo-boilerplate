package migrate

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cast"
	"golang-monorepo-boilerplate/core/config"
	"golang-monorepo-boilerplate/core/log"
	"gorm.io/gorm"
	"strings"
)

const DefaultMigrationsTable = "_migrations"

type (
	runnerDependencies interface {
		config.Provider
		log.Provider
	}
	Runner struct {
		db             *gorm.DB
		migrationsList MigrationsList
		tableName      string
		d              runnerDependencies
	}
	Migrations struct {
		File    string `gorm:"primaryKey;not null"`
		Applied bool   `gorm:"not null"`
	}
)

func (Migrations) TableName() string {
	return DefaultMigrationsTable
}

func (r *Runner) createMigrationsTable() error {
	if r.db.Migrator().HasTable(&Migrations{}) {
		return nil
	}
	err := r.db.Migrator().CreateTable(&Migrations{})
	if err != nil {
		return err
	}
	return nil
}

func (r *Runner) isMigrationApplied(db *gorm.DB, file string) bool {
	var exists bool

	err := db.Model(&Migrations{}).
		Select("count(*)").
		Where("file = ?", file).
		Find(&exists).
		Error

	return err == nil && exists
}

func (r *Runner) saveAppliedMigration(db *gorm.DB, file string) error {

	m := Migrations{
		File:    file,
		Applied: true,
	}

	return db.Create(&m).Error
}

func (r *Runner) saveRevertedMigration(db *gorm.DB, file string) error {
	return db.Where("file = ?", file).Delete(&Migrations{}).Error
}

func (r *Runner) lastAppliedMigrations(limit int) ([]string, error) {
	var migrations = make([]Migrations, 0, limit)
	var files = make([]string, 0, limit)

	err := r.db.Select("file").
		Where("applied = ?", true).
		Order("file DESC").
		Limit(limit).
		Find(&migrations).Error

	if err != nil {
		return nil, err
	}

	for _, migration := range migrations {
		files = append(files, migration.File)
	}

	return files, nil
}

func NewRunner(db *gorm.DB, migrationsList MigrationsList, d runnerDependencies) (*Runner, error) {
	runner := &Runner{
		db:             db,
		migrationsList: migrationsList,
		tableName:      DefaultMigrationsTable,
		d:              d,
	}

	if err := runner.createMigrationsTable(); err != nil {
		return nil, err
	}

	return runner, nil
}

func (r *Runner) Run(autoConfirm bool, args ...string) error {
	cmd := "up"
	if len(args) > 0 {
		cmd = args[0]
	}

	switch cmd {
	case "up":
		applied, err := r.Up()
		if err != nil {
			r.d.Logger().Error(err.Error())
			return err
		}

		if len(applied) == 0 {
			r.d.Logger().Info("No migrations to apply")
		} else {
			for _, file := range applied {
				r.d.Logger().Info(fmt.Sprintf("Applied migration %s", file))
			}
		}

		return nil
	case "down":
		toRevertCount := 1
		if len(args) > 1 {
			toRevertCount = cast.ToInt(args[1])
			if toRevertCount < 0 {
				// revert all applied migrations
				toRevertCount = len(r.migrationsList.Items())
			}
		}

		names, err := r.lastAppliedMigrations(toRevertCount)
		if err != nil {
			r.d.Logger().Error(err.Error())
			return err
		}

		if !autoConfirm {
			confirm := false
			prompt := &survey.Confirm{
				Message: fmt.Sprintf(
					"\n%v\nDo you really want to revert the last %d applied migration(s)?",
					strings.Join(names, "\n"),
					toRevertCount,
				),
			}
			survey.AskOne(prompt, &confirm)
			if !confirm {
				fmt.Println("The command has been cancelled")
				return nil
			}
		}

		reverted, err := r.Down(toRevertCount)
		if err != nil {
			r.d.Logger().Error(err.Error())
			return err
		}

		if len(reverted) == 0 {
			r.d.Logger().Info("No migrations to revert")
		} else {
			for _, file := range reverted {
				r.d.Logger().Info(fmt.Sprintf("Reverted migration %s", file))
			}
		}

		return nil
	case "history-sync":
		// todo implement history sync by looking at pocketbase logic
		// leaving it now because don't know what it does
		panic("not implemented")
	default:
		return fmt.Errorf("Unsupported command: %q\n", cmd)
	}
}

func (r *Runner) Up() ([]string, error) {
	applied := []string{}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		for _, m := range r.migrationsList.Items() {
			// skip applied
			if r.isMigrationApplied(tx, m.File) {
				continue
			}

			// ignore empty Up action
			if m.Up != nil {
				if err := m.Up(tx); err != nil {
					return fmt.Errorf("Failed to apply migration %s: %w", m.File, err)
				}
			}

			if err := r.saveAppliedMigration(tx, m.File); err != nil {
				return fmt.Errorf("Failed to save applied migration info for %s: %w", m.File, err)
			}

			applied = append(applied, m.File)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return applied, nil
}

func (r *Runner) Down(toRevertCount int) ([]string, error) {
	reverted := make([]string, 0, toRevertCount)

	names, appliedErr := r.lastAppliedMigrations(toRevertCount)
	if appliedErr != nil {
		return nil, appliedErr
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		for _, name := range names {
			for _, m := range r.migrationsList.Items() {
				if m.File != name {
					continue
				}

				// revert limit reached
				if toRevertCount-len(reverted) <= 0 {
					return nil
				}

				// ignore empty Down action
				if m.Down != nil {
					if err := m.Down(tx); err != nil {
						return fmt.Errorf("Failed to revert migration %s: %w", m.File, err)
					}
				}

				if err := r.saveRevertedMigration(tx, m.File); err != nil {
					return fmt.Errorf("Failed to save reverted migration info for %s: %w", m.File, err)
				}

				reverted = append(reverted, m.File)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return reverted, nil
}
