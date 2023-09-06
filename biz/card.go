package biz

import (
	"PrometheusAlert/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func ConstructAlterCard(message models.AlertContent, img string) (card string) {
	cardContent := &CardContent{
		Config: &CardConfig{
			WideScreenMode: true,
		},
		Header: &CardHeader{
			Template: "red",
			Title: &CardText{
				Tag:     "plain_text",
				Content: message.Receiver + "【报警】",
			},
		},
	}
	elements := make([]interface{}, 0)
	// card block 1
	element1 := &CardElement{
		Tag: "div",
		Fields: []*CardField{
			{
				IsShort: true,
				Text: &CardText{
					Content: "**🕐 触发时间：**\n" + message.Alerts[0].StartsAt,
					Tag:     "lark_md",
				},
			}, {
				IsShort: true,
				Text: &CardText{
					Content: "**🔢 告警事件 ID：**\n" + fmt.Sprintln(message.Alerts[0].Labels.Id),
					Tag:     "lark_md",
				},
			}, {
				IsShort: true,
				Text: &CardText{
					Content: "**📋 项目：**\n" + message.Alerts[0].Labels.Alertname,
					Tag:     "lark_md",
				},
			}, {
				IsShort: true,
				Text: &CardText{
					Content: "**👤 一级值班：**\n<at id=all>丛明明</at>",
					Tag:     "lark_md",
				},
			}, {
				IsShort: true,
				Text: &CardText{
					Content: "**👤 二级值班：**\n<at id=all>所有人</at>",
					Tag:     "lark_md",
				},
			},
		},
	}
	elements = append(elements, element1)

	// card block 2, image block
	element2 := &CardElement{
		Tag:    "img",
		ImgKey: img,
		Alt: &CardText{
			Content: " ",
			Tag:     "plain_text",
		},
		Title: &CardText{
			Content: "🔴 " + message.Alerts[0].Annotations.Summary,
			Tag:     "lark_md",
		},
	}
	elements = append(elements, element2)

	// card block 3, note block
	element3 := CardNote{
		Tag: "note",
	}
	element3Elements := make([]interface{}, 0)
	element3Elements = append(element3Elements, &CardText{
		Content: "🔴 " + message.Alerts[0].Annotations.Description,
		Tag:     "plain_text",
	})
	element3.Elements = element3Elements
	elements = append(elements, element3)

	// card action block
	element4 := &CardActionBlock{
		Tag: "action",
	}
	actions := make([]interface{}, 0)
	/*
			button := &CardButton{
				Tag: "button",
				Text: &CardText{
					Tag:     "plain_text",
					Content: "跟进处理",
				},
				Type:  "primary",
				Value: map[string]string{"key1": "value1"},
			}
		actions = append(actions, button)
	*/
	selectMenu := &CardSelectMenu{
		Tag: "select_static",
		PlaceHolder: &CardText{
			Content: "暂时屏蔽",
			Tag:     "plain_text",
		},
		Options: []*CardOption{
			{
				Text: &CardText{
					Content: "屏蔽10分钟",
					Tag:     "plain_text",
				},
				Value: "10",
			}, {
				Text: &CardText{
					Content: "屏蔽30分钟",
					Tag:     "plain_text",
				},
				Value: "30",
			}, {
				Text: &CardText{
					Content: "屏蔽1小时",
					Tag:     "plain_text",
				},
				Value: "60",
			}, {
				Text: &CardText{
					Content: "屏蔽1天",
					Tag:     "plain_text",
				},
				Value: "1440",
			}, {
				Text: &CardText{
					Content: "屏蔽3天",
					Tag:     "plain_text",
				},
				Value: "4320",
			},
			{
				Text: &CardText{
					Content: "屏蔽7天",
					Tag:     "plain_text",
				},
				Value: "10080",
			},
		},
		Value: map[string]string{"key": "value"},
	}
	actions = append(actions, selectMenu)
	element4.Actions = actions
	elements = append(elements, element4)

	// card split line
	element5 := &CardSplitLine{
		Tag: "hr",
	}
	elements = append(elements, element5)

	// card
	element6 := &CardElement{
		Tag: "div",
		Text: &CardText{
			Content: time.Now().Format("2006-01-02 15:04:05"),
			Tag:     "lark_md",
		},
	}
	elements = append(elements, element6)

	cardContent.Elements = elements

	cardBytes, err := json.Marshal(cardContent)
	if err != nil {
		logs.Error("failed to marshal")
		return ""
	}
	logs.Info("card_content: %v", string(cardBytes))
	return string(cardBytes)
}

