package main

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Digest struct {
	UserName  string `json:"username"`
	Password  string `json:"password"`
	Realm     string `json:"realm"`
	Qop       string `json:"qop"`
	Nonce     string `json:"nonce"`
	Uri       string `json:"uri"`
	Nc        string `json:"nc"`
	Cnonce    string `json:"cnonce"`
	Response  string `json:"Response"`
	LastNonce string
}

var digest Digest

//储存nonce 和 nc 对应关系
// var ncMap = make(map[string]string)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

//产生随机数
func RandomString() string {
	b := make([]byte, rand.Intn(10)+10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

//根据qop计算response
func ResponseQop(method string, hashMap map[string]string, digest Digest) (res string) {

	var A1, A2 string
	if hashMap["qop"] == "auth" {
		A1 = hashMap["username"][1:len(hashMap["username"])-1] + ":" + hashMap["realm"][1:len(hashMap["realm"])-1] + ":" + digest.Password
		A2 = method + ":" + hashMap["uri"][1:len(hashMap["uri"])-1]
	}
	res = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%x", md5.Sum([]byte(A1)))+":"+hashMap["nonce"][1:len(hashMap["nonce"])-1]+":"+hashMap["nc"]+":"+hashMap["cnonce"][1:len(hashMap["cnonce"])-1]+":"+hashMap["qop"]+":"+fmt.Sprintf("%x", md5.Sum([]byte(A2))))))
	return
}

//计算响应摘要
func ResponsePauth(hashMap map[string]string, digest Digest) (rspauth string) {
	var A1, A2 string
	if digest.Qop == "auth" {
		A1 = hashMap["username"][1:len(hashMap["username"])-1] + ":" + hashMap["realm"][1:len(hashMap["realm"])-1] + ":" + digest.Password
		A2 = ":" + hashMap["uri"][1:len(hashMap["uri"])-1]
	}

	rspauth = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%x", md5.Sum([]byte(A1)))+":"+hashMap["nonce"][1:len(hashMap["nonce"])-1]+":"+hashMap["nc"]+":"+hashMap["cnonce"][1:len(hashMap["cnonce"])-1]+":"+hashMap["qop"]+":"+fmt.Sprintf("%x", md5.Sum([]byte(A2))))))
	return
}

func SetAuthenticateHeader(c *gin.Context) {
	digest.Realm = "Digest test"
	digest.Qop = "auth"
	digest.Nonce = RandomString()
	c.Writer.Header().Set("WWW-Authenticate", `Digest realm="`+digest.Realm+`",qop="`+digest.Qop+`",nonce="`+digest.Nonce+`"`)
}

func DigestAuth(c *gin.Context, header http.Header) {
	auth := header.Get("Authorization")

	if auth == "" {
		SetAuthenticateHeader(c)
		c.String(401, "Unauthorized")
		return
	}

	auths := strings.SplitN(auth, " ", 2)
	str := strings.Split(auths[1], ", ")
	hashmap := make(map[string]string)
	for _, v := range str {
		split := strings.Split(v, ",")

		for _, s := range split {
			i := strings.Split(s, "=")
			hashmap[i[0]] = i[1]
		}
	}

	for k, v := range hashmap {
		fmt.Println(k, ":", v)
	}

	//验证nonce是否一致
	if hashmap["nonce"][1:len(hashmap["nonce"])-1] != digest.Nonce {
		fmt.Println("随机数：", hashmap["nonce"][1:len(hashmap["nonce"])-1], digest.Nonce)
		if digest.LastNonce == hashmap["nonce"][1:len(hashmap["nonce"])-1] {
			SetAuthenticateHeader(c)
			c.String(401, "Unauthorized")
		} else {
			digest.LastNonce = hashmap["nonce"][1 : len(hashmap["nonce"])-1]
			// c.Writer.Header
			c.String(402, "服务端随机数不匹配 error")
		}
		return
	}

	//判断同一随机数是否重复请求
	// if _, ok := ncMap[auths[3]]; ok {
	// 	fmt.Println("重复请求")
	// 	ncMap[auths[3]] = auths[4]
	// }

	// userName := auths[1]
	// password := "123456"
	// if auths[2] == "auth" {
	// 	rsp := md5(md5(userName + ":" + realm))
	// }
	method := c.Request.Method
	digest.Password = "123456"

	response := ResponseQop(method, hashmap, digest)
	//判断response值是否一致
	if hashmap["response"][1:len(hashmap["response"])-1] != response {
		// fmt.Println(hashmap["response"][1 : len(hashmap["response"])-1])
		// fmt.Println(response)
		if digest.LastNonce == hashmap["nonce"][1:len(hashmap["nonce"])-1] {
			SetAuthenticateHeader(c)
			c.String(401, "Unauthorized")
		} else {
			digest.LastNonce = hashmap["nonce"][1 : len(hashmap["nonce"])-1]
			// c.Writer.Header
			c.String(404, "认证失败")
		}
		return
	}

	//产生下一个新的随机数
	// has := md5.Sum([]byte(time.Now().String()))
	// nextnonce = fmt.Sprintf("%x", has)

	// //产生响应摘要
	// has = md5.Sum([]byte(time.Now().String()))
	// rspauth = fmt.Sprintf("%x", has)

	// cnonce = auths[6]
	// nc = auths[4]
	rspauth := ResponsePauth(hashmap, digest)
	c.Writer.Header().Set("Authorization-Info", `qop="`+hashmap["qop"]+`",rspauth="`+rspauth+`",cnonce="`+hashmap["cnonce"][1:len(hashmap["cnonce"])-1])
	c.String(200, "OK")

	// if nonce==header.Get()
}

func main() {
	r := gin.Default()
	r.GET("/index", func(c *gin.Context) {
		DigestAuth(c, c.Request.Header)
		// c.String(200, "hello world")
	})
	r.Run(":8090")
}
