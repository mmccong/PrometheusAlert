package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	createChatURL    = "https://open.feishu.cn/open-apis/im/v1/chats"
	inviteMembersURL = "https://open.feishu.cn/open-apis/im/v1/chats/%v/members"
	chatURL          = "https://open.feishu.cn/open-apis/im/v1/chats/%v"
)

var (
	UserA = "ou_a6e8e05ca57cd92676581d2518e3b0da"
	UserB = "ou_2f8e5575a29130a1c25d5e45fce20cf7"
	UserC = "ou_c88835a62fce5d13c1289e50171be675"
)

// CreateChatAndInviteUser creat a group with the robot as the group owner, and invite user to chat.
func CreateAlertChatAndInviteUser(ctx context.Context, token string, userIDList []string, body PrometheusAlertBody) (chatID string, err error) {
	if token == "" {
		token, err = GetTenantAccessToken(ctx)
		if err != nil {
			logs.Error("failed to get tenant access token")
			return "", err
		}
	}

	createChatReq := &CreateChatRequest{
		Name:        body.Alerts[0].Labels.Alertname,
		Description: body.Alerts[0].Annotations.Description,
		I18nNames: &I18nNames{
			ZhCn: "P" + body.Alerts[0].Labels.Level + ": 线上事故处理",
			EnUs: "P" + body.Alerts[0].Labels.Level + ": Online incident handling",
			JaJp: "P" + body.Alerts[0].Labels.Level + "：オンラインインシデント処理",
		},
	}

	createResp, err := createChatV1(ctx, token, createChatReq)
	if err != nil {
		logs.Error("failed to create chat")
		return "", err
	}

	openChatID := createResp.ChatId
	//userIDList := []string{UserA, UserB, UserC}
	inviteMembersRequest := &ChatMembersInviteRequest{
		IdList: userIDList,
	}

	inviteResp, err := chatMembersInvite(ctx, token, openChatID, inviteMembersRequest)
	if err != nil {
		logs.Error("failed to invited members to chat, chat_id: %v, user_id_list: %v", openChatID, userIDList)
		return "", err
	}
	if len(inviteResp.InvalidIDList) > 0 {
		logs.Info("invited member to chat find invalide user, invalied_ids: %v", inviteResp.InvalidIDList)
	}
	return openChatID, nil
}

// CreateChatAndInviteUser creat a group with the robot as the group owner, and invite user to chat.
func CreateChatAndInviteUser(ctx context.Context, token string) (chatID string, err error) {
	if token == "" {
		token, err = GetTenantAccessToken(ctx)
		if err != nil {
			logs.Error("failed to get tenant access token")
			return "", err
		}
	}

	createChatReq := &CreateChatRequest{
		Name:        "P0: 线上事故处理",
		Description: "线上紧急事故处理",
		I18nNames: &I18nNames{
			ZhCn: "P0: 线上事故处理",
			EnUs: "P0: Online incident handling",
			JaJp: "P0：オンラインインシデント処理",
		},
	}

	createResp, err := createChatV1(ctx, token, createChatReq)
	if err != nil {
		logs.Error("failed to create chat")
		return "", err
	}

	openChatID := createResp.ChatId
	userIDList := []string{UserA, UserB, UserC}
	inviteMembersRequest := &ChatMembersInviteRequest{
		IdList: userIDList,
	}

	inviteResp, err := chatMembersInvite(ctx, token, openChatID, inviteMembersRequest)
	if err != nil {
		logs.Error("failed to invited members to chat, chat_id: %v, user_id_list: %v", openChatID, userIDList)
		return "", err
	}
	if len(inviteResp.InvalidIDList) > 0 {
		logs.Info("invited member to chat find invalide user, invalied_ids: %v", inviteResp.InvalidIDList)
	}
	return openChatID, nil
}