func ConstructResolvedCard(message models.AlertContent, img string) (card string) {
	cardContent := &CardContent{
		Config: &CardConfig{
			WideScreenMode: true,
		},
		Header: &CardHeader{
			Template: "green",
			Title: &CardText{
				Tag:     "plain_text",
				Content: message.Receiver + "【已处理】",
			},
		},
	}
	elements := make([]interface{}, 0)
	// card block 1
	element1 := &CardElement{
		Tag: "div",
		Fields: []*CardField{
			{
				IsShort: true,
				Text: &CardText{
					Content: "**🕐 触发时间：**\n" + message.Alerts[0].StartsAt,
					Tag:     "lark_md",
				},
			}, {
				IsShort: true,
				Text: &CardText{
					Content: "**🔢 告警事件 ID：**\n" + message.Alerts[0].Fingerprint,
					Tag:     "lark_md",
				},
			}, {
				IsShort: true,
				Text: &CardText{
					Content: "**📋 项目：**\n" + message.Alerts[0].Labels.Alertname,
					Tag:     "lark_md",
				},
			}, {
				IsShort: true,
				Text: &CardText{
					Content: "**👤 一级值班：**\n<at id=all>丛明明</at>",
					Tag:     "lark_md",
				},
			}, {
				IsShort: true,
				Text: &CardText{
					Content: "**👤 二级值班：**\n<at id=all>所有人</at>",
					Tag:     "lark_md",
				},
			},
		},
	}
	elements = append(elements, element1)

	// card block 2, image block
	element2 := &CardElement{
		Tag:    "img",
		ImgKey: img,
		Alt: &CardText{
			Content: " ",
			Tag:     "plain_text",
		},
		Title: &CardText{
			Content: message.Alerts[0].Annotations.Description,
			Tag:     "lark_md",
		},
	}
	elements = append(elements, element2)

	// card block 3, note block
	element3 := CardNote{
		Tag: "note",
	}
	element3Elements := make([]interface{}, 0)
	element3Elements = append(element3Elements, &CardText{
		Content: "🔴 支付失败数  🔵 支付成功数",
		Tag:     "plain_text",
	})
	element3.Elements = element3Elements
	elements = append(elements, element3)

	// card split line
	element4 := &CardSplitLine{
		Tag: "hr",
	}
	elements = append(elements, element4)

	element5 := &CardNote{Tag: "note"}
	element5Elements := make([]interface{}, 0)
	element5Elements = append(element5Elements, &CardText{
		Tag:     "plain_text",
		Content: "✅ 李健已处理此报警",
	})
	element5.Elements = element5Elements
	elements = append(elements, element5)

	element6 := &CardElement{
		Tag: "div",
		Text: &CardText{
			Content: time.Now().Format("2006-01-02 15:04:05"),
			Tag:     "lark_md",
		},
	}
	elements = append(elements, element6)

	cardContent.Elements = elements
	cardBytes, err := json.Marshal(cardContent)

	if err != nil {
		logs.Error("failed to marshal")
		return ""
	}
	logs.Info("card_content: %v", string(cardBytes))
	return string(cardBytes)
}

