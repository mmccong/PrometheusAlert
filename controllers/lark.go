package controllers

import (
	"PrometheusAlert/biz"
	"PrometheusAlert/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"net/http"
	"strconv"
)

type LarkController struct {
	beego.Controller
}

type Message struct {
	ChatID  string `json:"chat_id"`
	Text    string `json:"text"`
	Title   string `json:"title"`
	UserIds string `json:"user_ids"`
}

type Event struct {
	OpenId        string `json:"open_id"`
	UserId        string `json:"user_id"`
	OpenMessageId string `json:"open_message_id"`
	TenantKey     string `json:"tenant_key"`
	Token         string `json:"token"`
	Action        Action `json:"action"`
}

type Action struct {
	Tag    string `json:"tag"`
	Option string `json:"option"`
	Value  Value  `json:"value"`
}
type Value struct {
	Value Key `json:"value"`
}
type Key struct {
	Key string `json:"key"`
}

type Challenge struct {
	Token     string `json:"token"`
	Type      string `json:"type"`
	Challenge string `json:"challenge"`
}
type Challenge1 struct {
	Challenge string `json:"challenge"`
}

type Card struct {
	Title   string   `json:"title"`
	Content []Action `json:"content"`
}

func (c *LarkController) GetIds() {
	emails := c.GetString("emails")
	logsign := "[" + LogsSign() + "]"
	logs.Info(logsign, string(c.Ctx.Input.RequestBody))
	user := models.LarkUserList{}
	logs.Info(logsign, string(c.Ctx.Input.RequestBody))
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &user)
	if err != nil {
		c.Abort("400")
	}
	c.Data["json"] = biz.GetUserIds(emails)
	logs.Info(logsign, c.Data["json"])
	c.ServeJSON()
}

func (c *LarkController) GetChatId() {
	ChatName := c.GetString("name")
	logsign := "[" + LogsSign() + "]"
	logs.Info(logsign, string(c.Ctx.Input.RequestBody))
	user := models.LarkUserList{}
	logs.Info(logsign, string(c.Ctx.Input.RequestBody))
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &user)
	if err != nil {
		c.Abort("400")
	}
	c.Data["json"] = biz.GetOneChatId(ChatName)
	logs.Info(logsign, c.Data["json"])
	c.ServeJSON()
}

func (c *LarkController) SendMessage() {
	var message Message
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &message)
	if err != nil {
		c.Abort("400")
	}
	logsign := "[" + LogsSign() + "]"
	//logs.Info(logsign, string(c.Ctx.Input.RequestBody))
	chatIds := beego.AppConfig.String("CHAT_ID")
	c.Data["json"] = biz.PostFeiShuAppToChat(message.Title, message.Text, chatIds, message.UserIds, logsign)
	logs.Info(logsign, c.Data["json"])
	c.ServeJSON()
}

func (c *LarkController) SendFeiShuCardMessage() {
	Chat := c.Input().Get("chat")
	var message models.AlertContent
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &message)
	if err != nil {
		c.Abort("400")
	}
	logSign := "[" + LogsSign() + "]"
	//c.Data["json"] = SendFeishuAppToChatWithTemplate(logSign)
	ctx := context.Background()
	token := biz.InitFeiShuAPP(logSign)
	chatId := biz.GetOneChatId(Chat)
	c.Data["json"] = biz.SendFeiShuAlertMessage(ctx, message, token, message.Status, chatId)
	logs.Info(logSign, c.Data["json"])
	c.ServeJSON()
}

func (c *LarkController) CreateChatAndInviteUser() {
	ctx := context.Background()
	var message Message
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &message)
	if err != nil {
		c.Abort("400")
	}
	logsign := "[" + LogsSign() + "]"
	//token, err := GetAccessToken(logsign)
	if err != nil {
		c.Abort("400")
	}
	c.Data["json"], err = biz.CreateChatAndInviteUser(ctx, "")
	if err != nil {
		c.Abort("400")
	}
	logs.Info(logsign, c.Data["json"])
	c.ServeJSON()
}

func (c *LarkController) GetOpenIds() {
	emails := c.GetString("emails")
	logsign := "[" + LogsSign() + "]"
	logs.Info(logsign, string(c.Ctx.Input.RequestBody))
	user := models.LarkUserList{}
	logs.Info(logsign, string(c.Ctx.Input.RequestBody))
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &user)
	if err != nil {
		c.Abort("400")
	}
	c.Data["json"] = biz.GetOpenIds(emails)
	logs.Info(logsign, c.Data["json"])
	c.ServeJSON()
}

func (c *LarkController) BotHook() {
	logsign := "[" + LogsSign() + "]"
	logs.Info(logsign, string(c.Ctx.Input.RequestBody))
	var challenge Challenge
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &challenge)
	fmt.Println("#################1")
	fmt.Println(challenge.Type)
	if err != nil {
		fmt.Println("#################2")
		c.Abort("400")
	}

	if challenge.Type == "url_verification" {
		fmt.Println("#################3")
		c.Data["json"] = Challenge1{Challenge: challenge.Challenge}
	} else if challenge.Type == "" {
		var event Event
		fmt.Println(c.Ctx.Input.RequestBody)
		err := json.Unmarshal(c.Ctx.Input.RequestBody, &event)
		if err != nil {
			logs.Error(logsign, "Error decoding event JSON", http.StatusBadRequest)
			return
		}
		//
		ActionData := event.Action
		num, _ := strconv.Atoi(ActionData.Option)
		// Process the card event
		fmt.Println("#################4")
		fmt.Println("Received card action. Tag:", event.Action.Tag, "Option:", ActionData.Option, "Value:", event.Action.Value)
		//SendTODoraemon
		err = biz.SendSilenceToDaemon(num, event.UserId, 6940)
		if err != nil {
			logs.Error(logsign, "Error sending message to Doraemon", http.StatusInternalServerError)
			return
		} else {
			logs.Info(logsign, "Send message to Doraemon success", http.StatusOK)
		}
		c.Data["json"] = event
	} else {
		fmt.Println("#################5")
		logs.Error(logsign, "Received unknown event type", challenge.Type)
		fmt.Println("Received unknown event type", challenge.Type)
	}

	c.ServeJSON()
}
