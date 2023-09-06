package biz

import (
	"PrometheusAlert/models"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	createMessageURL     = "https://open.feishu.cn/open-apis/im/v1/messages"
	uploadImageURL       = "https://open.feishu.cn/open-apis/im/v1/images"
	getMessageHistoryURL = "https://open.feishu.cn/open-apis/im/v1/messages"
)

func SendFeiShuAlertMessage(ctx context.Context, message models.AlertContent, token, status string, chatID string) error {
	var err error
	if token == "" {
		token, err = GetTenantAccessToken(ctx)
		if err != nil {
			logs.Error("failed to get tenant access token")
			return err
		}
	}

	var createResp *MessageItem
	var createReq *CreateMessageRequest
	switch status {
	case "firing1":
		content := "{\"text\":\"<at user_id=\\\"all\\\">所有人</at> 请注意，线上服务发生报警，请及时处理。 \\n服务负责人：<at user_id=\\\"ou_ba44c2d64d161c0f12d8548bef215311\\\">张三</at> \"}"
		createReq = genCreateMessageRequest(ctx, chatID, content, "text")
	case "warming":
		content := "{\"zh_cn\":{\"title\":\"线上服务报警通知！\",\"content\":[[{\"tag\":\"at\",\"user_id\":\"all\",\"user_name\":\"所有人\"},{\"tag\":\"text\",\"text\":\"请注意，线上服务发生报警，请及时处理。\"}],[{\"tag\":\"text\",\"text\":\"服务负责人：\"},{\"tag\":\"at\",\"user_id\":\"ou_ba44c2d64d161c0f12d8548bef215311\",\"user_name\":\"张三\"}]]}}"
		createReq = genCreateMessageRequest(ctx, chatID, content, "post")
	case "firing":
		image, err := UploadImage(token)
		if err != nil {
			logs.Error("failed to upload image")
			return err
		}
		cardContent := ConstructAlterCard(message, image.ImageKey)
		createReq = genCreateMessageRequest(ctx, chatID, cardContent, "interactive")
	case "resolved":
		image, err := UploadImage(token)
		if err != nil {
			logs.Error("failed to upload image")
			return err
		}
		cardContent := ConstructResolvedCard(message, image.ImageKey)
		createReq = genCreateMessageRequest(ctx, chatID, cardContent, "interactive")
	case "silence":
		image, err := UploadImage(token)
		if err != nil {
			logs.Error("failed to upload image")
			return err
		}
		cardContent := ConstructSilenceCard(message, image.ImageKey)
		createReq = genCreateMessageRequest(ctx, chatID, cardContent, "interactive")
	case "mistake":
		image, err := UploadImage(token)
		if err != nil {
			logs.Error("failed to upload image")
			return err
		}
		cardContent := ConstructSilenceCard(message, image.ImageKey)
		createReq = genCreateMessageRequest(ctx, chatID, cardContent, "interactive")
	case "ssl":
		msg := models.SSLContent{
			Title:   message.Title,
			Content: message.Content,
		}
		cardContent := ConstructSSLCertCard(msg)
		createReq = genCreateMessageRequest(ctx, chatID, cardContent, "interactive")
	case "markdown":
		msg := models.MarkdownContent{
			Title: message.Title,
			Text:  message.Content,
		}
		cardContent := ConstructMarkdownCard(msg)
		createReq = genCreateMessageRequest(ctx, chatID, cardContent, "interactive")
	default:
		return nil
	}

	createResp, err = SendMessage(ctx, token, createReq)
	if err != nil {
		logs.Error("send %v message failed, chat_id: %v", status, chatID)
		return err
	}

	msgID := createResp.MessageID
	logs.Info("succeed send alert message, msg_id: %v", msgID)
	return nil
}

