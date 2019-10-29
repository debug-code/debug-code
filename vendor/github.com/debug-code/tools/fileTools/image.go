package fileTools

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"
	"wechat-go/tools"
)

func Base64ToImage(img, toimg string) {

	ddd, _ := base64.StdEncoding.DecodeString(img)
	fmt.Println(string(ddd))
	saveTitle := tools.GetRandomStr(8, true, true, true)
	saveTitle += strconv.Itoa(int(time.Now().Unix()))

	err := ioutil.WriteFile(toimg, ddd, 0666)
	if err != nil {
		fmt.Println(err)
	}

}
