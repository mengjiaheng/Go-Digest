package utils

import (
	"crypto/md5"
	"demo1/core"
	"fmt"
)

func ResponseQop(hashMap map[string]string, digest core.Digest) (res string) {

	var A1, A2 string
	if len(digest.Qop) != 0 && digest.Qop[1:len(digest.Qop)-1] == "auth" {
		A1 = digest.UserName + ":" + hashMap["realm"][1:len(hashMap["realm"])-1] + ":" + digest.Password
		A2 = digest.Method + ":" + digest.Uri
	}

	// res = fmt.Sprintf("%x", md5.Sum([]byte(A1))) + ":" + hashMap["nonce"][1:len(hashMap["nonce"])-1] + ":" + digest.Nc + ":" + digest.Cnonce + ":" + digest.Qop + ":" + fmt.Sprintf("%x", md5.Sum([]byte(A2)))
	res = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%x", md5.Sum([]byte(A1)))+":"+hashMap["nonce"][1:len(hashMap["nonce"])-1]+":"+digest.Nc+":"+digest.Cnonce+":"+digest.Qop[1:len(digest.Qop)-1]+":"+fmt.Sprintf("%x", md5.Sum([]byte(A2))))))
	return
}
