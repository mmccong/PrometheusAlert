package biz

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
)

type SilenceBody struct {
	Duration int    `json:"duration"`
	Ids      []int  `json:"ids"`
	User     string `json:"user"`
}

type respBody struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func SendSilenceToDaemon(duration int, user string, ids int) error {
	// TODO
	//url := beego.AppConfig.String("DORAEMON_URL")
	url := "http://172.16.31.143:32000/api/v1/alerts"
	Cookie := "csrftoken=AHpxmCWHl2ZCSeRQ8xEkSRrLI1bvhOyu4CUAqrdbkADDs8f7l6u2685iHL79kvaG; username=admin; password=Es.9527; beegosessionID=82fb5c2fb8b2535ac0185b6eb8b197d0"

	requestBody := SilenceBody{
		Duration: duration,
		Ids:      []int{ids},
		User:     user,
	}
	// 将body转换为JSON格式
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		logs.Error("转换Json格式失败：%s", err)
		return err
	}
	logs.Info("Post Body requestBodyBytes:", string(requestBodyBytes))

	// 发送POST请求
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		logs.Error("创建POST请求失败：%s", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", Cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logs.Error("POST请求失败：%s", err)
		return err
	}

	defer resp.Body.Close()

	// Read and parse response body as JSON into a map
	var jsonData map[string]interface{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return err
	}

	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return err
	}

	// Process the parsed JSON data as a map
	fmt.Println("Response Data:", jsonData)

	if jsonData["code"] != 200 {
		logs.Error("Post Silence to Doraemon failed, response code is %d", jsonData["code"])
		return fmt.Errorf("Post Silence to Doraemon failed, response code is %d", jsonData["code"])
	}

	return nil
}
