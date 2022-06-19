package qrcode

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/zedisdog/sweetbean/sdk/wechat/mini/auth"
	"github.com/zedisdog/sweetbean/sdk/wechat/mini/qrcode/response"
)

const (
	qrCodeUnlimited = "https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=%s"
)

func NewQrCode(a *auth.Auth) *QrCode {
	return &QrCode{
		auth: a,
	}
}

type QrCode struct {
	auth *auth.Auth
}

func (q QrCode) GetUnlimited(sence map[string]string) (r response.QrCodeUnlimited, err error) {
	token, err := q.auth.GetAccessToken()
	if err != nil {
		return
	}
	resp, err := http.Get(fmt.Sprintf(qrCodeUnlimited, token))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(content, &r)
	if err == nil {
		return
	} else {
		r.Content = content
		err = nil
	}
	return
}
