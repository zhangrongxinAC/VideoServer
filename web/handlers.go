package main

import (
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/julienschmidt/httprouter"
)

type HomePage struct {
	Name string
}

type UserPage struct {
	Name string
}

func homeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cname, err1 := r.Cookie("username")
	sid, err2 := r.Cookie("session")

	if err1 != nil || err2 != nil {
		p := &HomePage{Name: "jacky"} // 对页面做注解
		t, e := template.ParseFiles("./templates/home.html")
		if e != nil {
			log.Printf("Parsing template home.html error: %s", e)
			return
		}

		t.Execute(w, p)
		return
	}

	if len(cname.Value) != 0 && len(sid.Value) != 0 {
		http.Redirect(w, r, "/userhome", http.StatusFound)
		return
	}

}

func userHomeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("Header", r.Header)
	cname, err1 := r.Cookie("username") // 读取cookie的字段
	_, err2 := r.Cookie("session")
	if err1 != nil || err2 != nil {
		http.Redirect(w, r, "/", http.StatusFound) // 重定向主页
		return
	}

	fname := r.FormValue("username")
	var p *UserPage
	if len(cname.Value) != 0 {
		p = &UserPage{Name: cname.Value} // 渲染名字
	} else if len(fname) != 0 {
		p = &UserPage{Name: fname}
	}
	t, e := template.ParseFiles("./templates/userhome.html")
	if e != nil {
		log.Printf("Parsing userhome.html error: %s", e)
		return
	}
	t.Execute(w, p)
}

func apiHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	log.Println("apiHandler Header", r.Header)
	if r.Method != http.MethodPost {
		re, _ := json.Marshal(ErrorRequestNotRecognized)
		io.WriteString(w, string(re))
		return
	}

	res, _ := ioutil.ReadAll(r.Body)
	apibody := &ApiBody{}
	if err := json.Unmarshal(res, apibody); err != nil {
		re, _ := json.Marshal(ErrorRequestBodyParseFailed)
		io.WriteString(w, string(re))
		return
	}
	request(apibody, w, r)
	defer r.Body.Close()
}

func proxyHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	u, _ := url.Parse("http://127.0.0.1:10002/")
	proxy := httputil.NewSingleHostReverseProxy(u) // 代理
	proxy.ServeHTTP(w, r)
}

func createUserProxyHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// 可以做一个网关，同一个ip不断地去访问，可以直接返回错误

	u, _ := url.Parse("http://127.0.0.1:10000/")   // 直接转发给后端，就不需要web去做解析
	proxy := httputil.NewSingleHostReverseProxy(u) // 代理
	proxy.ServeHTTP(w, r)
}
