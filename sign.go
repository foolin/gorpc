package gorpc
import (
	"fmt"
	"crypto/sha1"
	"crypto/hmac"
)


func makeSign(value, secret string) string  {
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(value))
	hash := mac.Sum(nil)
	return fmt.Sprintf("%x", hash)
}