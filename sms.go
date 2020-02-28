package smssdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type WelinkBase struct {
	Name  string `json:"sname"`
	Pwd   string `json:"spwd"`
	Prdid string `json:"sprdid"`
}

type WelinkSms struct {
	Name  string `json:"sname"`
	Pwd   string `json:"spwd"`
	Prdid string `json:"sprdid"`
	Dst   string `json:"sdst"`
	Msg   string `json:"smsg"`
}

type Remain struct {
	State  int64 `json:"State"`
	Remain int64 `json:"Remain"`
}

type SendResult struct {
	MsgState string `json:"MsgState"`
	State    int64  `json:"State"`
	MsgID    string `json:"MsgID"`
	Reserve  int64  `json:"Reserve"`
}

const apiGateWay = "http://api.51welink.com/"

func (s *WelinkSms) SendSMS() (*SendResult, error) {
	result, err := httpSend(apiGateWay+"json/sms/g_Submit", s)
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(result))
	data := SendResult{}
	if err = json.Unmarshal(result, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (b *WelinkBase) GetRemain() (*Remain, error) {
	result, err := httpSend(apiGateWay+"json/Query/GetRemain", b)
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(result))
	data := Remain{}
	if err = json.Unmarshal(result, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func httpSend(router string, param interface{}) ([]byte, error) {
	sendBody, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", router, bytes.NewBuffer(sendBody))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	httpClient := &http.Client{}
	httpClient.Timeout = 3 * time.Second
	response, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求错误:%d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