func createChatV1(ctx context.Context, token string, createChatRequest *CreateChatRequest) (*CreateChatRespBody, error) {
	var err error
	if token == "" {
		token, err = GetTenantAccessToken(ctx)
		if err != nil {
			logs.Error("failed to get tenant access token")
			return nil, err
		}
	}
	cli := &http.Client{}

	reqBytes, err := json.Marshal(createChatRequest)
	if err != nil {
		logs.Error("failed to marshal")
		return nil, err
	}
	req, err := http.NewRequest("POST", createChatURL, strings.NewReader(string(reqBytes)))
	if err != nil {
		logs.Error("new request failed")
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	q := req.URL.Query()
	q.Add("user_id_type", "open_id")
	req.URL.RawQuery = q.Encode()
	var logID string
	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("create chat failed, err=%v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("read body failed")
		return nil, err
	}
	if resp != nil && resp.Header != nil {
		logID = resp.Header.Get("x-tt-logid")
	}

	createChatResp := &CreateChatResponse{}
	err = json.Unmarshal(body, createChatResp)
	if err != nil {
		logs.Error("failed to unmarshal")
		return nil, err
	}
	if createChatResp.Code != 0 {
		logs.Warn("failed to create chat, code: %v, msg: %v, log_id: %v", createChatResp.Code, createChatResp.Message, logID)
		return nil, fmt.Errorf("create chat failed")
	}
	logs.Info("succeed create chat, chat_id: %v", createChatResp.Data.ChatId)
	return createChatResp.Data, nil
}

func chatMembersInvite(ctx context.Context, token string, chatID string, inviteRequest *ChatMembersInviteRequest) (*ChatMembersInviteRespBody, error) {
	var err error
	if token == "" {
		token, err = GetTenantAccessToken(ctx)
		if err != nil {
			logs.Error("failed to get tenant access token")
			return nil, err
		}
	}
	cli := &http.Client{}

	reqBytes, err := json.Marshal(inviteRequest)
	if err != nil {
		logs.Error("failed to marshal")
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf(inviteMembersURL, chatID), strings.NewReader(string(reqBytes)))
	if err != nil {
		logs.Error("get request failed")
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	q := req.URL.Query()
	q.Add("member_id_type", "open_id")
	req.URL.RawQuery = q.Encode()

	logID := ""
	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("invite members to chat failed, err=%v", err)
	}
	if resp != nil && resp.Header != nil {
		logID = resp.Header.Get("x-tt-logid")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("read body failed")
		return nil, err
	}

	inviteMemberResp := &ChatMembersInviteResponse{}
	err = json.Unmarshal(body, inviteMemberResp)
	if err != nil {
		logs.Error("failed to unmarshal")
		return nil, err
	}
	if inviteMemberResp.Code != 0 {
		logs.Warn("invite chatter failed, code: %v, msg: %v, log_id: %v", inviteMemberResp.Code, inviteMemberResp.Message, logID)
		return nil, fmt.Errorf("invite chatter failed")
	}
	logs.Info("succeed invited members to chat, resp: %v, log_id: %v", inviteMemberResp, logID)

	return inviteMemberResp.Data, nil
}

func UpdateChat(ctx context.Context, token, chatID string, updateChatReq *UpdateChatRequest) (*UpdateChatResponse, error) {
	var err error
	if token == "" {
		token, err = GetTenantAccessToken(ctx)
		if err != nil {
			logs.Error("failed to get tenant access token")
			return nil, err
		}
	}
	cli := &http.Client{}

	reqBytes, err := json.Marshal(updateChatReq)
	if err != nil {
		logs.Error("failed to marshal")
		return nil, err
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf(chatURL, chatID), strings.NewReader(string(reqBytes)))
	if err != nil {
		logs.Error("get request failed")
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	q := req.URL.Query()
	q.Add("user_id_type", "open_id")
	req.URL.RawQuery = q.Encode()

	logID := ""
	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("update chat failed, err=%v", err)
	}
	logID = resp.Header.Get("x-tt-logid")
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("read body failed")
		return nil, err
	}
	logs.Info("body: %v", string(body))
	updateChatResp := &UpdateChatResponse{}
	err = json.Unmarshal(body, updateChatResp)
	if err != nil {
		logs.Error("failed to unmarshal")
		return nil, err
	}

	if updateChatResp.Code != 0 {
		logs.Warn("failed to create chat, code: %v, msg: %v, log_id: %v", updateChatResp.Code, updateChatResp.Message, logID)
		return nil, fmt.Errorf("update chat failed")
	}

	logs.Info("succeed update chat, log_id: %v", logID)
	return updateChatResp, nil
}

func GetChatInfo(ctx context.Context, token, chatID string) (*GetChatInfoResponseBody, error) {
	var err error
	if token == "" {
		token, err = GetTenantAccessToken(ctx)
		if err != nil {
			logs.Error("failed to get tenant access token")
			return nil, err
		}
	}
	cli := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf(chatURL, chatID), nil)
	if err != nil {
		logs.Error("get request failed")
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	q := req.URL.Query()
	q.Add("user_id_type", "open_id")
	req.URL.RawQuery = q.Encode()

	logID := ""
	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("update chat failed, err=%v", err)
	}
	logID = resp.Header.Get("x-tt-logid")
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("read body failed")
		return nil, err
	}
	logs.Info("body: %v", string(body))
	getChatResp := &GetChatInfoResponse{}
	err = json.Unmarshal(body, getChatResp)
	if err != nil {
		logs.Error("failed to unmarshal")
		return nil, err
	}

	if getChatResp.Code != 0 {
		logs.Warn("failed to create chat, code: %v, msg: %v, log_id: %v", getChatResp.Code, getChatResp.Message, logID)
		return nil, fmt.Errorf("update chat failed")
	}

	logs.Info("succeed update chat, log_id: %v", logID)
	return getChatResp.Data, nil
}
