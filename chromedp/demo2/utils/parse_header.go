package utils

import (
	"net/http"
	"strings"
)

func ParseAuthorization(header http.Header) map[string]string {
	auth := header.Get("Www-Authenticate")

	//如果认证方式为Digest

	if strings.Contains(auth, "Digest") {

		auths := strings.SplitN(auth, " ", 2)
		hashMap := StringToMap(auths[1])

		return hashMap
	} else {
		return nil
	}
}
