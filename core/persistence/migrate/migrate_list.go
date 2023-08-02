package migrate

import (
	"gorm.io/gorm"
	"path/filepath"
	"runtime"
	"sort"
)

type Migration struct {
	File string
	Up   func(db *gorm.DB) error
	Down func(db *gorm.DB) error
}

type MigrationsList struct {
	list []*Migration
}

func (l *MigrationsList) Item(index int) *Migration {
	return l.list[index]
}

func (l *MigrationsList) Items() []*Migration {
	return l.list
}

func (l *MigrationsList) Register(
	up func(db *gorm.DB) error,
	down func(db *gorm.DB) error,
	optFilename ...string,
) {
	var file string
	if len(optFilename) > 0 {
		file = optFilename[0]
	} else {
		_, path, _, _ := runtime.Caller(1)
		file = filepath.Base(path)
	}

	l.list = append(l.list, &Migration{
		File: file,
		Up:   up,
		Down: down,
	})

	sort.Slice(l.list, func(i int, j int) bool {
		return l.list[i].File < l.list[j].File
	})
}
