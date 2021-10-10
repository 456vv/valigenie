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
//        "gwCommands": [
//            {
//                "commandDomain": "AliGenie.Speaker",
//                "commandName": "Speak",
//                "payload": {
//                  	"type": "text",             //默认 text
//                    	"text": "这是语音播放的TTS需要用到的Command",
//                  	"expectSpeech": true,       //是否开麦,默认 false
//                  	"needLight": true,          //开麦时是否需要灯光提示用户
//                  	"needVoice": true,          //开麦时是否需要声音提示用户
//                  	"wakeupType": "continuity"  // 如果需要开麦，这里需要设置为 continuity
//                }
//            }
//        ]
//    }
//}

//追问的意图名称和ID，ResultType需在设置为"ASK_INF"
type ResponseValueAskedInfo struct{
	ParameterName	string					`json:"parameterName"`		//1，询问的参数名（非实体名）
	IntentId		int64					`json:"intentId"`			//1，意图ID，从请求参数中可以获得
}

//根据音频素材id播放音频素材
type ResponseValueAction struct{
	Name			string					`json:"name"`				//1，Action名，该名字必须设置为“audioPlayGenieSource”
	Properties		map[string]string		`json:"properties"`			//1，Action中的信息字段，“audioGenieId”的key必须设置，标示播放的开放平台存储的音频ID
}

//gwCommands 字段是 V3.0 SDK 中一个特殊的字段，在响应数据中携带了gwCommands 后，原 reply、actions 字段会被忽略。
//自定义技能中使用 TPL 模板，需要在响应数据中 gwCommands 字段里携带页面展示需要的数据
type ResponseGwCommands struct{
	CommandDomain	string					`json:"commandDomain"`		//1,指令命名空间
	CommandName		string					`json:"commandName"`		//1,指令名称
	Payload			map[string]interface{}	`json:"payload"`			//1,指令数据
}

//确认信息，ResultType需在设置为"CONFIRM"
type ResponseConfirmParaInfo struct{
	ConfirmParameterName	string			`json:"confirmParameterName"`	//1,用户表达匹配到此参数，表示"确定"意思
	DenyParameterName		string			`json:"denyParameterName"`		//1,用户表达匹配到此参数，表示"否定"意思
}

type ResponseValue struct{
	Reply			string					`json:"reply"`						//1，回复播报语句
	ResultType		string					`json:"resultType"`					//1，回复时的状态标识（ASK_INF：信息获取，例如“请问从哪个城市出发”，在此状态下，用户说的下一句话优先进入本意图进行有效信息抽取 RESULT：正常完成交互的阶段并给出回复 CONFIRM：期待确认）
	Properties		map[string]string		`json:"properties,omitempty"`		//0，生成回复语句时携带的额外信息
	AskedInfos		[]ResponseValueAskedInfo`json:"askedInfos,omitempty"`		//0，在ASK_INF状态下，必须设置本次追问的具体参数名（开发者平台意图参数下配置的参数信息）	在ASK_INF状态下必须
	Actions			[]ResponseValueAction	`json:"actions,omitempty"`			//0，播控类信息，支持播放音频素材和 TTS 文本
	GwCommands		[]ResponseGwCommands	`json:"gwCommands,omitempty"`		//0，最新版响应协议定义的command结构
	ConfirmParaInfo	ResponseConfirmParaInfo	`json:"confirmParaInfo,omitempty"`	//0，resultType: CONFIRM状态下可以携带的匹配用户[肯定]和[否定]回答的参数名称。
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















