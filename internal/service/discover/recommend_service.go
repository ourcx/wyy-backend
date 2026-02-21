package service

import (
	"wyy/internal/repo/discover"
)

type RecommendService struct {
	RecommendRepo *repo.RecommendRepo
}

func NewRecommendService(RecommendRepo *repo.RecommendRepo) *RecommendService {
	return &RecommendService{RecommendRepo: RecommendRepo}
}
