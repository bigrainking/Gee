package main

// 使用New()创建 gee 的实例，使用 GET()方法添加路由，最后使用Run()启动Web服务。
// 这里的路由，只是静态路由，不支持/hello/:name这样的动态路由，动态路由我们将在下一次实现。

import (
	"fmt"
	"log"
	"net/http"

	"Gee/day1-httpbase/base3/gee"
	// "gee"
)

func main() {
	// 创建一个Gee实例
	r := gee.New()
	// Get方法添加路由
	r.GET("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	})
	r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	})

	// 启动服务器
	log.Fatal(r.Run(":9999"))
}
