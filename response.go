package valigenie

import (
	"encoding/json"
	"log"
	"io"
	"bytes"
	"strings"
	"fmt"
)
//{
//    "returnCode": "0",
//    "returnErrorSolution": "",
//    "returnMessage": "",
//    "returnValue": {
//        "reply": "哈哈哈",
//        "resultType": "RESULT",
//        "actions": [
//            {
//                "name": "audioPlayGenieSource",
//                "properties": {
//                    "audioGenieId": "123"
//                }
//            }
//        ],
//        "properties": {},
//        "executeCode": "SUCCESS",
//        "msgInfo": ""
//    }
//}
type ResponseValueAskedInfo struct{
	ParameterName	string					`json:"parameterName"`		//1，询问的参数名（非实体名）
	IntentId		int64					`json:"intentId"`			//1，意图ID，从请求参数中可以获得
}
type ResponseValueAction struct{
	Name			string					`json:"name"`				//1，Action名，该名字必须设置为“audioPlayGenieSource”
	Properties		map[string]string		`json:"properties"`			//1，Action中的信息字段，“audioGenieId”的key必须设置，标示播放的开放平台存储的音频ID
}

type ResponseValue struct{
	Reply			string					`json:"reply"`						//1，回复播报语句
	ResultType		string					`json:"resultType"`					//1，回复时的状态标识（ASK_INF：信息获取，例如“请问从哪个城市出发”，在此状态下，用户说的下一句话优先进入本意图进行有效信息抽取 RESULT：正常完成交互的阶段并给出回复 CONFIRM：期待确认）
	Properties		map[string]string		`json:"properties,omitempty"`		//0，生成回复语句时携带的额外信息
	AskedInfos		[]ResponseValueAskedInfo`json:"askedInfos,omitempty"`		//0，在ASK_INF状态下，必须设置本次追问的具体参数名（开发者平台意图参数下配置的参数信息）	在ASK_INF状态下必须
	Actions			[]ResponseValueAction	`json:"actions,omitempty"`			//0，播控类信息，目前只支持播放音频
	ExecuteCode		string					`json:"executeCode"`				//1，“SUCCESS”代表执行成功；“PARAMS_ERROR”代表接收到的请求参数出错；“EXECUTE_ERROR”代表自身代码有异常；“REPLY_ERROR”代表回复结果生成出错
}


//,omitempty
type Response struct{
	req						*Request
	ReturnCode				string			`json:"returnCode"`							//1,“0”默认表示成功，其他不成功的字段自己可以确定
	ReturnErrorSolution		string			`json:"returnErrorSolution,omitempty"`		//0，出错时解决办法的描述信息
	ReturnMessage			string			`json:"returnMessage,omitempty"`			//0，返回执行成功的描述信息
	ReturnValue				ResponseValue	`json:"returnValue"`						//1，意图理解后的执行结果
}

//追问
//reply string							回复
//askedInfos []string					询问的参数名
//args... interface{}					参数
func (T *Response) Askf(reply string, askedInfos []string, args... interface{}){
	T.Ask(fmt.Sprintf(reply, args...), askedInfos)
}
func (T *Response) Ask(reply string, askedInfos []string){
	T.ReturnValue.Reply				= reply
	T.ReturnValue.ResultType		="ASK_INF"
	T.ReturnValue.ExecuteCode		="SUCCESS"
	T.ReturnCode					= "0"
	
	if askedInfos == nil {
		askedInfos = strings.Split(T.req.WordsPair(), ",")
	}
	for _, va := range askedInfos {
		askInfo := ResponseValueAskedInfo{
			ParameterName	: va,
			IntentId		: T.req.IntentId,
		}
		T.ReturnValue.AskedInfos = append(T.ReturnValue.AskedInfos, askInfo)
	}
}

//确认
//reply string							回复
//args... interface{}					参数
func (T *Response) Confirmf(reply string, args... interface{}){
	T.Confirm(fmt.Sprintf(reply, args...))
}
func (T *Response) Confirm(reply string){
	T.ReturnValue.Reply				= reply
	T.ReturnValue.ResultType		="CONFIRM"
	T.ReturnValue.ExecuteCode		="SUCCESS"
	T.ReturnCode					= "0"
}

//结果
//reply string							回复
//args... interface{}					参数
func (T *Response) Resultf(reply string, args... interface{}){
	T.Result(fmt.Sprintf(reply, args...))
}
func (T *Response) Result(reply string){
	T.ReturnValue.Reply				= reply
	T.ReturnValue.ResultType		="RESULT"
	T.ReturnValue.ExecuteCode		="SUCCESS"
	T.ReturnCode					= "0"
}

//错误
//errStr, code string		文本，错误代码
func (T *Response) Error(errStr, code string) {
	var err string
	switch code {
	case "400":
		err	= "PARAMS_ERROR"
	case "500","":
		err	= "EXECUTE_ERROR"
	case "404":
		err = "REPLY_ERROR"
	}
	T.ReturnValue.ResultType	="RESULT"
	T.ReturnValue.ExecuteCode	= err
	T.ReturnCode				= code
	T.ReturnErrorSolution		= errStr
}

//写入到w
//w io.Writer	写入
func (T *Response) WriteTo(w io.Writer) (n int64, err error) {
	if T.ReturnCode == "" {
		T.ReturnValue.Reply				= "内部错误，这是程序员忘了调用回复接口造成的。请联系官方修复。谢谢！"
		T.ReturnValue.ResultType		="RESULT"
		T.ReturnValue.ExecuteCode		="SUCCESS"
		T.ReturnCode					= "0"
	}
	buf := bytes.NewBuffer(nil)
	buf.Grow(1024)
	err = json.NewEncoder(buf).Encode(T)
	if err != nil {
		w.Write([]byte(`{"returnCode":"-1", "returnMessage":{"ExecuteCode":"REPLY_ERROR"}}`))
		log.Println(err)
		return 0, err
	}
	return buf.WriteTo(w)
}















