package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/gocolly/colly"
)

type pageInfo struct {
	StatusCode int
	Links      map[string]int
}

type Digest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Realm    string `json:"realm"`
	Qop      string `json:"qop"`
	Nonce    string `json:"nonce"`
	Uri      string `json:"uri"`
	Nc       string `json:"nc"`
	Cnonce   string `json:"cnonce"`
	Response string `json:"Response"`
}

var digest Digest

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

//产生随机数
func RandomString() string {
	b := make([]byte, rand.Intn(10)+10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

//储存nonce 和 nc 对应关系
var ncMap = make(map[string]string)

//根据qop计算response
func ResponseQop(method string, hashMap map[string]string, digest Digest) (res string) {

	var A1, A2 string
	if digest.Qop == "auth" {
		A1 = digest.UserName + ":" + hashMap["realm"][1:len(hashMap["realm"])-1] + ":" + digest.Password
		A2 = method + ":" + digest.Uri
	}

	res = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%x", md5.Sum([]byte(A1)))+":"+hashMap["nonce"][1:len(hashMap["nonce"])-1]+":"+digest.Nc+":"+digest.Cnonce+":"+digest.Qop+":"+fmt.Sprintf("%x", md5.Sum([]byte(A2))))))
	// res = fmt.Sprintf("%x", md5.Sum([]byte(A1))) + ":" + hashMap["nonce"][1:len(hashMap["nonce"])-1] + ":" + digest.Nc + ":" + digest.Cnonce + ":" + digest.Qop + ":" + fmt.Sprintf("%x", md5.Sum([]byte(A2)))
	return
}

func handler(w http.ResponseWriter, r *http.Request) {

	URL := r.URL.Query().Get("url")
	if URL == "" {
		log.Println("missing URL argument")
		return
	}
	log.Println("visiting", URL)

	c := colly.NewCollector()

	p := &pageInfo{Links: make(map[string]int)}

	// extract status code
	c.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode)
		p.StatusCode = r.StatusCode
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("error:", r.StatusCode, err)

		if r.StatusCode == 401 {
			auth := r.Headers.Get("WWW-Authenticate")
			auths := strings.SplitN(auth, " ", 2)

			// fmt.Println(auths[1])
			//如果认证方式为Digest
			if auths[0] == "Digest" {

				str := strings.Split(auths[1], ",")
				hashMap := make(map[string]string)
				for _, v := range str {
					split := strings.Split(v, ",")

					for _, s := range split {
						i := strings.Split(s, "=")
						hashMap[i[0]] = i[1]
					}
				}

				for k, v := range hashMap {
					fmt.Println(k, ":", v)
				}

				// fmt.Println(digest)
				// var digest Digest

				//获得请求方法与uri
				method := r.Request.Method
				digest.Uri = r.Request.URL.RequestURI()
				digest.Qop = "auth"

				//产生客户端随机数与请求计数器
				digest.Cnonce = RandomString()
				digest.Nc = "0000001"

				//设置账号密码
				digest.UserName = "mengjiaheng"
				digest.Password = "123456"

				digest.Response = ResponseQop(method, hashMap, digest)
				c1 := colly.NewCollector()
				//第二次请求开始
				c1.OnRequest(func(r *colly.Request) {
					// Request头部设定
					r.Headers.Set("Authorization", `Digest username="`+digest.UserName+`",realm="`+hashMap["realm"][1:len(hashMap["realm"])-1]+`",qop=`+hashMap["qop"][1:len(hashMap["qop"])-1]+`,nonce="`+hashMap["nonce"][1:len(hashMap["nonce"])-1]+`",uri="`+digest.Uri+`",nc=`+digest.Nc+`,cnonce="`+digest.Cnonce+`",response="`+digest.Response+`"`)

				})
				// extract status code
				c1.OnResponse(func(r *colly.Response) {
					authTemp := r.Headers.Get("Authorization-Info")
					fmt.Println("第二次响应头：", authTemp)
					log.Println("response received", r.StatusCode)
					p.StatusCode = r.StatusCode
				})
				c1.OnError(func(r *colly.Response, err error) {
					log.Println("error:", r.StatusCode, err)
				})

				c1.Visit(URL)
				// dump results
				b, err := json.Marshal(p)
				if err != nil {
					log.Println("failed to serialize response:", err)
					return
				}
				w.Header().Add("Content-Type", "application/json")
				w.Write(b)
			}
		}
		p.StatusCode = r.StatusCode
	})

	c.Visit(URL)

	// dump results
	b, err := json.Marshal(p)
	if err != nil {
		log.Println("failed to serialize response:", err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func main() {
	// example usage: curl -s 'http://127.0.0.1:7171/?url=http://go-colly.org/'
	addr := ":7171"

	http.HandleFunc("/", handler)

	log.Println("listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