func ConstructSSLCertCard(message models.SSLContent) (card string) {
	/*
		var msg []models.SSLContentText

		err := json.Unmarshal([]byte(message.Content), &msg)
		if err != nil {
			logs.Error("failed to unmarshal")
			return ""
		}
	*/
	cardContent := &CardContent{
		Config: &CardConfig{
			WideScreenMode: true,
		},
		Header: &CardHeader{
			Template: "red",
			Title: &CardText{
				Tag:     "plain_text",
				Content: message.Title,
			},
		},
	}

	elements := make([]interface{}, 0)
	element1 := &CardElement{
		Tag: "div",
		Fields: []*CardField{
			{
				IsShort: true,
				Text: &CardText{
					Content:   "**二级域名**",
					Tag:       "lark_md",
					TextAlign: "center",
				},
			},
		},
	}
	elements = append(elements, element1)

	element2 := &CardElement{
		Tag: "div",
		Fields: []*CardField{
			{
				IsShort: true,
				Text: &CardText{
					Content:   message.Content,
					Tag:       "lark_md",
					TextAlign: "center",
				},
			},
		},
	}
	elements = append(elements, element2)

	/*
		element1 := &CardElement{
			Tag: "column",
			Fields: []*CardField{
				{
					IsShort: true,
					Text: &CardText{
						Content: message.Content,
						Tag:     "lark_md",
					},
				},
			},
		}
		elements = append(elements, element1)



			element2 := &CardElement{
				Tag: "div",
				Fields: []*CardField{
					{
						IsShort: true,
						Text: &CardText{
							Content:   "**颁发日期**",
							Tag:       "markdown",
							TextAlign: "center",
						},
					},
				},
			}
			elements = append(elements, element2)

			element3 := &CardElement{
				Tag: "div",
				Fields: []*CardField{
					{
						IsShort: true,
						Text: &CardText{
							Content:   "**截止日期**",
							Tag:       "lark_md",
							TextAlign: "center",
						},
					},
				},
			}
			elements = append(elements, element3)

			element4 := &CardElement{
				Tag: "div",
				Fields: []*CardField{
					{
						IsShort: true,
						Text: &CardText{
							Content:   "**剩余天数**",
							Tag:       "lark_md",
							TextAlign: "center",
						},
					},
				},
			}

			elements = append(elements, element4)

			for _, v := range msg {
				//element := "element" + strconv.Itoa(index)
				element21 := &CardElement{
					Tag: "div",
					Fields: []*CardField{
						{
							IsShort: true,
							Text: &CardText{
								Content:   v.Domain,
								Tag:       "lark_md",
								TextAlign: "center",
							},
						},
					},
				}
				elements = append(elements, element21)

				element22 := &CardElement{
					Tag: "div",
					Fields: []*CardField{
						{
							IsShort: true,
							Text: &CardText{
								Content:   v.CreateTime,
								Tag:       "lark_md",
								TextAlign: "center",
							},
						},
					},
				}
				elements = append(elements, element22)

				element23 := &CardElement{
					Tag: "div",
					Fields: []*CardField{
						{
							IsShort: true,
							Text: &CardText{
								Content:   v.UpdateTime,
								Tag:       "lark_md",
								TextAlign: "center",
							},
						},
					},
				}
				elements = append(elements, element23)

				element24 := &CardElement{
					Tag: "div",
					Fields: []*CardField{
						{
							IsShort: true,
							Text: &CardText{
								Content:   strconv.Itoa(v.ExpireTime),
								Tag:       "lark_md",
								TextAlign: "center",
							},
						},
					},
				}
				elements = append(elements, element24)
			}
	*/
	// card split line
	element5 := &CardSplitLine{
		Tag: "hr",
	}
	elements = append(elements, element5)
	//card send time
	element6 := &CardElement{
		Tag: "div",
		Text: &CardText{
			Content: time.Now().Format("2006-01-02 15:04:05"),
			Tag:     "lark_md",
		},
	}
	elements = append(elements, element6)
	cardContent.Elements = elements
	cardBytes, err := json.Marshal(cardContent)
	if err != nil {
		logs.Error("failed to marshal")
		return ""
	}
	logs.Info("card_content: %v", string(cardBytes))
	return string(cardBytes)
}

func ConstructMarkdownCard(message models.MarkdownContent) (card string) {
	cardContent := &CardContent{
		Config: &CardConfig{
			WideScreenMode: true,
		},
		Header: &CardHeader{
			Template: "orange",
			Title: &CardText{
				Tag:     "plain_text",
				Content: message.Title + "数据库审计",
			},
		},
	}

	elements := make([]interface{}, 0)
	element1 := &CardElement{
		Tag: "div",
		Fields: []*CardField{
			{
				IsShort: true,
				Text: &CardText{
					Content:   message.Text,
					Tag:       "lark_md",
					TextAlign: "center",
				},
			},
		},
	}
	elements = append(elements, element1)

	// card split line
	element5 := &CardSplitLine{
		Tag: "hr",
	}
	elements = append(elements, element5)
	//card send time
	element6 := &CardElement{
		Tag: "div",
		Text: &CardText{
			Content: time.Now().Format("2006-01-02 15:04:05"),
			Tag:     "lark_md",
		},
	}
	elements = append(elements, element6)
	cardContent.Elements = elements
	cardBytes, err := json.Marshal(cardContent)
	if err != nil {
		logs.Error("failed to marshal")
		return ""
	}
	logs.Info("card_content: %v", string(cardBytes))
	return string(cardBytes)
}

