package utils

import (
	"crypto/md5"
	"fmt"
)

// Md5
// @param str 字符串
func Md5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}
