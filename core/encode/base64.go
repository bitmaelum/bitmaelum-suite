package encode

import "encoding/base64"

// Encode data into base64
func Encode(src []byte) []byte {
    dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))

    base64.StdEncoding.Encode(dst, src)

    return dst
}

// Decode base64 back into bytes
func Decode(src []byte) ([]byte, error) {
    dst := make([]byte, base64.StdEncoding.DecodedLen(len(src)))

    _, err := base64.StdEncoding.Decode(dst, src)
    return dst, err
}
