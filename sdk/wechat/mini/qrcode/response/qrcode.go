package response

import "github.com/zedisdog/sweetbean/sdk/wechat/mini/common"

type QrCodeUnlimited struct {
	common.ErrorResponse
	Content []byte
}
