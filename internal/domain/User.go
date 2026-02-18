package domain

type User struct {
	Email    string `gorm:"primaryKey"`
	Password string
	Name     string
	ID       int64
}
