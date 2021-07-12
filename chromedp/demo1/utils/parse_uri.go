package utils

import "strings"

//解析uri
func ParseURI(URI string) (uri string) {

	str := strings.Split(URI, "/")
	for i := 3; 2 < i && i < len(str); i++ {
		uri = uri + "/" + str[i]
	}
	return
}
