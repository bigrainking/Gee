// 用标准库启动web服务器

package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// 注册hander到Addr
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/hello", helloHandler)
	// 监听端口，并启动服务
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// 自定义handler处理器
func indexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "url = %q", req.URL.Path)
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "url = %q", req.URL.Path)
}
