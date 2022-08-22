package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

// 1. 提供访问POSTFORM参数 QUERY查询参数
// 2. 提供了快速构造String/Data/JSON/HTML响应的方法
// 3. 提供访问 req w
// 4. 提供快速访问属性request Path Method
type Context struct {
	// origin objects// 3. 提供访问 req w
	Writer http.ResponseWriter
	Req    *http.Request
	// request info // 4.提供快速访问属性request Path Method
	Path   string
	Method string
	// response info
	StatusCode int
}

// 创建
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

// 1.1 Context提供查询request中Form表单某个key的值的功能
// 获取request提交body的表单中key对应的value
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// 1.2 Context提供访问Query参数：获取request中GET请求的query参数
func (c *Context) Query(key string) string {
	fmt.Println(c.Req.URL.Query().Get(key))
	return c.Req.URL.Query().Get(key)
}

// 2.1 快速构造JSON返回数据
// 传入需要返回的statuscode和需要返回的数据
func (c *Context) JSON(code int, obj interface{}) {
	// 1. 设置头：构造时直接设置好,赋值给Context
	c.SetHeader("Content-Type", "application/json")
	// 2. 设置返回状态值：设置statuscode
	c.SetStatusCode(code)
	// 3. 转换返回格式
	encoder := json.NewEncoder(c.Writer) //创建返回值的JSON空结构体
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}
func (c *Context) SetStatusCode(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// 2.2 快速构造String响应
// 返回string类型到Response
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.SetStatusCode(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values)))
}

func (c *Context) Data(code int, data []byte) {
	c.SetStatusCode(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.SetStatusCode(code)
	c.Writer.Write([]byte(html))
}
