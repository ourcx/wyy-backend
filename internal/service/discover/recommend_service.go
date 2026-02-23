package service

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"

	"wyy/internal/repo/discover"
)

type RecommendService struct {
	RecommendRepo *repo.RecommendRepo
}

func NewRecommendService(RecommendRepo *repo.RecommendRepo) *RecommendService {
	return &RecommendService{RecommendRepo: RecommendRepo}
}

// 推荐核心模块
// RecommendItem 召回阶段返回的候选项目
type RecommendItem struct {
	SongID string
	Score  float64 // 算法原始得分
	Reason string  // 推荐理由（可选）
}

// RecommendRequest 推荐请求参数（可根据需要扩展）
type RecommendRequest struct {
	UserID string
	Size   int // 期望召回数量
	// 其他上下文：时间、设备、位置等
}

type Recommender interface {
	// Recommend 返回候选歌曲列表
	Recommend(ctx context.Context, req *RecommendRequest) ([]*RecommendItem, error)
}
type UserBasedCFRecommender struct {
	userActionRepo repo.UserActionRepo
	cacheRepo      repo.CacheRepo
	songRepo       repo.SongRepo
	// 配置：相似度计算方式、邻居数量等
	topK                int     // 相似用户数量
	similarityThreshold float64 // 相似度阈值
}

func NewUserBasedCFRecommender(userActionRepo repo.UserActionRepo, cacheRepo repo.CacheRepo, songRepo repo.SongRepo, topK int, similarityThreshold float64) *UserBasedCFRecommender {
	return &UserBasedCFRecommender{
		userActionRepo:      userActionRepo,
		cacheRepo:           cacheRepo,
		songRepo:            songRepo,
		topK:                topK,
		similarityThreshold: similarityThreshold,
	}
}

