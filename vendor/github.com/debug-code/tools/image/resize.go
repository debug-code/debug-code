package image

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"strings"

	"github.com/nfnt/resize"
	"golang.org/x/image/bmp"
)

func getByteFormat(bytes []byte) string {
	if len(bytes) < 4 {
		return ""
	}
	if bytes[0] == 0x89 && bytes[1] == 0x50 && bytes[2] == 0x4E && bytes[3] == 0x47 {
		return "png"
	}
	if bytes[0] == 0xFF && bytes[1] == 0xD8 {
		return "jpg"
	}
	if bytes[0] == 0x47 && bytes[1] == 0x49 && bytes[2] == 0x46 && bytes[3] == 0x38 {
		return "gif"
	}
	if bytes[0] == 0x42 && bytes[1] == 0x4D {
		return "bmp"
	}
	return ""
}

// FromBase64 从base64中读取图片
func FromBase64(src string) (image.Image, error) {
	bs, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return nil, err
	}
	io := bytes.NewReader(bs)
	t := getByteFormat(bs)
	switch t {
	case "png":
		return png.Decode(io)
	case "jpg":
		return jpeg.Decode(io)
	case "gif":
		return gif.Decode(io)
	case "bmp":
		return bmp.Decode(io)
	default:
		return nil, errors.New("no format image")
	}
}

func resizeBase64Src(src string, newWidth int) (string, error) {
	n := strings.Index(src, "base64,")
	if n > -1 {
		src = string(src[n+7:])
	}
	img, err := FromBase64(src)
	if err != nil {
		return "", err
	}
	width := img.Bounds().Dy()
	if width > newWidth {
		img = resize.Resize(uint(newWidth), 0, img, resize.Lanczos3)
	}
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, img, nil)
	if err != nil {
		return "", err
	}
	ioutil.WriteFile(fmt.Sprintf("D:/test_%d.jpg", width), buf.Bytes(), 667)
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// GenBase64ThumbSrc 生成缩略图
func GenBase64ThumbSrc(src string) (string, error) {
	return resizeBase64Src(src, 150)
}

// ResizeBase64Src 修改大小
//data:image/jpeg;base64,
func ResizeBase64Src(src string) (string, error) {
	return resizeBase64Src(src, 1000)
}

// Base64Src ...
func Base64Src(src string) string {
	n := strings.Index(src, "base64,")
	if n > -1 {
		return src
	}
	return "data:image/jpeg;base64," + src
}
