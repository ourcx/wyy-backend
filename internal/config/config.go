package config

import "fmt"

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	// 其他模块配置...
}

type ServerConfig struct {
	Port int
	Mode string
}

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	// 连接池参数等
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
}

// 可以添加辅助方法，比如生成 DSN
func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true",
		d.User, d.Password, d.Host, d.Port, d.DBName)
}
