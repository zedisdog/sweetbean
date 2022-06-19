package qrcode

import (
	"encoding/json"
	"fmt"

	"github.com/zedisdog/sweetbean/net/http"
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

type Color struct {
	R int `json:"r"`
	G int `json:"g"`
	B int `json:"b"`
}

type QrCodeUnlimitedOptions struct {
	Page       string `json:"page"`
	CheckPath  bool   `json:"check_page"`
	EnvVersion string `json:"env_version"`
	Width      int    `json:"width"`
	AutoColor  bool   `json:"auto_color"`
	LineColor  Color  `json:"line_color"`
	IsHyaline  bool   `json:"is_hyaline"`
}

type QrCode struct {
	auth *auth.Auth
}

func (q QrCode) GetUnlimited(sence map[string]string, setters ...func(*QrCodeUnlimitedOptions)) (r response.QrCodeUnlimited, err error) {
	options := QrCodeUnlimitedOptions{
		Page:       "",
		CheckPath:  true,
		EnvVersion: "release",
		Width:      430,
		AutoColor:  false,
		LineColor: Color{
			G: 0,
			R: 0,
			B: 0,
		},
		IsHyaline: false,
	}
	for _, set := range setters {
		set(&options)
	}
	token, err := q.auth.GetAccessToken()
	if err != nil {
		return
	}
	resp, err := http.PostJSON(fmt.Sprintf(qrCodeUnlimited, token), options)
	if err != nil {
		return
	}
	err = json.Unmarshal(resp, &r)
	if err == nil {
		return
	} else {
		r.Content = resp
		err = nil
	}
	return
}
