package main

import (
	"gin-blog-hufeng/routes"
	"golang.org/x/sync/errgroup"
	"log"
)

var g errgroup.Group

func main() {
	// 初始化全局变量 g
	routes.InitGlobalVariable()

	// 前台服务接口

	// 后台服务接口
	g.Go(func() error {
		return routes.BackendServer().ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
