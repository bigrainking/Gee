package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Context对象
type Context struct {
	// origin
	Req    *http.Request
	Writer http.ResponseWriter
	// 常用属性
	Method string
	Path   string
	Params map[string]string // 记录当前请求内容
	// 状态码
	StatusCode int
}

// 构造函数 : 不对外
func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Req:    r,
		Writer: w,
		Method: r.Method,
		Path:   r.URL.Path}
}

// 为JSON方便简洁:创建装JSON数据的结构体,对外提供
type H map[string]interface{}

// 1. 获取request参数:query postForm
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) PostForm(key string) string {
	return c.Req.PostFormValue(key)
}

func (c *Context) Param(key string) string {
	return c.Params[key]
}

// 2. 构造响应请求JSON HTML String Data
// 		1. 设置格式 2. 设置statuscode 3. 数据填充
//		下面每个的Response类型不同
func (c *Context) JSON(status int, data H) {
	// 设置头部信息Json
	c.SetHeader("Content-Type", "application/json")
	// 设置状态码
	c.SetStatus(status)
	// data参数填充
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(data); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) String(status int, form string, values ...interface{}) {
	c.SetStatus(status)
	c.SetHeader("Content-Type", "text/plain")
	// s := []interface{}(values)
	c.Writer.Write([]byte(fmt.Sprintf(form, values...)))
}

func (c *Context) HTML(status int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.SetStatus(status)
	c.Writer.Write([]byte(html))
}

// Response没有类型
func (c *Context) Data(status int, data []byte) {
	c.SetStatus(status)
	c.Writer.Write(data) //直接就是对应形式不用转换
}

// 头部设置信息
func (c *Context) SetHeader(key, value string) {
	// 头部参数设置
	c.Writer.Header().Set(key, value)
}

func (c *Context) SetStatus(code int) {
	c.StatusCode = code //context对象中status
	// 状态码设置
	c.Writer.WriteHeader(code) // 构造的响应中设置
}
