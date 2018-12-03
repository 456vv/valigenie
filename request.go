package valigenie
	

/*
{
	"token": "1:b0c6c5f881196f4def862f9cb8e1221c16e3fe8ee7ef9c0e01e0fadc149ded0bf4a7b2e54fb48834c99bf5b56f4a857f979fb2afdcd0fae502e23a3356579d44",
	"sessionId": "e3abf4d6-5e3c-4a2d-b2e4-0ff88b829021",
	"utterance": "雀禾家居 你好",
	"skillId": 16734,
	"skillName": "雀禾家居",
	"intentId": 21492,
	"intentName": "SwitchIntent",
	"requestData": {},
	"slotEntities": [{
		"intentParameterId": 21261,//这个不是所有请求都有，注意注意
		"intentParameterName": "any",
		"originalValue": "一岁",
		"standardValue": "{\"unit\":\"岁\",\"value\":\"1\"}",
		"liveTime": 0,
		"createTimeStamp": 1528805442521,
		"slotValue": "{\"unit\":\"岁\",\"value\":\"1\"}"
	}],
	"botId": 29758,
	"domainId": 18628,
	"requestId": "20180612194925775-1062531537"
}

*/

type RequestSlotEntitie struct{
	IntentParameterId		int64		`json:"intentParameterId"`				//1,意图参数ID
	IntentParameterName		string		`json:"intentParameterName"`			//1,意图参数名
	OriginalValue			string		`json:"originalValue"`					//1,原始句子中抽取出来的未做处理的slot值
	StandardValue			string		`json:"standardValue"`					//1,原始slot归一化后的值
	LiveTime				int64		`json:"liveTime"`						//1,该slot已生存时间（会话轮数）
	CreateTimeStamp			int64		`json:"createTimeStamp"`				//1,该slot产生的时间点
	SlotName				string		`json:"slotName"`						//1,slot名称
	SlotValue				string		`json:"slotValue"`						//1,slot值
}
	
type Request struct {
	SessionId	string					`json:"sessionId"`						//1,会话ID，session内的对话此ID相同
	BotId		int64					`json:"botId"`							//1,应用ID，来自于创建的应用或者技能
	Utterance	string					`json:"utterance"`						//1,用户输入语句
	SkillId		int64					`json:"skillId"`						//1,技能ID
	SkillName	string					`json:"skillName"`						//1,技能名称
	IntentName	string					`json:"intentName"`						//1,意图名称
	Token		string					`json:"token"`							//0,技能鉴权token，可以不需要，如果有安全需求需要配置	
	RequestData	map[string]string		`json:"requestData"`					//0,业务请求附带参数,来自于设备调用语义理解服务额外携带的信息，只做透传
	SlotEntities []RequestSlotEntitie	`json:"slotEntities"`					//1,从用户语句中抽取出的slot参数信息
	DomainId	int64					`json:"domainId"`						//1,领域ID
	IntentId	int64					`json:"intentId"`						//1,意图ID
	RequestId	string					`json:"requestId"`						//1,请求Id
}

//参数名
//	[]string	参数名
func (T *Request) IntentParameterNames() []string {
	var names []string
	for _, slot  := range T.SlotEntities {
		names = append(names, slot.IntentParameterName)
	}
	return names
}

//原值
//	name string		参数名
//	string			值
func (T *Request) OriginalValue(name string) string {
	for _, slot  := range T.SlotEntities {
		if slot.IntentParameterName == name {
			return slot.OriginalValue
		}
	}
	return ""
}

//Slot原值
//	map[string]string	参数名/原值
func (T *Request) SlotOriginalValues() map[string]string {
	m := make(map[string]string)
	for _, slot  := range T.SlotEntities {
		m[slot.IntentParameterName] = slot.OriginalValue
	}
	return m
}

//词对
//	string	词对字符
func (T *Request) WordsPair() string {
	var wordsPair string
	for _, slot  := range T.SlotEntities {
		if wordsPair != "" {
			wordsPair+=","
		}
		wordsPair+=slot.IntentParameterName
	}
	return wordsPair
}




