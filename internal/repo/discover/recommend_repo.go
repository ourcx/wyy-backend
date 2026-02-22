package repo

import "gorm.io/gorm"

type RecommendRepo struct {
	db *gorm.DB
}

func NewRecommendRepo(db *gorm.DB) *RecommendRepo {
	return &RecommendRepo{
		db: db,
	}
}

//推荐系统核心设计
