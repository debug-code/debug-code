package httpRequest

import (
	"bytes"
	"fmt"

	"io/ioutil"
	"net/http"
	"net/url"
)

func Post(link string, body []byte, paramMap ...map[string]string) (string, error) {
	var urlPath string
	var err error
	urlPath, err = setUrl(link, paramMap...)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", urlPath, bytes.NewBuffer(body)) // 新建请求
	if err != nil {
		fmt.Println("NewRequest", err)
		return "", err
	}

	if len(paramMap) > 1 {
		for k, v := range paramMap[1] {
			req.Header.Set(k, v)
		}
	}

	//beego.Info(req.Header)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	//beego.Info("post resp", resp)

	defer resp.Body.Close()

	bodys, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodys), nil
}

func Put(link string, body []byte, paramMap ...map[string]string) (string, error) {
	var urlPath string
	var err error
	urlPath, err = setUrl(link, paramMap...)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("PUT", urlPath, bytes.NewBuffer(body)) // 新建请求
	if err != nil {
		fmt.Println("NewRequest", err)
		return "", err
	}

	if len(paramMap) > 1 {
		for k, v := range paramMap[1] {
			req.Header.Set(k, v)
		}
	}

	//beego.Info(req.Header)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	//beego.Info("post resp", resp)

	defer resp.Body.Close()

	bodys, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodys), nil
}

func Get(link string, paramMap ...map[string]string) (string, error) {
	var urlPath string
	var err error
	//client := &http.Client{}
	urlPath, err = setUrl(link, paramMap...)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("GET", urlPath, nil) // 新建请求
	if err != nil {
		fmt.Println("err1", err)
		return "", err
	}
	if len(paramMap) > 1 {
		//urlPath, err = setUrl(link, paramMap[1])
		for k, v := range paramMap[1] {
			req.Header.Set(k, v)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	//resp, err := client.Do(req)
	if err != nil { 
		return "", err
	}
	defer resp.Body.Close()

	bodys, _ := ioutil.ReadAll(resp.Body)

	return string(bodys), nil
}

func Delete(link string, paramMap ...map[string]string) (string, error) {
	var urlPath string
	var err error
	//client := &http.Client{}
	urlPath, err = setUrl(link, paramMap...)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("DELETE", urlPath, nil) // 新建请求
	if err != nil {
		fmt.Println("err1", err)
		return "", err
	}
	if len(paramMap) > 1 {
		//urlPath, err = setUrl(link, paramMap[1])
		for k, v := range paramMap[1] {
			req.Header.Set(k, v)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	//resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodys, _ := ioutil.ReadAll(resp.Body)

	return string(bodys), nil
}

func GetEasy(link string, paramMap ...map[string]string) (string, error) {
	urlPath, err := setUrl(link, paramMap...)
	if err != nil {
		return "", err
	}

	resp, err := http.Get(urlPath)
	if err != nil {
		return "", err
	}
	//defer resp.Body.Close()
	s, _ := ioutil.ReadAll(resp.Body)

	return string(s), nil
}

func setUrl(link string, paramMap ...map[string]string) (string, error) {
	params := url.Values{}
	if len(paramMap) > 0 {
		for k, v := range paramMap[0] {
			params.Set(k, v)
		}
	}

	Url, err := url.Parse(link)
	if err != nil {
		return "", err

	}
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	return urlPath, nil
}
