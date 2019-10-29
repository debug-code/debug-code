package tools

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"git.xrewin.com/go/beego"
	"github.com/boltdb/bolt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const ASCII_LOWERCASE = "abcdefghijklmnopqrstuvwxyz"
const ASCII_UPPERCASE = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const ASCII_LETTERS = ASCII_LOWERCASE + ASCII_UPPERCASE
const DIGITS = "0123456789"

//将字符串数组拼接成一个字符串
func ArrToString(arr []string) string {
	str := strings.Trim(fmt.Sprint(arr), "[]")
	str = strings.Replace(str, " ", "", -1)
	return str
}

//将[]byte 转换成字符  例子：\x34\xaa
func BytesToHexString(bt []byte, sign bool) string {
	var re string
	if sign {
		for k := range bt {
			//fmt.Println(k)
			re += `\x` + hex.EncodeToString(bt[k:k+1])
		}
	} else {
		for k := range bt {
			//fmt.Println(k)
			re += hex.EncodeToString(bt[k : k+1])
		}
	}
	return re
}

func HexStringToBytes(hs string, sign bool) ([]byte, error) {
	var str string
	if sign {
		//fmt.Println(hs)
		str = strings.Replace(hs, `\x`, "", -1)
		//fmt.Println(str)

	} else {
		str = hs
	}
	re, err := hex.DecodeString(str)
	if err != nil {
		//fmt.Println(err)
		return nil, err
	}
	return re, nil
}

//？？？表示将32位的主机字节顺序转化为32位的网络字节顺序
//int to byte
func SockHtonl(input int) []byte {

	reBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(reBytes, uint32(input))

	return reBytes
}

//byte to int
func SockNonhl(inputBytes []byte) uint32 {

	//reBytes := make([]byte, 4)
	reInt := binary.BigEndian.Uint32(inputBytes)

	return reInt
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
func GetRandomStrWithSrc(num int, args ...string) string {
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

//补位
func Pkcs7Encode(text []byte, blockSize int) []byte {
	//计算需要填充的位数
	needSize := blockSize - (len(text) % blockSize)
	if needSize == 0 {
		needSize = blockSize
	}
	//获得补位字符
	//Repeat()函数的功能是把切片b复制count个,然后合成一个新的字节切片返回.
	//func Repeat(b[]byte,count int) []byte
	pad := bytes.Repeat([]byte{byte(needSize)}, needSize)
	return append(text, pad...)
}

//
func Pkcs7Decode(text []byte) []byte {

	length := len(text)
	unpadding := int(text[length-1])

	return text[:(length - unpadding)]
}

//AES 加密
func AesEncrypt(encodeStr string, key []byte) (string, error) {
	encodeBytes := []byte(encodeStr)
	//根据key 生成密文
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	encodeBytes = Pkcs7Encode(encodeBytes, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:16])
	crypted := make([]byte, len(encodeBytes))
	blockMode.CryptBlocks(crypted, encodeBytes)
	// 使用BASE64对加密后的字符串进行编码
	//fmt.Println(crypted)
	//base64.StdEncoding.EncodeToString(crypted)
	return base64.StdEncoding.EncodeToString(crypted), nil
}

func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)

	//fmt.Println("削减")
	//fmt.Println(origData)
	//fmt.Println(string(origData))
	origData = Pkcs7Decode(origData)
	return origData, nil

}

//bolt

func Save(key string, val string) error {
	db, err := bolt.Open("mydb.db", 0600, nil)
	if err != nil {
		fmt.Println("open", err)
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		//b := tx.Bucket([]byte("test"))
		b, err := tx.CreateBucketIfNotExists([]byte("test"))
		//dd := 1
		err = b.Put([]byte(key), []byte(val))

		return err

	})
	return err
}

func GetKey(key string) (string, error) {
	db, err := bolt.Open("mydb.db", 0600, nil)
	if err != nil {
		fmt.Println("open", err)
	}
	defer db.Close()
	var val string
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("test"))
		val = string(b.Get([]byte(key)))
		return err

	})
	return val, err
}

//emoji 表情处理
//解码
func UnicodeEmojiDecode(s string) string {
	//emoji表情的数据表达式
	re := regexp.MustCompile("\\[[\\\\u0-9a-zA-Z]+\\]")
	//提取emoji数据表达式
	reg := regexp.MustCompile("\\[\\\\u|]")
	src := re.FindAllString(s, -1)
	for i := 0; i < len(src); i++ {
		e := reg.ReplaceAllString(src[i], "")
		p, err := strconv.ParseInt(e, 16, 32)
		if err == nil {
			s = strings.Replace(s, src[i], string(rune(p)), -1)
		}
	}
	return s
}

