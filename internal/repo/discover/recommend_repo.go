package repo

import (
	"context"

	"gorm.io/gorm"
)

type RecommendRepo struct {
	db *gorm.DB
}

func NewRecommendRepo(db *gorm.DB) *RecommendRepo {
	return &RecommendRepo{
		db: db,
	}
}

// 推荐系统核心设计
// PlayRecord 代表一次播放记录（可能包含时间、时长等）
type PlayRecord struct {
	UserID   string
	SongID   string
	PlayedAt int64 // 时间戳
	Duration int   // 播放时长（秒）
}

// UserAction 代表用户的完整行为（可用于评分计算）
type UserAction struct {
	UserID    string
	SongID    string
	Action    string  // "play", "like", "skip", "rate"
	Value     float64 // 如果是评分则为分数，否则为 0/1
	Timestamp int64
}

type UserActionRepo interface {
	// 获取用户最近的播放记录（用于实时个性化召回）
	GetRecentPlays(ctx context.Context, userID string, limit int) ([]*PlayRecord, error)

	// 获取用户的所有历史行为（用于离线计算或协同过滤）
	GetAllUserActions(ctx context.Context, userID string) ([]*UserAction, error)

	// 批量获取多个用户的行为（用于 User-Based CF）
	BatchGetUserActions(ctx context.Context, userIDs []string) (map[string][]*UserAction, error)

	// 插入新的用户行为（实时写入，可能异步）
	InsertAction(ctx context.Context, action *UserAction) error

	// 获取用户对特定歌曲的行为（用于判断是否已听过）
	GetUserActionOnSong(ctx context.Context, userID, songID string) (*UserAction, error)
}
type Song struct {
	ID          string
	Name        string
	Artist      string
	Album       string
	Tags        []string // 风格标签
	Duration    int
	PublishTime int64
	Features    []float64 // 音频特征向量（可选）
}

type SongRepo interface {
	// 根据 ID 批量获取歌曲信息（用于补充推荐结果详情）
	GetSongs(ctx context.Context, songIDs []string) ([]*Song, error)

	// 获取某标签下的热门歌曲（基于内容召回）
	GetTopSongsByTag(ctx context.Context, tag string, limit int) ([]*Song, error)

	// 获取相似歌曲（基于音频特征或协同过滤结果）
	GetSimilarSongs(ctx context.Context, songID string, limit int) ([]*Song, error)

	// 获取艺人的代表歌曲
	GetSongsByArtist(ctx context.Context, artist string, limit int) ([]*Song, error)
}
type UserProfile struct {
	UserID           string
	PreferredTags    map[string]float64 // 标签权重
	PreferredArtists []string
	Vector           []float64 // 隐向量（Embedding）
	UpdateTime       int64
}

type UserProfileRepo interface {
	// 获取用户画像
	GetUserProfile(ctx context.Context, userID string) (*UserProfile, error)

	// 更新用户画像（离线任务调用）
	UpdateUserProfile(ctx context.Context, profile *UserProfile) error
}

type CacheRepo interface {
	// 获取全局热门歌曲（用于冷启动、兜底）
	GetGlobalHotSongs(ctx context.Context, limit int) ([]*Song, error)

	// 获取最新上架歌曲
	GetNewestSongs(ctx context.Context, limit int) ([]*Song, error)

	// 获取用户的协同过滤相似用户列表（缓存计算结果）
	GetSimilarUsers(ctx context.Context, userID string) ([]string, error)

	// 设置相似用户列表
	SetSimilarUsers(ctx context.Context, userID string, similarUsers []string, ttl int64) error
}
type Repository struct {
	UserAction  UserActionRepo
	Song        SongRepo
	UserProfile UserProfileRepo
	Cache       CacheRepo
}
