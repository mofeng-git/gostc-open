package sqlite

import (
	//"github.com/glebarez/sqlite"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/ncruces/go-sqlite3/gormlite"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io"
	"log"
	"os"
	"time"
)

type Sqlite struct {
	db  *gorm.DB
	log *lumberjack.Logger
}

func (impl *Sqlite) GetDB() *gorm.DB {
	return impl.db
}

func (impl *Sqlite) AutoMigrate(table ...any) error {
	return impl.GetDB().AutoMigrate(table...)
}

func (impl *Sqlite) Close() {
	if impl.log != nil {
		_ = impl.log.Close()
	}
	impl.db.Exec("VACUUM;")
	d, err := impl.db.DB()
	if err != nil {
		return
	}
	_ = d.Close()
}

func NewDB(dbFile, logLevel string, toFile string, console bool) (*Sqlite, error) {
	var impl Sqlite
	gormConfig := gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 关闭自动建表的外键约束
	}
	gormConfig.NamingStrategy = schema.NamingStrategy{
		SingularTable: true,
	}

	loggerLevel := logger.Silent
	switch logLevel {
	case "info", "INFO", "Info":
		loggerLevel = logger.Info
	case "warn", "WARN", "Warn":
		loggerLevel = logger.Warn
	case "error", "ERROR", "Error":
		loggerLevel = logger.Error
	}
	logConfig := logger.Config{
		SlowThreshold:             time.Second, // 慢 SQL 阈值
		LogLevel:                  loggerLevel, // 日志级别
		IgnoreRecordNotFoundError: false,       // 忽略ErrRecordNotFound（记录未找到）错误
		Colorful:                  false,       // 禁用彩色打印
	}

	var writers []io.Writer
	if console {
		writers = append(writers, os.Stdout)
	}
	if toFile != "" {
		dbLog := &lumberjack.Logger{
			Filename:   toFile,
			MaxSize:    100,
			MaxAge:     30,
			MaxBackups: 10,
			LocalTime:  true,
			Compress:   true,
		}
		writers = append(writers, dbLog)
		impl.log = dbLog
	}
	gormConfig.Logger = logger.New(
		log.New(io.MultiWriter(writers...), "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logConfig,
	)
	d, err := gorm.Open(gormlite.Open(dbFile), &gormConfig)
	if err != nil {
		return nil, err
	}
	db, err := d.DB()
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	impl.db = d
	return &impl, nil
}

// 内存数据库
func NewMemoryDB(logLevel string, toFile string, console bool) (*Sqlite, error) {
	var impl Sqlite
	gormConfig := gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 关闭自动建表的外键约束
	}
	gormConfig.NamingStrategy = schema.NamingStrategy{
		SingularTable: true,
	}

	loggerLevel := logger.Silent
	switch logLevel {
	case "info", "INFO", "Info":
		loggerLevel = logger.Info
	case "warn", "WARN", "Warn":
		loggerLevel = logger.Warn
	case "error", "ERROR", "Error":
		loggerLevel = logger.Error
	}
	logConfig := logger.Config{
		SlowThreshold:             time.Second, // 慢 SQL 阈值
		LogLevel:                  loggerLevel, // 日志级别
		IgnoreRecordNotFoundError: false,       // 忽略ErrRecordNotFound（记录未找到）错误
		Colorful:                  false,       // 禁用彩色打印
	}

	var writers []io.Writer
	if console {
		writers = append(writers, os.Stdout)
	}
	if toFile != "" {
		dbLog := &lumberjack.Logger{
			Filename:   toFile,
			MaxSize:    100,
			MaxAge:     30,
			MaxBackups: 10,
			LocalTime:  true,
			Compress:   true,
		}
		writers = append(writers, dbLog)
		impl.log = dbLog
	}
	gormConfig.Logger = logger.New(
		log.New(io.MultiWriter(writers...), "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logConfig,
	)
	d, err := gorm.Open(gormlite.Open(":memory:"), &gormConfig)
	if err != nil {
		return nil, err
	}
	db, err := d.DB()
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	impl.db = d
	return &impl, nil
}
