package main

import (
	"net/http"

	"Gee/day2-context/gee"
)

func main() {
	r := gee.New()
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})
	r.GET("/hello", func(c *gee.Context) {
		// expect /hello?name=geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path, c.Method)
	})

	r.POST("/login", func(c *gee.Context) {
		// 构造返回的JSON数据
		c.JSON(http.StatusOK, gee.H{
			// 第一条数据是 "username":获取request中上传的username ： 通过本次会话的Context获取
			"username": c.PostForm("username"), //c.PostForm("username")获取请求request中的username
			"password": c.PostForm("password"),
		})
	})

	r.Run(":9999")
}