func (r *UserBasedCFRecommender) Recommend(ctx context.Context, req *RecommendRequest) ([]*RecommendItem, error) {
	// 1. 获取目标用户的行为
	actions, err := r.userActionRepo.GetAllUserActions(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	// 2. 查找相似用户（可从缓存或实时计算）
	similarUsers, err := r.cacheRepo.GetSimilarUsers(ctx, req.UserID)
	if err != nil || len(similarUsers) == 0 {
		// 实时计算相似度（利用 goroutine 并发）
		similarUsers = r.calcSimilarUsers(ctx, req.UserID, actions)
		// 异步缓存
		go r.cacheRepo.SetSimilarUsers(ctx, req.UserID, similarUsers, 3600)
	}
	// 3. 聚合相似用户听过的歌曲
	candidates := r.aggregateSongsFromUsers(ctx, similarUsers, req.Size)
	return candidates, nil
}

// calcSimilarUsers 计算相似用户
func (r *UserBasedCFRecommender) calcSimilarUsers(ctx context.Context, userID string, userActions []*repo.UserAction) []string {
	// 获取所有用户的行为（实际应用中应该从数据库或缓存中获取）
	// 这里简化处理，实际应该调用 BatchGetUserActions
	// 为了演示，我们假设有一些相似用户
	// TODO: 实现真实的相似度计算算法（如余弦相似度、Jaccard 相似度等）

	// 构建用户的歌曲集合
	userSongSet := make(map[string]bool)
	for _, action := range userActions {
		userSongSet[action.SongID] = true
	}

	// 模拟相似用户（实际应用中应该从数据库获取所有用户并计算相似度）
	// 这里返回一些假数据
	similarUsers := []string{
		"user_001",
		"user_002",
		"user_003",
	}

	// 限制返回数量
	if len(similarUsers) > r.topK {
		similarUsers = similarUsers[:r.topK]
	}

	return similarUsers
}

// aggregateSongsFromUsers 从相似用户中聚合歌曲
func (r *UserBasedCFRecommender) aggregateSongsFromUsers(ctx context.Context, similarUsers []string, size int) []*RecommendItem {
	// 获取相似用户的行为
	actionsMap, err := r.userActionRepo.BatchGetUserActions(ctx, similarUsers)
	if err != nil {
		return nil
	}

	// 统计歌曲出现次数和得分
	songScores := make(map[string]float64)
	songCount := make(map[string]int)

	for _, userID := range similarUsers {
		actions := actionsMap[userID]
		for _, action := range actions {
			// 累加得分（可以根据行为类型加权）
			weight := 1.0
			switch action.Action {
			case "like":
				weight = 2.0
			case "play":
				weight = 1.0
			case "skip":
				weight = 0.5
			}
			songScores[action.SongID] += action.Value * weight
			songCount[action.SongID]++
		}
	}

	// 转换为 RecommendItem 列表
	items := make([]*RecommendItem, 0, len(songScores))
	for songID, score := range songScores {
		items = append(items, &RecommendItem{
			SongID: songID,
			Score:  score / float64(songCount[songID]), // 使用平均分
			Reason: "协同过滤推荐",
		})
	}

	// 按得分排序
	sortByScore(items)

	// 限制返回数量
	if len(items) > size {
		items = items[:size]
	}

	return items
}

type Filter interface {
	Filter(ctx context.Context, userID string, items []*RecommendItem) ([]*RecommendItem, error)
}

// Ranker 排序器接口
type Ranker interface {
	Rank(ctx context.Context, userID string, items []*RecommendItem) ([]*RecommendItem, error)
}

// ScoreBasedRanker 基于得分的排序器
type ScoreBasedRanker struct{}

func NewScoreBasedRanker() *ScoreBasedRanker {
	return &ScoreBasedRanker{}
}

func (r *ScoreBasedRanker) Rank(ctx context.Context, userID string, items []*RecommendItem) ([]*RecommendItem, error) {
	// 按得分降序排序
	sortByScore(items)
	return items, nil
}

// Mixer 混排器接口
type Mixer interface {
	Mix(ctx context.Context, itemGroups [][]*RecommendItem, size int) ([]*RecommendItem, error)
}

// SimpleMixer 简单混排器（直接合并）
type SimpleMixer struct{}

func NewSimpleMixer() *SimpleMixer {
	return &SimpleMixer{}
}

func (m *SimpleMixer) Mix(ctx context.Context, itemGroups [][]*RecommendItem, size int) ([]*RecommendItem, error) {
	// 合并所有组
	allItems := make([]*RecommendItem, 0)
	for _, group := range itemGroups {
		allItems = append(allItems, group...)
	}

	// 按得分排序
	sortByScore(allItems)

	// 截断到指定大小
	if len(allItems) > size {
		allItems = allItems[:size]
	}

	return allItems, nil
}

// WeightedMixer 加权混排器（按比例混排）
type WeightedMixer struct {
	weights []float64 // 每个召回源的权重
}

func NewWeightedMixer(weights []float64) *WeightedMixer {
	return &WeightedMixer{weights: weights}
}

func (m *WeightedMixer) Mix(ctx context.Context, itemGroups [][]*RecommendItem, size int) ([]*RecommendItem, error) {
	if len(itemGroups) == 0 {
		return []*RecommendItem{}, nil
	}

	// 计算每个召回源应该返回的数量
	weights := m.weights
	if len(weights) != len(itemGroups) {
		// 如果权重数量不匹配，使用均等权重
		weights = make([]float64, len(itemGroups))
		for i := range weights {
			weights[i] = 1.0
		}
	}

	totalWeight := 0.0
	for _, w := range weights {
		totalWeight += w
	}

	// 计算每个源的配额
	quotas := make([]int, len(itemGroups))
	for i, w := range weights {
		quotas[i] = int(float64(size) * w / totalWeight)
	}

	// 从每个源中取出指定数量的项目
	result := make([]*RecommendItem, 0, size)
	for i, group := range itemGroups {
		// 按得分排序
		sortedGroup := make([]*RecommendItem, len(group))
		copy(sortedGroup, group)
		sortByScore(sortedGroup)

		// 取出配额数量
		count := quotas[i]
		if count > len(sortedGroup) {
			count = len(sortedGroup)
		}
		result = append(result, sortedGroup[:count]...)
	}

	return result, nil
}

// 已听过滤实现
type ListenedFilter struct {
	userActionRepo repo.UserActionRepo
}

func NewListenedFilter(userActionRepo repo.UserActionRepo) *ListenedFilter {
	return &ListenedFilter{
		userActionRepo: userActionRepo,
	}
}

func (f *ListenedFilter) Filter(ctx context.Context, userID string, items []*RecommendItem) ([]*RecommendItem, error) {
	// 获取用户已听歌曲ID集合
	playedSet, err := f.getPlayedSet(ctx, userID)
	if err != nil {
		return nil, err
	}
	filtered := make([]*RecommendItem, 0, len(items))
	for _, item := range items {
		if !playedSet[item.SongID] {
			filtered = append(filtered, item)
		}
	}
	return filtered, nil
}

// getPlayedSet 获取用户已听歌曲集合
func (f *ListenedFilter) getPlayedSet(ctx context.Context, userID string) (map[string]bool, error) {
	actions, err := f.userActionRepo.GetAllUserActions(ctx, userID)
	if err != nil {
		return nil, err
	}

	playedSet := make(map[string]bool)
	for _, action := range actions {
		playedSet[action.SongID] = true
	}

	return playedSet, nil
}

type RecommendationService struct {
	recallers []Recommender // 多路召回器
	ranker    Ranker
	filters   []Filter
	mixer     Mixer
	songRepo  repo.SongRepo // 用于获取详情
}

func NewRecommendationService(recallers []Recommender, ranker Ranker, filters []Filter, mixer Mixer, songRepo repo.SongRepo) *RecommendationService {
	return &RecommendationService{
		recallers: recallers,
		ranker:    ranker,
		filters:   filters,
		mixer:     mixer,
		songRepo:  songRepo,
	}
}

func (s *RecommendationService) GetRecommendations(ctx context.Context, userID string, size int) ([]*repo.Song, error) {
	// 1. 并发执行多路召回
	var mu sync.Mutex
	allCandidates := make([]*RecommendItem, 0)
	eg, ctx := errgroup.WithContext(ctx)
	for _, recaller := range s.recallers {
		recaller := recaller
		eg.Go(func() error {
			items, err := recaller.Recommend(ctx, &RecommendRequest{UserID: userID, Size: size})
			if err != nil {
				return err
			}
			mu.Lock()
			allCandidates = append(allCandidates, items...)
			mu.Unlock()
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	// 2. 去重（保留得分高的）
	allCandidates = s.deduplicate(allCandidates)

	// 3. 排序
	ranked, err := s.ranker.Rank(ctx, userID, allCandidates)
	if err != nil {
		return nil, err
	}

	// 4. 依次应用过滤器
	for _, filter := range s.filters {
		ranked, err = filter.Filter(ctx, userID, ranked)
		if err != nil {
			return nil, err
		}
	}

	// 5. 混排并截断
	finalItems, err := s.mixer.Mix(ctx, [][]*RecommendItem{ranked}, size)
	if err != nil {
		return nil, err
	}

	// 6. 获取歌曲详情
	songIDs := extractSongIDs(finalItems)
	songs, err := s.songRepo.GetSongs(ctx, songIDs)
	if err != nil {
		return nil, err
	}

	// 7. 按 finalItems 顺序返回歌曲
	return sortSongsByOrder(songs, finalItems), nil
}

// 辅助方法
func (s *RecommendationService) deduplicate(items []*RecommendItem) []*RecommendItem {
	seen := make(map[string]bool)
	result := make([]*RecommendItem, 0, len(items))
	for _, item := range items {
		if !seen[item.SongID] {
			seen[item.SongID] = true
			result = append(result, item)
		}
	}
	return result
}

// extractSongIDs 从推荐项中提取歌曲ID
func extractSongIDs(items []*RecommendItem) []string {
	songIDs := make([]string, 0, len(items))
	for _, item := range items {
		songIDs = append(songIDs, item.SongID)
	}
	return songIDs
}

// sortSongsByOrder 按照推荐项的顺序排序歌曲
func sortSongsByOrder(songs []*repo.Song, items []*RecommendItem) []*repo.Song {
	// 构建歌曲ID到歌曲的映射
	songMap := make(map[string]*repo.Song)
	for _, song := range songs {
		songMap[song.ID] = song
	}

	// 按照推荐项的顺序返回歌曲
	result := make([]*repo.Song, 0, len(items))
	for _, item := range items {
		if song, exists := songMap[item.SongID]; exists {
			result = append(result, song)
		}
	}

	return result
}

// sortByScore 按得分降序排序推荐项
func sortByScore(items []*RecommendItem) {
	for i := 0; i < len(items)-1; i++ {
		for j := i + 1; j < len(items); j++ {
			if items[i].Score < items[j].Score {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
}
