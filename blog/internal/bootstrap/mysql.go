package bootstrap

import (
	"database/sql"
	"time"

	"blog/internal/config"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

func InitMySQL(cfg config.MySQLConfig, log *zap.Logger) (*gorm.DB, *sql.DB, error) {
	gormLogger := glogger.New(
		zap.NewStdLog(log.Named("gorm")),
		glogger.Config{
			SlowThreshold:             500 * time.Millisecond,
			LogLevel:                  glogger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{Logger: gormLogger})
	if err != nil {
		return nil, nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}
	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetimeSec > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetimeSec) * time.Second)
	}
	return db, sqlDB, nil
}
