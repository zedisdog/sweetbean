package official

import "github.com/zedisdog/sweetbean/sdk/wechat/official/oauth"

func NewClient(appID string, secret string) *Client {
	return &Client{
		appID:  appID,
		secret: secret,
	}
}

type Client struct {
	appID  string
	secret string
	oauth  *oauth.Oauth
}

func (c Client) Oauth() *oauth.Oauth {
	if c.oauth == nil {
		c.oauth = oauth.NewOauth(c.appID, c.secret)
	}
	return c.oauth
}
