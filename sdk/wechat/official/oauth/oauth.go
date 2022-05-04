package oauth

import "fmt"

const (
	SnsapiBase     = "snsapi_base"
	SnsapiUserinfo = "snsapi_userinfo"
	oauthUrlTmpl   = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s#wechat_redirect"
)

func NewOauth(appID string, secret string) *Oauth {
	return &Oauth{
		appID:  appID,
		secret: secret,
	}
}

type Oauth struct {
	appID  string
	secret string
}

func (o Oauth) GenRedirectUrl(callbackUrl string, options ...func(*redirectOptions)) string {
	opt := &redirectOptions{
		scope: SnsapiBase,
	}
	for _, set := range options {
		set(opt)
	}
	return fmt.Sprintf(oauthUrlTmpl, o.appID, callbackUrl, opt.scope, opt.State())
}
