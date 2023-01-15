package qr

import "github.com/skip2/go-qrcode"

func EncodeStringToPng(str string) ([]byte, error) {
	return qrcode.Encode(str, qrcode.Medium, 256)
}
