package service

import (
	"errors"
	"wyy/internal/domain"
	"wyy/internal/repo"

	"golang.org/x/crypto/bcrypt" // 用于密码加密，需 go get
)

type UserService struct {
	userRepo *repo.UserRepo
}

func NewUserService(userRepo *repo.UserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

// Register 用户注册业务
func (s *UserService) Register(name, email, password string) (*domain.User, error) {
	// 1. 检查邮箱是否已存在
	existing, _ := s.userRepo.GetByEmail(email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	// 2. 加密密码
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 3. 创建领域对象
	user := &domain.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPwd),
	}

	// 4. 保存到数据库
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

// Login 用户登录
func (s *UserService) Login(email, password string) (*domain.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}
	return user, nil
}

// GetUser 获取用户信息
func (s *UserService) GetUser(id int64) (*domain.User, error) {
	return s.userRepo.GetByID(id)
}
