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
