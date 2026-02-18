# wyy-backend 小海音乐后端

小海音乐是一个仿照网易云的项目，包括前端 react 编写和后端 golang 编写

使用 swagger 生成文档要在根目录下运行swag init -g cmd/app/main.go -o cmd/app/docs --parseDependency --parseInternal  
生成的 api 文档在http://localhost:8080/swagger/index.html#/这个链接下面