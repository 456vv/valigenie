package valigenie


import (
	"net/http"
)

// 程序
type Application struct{
	HandleFunc				http.HandlerFunc							// 处理函数
	ValidReqTimestamp		int											// 有效时间，秒为单位
}