func SendAlertMessage(ctx context.Context, token, msgType string, chatID string) error {
	var err error
	if token == "" {
		token, err = GetTenantAccessToken(ctx)
		if err != nil {
			logs.Error("failed to get tenant access token")
			return err
		}
	}

	var createResp *MessageItem
	var createReq *CreateMessageRequest
	var message models.AlertContent
	switch msgType {
	case "text":
		content := "{\"text\":\"<at user_id=\\\"all\\\">所有人</at> 请注意，线上服务发生报警，请及时处理。 \\n服务负责人：<at user_id=\\\"ou_ba44c2d64d161c0f12d8548bef215311\\\">张三</at> \"}"
		createReq = genCreateMessageRequest(ctx, chatID, content, msgType)
	case "post":
		content := "{\"zh_cn\":{\"title\":\"线上服务报警通知！\",\"content\":[[{\"tag\":\"at\",\"user_id\":\"all\",\"user_name\":\"所有人\"},{\"tag\":\"text\",\"text\":\"请注意，线上服务发生报警，请及时处理。\"}],[{\"tag\":\"text\",\"text\":\"服务负责人：\"},{\"tag\":\"at\",\"user_id\":\"ou_ba44c2d64d161c0f12d8548bef215311\",\"user_name\":\"张三\"}]]}}"
		//content := "{\"zh_cn\":[{\"tag\":\"column_set\",\"flex_mode\":\"none\",\"background_style\":\"default\",\"columns\":[{\"tag\":\"column\",\"width\":\"weighted\",\"weight\":1,\"vertical_align\":\"top\",\"elements\":[{\"tag\":\"div\",\"text\":{\"content\":\"**🔴 报警服务：**\\nQA 7\",\"tag\":\"lark_md\"}}]},{\"tag\":\"column\",\"width\":\"weighted\",\"weight\":1,\"vertical_align\":\"top\",\"elements\":[{\"tag\":\"div\",\"text\":{\"content\":\"**🕐 时间：**\\n2023-02-23 20:17:51\",\"tag\":\"lark_md\"}}]}]},{\"tag\":\"column_set\",\"flex_mode\":\"none\",\"background_style\":\"default\",\"columns\":[{\"tag\":\"column\",\"width\":\"weighted\",\"weight\":1,\"vertical_align\":\"top\",\"elements\":[{\"tag\":\"div\",\"text\":{\"content\":\"**👤 一级值班：**\\n[@王冰](https://open.feishu.cn/document/ugTN1YjL4UTN24CO1UjN/uUzN1YjL1cTN24SN3UjN?from=mcb)\",\"tag\":\"lark_md\"}}]},{\"tag\":\"column\",\"width\":\"weighted\",\"weight\":1,\"vertical_align\":\"top\",\"elements\":[{\"tag\":\"markdown\",\"content\":\"**👤 二级值班：**\\n[@李天天](https://open.feishu.cn/document/ugTN1YjL4UTN24CO1UjN/uUzN1YjL1cTN24SN3UjN?from=mcb)\"}]}]},{\"tag\":\"div\",\"text\":{\"content\":\"支付方式 支付成功率低于50%\",\"tag\":\"plain_text\"}},{\"alt\":{\"content\":\"\",\"tag\":\"plain_text\"},\"img_key\":\"img_v2_8b2ebeaf-c97c-411d-a4dc-4323e8cba10g\",\"tag\":\"img\"},{\"elements\":[{\"content\":\"🔴 支付失败数  🔵 支付成功数\",\"tag\":\"plain_text\"}],\"tag\":\"note\"},{\"actions\":[{\"tag\":\"button\",\"text\":{\"content\":\"跟进处理\",\"tag\":\"plain_text\"},\"type\":\"primary\",\"value\":{\"sloved\":\"user\"}},{\"options\":[{\"text\":{\"content\":\"屏蔽10分钟\",\"tag\":\"plain_text\"},\"value\":\"1\"},{\"text\":{\"content\":\"屏蔽30分钟\",\"tag\":\"plain_text\"},\"value\":\"2\"},{\"text\":{\"content\":\"屏蔽1小时\",\"tag\":\"plain_text\"},\"value\":\"3\"},{\"text\":{\"content\":\"屏蔽24小时\",\"tag\":\"plain_text\"},\"value\":\"4\"}],\"placeholder\":{\"content\":\"暂时屏蔽报警\",\"tag\":\"plain_text\"},\"tag\":\"select_static\",\"value\":{\"key\":\"value\"}}],\"tag\":\"action\"},{\"tag\":\"hr\"},{\"tag\":\"div\",\"text\":{\"content\":\"🙋🏼 [我要反馈误报](https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/im-v1/message-development-tutorial/introduction?from=mcb) | 📝 [录入报警处理过程](https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/im-v1/message-development-tutorial/introduction?from=mcb)\",\"tag\":\"lark_md\"}}]}"
		createReq = genCreateMessageRequest(ctx, chatID, content, msgType)
	case "interactive":
		image, err := UploadImage(token)
		if err != nil {
			logs.Error("failed to upload image")
			return err
		}
		cardContent := ConstructAlterCard(message, image.ImageKey)
		createReq = genCreateMessageRequest(ctx, chatID, cardContent, msgType)
	default:
		return nil
	}

	createResp, err = SendMessage(ctx, token, createReq)
	if err != nil {
		logs.Error("send %v message failed, chat_id: %v", msgType, chatID)
		return err
	}

	msgID := createResp.MessageID
	logs.Info("succeed send alert message, msg_id: %v", msgID)
	return nil
}

