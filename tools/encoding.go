package tools

import "encoding/base64"

func Base64Encode(src []byte) (buf []byte) {
	buf = make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(buf, src)
	return
}

func Base64Decode(src string) (dbuf []byte, err error) {
	dbuf = make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	n, err := base64.StdEncoding.Decode(dbuf, []byte(src))
	return dbuf[:n], err
}
