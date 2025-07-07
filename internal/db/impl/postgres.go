package impl

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"medods_test_task/internal/db/intf"
)

type PostgresDB struct {
	db *gorm.DB
}

func NewPostgresDB() intf.Database {
	return &PostgresDB{}
}

func (p *PostgresDB) Connect(dsn string) error {
	var err error
	p.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	return err
}

func (p *PostgresDB) DB() *gorm.DB {
	return p.db
}

func (p *PostgresDB) Migrate(models ...interface{}) error {
	return p.db.AutoMigrate(models...)
}

func (p *PostgresDB) Close() error {
	db, err := p.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