func SendMessage(ctx context.Context, token string, createReq *CreateMessageRequest) (*MessageItem, error) {
	var err error
	if token == "" {
		token, err = GetTenantAccessToken(ctx)
		if err != nil {
			logs.Error("failed to get tenant access token")
			return nil, err
		}
	}
	cli := &http.Client{}

	reqBytes, err := json.Marshal(createReq)
	if err != nil {
		logs.Error("failed to marshal")
		return nil, err
	}
	req, err := http.NewRequest("POST", createMessageURL, strings.NewReader(string(reqBytes)))
	if err != nil {
		logs.Error("new request failed")
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	q := req.URL.Query()
	q.Add("receive_id_type", "chat_id")
	req.URL.RawQuery = q.Encode()

	var logID string
	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("create message failed, err=%v", err)
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

	createMessageResp := &CreateMessageResponse{}
	err = json.Unmarshal(body, createMessageResp)
	if err != nil {
		logs.Error("failed to unmarshal")
		return nil, err
	}
	if createMessageResp.Code != 0 {
		logs.Warn("failed to create message, code: %v, msg: %v, log_id: %v", createMessageResp.Code, createMessageResp.Message, logID)
		return nil, fmt.Errorf("create message failed")
	}
	logs.Info("succeed create message, msg_id: %v", createMessageResp.Data.MessageID)
	return createMessageResp.Data, nil
}

func genCreateMessageRequest(ctx context.Context, chatID, content, msgType string) *CreateMessageRequest {
	return &CreateMessageRequest{
		ReceiveID: chatID,
		Content:   content,
		MsgType:   msgType,
	}
}

func GetChatAllMessageAndReview(ctx context.Context, token, chatID string) error {
	var err error
	if token == "" {
		token, err = GetTenantAccessToken(ctx)
		if err != nil {
			logs.Error("failed to get tenant access token")
			return err
		}
	}
	start := "0"
	end := fmt.Sprintf("%v", time.Now().Unix())

	pwd, _ := os.Getwd()
	parent := filepath.Dir(pwd)
	path := parent + fmt.Sprintf("/resource/download/chat_%v_history.txt", chatID)

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logs.Error("open file failed")
		return err
	}
	defer file.Close()

	write := bufio.NewWriter(file)
	write.WriteString(fmt.Sprintf("chat(%v) history\n", chatID))

	hasMore := true
	pageToken := ""
	for {
		if !hasMore {
			break
		}
		getMsgResp, err := GetChatMessageHistory(ctx, token, chatID, start, end, pageToken, "10")
		if err != nil {
			logs.Error("failed to get chat message")
			break
		}

		if len(getMsgResp.Items) > 0 {
			for _, item := range getMsgResp.Items {
				senderID := item.Sender.ID
				createTime := item.CreateTime
				intCreateTime, err := strconv.ParseInt(createTime, 10, 64)
				if err != nil {
					continue
				}

				createTime = fmt.Sprintf("%v", time.Unix(intCreateTime/1000, 0))
				content := item.Body.Content
				str := fmt.Sprintf("chatter(%v) at (%v) send: %v", senderID, createTime, content)
				write.WriteString(str + "\n")
			}
			write.Flush()
		}
		pageToken = getMsgResp.PageToken
		hasMore = getMsgResp.HasMore
		time.Sleep(50 * time.Millisecond)
	}

	return nil
}

func GetChatMessageHistory(ctx context.Context, token, chatID string, start, end, pageToken, pageSize string) (*GetMessageHistoryBody, error) {
	var err error
	if token == "" {
		token, err = GetTenantAccessToken(ctx)
		if err != nil {
			logs.Error("failed to get tenant access token")
			return nil, err
		}
	}
	cli := &http.Client{}

	req, err := http.NewRequest("GET", getMessageHistoryURL, nil)
	if err != nil {
		logs.Error("new request failed")
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	q := req.URL.Query()
	q.Add("container_id_type", "chat")
	q.Add("container_id", chatID)
	q.Add("start_time", start)
	q.Add("end_time", end)
	q.Add("page_token", pageToken)
	q.Add("page_size", pageSize)
	req.URL.RawQuery = q.Encode()

	var logID string
	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get message failed, err=%v", err)
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

	getMessageResp := &GetMessageHistoryResponse{}
	err = json.Unmarshal(body, getMessageResp)
	if err != nil {
		logs.Error("failed to unmarshal")
		return nil, err
	}
	if getMessageResp.Code != 0 {
		logs.Warn("failed to get message, code: %v, msg: %v, log_id: %v", getMessageResp.Code, getMessageResp.Message, logID)
		return nil, fmt.Errorf("get message hitory failed")
	}

	return getMessageResp.Data, nil
}
