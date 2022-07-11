package main

import (
	"todo_list/conf"
	"todo_list/routes"
)

func main() {
	conf.Init() //看来想调用方法还是要大写，这个是初始化配置
	//转载路由 swag init -g common.go
	r := routes.NewRouter()
	//配置web服务端口，3000端口
	_ = r.Run(conf.HttpPort)
}
