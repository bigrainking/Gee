package main

import (
	"Gee/day6-template/gee"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := gee.New()
	r.Use(gee.Logger())
	// 1. 设置模板自定义函数
	t := template.FuncMap{
		"FormatAsDate": FormatAsDate,
	}
	r.SetFuncMap(t)
	// 2. 模板加载：将template/下所有模板文件加载
	r.LoadHTMLGlob("templates/*") //当前问价路径是main.go所在文件夹
	// 设置静态路径映射：将磁盘上的文件映射到路径上
	r.Static("/assets", "./static")

	// 3. 渲染： 将内容结构体输出给对应模板
	stu1 := &student{Name: "Geektutu", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "css.html", nil) //输出HTML类型的输出
	})
	r.GET("/students", func(c *gee.Context) {
		c.HTML(http.StatusOK, "arr.html", gee.H{
			"title":  "gee",
			"stuArr": [2]*student{stu1, stu2}, // 声明创建一个student类型的，类似[]int{}
		})
	})
	r.GET("/date", func(c *gee.Context) {
		c.HTML(http.StatusOK, "custom_func.html", gee.H{
			"title": "gee",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	r.Run(":9999")
}