func ConstructSilenceCard(message models.AlertContent, img string) (card string) {
	cardContent := &CardContent{
		Config: &CardConfig{
			WideScreenMode: true,
		},
		Header: &CardHeader{
			Template: "grey",
			Title: &CardText{
				Tag:     "plain_text",
				Content: message.Receiver + "【已屏蔽报警 1 小时】",
			},
		},
	}
	elements := make([]interface{}, 0)
	// card block 1
	element1 := &CardElement{
		Tag: "div",
		Fields: []*CardField{
			{
				IsShort: true,
				Text: &CardText{
					Content: "**🕐 触发时间：**\n" + message.Alerts[0].StartsAt,
					Tag:     "lark_md",
				},
			}, {
				IsShort: true,
				Text: &CardText{
					Content: "**🔢 告警事件 ID：**\n" + message.Alerts[0].Fingerprint,
					Tag:     "lark_md",
				},
			}, {
				IsShort: true,
				Text: &CardText{
					Content: "**📋 项目：**\n" + message.Alerts[0].Labels.Alertname,
					Tag:     "lark_md",
				},
			}, {
				IsShort: true,
				Text: &CardText{
					Content: "**👤 一级值班：**\n<at id=all>丛明明</at>",
					Tag:     "lark_md",
				},
			}, {
				IsShort: true,
				Text: &CardText{
					Content: "**👤 二级值班：**\n<at id=all>所有人</at>",
					Tag:     "lark_md",
				},
			},
		},
	}
	elements = append(elements, element1)

	// card block 2, image block
	element2 := &CardElement{
		Tag:    "img",
		ImgKey: img,
		Alt: &CardText{
			Content: " ",
			Tag:     "plain_text",
		},
		Title: &CardText{
			Content: message.Alerts[0].Annotations.Description,
			Tag:     "lark_md",
		},
	}
	elements = append(elements, element2)

	// card block 3, note block
	element3 := CardNote{
		Tag: "note",
	}
	element3Elements := make([]interface{}, 0)
	element3Elements = append(element3Elements, &CardText{
		Content: "🔴 支付失败数  🔵 支付成功数",
		Tag:     "plain_text",
	})
	element3.Elements = element3Elements
	elements = append(elements, element3)

	// card action block
	element4 := &CardActionBlock{
		Tag: "action",
	}
	actions := make([]interface{}, 0)
	button := &CardButton{
		Tag: "button",
		Text: &CardText{
			Tag:     "plain_text",
			Content: "跟进处理",
		},
		Type:  "primary",
		Value: map[string]string{"continue": "true"},
	}
	actions = append(actions, button)
	button2 := &CardButton{
		Tag: "button",
		Text: &CardText{
			Tag:     "plain_text",
			Content: "取消屏蔽",
		},
		Type:  "primary",
		Value: map[string]string{"unsilence": "true"},
	}
	actions = append(actions, button2)
	element4.Actions = actions
	elements = append(elements, element4)

	// card split line
	element5 := &CardSplitLine{
		Tag: "hr",
	}
	elements = append(elements, element5)

	// card
	element6 := &CardElement{
		Tag: "div",
		Text: &CardText{
			Content: "🙋🏼 [我要反馈误报](https://open.feishu.cn/) | 📝 [录入报警处理过程](https://open.feishu.cn/)",
			Tag:     "lark_md",
		},
	}
	elements = append(elements, element6)

	element7 := &CardElement{
		Tag: "div",
		Text: &CardText{
			Content: time.Now().Format("2006-01-02 15:04:05"),
			Tag:     "lark_md",
		},
	}
	elements = append(elements, element7)

	cardContent.Elements = elements

	cardBytes, err := json.Marshal(cardContent)
	if err != nil {
		logs.Error("failed to marshal")
		return ""
	}
	logs.Info("card_content: %v", string(cardBytes))
	return string(cardBytes)
}

func UploadImage(token string) (*UploadImageResponseBody, error) {
	cli := &http.Client{}

	pwd, _ := os.Getwd()
	parent := filepath.Dir(pwd)
	fmt.Println(parent)
	path := parent + "/PrometheusAlert/resource/upload/alert.png"
	logs.Info("path: %v", path)
	image, err := os.Open(path)
	if err != nil {
		logs.Error("failed to open image")
		return nil, err
	}
	defer image.Close()

	buffer := &bytes.Buffer{}
	write := multipart.NewWriter(buffer)
	w, err := write.CreateFormFile("image", filepath.Base(path))
	if err != nil {
		logs.Error("failed to create form file")
		return nil, err
	}
	_, err = io.Copy(w, image)
	if err != nil {
		logs.Error("copy image failed")
		return nil, err
	}
	params := make(map[string]string)
	params["image_type"] = "message"
	for k, v := range params {
		err = write.WriteField(k, v)
		if err != nil {
			return nil, err
		}
	}

	err = write.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uploadImageURL, buffer)
	if err != nil {
		logs.Error("new request failed")
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", write.FormDataContentType())

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
	uploadImageResp := &UploadImageResponse{}
	err = json.Unmarshal(body, uploadImageResp)
	if err != nil {
		logs.Error("failed to unmarshal")
		return nil, err
	}
	if uploadImageResp.Code != 0 {
		logs.Warn("failed to upload image, code: %v, msg: %v, log_id: %v", uploadImageResp.Code, uploadImageResp.Message, logID)
		return nil, fmt.Errorf("create image failed")
	}
	return uploadImageResp.Data, nil
}
