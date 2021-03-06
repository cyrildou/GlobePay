package sign

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
	"strings"
	"time"
)

type Sign struct {
	PartnerCode    string //商户编码
	Time           int64  //UTC毫秒时间戳
	NonceStr       string //随机字符串(建议长度在10到30位)
	CredentialCode string //系统为商户分配的开发校验码
	Sign           string //签名
}

func createRandomString(len int) string {
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(str[randomInt.Int64()])
	}
	return container
}

func NewSign(partnerCode string, credentialCode string) *Sign {
	var sign Sign
	sign.PartnerCode = partnerCode
	sign.CredentialCode = credentialCode
	return &sign
}

func (sign *Sign) GenSign() {

	sign.Time = time.Now().UTC().UnixNano() / 1e6 //UTC毫秒时间戳
	sign.NonceStr = createRandomString(15)

	signStr := fmt.Sprintf("%v&%v&%v&%v", sign.PartnerCode, sign.Time, sign.NonceStr, sign.CredentialCode)

	sum := sha256.Sum256([]byte(signStr))
	sign.Sign = strings.ToLower(fmt.Sprintf("%x", sum))
}

func (sign *Sign) GenSignByParams(t int64, nonceStr string) string {

	sign.Time = t
	sign.NonceStr = nonceStr

	signStr := fmt.Sprintf("%v&%v&%v&%v", sign.PartnerCode, sign.Time, sign.NonceStr, sign.CredentialCode)

	sum := sha256.Sum256([]byte(signStr))
	return strings.ToLower(fmt.Sprintf("%x", sum))
}

func (sign *Sign) GenSignURL(host string) string {

	sign.GenSign()
	query := fmt.Sprintf("time=%v&nonce_str=%v&sign=%v", sign.Time, sign.NonceStr, sign.Sign)
	return host + "?" + query
}
