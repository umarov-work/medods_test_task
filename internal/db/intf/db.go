package intf

import (
	"gorm.io/gorm"
)

type Database interface {
	Connect(dsn string) error
	DB() *gorm.DB
	Migrate(models ...interface{}) error
	Close() error
}
