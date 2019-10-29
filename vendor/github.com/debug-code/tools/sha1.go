package tools

import (
	"crypto/sha1"
	"encoding/hex"
)


//sha加密
func SHA1(str string)(string, error){


	// 产生一个散列值得方式是 sha1.New()，sha1.Write(bytes)，
	// 然后 sha1.Sum([]byte{})。这里我们从一个新的散列开始。
	h := sha1.New()
	// 写入要处理的字节。如果是一个字符串，
	// 需要使用[]byte(s) 来强制转换成字节数组。
	_, err := h.Write([]byte(str))
	if err != nil{
		return "", err
	}
	// 这个用来得到最终的散列值的字符切片。
	// Sum 的参数可以用来都现有的字符切片追加额外的字节切片：一般不需要要。
	bs := h.Sum(nil)
	return hex.EncodeToString(bs[:]), nil
}