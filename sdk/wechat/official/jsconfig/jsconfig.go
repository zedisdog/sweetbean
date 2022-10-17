package jsconfig

import (
	"fmt"
	"time"

	"github.com/zedisdog/sweetbean/errx"
	"github.com/zedisdog/sweetbean/sdk/wechat/official/ticket"
	"github.com/zedisdog/sweetbean/tools"
)

// debug: true, // 开启调试模式,调用的所有 api 的返回值会在客户端 alert 出来，若要查看传入的参数，可以在 pc 端打开，参数信息会通过 log 打出，仅在 pc 端时才会打印。
//
//	appId: '', // 必填，公众号的唯一标识
//	timestamp: , // 必填，生成签名的时间戳
//	nonceStr: '', // 必填，生成签名的随机串
//	signature: '',// 必填，签名
//	jsApiList: [] // 必填，需要使用的 JS 接口列表
type Config struct {
	AppID     string `json:"appId"`
	Timestamp int64  `json:"timestamp"`
	NonceStr  string `json:"nonceStr"`
	Signature string `json:"signature"`
}

func NewJsConfig(appID string, ticket *ticket.Ticket) *JsConfig {
	return &JsConfig{
		appID:  appID,
		ticket: ticket,
	}
}

type JsConfig struct {
	appID  string
	ticket *ticket.Ticket
}

func (j JsConfig) Gen(url string) (config Config, err error) {
	ticket, err := j.ticket.GetTicket()
	if err != nil {
		err = errx.Wrap(err, "gen config failed")
		return
	}

	config.AppID = j.appID
	config.NonceStr = tools.RandString(6)
	config.Timestamp = time.Now().Unix()
	text := fmt.Sprintf(
		"jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s",
		ticket,
		config.NonceStr,
		config.Timestamp,
		url,
	)
	config.Signature = tools.Sha1(text)

	return
}