//转码
func UnicodeEmojiCode(s string) string {
	ret := ""
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		if len(string(rs[i])) == 4 {
			u := `[\u` + strconv.FormatInt(int64(rs[i]), 16) + `]`
			ret += u

		} else {
			ret += string(rs[i])
		}
	}
	return ret
}

func GetFileSize(filename string) int64 {
	var result int64
	filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	return result
}

//Post file

//func PostFile(filename string, targetUrl string) (string, error) {
//	bodyBuf := &bytes.Buffer{}
//	bodyWriter := multipart.NewWriter(bodyBuf)
//
//	//关键的一步操作
//	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
//	if err != nil {
//		fmt.Println("error writing to buffer")
//		return "", err
//	}
//
//	//打开文件句柄操作
//	fh, err := os.Open(filename)
//	if err != nil {
//		fmt.Println("error opening file")
//		return "", err
//	}
//	defer fh.Close()
//
//	//iocopy
//	_, err = io.Copy(fileWriter, fh)
//	if err != nil {
//		return "", err
//	}
//
//	contentType := bodyWriter.FormDataContentType()
//	bodyWriter.Close()
//
//	resp, err := http.Post(targetUrl, contentType, bodyBuf)
//	if err != nil {
//		return "", err
//	}
//	defer resp.Body.Close()
//	resp_body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return "", err
//	}
//	return string(resp_body), nil
//}

func PostFile(link string, params map[string]string, name, path string) (*http.Request, error) {
	fp, err := os.Open(path) // 打开文件句柄
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	body := &bytes.Buffer{}                                       // 初始化body参数
	writer := multipart.NewWriter(body)                           // 实例化multipart
	part, err := writer.CreateFormFile(name, filepath.Base(path)) // 创建multipart 文件字段
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, fp) // 写入文件数据到multipart
	for key, val := range params {
		_ = writer.WriteField(key, val) // 写入body中额外参数，比如七牛上传时需要提供token
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", link, body) // 新建请求
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "multipart/form-data") // 设置请求头,!!!非常重要，否则远端无法识别请求
	return req, nil
}

//cmd
func CatPic(from string, to string) {
	//args := []string{"-ss", "00:00:02", "-i", from, to, "-s", w + "x" + h}
	args := []string{"-ss", "00:00:02", "-i", from, to}
	Lcmd("ffmpeg", args)
}
func Lcmd(command string, args []string) string {
	//args:=[]string{"network","ls"}
	out, err := exec.Command(command, args...).Output()
	if err != nil {
		beego.Error(err)
	}
	//fmt.Printf("%s",string(out))
	//fmt.Println(string(out))

	return string(out)
}

//str to json
func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

//percentage of two int

func Percentage(a, b int) string {
	if b == 0 {
		return "--"
	}

	c := float64(a) / float64(b) * 100

	return strconv.Itoa(int(c)) + "%"

}

func Percent(a, b int) string {
	res := "--"
	if b == 0 || a == b {
		return res
	}

	if a > b {
		res = strconv.Itoa(int(float64(b-a)/float64(b)*100)) + "%"
	} else {
		res = strconv.Itoa(int(float64(a-b)/float64(b)*100)) + "%"
	}
	return res
}

func Percent0(a, b int) string {
	res := "0%"
	if b == 0 || a == b {
		return res
	}
	res = strconv.Itoa(int(float64(a)/float64(b)*100)) + "%"
	return res
}

func Percent1(a, b int) string {
	res := "0"
	if b == 0 || a == b {
		return res
	}
	res = strconv.Itoa(int(float64(a) / float64(b) * 100))
	return res
}
func TimeSub(t1, t2 time.Time) int {
	t1 = t1.UTC().Truncate(24 * time.Hour)
	t2 = t2.UTC().Truncate(24 * time.Hour)
	return int(t1.Sub(t2).Hours() / 24)
}

func SortSIMapByvalue(maps map[string]int, sort bool) (keys []string, values []int) {

	for k, v := range maps {
		keys = append(keys, k)
		values = append(values, v)
	}

	arrylength := len(values)
	for i := 0; i < arrylength; i++ {
		min := i
		for j := i + 1; j < arrylength; j++ {
			if sort {
				if values[j] < values[min] {
					min = j
				}
			} else {
				if values[j] > values[min] {
					min = j
				}
			}

		}
		t := values[i]
		values[i] = values[min]
		values[min] = t

		tt := keys[i]
		keys[i] = keys[min]
		keys[min] = tt

	}
	return

}

func DivisionDIntTOInt(a, b int) int {
	if b == 0 || a == 0 {
		return 0
	}
	return int(float64(a) / float64(b))
}

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
