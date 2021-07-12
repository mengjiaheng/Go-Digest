package core

import (
	"context"
	"demo1/utils"
	"fmt"
	"strings"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type Digest struct {
	UserName string
	Password string
	Method   string
	Realm    string
	Qop      string
	Nonce    string
	Uri      string
	Nc       string
	Cnonce   string
	Response string
}

var digest Digest

//设置请求头
func LoadHeaders(auth string) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		return network.SetExtraHTTPHeaders(
			network.Headers{"Authorization": auth}).Do(ctx)

	}
}

//解析响应头数据
func ParseAuthorization(headers network.Headers) string {
	auth := headers["Www-Authenticate"].(string)

	auths := strings.SplitN(auth, " ", 2)

	//如果认证方式为Digest
	if auths[0] == "Digest" {
		hashMap := utils.StringMap(auths[1])

		digest.Qop = hashMap["qop"]

		//产生客户端随机数与请求计数器
		digest.Cnonce = utils.RandomString()
		digest.Nc = "0000001"

		//设置账号密码
		digest.UserName = "mengjiaheng"
		digest.Password = "123456"

		//
		var authCha network.AuthChallenge
		authCha.Scheme = "Digest"
		authCha.Realm = digest.Realm
		authCha.Origin = "http://127.0.0.1:8090/index"

		var authChaRes network.AuthChallengeResponse
		authChaRes.Response = "ProvideCredentials"
		authChaRes.Username = "mengjiaheng"
		authChaRes.Password = "123456"

		json, _ := authChaRes.MarshalJSON()
		fmt.Println(json)
		authChaRes.Response.String()
		digest.Response = ResponseQop(hashMap, digest)
		auth = `Digest username="` + digest.UserName + `",realm="` + hashMap["realm"][1:len(hashMap["realm"])-1] + `",qop=` + hashMap["qop"][1:len(hashMap["qop"])-1] + `,nonce="` + hashMap["nonce"][1:len(hashMap["nonce"])-1] + `",uri="` + digest.Uri + `",nc=` + digest.Nc + `,cnonce="` + digest.Cnonce + `",response="` + digest.Response + `"`
		return auth
	}
	return ""
}
