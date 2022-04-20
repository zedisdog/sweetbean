package mini

import "github.com/zedisdog/sweetbean/sdk/wechat/mini/sns"

func NewClient(appID string, secret string) *Client {
	return &Client{
		appID:  appID,
		secret: secret,
		Sns:    sns.NewSns(appID, secret),
	}
}

type Client struct {
	appID  string
	secret string
	*sns.Sns
}
