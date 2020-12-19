package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type middleWareHandler struct {
	r *httprouter.Router
	l *ConnLimiter
}

func NewMiddleWareHandler(r *httprouter.Router, cc int) http.Handler {
	m := middleWareHandler{}
	m.r = r
	m.l = NewConnLimiter(cc) // 限制数量
	return m
}

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()
	router.GET("/videos/:vid-id", streamHandler)  // 观看文件
	router.POST("/upload/:vid-id", uploadHandler) // 处理上传文件
	router.GET("/testpage", testPageHandler)
	return router
}

func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 要测试每个请求的响应时间
	// 下了锚点
	if !m.l.GetConn() { // 先做流控检测
		sendErrorResponse(w, http.StatusTooManyRequests, "Too many requests.")
		return
	}
	// 获取起始时间
	m.r.ServeHTTP(w, r) //调用第三方
	// 获取当前时间
	// 获取当前时间-获取起始时间
	// 取消锚点
	defer m.l.ReleaseConn()
}

func main() {
	r := RegisterHandlers()
	mh := NewMiddleWareHandler(r, 2)
	err := http.ListenAndServe(":10002", mh)
	if err != nil {
		log.Println(err)
	}
	log.Println("main fnish")
}
