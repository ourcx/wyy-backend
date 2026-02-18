// @title			Wyy API
// @version		1.0
// @description	音乐项目 API 文档
// @host			localhost:8080
// @BasePath		/api
package main

import (
	"log"
)

func main() {
	application, err := NewApp("configs/config.yaml")
	if err != nil {
		log.Fatalf("failed to init app: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
