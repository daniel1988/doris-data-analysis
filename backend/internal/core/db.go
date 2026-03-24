package core

import (
	"fmt"
	"time"

	"gitee.com/dmp_admin_v2/backend/internal/common"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	ProjectDbPrefix = "dmp_"
)

var (
	centerDB *gorm.DB
	dorisDBs = make(map[string]*gorm.DB)
)

// GetProjectDatabaseName 获取项目数据库名称
func GetProjectDatabaseName(projectAlias string) string {
	return ProjectDbPrefix + projectAlias
}

func Init(configPath string) error {
	cfg, err := LoadConfig(configPath)
	if err != nil {
		return err
	}

	// 初始化日志
	common.InitLogger(cfg.Logger)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%v&loc=%s",
		cfg.Doris.User, cfg.Doris.Password, cfg.Doris.Host, cfg.Doris.Port, cfg.Doris.Database,
		cfg.Doris.Charset, cfg.Doris.ParseTime, cfg.Doris.Loc)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect center db: %w", err)
	}

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(cfg.Doris.MaxIdleConns)
		sqlDB.SetMaxOpenConns(cfg.Doris.MaxOpenConns)
		if duration, err := time.ParseDuration(cfg.Doris.ConnMaxLifetime); err == nil {
			sqlDB.SetConnMaxLifetime(duration)
		}
	}

	centerDB = db
	return nil
}

func GetProjectCenter() *gorm.DB {
	if centerDB == nil {
		// 如果没有显式初始化，则尝试加载默认配置初始化一次
		if err := Init(""); err != nil {
			common.Logger.Errorf("Failed to auto init db: %v", err)
			return nil
		}
	}
	return centerDB
}

func GetProjectDoris(projectAlias string) (*gorm.DB, error) {
	if db, ok := dorisDBs[projectAlias]; ok {
		return db, nil
	}

	// 确保 centerDB 和 GlobalConfig 已初始化
	if GetProjectCenter() == nil {
		return nil, fmt.Errorf("failed to initialize project center")
	}

	// 使用全局配置连接到项目的数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%v&loc=%s",
		GlobalConfig.Doris.User, GlobalConfig.Doris.Password, GlobalConfig.Doris.Host, GlobalConfig.Doris.Port,
		GetProjectDatabaseName(projectAlias), // 使用 dmp_项目别名 作为数据库名称
		GlobalConfig.Doris.Charset, GlobalConfig.Doris.ParseTime, GlobalConfig.Doris.Loc)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect doris for project %s: %w", projectAlias, err)
	}

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(GlobalConfig.Doris.MaxIdleConns)
		sqlDB.SetMaxOpenConns(GlobalConfig.Doris.MaxOpenConns)
		if duration, err := time.ParseDuration(GlobalConfig.Doris.ConnMaxLifetime); err == nil {
			sqlDB.SetConnMaxLifetime(duration)
		}
	}

	dorisDBs[projectAlias] = db
	return db, nil
}
