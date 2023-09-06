package biz

import (
	"PrometheusAlert/models"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type DashboardJson struct {
	Telegram        int `json:"telegram"`
	Smoordx         int `json:"smoordx"`
	Smoordh         int `json:"smoordh"`
	Alydx           int `json:"alydx"`
	Alydh           int `json:"alydh"`
	Bdydx           int `json:"bdydx"`
	Bark            int `json:"bark"`
	Dingding        int `json:"dingding"`
	Email           int `json:"email"`
	Feishu          int `json:"feishu"`
	Hwdx            int `json:"hwdx"`
	Rlydx           int `json:"rlydx"`
	Ruliu           int `json:"ruliu"`
	Txdx            int `json:"txdx"`
	Txdh            int `json:"txdh"`
	Webhook         int `json:"webhook"`
	Weixin          int `json:"weixin"`
	Workwechat      int `json:"workwechat"`
	Voice           int `json:"voice"`
	Zabbix          int `json:"zabbix"`
	Grafana         int `json:"grafana"`
	Graylog         int `json:"graylog"`
	Prometheus      int `json:"prometheus"`
	Prometheusalert int `json:"prometheusalert"`
	Aliyun          int `json:"prometheusalert"`
}

var ChartsJson DashboardJson

func InitFeiShuAPP(logsign string) string {
	open := beego.AppConfig.String("open-feishuapp")
	if open != "1" {
		logs.Info(logsign, "[feishuapp]", "飞书APP接口未配置未开启状态,请先配置open-feishuapp为1")
		return "飞书APP接口未配置未开启状态,请先配置open-feishuapp为1"
	}
	ctx := context.Background()
	token, err := GetTenantAccessToken(ctx)
	if err != nil {
		logs.Error("failed to get tenant access token")
		return "failed to get tenant access token"
	}
	return token
}

func GetUserIds(users string) string {
	// 创建 Client
	// 如需SDK自动管理租户Token的获取与刷新，可调用lark.WithEnableTokenCache(true)进行设置
	ctx := context.Background()
	token, err := GetTenantAccessToken(ctx)
	if err != nil {
		logs.Error("failed to get tenant access token")
		return "failed to get tenant access token"
	}

	appId := beego.AppConfig.String("FEISHU_APPID")
	appSecret := beego.AppConfig.String("FEISHU_APPSECRET")
	client := lark.NewClient(appId, appSecret, lark.WithEnableTokenCache(false))

	// 创建请求对象
	req := larkcontact.NewBatchGetIdUserReqBuilder().
		UserIdType(`user_id`).
		Body(larkcontact.NewBatchGetIdUserReqBodyBuilder().
			//Emails(users).
			Emails(strings.Split(users, ",")).
			//Mobiles([]string{`13812345678`, `13812345679`}).
			Build()).
		Build()

	// 发起请求
	// 如开启了SDK的Token管理功能，就无需在请求时调用larkcore.WithTenantAccessToken("-xxx")来手动设置租户Token了
	resp, err := client.Contact.User.BatchGetId(context.Background(), req, larkcore.WithTenantAccessToken(token))

	// 处理错误
	if err != nil {
		return ""
	}

	// 服务端错误处理
	if !resp.Success() {
		return ""
	}

	// 业务处理
	resBy, err := json.Marshal(resp.Data.UserList)
	if err != nil {
		return "err"
	}

	var newData []models.LarkUser
	jsonRes := json.Unmarshal(resBy, &newData)
	if jsonRes != nil {
		return "err"
	}

	var UidSlice []string
	for _, value := range newData {
		UidSlice = append(UidSlice, value.UserId)
	}
	result := strings.Join(UidSlice, ",")
	return result
}

func GetOpenIds(users string) string {
	// 创建 Client
	ctx := context.Background()
	token, err := GetTenantAccessToken(ctx)
	if err != nil {
		logs.Error("failed to get tenant access token")
		return "failed to get tenant access token"
	}
	// 如需SDK自动管理租户Token的获取与刷新，可调用lark.WithEnableTokenCache(true)进行设置
	appId := beego.AppConfig.String("FEISHU_APPID")
	appSecret := beego.AppConfig.String("FEISHU_APPSECRET")
	client := lark.NewClient(appId, appSecret, lark.WithEnableTokenCache(false))

	// 创建请求对象
	req := larkcontact.NewBatchGetIdUserReqBuilder().
		Body(larkcontact.NewBatchGetIdUserReqBodyBuilder().
			//Emails(users).
			Emails(strings.Split(users, ",")).
			//Mobiles([]string{`13812345678`, `13812345679`}).
			Build()).
		Build()

	// 发起请求
	// 如开启了SDK的Token管理功能，就无需在请求时调用larkcore.WithTenantAccessToken("-xxx")来手动设置租户Token了
	resp, err := client.Contact.User.BatchGetId(context.Background(), req, larkcore.WithTenantAccessToken(token))

	// 处理错误
	if err != nil {
		return ""
	}

	// 服务端错误处理
	if !resp.Success() {
		return ""
	}

	// 业务处理
	resBy, err := json.Marshal(resp.Data.UserList)
	if err != nil {
		return "err"
	}

	var newData []models.LarkUser
	jsonRes := json.Unmarshal(resBy, &newData)
	if jsonRes != nil {
		return "err"
	}

	var UidSlice []string
	for _, value := range newData {
		UidSlice = append(UidSlice, value.UserId)
	}
	result := strings.Join(UidSlice, ",")
	return result
}

func GetOneChatId(name string) string {
	// 创建 Client
	ctx := context.Background()
	token, err := GetTenantAccessToken(ctx)
	if err != nil {
		logs.Error("failed to get tenant access token")
		return "failed to get tenant access token"
	}
	// 如需SDK自动管理租户Token的获取与刷新，可调用lark.WithEnableTokenCache(true)进行设置
	appId := beego.AppConfig.String("FEISHU_APPID")
	appSecret := beego.AppConfig.String("FEISHU_APPSECRET")
	client := lark.NewClient(appId, appSecret, lark.WithEnableTokenCache(false))

	// 创建请求对象
	req := larkim.NewListChatReqBuilder().
		PageSize(20).
		Build()

	// 发起请求
	// 如开启了SDK的Token管理功能，就无需在请求时调用larkcore.WithTenantAccessToken("-xxx")来手动设置租户Token了
	resp, err := client.Im.Chat.List(context.Background(), req, larkcore.WithTenantAccessToken(token))

	// 处理错误
	if err != nil {
		return fmt.Sprintln("err: v%", err)
	}

	// 服务端错误处理
	if !resp.Success() {
		logs.Error("GetChatId: %s", err)
		return "err"
	}

	// 业务处理
	resBy, err := json.Marshal(resp.Data.Items)
	if err != nil {
		logs.Error("GetChatId: %s", err)
		return "err"
	}

	var newData []models.LarkItem
	jsonRes := json.Unmarshal(resBy, &newData)
	if jsonRes != nil {
		logs.Error("GetChatId: %s", jsonRes)
		return "err"
	}
	var ChatId string
	for _, value := range newData {
		if name == value.Name {
			ChatId = value.ChatId
		}
	}
	return ChatId
}

func PostFeiShuAppToChat(title, text, chatIds, userIds, logsign string) string {
	ctx := context.Background()
	token, err := GetTenantAccessToken(ctx)
	if err != nil {
		logs.Error("failed to get tenant access token")
		return "failed to get tenant access token"
	}

	open := beego.AppConfig.String("open-feishuapp")
	if open != "1" {
		logs.Info(logsign, "[feishuapp]", "飞书APP接口未配置未开启状态,请先配置open-feishuapp为1")
		return "飞书APP接口未配置未开启状态,请先配置open-feishuapp为1"
	}
	var color string
	if strings.Count(text, "\"status\":\"firing\"") > 0 && strings.Count(text, "firing") > 0 {
		color = "orange"
	} else if strings.Count(text, "resolved") > 0 {
		color = "green"
	} else {
		color = "red"
	}

	var ids string
	var atsStr string
	var SendContent string
	if len(userIds) > 0 {
		if strings.Count(userIds, "@") > 0 {
			ids = GetUserIds(userIds)
		}
		for _, value := range strings.Split(ids, ",") {
			at := fmt.Sprintf("<at id=%v></at> ", value)
			atsStr += at
		}
		SendContent = text + " " + atsStr
	} else {
		SendContent = text
	}

	var result []byte
	if chatIds != "" {
		ReceiveIds := strings.Split(chatIds, ",")
		fsAppContent :=
			&models.FSAPPCards{
				FSAPPConfig: models.FSAPPConf{
					WideScreenMode: true,
					EnableForward:  true,
				},
				FSAPPHeader: models.FSAPPHeaders{
					FSAPPTitle: models.FSAPPTitles{
						Content: title,
						Tag:     "plain_text",
					},
					Template: color,
				},
				FSAPPElements: []models.FSAPPElement{
					models.FSAPPElement{
						Tag: "div",
						Text: models.Te{
							Content: SendContent,
							Tag:     "lark_md",
						},
					},
					{
						Tag: "hr",
					},
					{
						Tag: "note",
						FSAPPElements: []models.FSAPPElement{
							{
								Content: title,
								Tag:     "lark_md",
							},
						},
					},
				},
			}
		contentByte, _ := json.Marshal(fsAppContent)
		for _, ReceiveId := range ReceiveIds {
			u := models.FSContentAPP{
				MsgType:      "interactive",
				ReceiveId:    ReceiveId,
				FSAPPContent: string(contentByte),
			}
			var ReceiveType string
			if strings.Contains(ReceiveId, "ou_") {
				ReceiveType = "open_id"
			} else if strings.Contains(ReceiveId, "on_") {
				ReceiveType = "union_id"
			} else if strings.Contains(ReceiveId, "oc_") {
				ReceiveType = "chat_id"
			} else if strings.Contains(ReceiveId, "@") {
				ReceiveType = "email"
			} else {
				ReceiveType = "user_id"
			}
			b := new(bytes.Buffer)
			json.NewEncoder(b).Encode(u)
			logs.Info(logsign, "[feishuapp]", b)
			var tr *http.Transport
			if proxyUrl := beego.AppConfig.String("proxy"); proxyUrl != "" {
				proxy := func(_ *http.Request) (*url.URL, error) {
					return url.Parse(proxyUrl)
				}
				tr = &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
					Proxy:           proxy,
				}
			} else {
				tr = &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				}
			}
			client := &http.Client{Transport: tr}
			FSUrl := fmt.Sprintf("https://open.feishu.cn/open-apis/im/v1/messages?receive_id_type=%s", ReceiveType)
			req, err := http.NewRequest("POST", FSUrl, b)
			if err != nil {
				logs.Error(logsign, "[feishuapp]", title+": "+err.Error())
			}
			req.Header.Set("Content-Type", "application/json; charset=utf-8")
			req.Header.Set("Authorization", "Bearer "+token)
			resp, err := client.Do(req)
			if err != nil {
				logs.Error(logsign, "[feishuapp]", err.Error())
			}
			defer resp.Body.Close()
			result, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				logs.Error(logsign, "[feishuapp]", title+": "+err.Error())
			}
			models.AlertToCounter.WithLabelValues("feishuapp").Add(1)
			ChartsJson.Feishu += 1
			logs.Info(logsign, "[feishuapp]", title+": "+string(result))
			//return string(result)
		}
	}
	return string(result)
}

func SendFeishuAppToChatWithTemplate(logsign string) string {
	ctx := context.Background()
	token, err := GetTenantAccessToken(ctx)
	if err != nil {
		logs.Error("failed to get tenant access token")
		return "failed to get tenant access token"
	}

	url := "https://open.feishu.cn/open-apis/im/v1/messages?receive_id_type=chat_id"
	method := "POST"
	payload := strings.NewReader(`{
	"receive_id": "oc_0d375e3e1d3db61cd523a70808bfc87d",
	"msg_type": "interactive",
	"content": "{\"type\": \"template\",\"data\": { \"template_id\": \"ctp_AAqk3LgxmZ1D\"} }",
	"msg_type": "interactive"
}`)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fmt.Println(string(body))
	return ""
}
