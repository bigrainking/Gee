package main

import (
	"fmt"
	"log"
	"net/http"
)

// 自定义实例，实现第二个参数的hander，让所有request进入该实例

// 结构体对象
type Engine struct{}

// 实现ServerHTTP
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		fmt.Fprintf(w, "url = %q", req.URL.Path)
	case "/hello":
		fmt.Fprintf(w, "url = %q", req.URL.Path)
	default:
		fmt.Fprintf(w, "404 not found!")
	}
}

func main() {
	engine := new(Engine) //创建一个实例对象
	log.Fatal(http.ListenAndServe(":8080", engine))
}
