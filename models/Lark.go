package models

type LarkUser struct {
	Email  string `json:"email"`
	UserId string `json:"user_id"`
}

type LarkUserList struct {
	UserList []LarkUser
}

type LarkItem struct {
	Avatar      string `json:"avatar"`
	ChatId      string `json:"chat_id"`
	Description string `json:"description"`
	//External    string `json:"external"`
	Name        string `json:"name"`
	OwnerId     string `json:"owner_id"`
	OwnerIdType string `json:"owner_id_type"`
	TenantKey   string `json:"tenant_key"`
}

type FSAPPConf struct {
	WideScreenMode bool `json:"wide_screen_mode"`
	EnableForward  bool `json:"enable_forward"`
}

type FSAPPTe struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type FSAPPElement struct {
	Tag           string         `json:"tag"`
	Text          Te             `json:"text"`
	Content       string         `json:"content"`
	FSAPPElements []FSAPPElement `json:"elements"`
}

type FSAPPTitles struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type FSAPPHeaders struct {
	FSAPPTitle FSAPPTitles `json:"title"`
	Template   string      `json:"template"`
}

type FSAPPCards struct {
	FSAPPConfig   FSAPPConf      `json:"config"`
	FSAPPElements []FSAPPElement `json:"elements"`
	FSAPPHeader   FSAPPHeaders   `json:"header"`
}

type FSContentAPP struct {
	MsgType      string `json:"msg_type"`
	ReceiveId    string `json:"receive_id"` //用户传入的ID，可以是 open_id、user_id、union_id、email、chat_id
	FSAPPContent string `json:"content"`
}

type Conf struct {
	WideScreenMode bool `json:"wide_screen_mode"`
	EnableForward  bool `json:"enable_forward"`
}

type Te struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type Element struct {
	Tag      string    `json:"tag"`
	Text     Te        `json:"text"`
	Content  string    `json:"content"`
	Elements []Element `json:"elements"`
}

type Titles struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type Headers struct {
	Title    Titles `json:"title"`
	Template string `json:"template"`
}

type Cards struct {
	Config   Conf      `json:"config"`
	Elements []Element `json:"elements"`
	Header   Headers   `json:"header"`
}

type FSMessagev2 struct {
	MsgType string `json:"msg_type"`
	Email   string `json:"email"` //@所使用字段
	Card    Cards  `json:"card"`
}

type TenantAccessMeg struct {
	AppId     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

type TenantAccessResp struct {
	Code              int    `json:"code"`
	Msg               string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
}
