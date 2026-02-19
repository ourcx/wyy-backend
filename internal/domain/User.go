package domain

// User 用户的类型
type User struct {
	Email    string `gorm:"primaryKey"`
	Password string
	Name     string
	ID       int64
}

// UserRegister 用户注册传入的类型
type UserRegister struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UserLogin 用户登录的传入的类型
type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
