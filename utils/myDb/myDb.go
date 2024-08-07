package myDb

import (
	"database/sql"
	"go-instaloader/config"
	"go-instaloader/utils/rlog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	syslog "log"
	"time"
)

var (
	db    *gorm.DB
	sqlDB *sql.DB
)

func GetDb() *gorm.DB {
	return db
}

func openDB(dsn string, config *gorm.Config, maxIdleConns, maxOpenConns int) (err error) {
	if config == nil {
		config = &gorm.Config{}
	}

	if config.NamingStrategy == nil {
		config.NamingStrategy = schema.NamingStrategy{
			TablePrefix:   "t_",
			SingularTable: true,
		}
	}

	if db, err = gorm.Open(mysql.Open(dsn), config); err != nil {
		rlog.Infof("opens database failed: %v", err.Error())
		return
	}

	if sqlDB, err = db.DB(); err == nil {
		sqlDB.SetMaxIdleConns(maxIdleConns)
		sqlDB.SetMaxOpenConns(maxOpenConns)
	} else {
		rlog.Info(err)
	}

	//rlog.Info("Successfully connected to mysql service")
	return
}

func DBInit() {
	gormConf := &gorm.Config{}
	newLogger := logger.New(
		syslog.New(rlog.GetLogger().GetWriter(), "\r\n[db]", rlog.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		})
	gormConf.Logger = newLogger
	err := openDB(config.Instance.MySqlUrl, gormConf,
		config.Instance.MySqlMaxIdle, config.Instance.MySqlMaxOpen)
	if err != nil {
		rlog.Info(err)
		panic(err)
	}
	if config.Instance.ShowSql {
		db.Logger = db.Logger.LogMode(logger.Info)
	} else {
		db.Logger = db.Logger.LogMode(logger.Warn)
	}
	rlog.Debug("mySQL connection established")
}
