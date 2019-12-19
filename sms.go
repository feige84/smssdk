package smssdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
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

var errInfo = make(map[int64]string)

func init() {
	errInfo[0] = "提交成功"
	errInfo[-99] = "异常"
	errInfo[101] = "提交参数不可为空，或参数格式错误"
	errInfo[102] = "发送时间格式不正确,正确格式为yyyy-MM-ddHH:mm:ss"
	errInfo[201] = "Ems格式转换错误"
	errInfo[202] = "Tms内容异常"
	errInfo[203] = "Mix格式彩信增加Smil文件错误"
	errInfo[1007] = "错误的信息类型：XXX"
	errInfo[1008] = "超过最大并发提交用户数XXX"
	errInfo[1009] = "号码为空或超过最大提交号码个数XXX"
	errInfo[1010] = "信息内容为空或超过最大信息字节长度XXX"
	errInfo[1011] = "超过最大企业号码长度:XXX,或企业号码不包含:XXX"
	errInfo[1013] = "账号密码不正确或账号状态异常"
	errInfo[1014] = "账户提交方式不正确或Ip受限"
	errInfo[1015] = "提交速度受限:XXX条/秒，指用户的提交速度大于规定的单个用户的最大提交速度(默认5000次/秒)"
	errInfo[1016] = "产品不存在或未开启:XXX"
	errInfo[1017] = "提交信息类型与产品信息类型不符合"
	errInfo[1018] = "超过产品发送时段:Begin=XXX,End=XXX,Curr=XXX"
	errInfo[1019] = "提交彩信必须有标题XXX"
	errInfo[1020] = "提交短信不可超过XXX个字"
	errInfo[1021] = "提交彩信不可超过XXXK"
	errInfo[1022] = "消息拆分失败，指长短信在进行拆分后的信息条数等于0或者大于8 "
	errInfo[1023] = "无效计费条数"
	errInfo[1025] = "Account:XXX 余额不足或者计费异常"
	errInfo[1026] = "提交至调度中心失败"
	errInfo[1027] = "提交成功，信息保存失败"
	errInfo[1028] = "账户%s无对应的产品%d"
	errInfo[1029] = "扩展产品%d不可提交多个号码"
	errInfo[1031] = "提交时间[%s]+31天>定时发送时间[%s]>提交时间[%s] 规则不成立"
	errInfo[1032] = "自由签名的产品101161801,签名格式不正确"
	errInfo[1033] = "产品%d的正则签名%s配置有误"
	errInfo[1035] = "强制签名的产品%d,签名格式不正确"
	errInfo[1037] = "未成功加载账户强制签名报备模块"
	errInfo[1038] = "强制签名的产品%d,签名%s未报备"
	errInfo[1039] = "未成功加载白名单模块"
	errInfo[1040] = "消息提交成功，但消息编号生产失败"
	errInfo[1041] = "未成功加载账户内容模板模块"
	errInfo[1042] = "内容不符合模板"
	errInfo[1043] = "未成功加载账户安全登录模块"
}

func (s *WelinkSms) SendSMS() (int64, string) {
	result, err := httpSend("http://api.51welink.com/json/sms/g_Submit", s)
	if err != nil {
		return -99, "异常"
	}
	state := gjson.GetBytes(result, "State")
	msgState := gjson.GetBytes(result, "MsgState")
	if state.Exists() {
		return state.Int(), msgState.String()
	} else {
		return -99, "异常"
	}
}

func (b *WelinkBase) GetRemain() (int64, bool) {
	result, err := httpSend("http://api.51welink.com/json/Query/GetRemain", b)
	if err != nil {
		return -99, false
	}
	state := gjson.GetBytes(result, "State")
	if state.Exists() {
		if state.Int() == 0 {
			return gjson.GetBytes(result, "Remain").Int(), true
		}
	}
	return 0, false
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
