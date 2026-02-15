package repo

import (
	"gorm.io/gorm"
)

type MetaRepo struct {
	db *gorm.DB
}

// NewMetaRepo 创建 MetaRepo 实例
func NewMetaRepo(db *gorm.DB) *MetaRepo {
	return &MetaRepo{db: db}
}

// GetAllTables 返回数据库中的所有表名
func (r *MetaRepo) GetAllTables() ([]string, error) {
	tables, err := r.db.Migrator().GetTables()
	if err != nil {
		return nil, err
	}
	return tables, nil
}
