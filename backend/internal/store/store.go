package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/glebarez/sqlite"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"

	"html-ppt/backend/internal/config"
)

// Store 是数据库层的统一入口。目前只有 GORM 句柄，
// 后续按需往这里加仓储方法（如 sessions.Save / users.FindByAPIKey）。
type Store struct {
	DB *gorm.DB
}

// Open 连接 MySQL：库不存在则自动创建，并对占位模型执行 AutoMigrate。
func Open(ctx context.Context, cfg config.DB) (*Store, error) {
	// 第一步：不带库名连上服务器，建库（这样用户拿到手就能跑，不用手动 CREATE DATABASE）
	server, err := sql.Open("mysql", dsn(cfg, false))
	if err != nil {
		return nil, err
	}
	_, execErr := server.ExecContext(ctx,
		fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", cfg.Name))
	closeErr := server.Close()
	if execErr != nil {
		return nil, fmt.Errorf("create database %s: %w", cfg.Name, execErr)
	}
	if closeErr != nil {
		return nil, closeErr
	}

	// 第二步：正式连接 + 建表
	db, err := gorm.Open(gormmysql.Open(dsn(cfg, true)), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.WithContext(ctx).AutoMigrate(&User{}, &Deck{}, &ChatSession{}); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}
	return &Store{DB: db}, nil
}

// Ping 供健康检查使用。
func (s *Store) Ping(ctx context.Context) error {
	sqlDB, err := s.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

// OpenMemory 内存 SQLite（纯 Go 驱动，无 cgo），专供测试：
// 与 MySQL 同一套模型与 AutoMigrate，测试里测的表结构和生产一致。
// 不放进 _test.go 是因为使用方是其他包的测试（跨包不可见 _test.go 导出）。
func OpenMemory(ctx context.Context) (*Store, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open memory sqlite: %w", err)
	}
	if err := db.WithContext(ctx).AutoMigrate(&User{}, &Deck{}, &ChatSession{}); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}
	return &Store{DB: db}, nil
}

func dsn(cfg config.DB, withDB bool) string {
	c := mysql.Config{
		User:                 cfg.User,
		Passwd:               cfg.Password,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		AllowNativePasswords: true,
		ParseTime:            true,
		Loc:                  time.Local,
		Collation:            "utf8mb4_unicode_ci",
	}
	if withDB {
		c.DBName = cfg.Name
	}
	return c.FormatDSN()
}
