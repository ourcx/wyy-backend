package repo

import (
	"wyy/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewDB 根据配置创建数据库连接
func NewDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := cfg.DSN()
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
