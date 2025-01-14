package biz

import (
	"PrometheusAlert/models"
	"context"
	"github.com/astaxie/beego/logs"
	//"logs"
	"testing"
	"time"
)

func TestSendAlertMessage(t *testing.T) {
	ctx := context.Background()
	token, err := GetTenantAccessToken(ctx)
	if err != nil {
		logs.Error("failed to get tenant access token")
		t.Fail()
		return
	}

	chatID, err := CreateChatAndInviteUser(ctx, token)
	if err != nil {
		logs.Error("failed to create chat")
		t.Fail()
		return
	}
	// Ensure successful information synchronization of chat
	time.Sleep(3 * time.Second)

	msgTpyes := []string{"text", "post", "interactive"}
	for _, msgType := range msgTpyes {
		err := SendAlertMessage(ctx, token, msgType, chatID)
		if err != nil {
			logs.Error("send %v message failed", msgType)
			continue
		}
	}
	logs.Info("succeed create chat and send msg")
}

func TestUploadImage(t *testing.T) {
	ctx := context.Background()
	token, err := GetTenantAccessToken(ctx)
	if err != nil {
		logs.Error("failed to get tenant access token")
		t.Fail()
		return
	}
	uploadImage, err := UploadImage(token)
	if err != nil {
		logs.Error("failed to upload image")
		t.Fail()
		return
	}
	logs.Info("succeed upload image, image_key: %v", uploadImage.ImageKey)
	var message models.AlertContent
	card := ConstructAlterCard(message, uploadImage.ImageKey)
	logs.Info("card: %v", card)

}

// todo remove
func TestSendPostMessage(t *testing.T) {
	ctx := context.Background()
	token, err := GetTenantAccessToken(ctx)
	if err != nil {
		logs.Error("failed to get tenant access token")
		t.Fail()
		return
	}
	createReq := &CreateMessageRequest{
		ReceiveID: "oc_95ac7aa44555d1e947f6cb8203dbebf4",
		MsgType:   "post",
		Content:   "{\"zh_cn\":{\"title\":\"我是一个标题\",\"content\":[[{\"tag\":\"text\",\"text\":\"第一行 :\"},{\"tag\":\"a\",\"href\":\"http://www.feishu.cn\",\"text\":\"超链接\"},{\"tag\":\"at\",\"user_id\":\"ou_1avnmsbv3k45jnk34j5\",\"user_name\":\"tom\"}],[{\"tag\":\"img\",\"image_key\":\"img_v2_0cc066b0-d406-4f96-af9e-bb551984729j\",\"width\":600,\"height\":300}],[{\"tag\":\"text\",\"text\":\"第二行:\"},{\"tag\":\"text\",\"text\":\"文本测试\"}],[{\"tag\":\"img\",\"image_key\":\"img_v2_932dd3b1-112e-4e0b-bbb4-e8ff505fbb8j\",\"width\":300,\"height\":200}]]}}",
	}

	SendMessage(ctx, token, createReq)
}

func TestGetChatAllMessageAndReview(t *testing.T) {
	ctx := context.Background()
	token, err := GetTenantAccessToken(ctx)
	if err != nil {
		logs.Error("failed to get tenant access token")
		t.Fail()
		return
	}
	chatID := "oc_8691a127fda5570eacc05628e90ca04a"
	err = GetChatAllMessageAndReview(ctx, token, chatID)
	if err != nil {
		logs.Error("failed to get chat all message ")
	}
}
