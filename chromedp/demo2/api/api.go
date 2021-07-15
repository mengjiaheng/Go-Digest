package api

import (
	"demo2/utils"
	"fmt"
	"net/http"
	"strings"
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

//计算客户端response
func ResponseQop(digest Digest) (res string) {

	var A1, A2 string

	A1 = fmt.Sprintf("%s:%s:%s", digest.UserName, digest.Realm, digest.Password)

	if len(digest.Qop) != 0 && digest.Qop == "auth-int" {
		A2 = fmt.Sprintf("%s:%s:%s", digest.Method, digest.Uri, utils.Md5("request-entity-body"))

	} else {
		A2 = fmt.Sprintf("%s:%s", digest.Method, digest.Uri)
		digest.Qop = "auth"
	}

	res = fmt.Sprintf("%s:%s:%s:%s:%s:%s", utils.Md5(A1), digest.Nonce, digest.Nc, digest.Cnonce, digest.Qop, utils.Md5(A2))
	fmt.Println("加密前：", res)
	res = utils.Md5(res)
	fmt.Println("加密后：", res)

	return
}

func ClientRequest(url string) string {

	//创建request 请求
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("New Request:", err)
		return ""
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("request err:", err)
		return ""
	}

	var hashMap map[string]string

	//如果需要认证
	if response.StatusCode == 401 {
		hashMap = utils.ParseAuthorization(response.Header)
	} else {
		fmt.Println(response.StatusCode)
		return ""
	}

	var digest Digest
	digest.Qop = strings.Trim(hashMap["qop"], `"`)
	digest.Realm = strings.Trim(hashMap["realm"], `"`)
	digest.Nonce = strings.Trim(hashMap["nonce"], `"`)

	digest.Method = request.Method
	digest.Uri = request.URL.Path
	digest.Cnonce = utils.RandomString()
	digest.Nc = "0000001"
	digest.UserName = "mengjiaheng"
	digest.Password = "123456"
	digest.Response = ResponseQop(digest)

	auth := fmt.Sprintf(`"Digest username="%s",realm="%s",qop=%s,nonce="%s",uri="%s",nc=%s,cnonce="%s",response="%s"`, digest.UserName, digest.Realm, digest.Qop, digest.Nonce, digest.Uri, digest.Nc, digest.Cnonce, digest.Response)
	return auth
}
