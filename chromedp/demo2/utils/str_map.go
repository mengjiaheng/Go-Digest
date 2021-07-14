package utils

import "strings"

func StringToMap(str string) map[string]string {
	sli := strings.Split(str, ",")
	hashMap := make(map[string]string)
	for _, v := range sli {
		split := strings.Split(v, ",")

		for _, s := range split {
			i := strings.Split(s, "=")
			hashMap[i[0]] = i[1]
		}
	}
	return hashMap
}
