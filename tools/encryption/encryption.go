package encryption

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"
)

const ASCII_LOWERCASE = "abcdefghijklmnopqrstuvwxyz"
const ASCII_UPPERCASE = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const ASCII_LETTERS = ASCII_LOWERCASE + ASCII_UPPERCASE
const DIGITS = "0123456789"

func StringToSha256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}
func StringToMd5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

//生成随机字符串范围包括大小写字母，数字
func GetRandomStr(num int, lower bool, upper bool, digits bool) string {
	var re string
	var list string
	if lower {
		list += ASCII_LOWERCASE
	}
	if upper {
		list += ASCII_UPPERCASE
	}
	if digits {
		list += DIGITS
	}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < num; i++ {
		index := rand.Intn(len(list) - 1)
		re += string(list[index])
	}

	return re
}
func GetRandomStrWithStrings(num int, args ...string) string {
	var re string
	var list string

	if len(args) > 0 {
		for _, v := range args {
			list += v
		}
	}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < num; i++ {
		index := rand.Intn(len(list) - 1)
		re += string(list[index])
	}

	return re
}
