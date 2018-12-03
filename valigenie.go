package valigenie

import (
	"sync"
	"net/http"
	"encoding/json"
	"fmt"
)

// aligenie
type Aligenie struct {
	AppIdAttr 		string												// 属性id

	apps	map[string]*Application										// app集
	m		sync.Mutex													// 锁
}

//设置APP
//	id string			id名称
//	app *Application	app配置
func (T *Aligenie) SetApp(id string, app *Application){
	T.m.Lock()
	defer T.m.Unlock()
	if T.apps == nil {
		T.apps = make(map[string]*Application)
	}
	if app == nil {
		delete(T.apps, id)
		return
	}
	T.apps[id]	= app
}

//服务处理
//	w http.ResponseWriter	http响应
//	r *http.Request			http请求
func (T *Aligenie) ServeHTTP(w http.ResponseWriter, r *http.Request){
	T.m.Lock()
	var (
		query		= r.URL.Query()
		argId		= query.Get(T.AppIdAttr)
		res 		= &Response{}
	)
	if argId == "" {
		argId = r.Header.Get(T.AppIdAttr)
	}
	app, ok := T.apps[argId]
	T.m.Unlock()

	if !ok {
		res.Error(fmt.Sprintf("valigenie: 参数验证不通过error(%s=%v)", T.AppIdAttr, argId), "400")
		res.WriteTo(w)
		return
	}

	defer r.Body.Close()
	req := new(Request)
	res.req = req
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		res.Error(fmt.Sprintf("valigenie: 数据验证不通过error(%v)", err), "400")
		res.WriteTo(w)
		return
	}

	intent := &Intent{
			Response			: res,
			Request				: req,
			App					: app,
		}
	intent.ServeHTTP(w, r)
}














