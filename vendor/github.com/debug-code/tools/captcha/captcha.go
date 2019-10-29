package captcha

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/lifei6671/gocaptcha"
	"image/jpeg"
	"strconv"
	"time"
	"wechat-go/tools"
)

const ASCII_LOWERCASE = "abcdefghjkmnpqrstuvwxyz"
const ASCII_UPPERCASE = "ABCDEFGHJKMNPQRSTUVWXYZ"
const ASCII_LETTERS = ASCII_LOWERCASE + ASCII_UPPERCASE
const DIGITS = "23456789"

type Captcha struct {
	Id     string
	Key    string
	Base64 string
}

func GetCapcha() Captcha {

	C := Captcha{}

	C.Key = tools.GetRandomStrWithSrc(6, DIGITS, ASCII_LETTERS)

	err := gocaptcha.ReadFonts("tools/captcha/fonts", ".ttf")
	if err != nil {
		fmt.Println(err)
		return C
	}
	//初始化一个验证码对象
	captchaImage := gocaptcha.NewCaptchaImage(240, 80, gocaptcha.RandLightColor())

	//画上三条随机直线
	captchaImage.DrawLine(3)

	//画边框
	captchaImage.DrawBorder(gocaptcha.ColorToRGB(0x17A7A7A))

	//画随机噪点
	captchaImage.DrawNoise(gocaptcha.CaptchaComplexHigh)

	//画随机文字噪点
	captchaImage.DrawTextNoise(gocaptcha.CaptchaComplexLower)

	captchaImage.DrawText(C.Key)

	img := captchaImage.GetImage()

	buffer := new(bytes.Buffer)
	//
	//ww := bufio.NewWriter(w)
	err = jpeg.Encode(buffer, img, nil)
	if err != nil {
		fmt.Println(err)
	}

	Base64 := base64.StdEncoding.EncodeToString(buffer.Bytes())
	if err != nil {
		fmt.Println(err)
		return C
	}

	C.Base64 = "data:image/jpg;base64," + Base64

	C.Id = tools.StringToSha256(C.Base64 + strconv.Itoa(int(time.Now().Unix())))

	return C

}
