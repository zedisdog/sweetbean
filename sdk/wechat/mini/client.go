package mini

import (
	"github.com/zedisdog/sweetbean/sdk/wechat/mini/auth"
	"github.com/zedisdog/sweetbean/sdk/wechat/mini/qrcode"
	"github.com/zedisdog/sweetbean/sdk/wechat/mini/sns"
)

func NewClient(appID string, secret string) *Client {
	auth := auth.NewAuth(appID, secret)
	return &Client{
		appID:  appID,
		secret: secret,
		Sns:    sns.NewSns(appID, secret),
		Auth:   auth,
		QrCode: qrcode.NewQrCode(auth),
	}
}

type Client struct {
	appID  string
	secret string
	*sns.Sns
	*auth.Auth
	*qrcode.QrCode
}
