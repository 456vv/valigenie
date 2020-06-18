# valigenie [![Build Status](https://travis-ci.org/456vv/valigenie.svg?branch=master)](https://travis-ci.org/456vv/valigenie)
golang valigenie，天猫精灵自定义版本

# **列表：**
```go
aligenie.go============================================================================================================================================================================================
type Aligenie struct {                                                                  // aligenie
    AppIdAttr         string                                                            // 属性id
        
    apps    map[string]*Application                                                     // app集
    m        sync.Mutex                                                                 // 锁
}           
    func (T *Aligenie) SetApp(id string, app *Application)                              // 设置APP
    func (T *Aligenie) ServeHTTP(w http.ResponseWriter, r *http.Request)                // 服务处理
application.go============================================================================================================================================================================================
type Application struct{                                                                // 程序
    HandleFunc                http.HandlerFunc                                          // 处理函数
    ValidReqTimestamp        int                                                        // 有效时间，秒为单位
}               

intent.go============================================================================================================================================================================================
type Intent struct{                                                                     // 意图
    Response        *Response                                                           // 响应
    Request         *Request                                                            // 请求
    App                 *Application                                                    // app
}       
    func (T *Intent) ServeHTTP(w http.ResponseWriter, r *http.Request)                  // 服务处理
request.go============================================================================================================================================================================================
type RequestSlotEntitie struct{
    IntentParameterId        int64        `json:"intentParameterId"`                    //1,意图参数ID
    IntentParameterName        string        `json:"intentParameterName"`               //1,意图参数名
    OriginalValue            string        `json:"originalValue"`                       //1,原始句子中抽取出来的未做处理的slot值
    StandardValue            string        `json:"standardValue"`                       //1,原始slot归一化后的值
    LiveTime                int64        `json:"liveTime"`                              //1,该slot已生存时间（会话轮数）
    CreateTimeStamp            int64        `json:"createTimeStamp"`                    //1,该slot产生的时间点
    SlotName                string        `json:"slotName"`                             //1,slot名称
    SlotValue                string        `json:"slotValue"`                           //1,slot值
}
    
type Request struct {
    SessionId    string                    `json:"sessionId"`                           //1,会话ID，session内的对话此ID相同
    BotId        int64                    `json:"botId"`                                //1,应用ID，来自于创建的应用或者技能
    Utterance    string                    `json:"utterance"`                           //1,用户输入语句
    SkillId        int64                    `json:"skillId"`                            //1,技能ID
    SkillName    string                    `json:"skillName"`                           //1,技能名称
    IntentName    string                    `json:"intentName"`                         //1,意图名称
    Token        string                    `json:"token"`                               //0,技能鉴权token，可以不需要，如果有安全需求需要配置
    RequestData    map[string]string        `json:"requestData"`                        //0,业务请求附带参数,来自于设备调用语义理解服务额外携带的信息，只做透传
    SlotEntities []RequestSlotEntitie    `json:"slotEntities"`                          //1,从用户语句中抽取出的slot参数信息
    DomainId    int64                    `json:"domainId"`                              //1,领域ID
    IntentId    int64                    `json:"intentId"`                              //1,意图ID
}   

    func (T *Request) IntentParameterNames() []string                                   // 意图参数名
    func (T *Request) WordsPair() string                                                // 词对
    func (T *Request) OriginalValue(name string) string                                 // 原值
    func (T *Request) SlotOriginalValues() map[string]string                            // Slot原值
response.go============================================================================================================================================================================================
type ResponseValueAskedInfo struct{
    ParameterName    string                    `json:"parameterName"`                   //1，询问的参数名（非实体名）
    IntentId        int64                    `json:"intentId"`                          //1，意图ID，从请求参数中可以获得
}
type ResponseValueAction struct{
    Name            string                    `json:"name"`                             //1，Action名，该名字必须设置为“audioPlayGenieSource”
    Properties        map[string]string        `json:"properties"`                      //1，Action中的信息字段，“audioGenieId”的key必须设置，标示播放的开放平台存储的音频ID
}

type ResponseValue struct{
    Reply            string                    `json:"reply"`                           //1，回复播报语句
    ResultType        string                    `json:"resultType"`                     //1，回复时的状态标识（ASK_INF：信息获取，例如“请问从哪个城市出发”，在此状态下，用户说的下一句话优先进入本意图进行有效信息抽取 RESULT：正常完成交互的阶段并给出回复 CONFIRM：期待确认）
    Properties        map[string]string        `json:"properties,omitempty"`            //0，生成回复语句时携带的额外信息
    AskedInfos        []ResponseValueAskedInfo`json:"askedInfos,omitempty"`             //0，在ASK_INF状态下，必须设置本次追问的具体参数名（开发者平台意图参数下配置的参数信息）    在ASK_INF状态下必须
    Actions            []ResponseValueAction    `json:"actions,omitempty"`              //0，播控类信息，目前只支持播放音频
    ExecuteCode        string                    `json:"executeCode"`                   //1，“SUCCESS”代表执行成功；“PARAMS_ERROR”代表接收到的请求参数出错；“EXECUTE_ERROR”代表自身代码有异常；“REPLY_ERROR”代表回复结果生成出错
}


//,omitempty
type Response struct{
    req                        *Request
    ReturnCode                string            `json:"returnCode"`                     //1,“0”默认表示成功，其他不成功的字段自己可以确定
    ReturnErrorSolution        string            `json:"returnErrorSolution,omitempty"` //0，出错时解决办法的描述信息
    ReturnMessage            string            `json:"returnMessage,omitempty"`         //0，返回执行成功的描述信息
    ReturnValue                ResponseValue    `json:"returnMessage"`                  //1，意图理解后的执行结果
}

    func (T *Response) Ask(reply string, askedInfos []string)                           //追问
    func (T *Response) Confirm(reply string)                                            //确认
    func (T *Response) Result(reply string)                                             //结果
    func (T *Response) Error(errStr, code string)                                       //错误
    func (T *Response) WriteTo(w io.Writer)                                             //写入到w

responseWriter.go============================================================================================================================================================================================
type ResponseWriter interface{
    Ask(reply string, askedInfos []string)
    Confirm(reply string)
    Result(reply string)
    Error(errStr, code string)
    WriteTo(w io.Writer)
}
```