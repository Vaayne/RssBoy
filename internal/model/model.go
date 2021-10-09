package model

import (
	"database/sql"
	"time"

	"github.com/indes/flowerss-bot/internal/config"
	"github.com/indes/flowerss-bot/internal/log"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

var (
	db    *gorm.DB
	sqlDB *sql.DB
)

// InitDB init db object
func InitDB() {
	connectDB()
	configDB()
	updateTable()
}

func configDB() {
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
}

func updateTable() {
	createOrUpdateTable(&Subscribe{})
	createOrUpdateTable(&User{})
	createOrUpdateTable(&Source{})
	createOrUpdateTable(&Option{})
	createOrUpdateTable(&Content{})
}

// connectDB connect to db
func connectDB() {
	if config.RunMode == config.TestMode {
		return
	}

	logger := zapgorm2.New(log.Logger)
	gormConfig := &gorm.Config{
		Logger: logger,
	}
	var err error
	if config.EnableMysql {
		db, err = gorm.Open(mysql.Open(config.Mysql.GetMySQLConnectingString()), gormConfig)
	} else if config.EnablePostgreSQL {
		db, err = gorm.Open(
			postgres.New(postgres.Config{
				DSN:                  config.PostgreSQL.GetPostgreSQLConnectingString(),
				PreferSimpleProtocol: true,
			}),
			gormConfig,
		)
	} else {
		db, err = gorm.Open(sqlite.Open(config.SQLitePath), gormConfig)
	}
	if err != nil {
		zap.S().Fatalf("connect db failed, err: %+v", err)
	}
	sqlDB, err = db.DB()
	if err != nil {
		zap.S().Fatalf("get sql db failed, err: %+v", err)
	}
}

// Disconnect disconnects from the database.
func Disconnect() {
	sqlDB.Close()
}

// createOrUpdateTable create table or Migrate table
func createOrUpdateTable(model interface{}) {

	if !db.Migrator().HasTable(model) {
		db.Migrator().CreateTable(model)
	} else {
		db.AutoMigrate(model)
	}
}

//EditTime timestamp
type EditTime struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}
