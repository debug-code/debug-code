package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

func SendEmail(to []string, from, theme, body string) error {
	auth := smtp.PlainAuth("", "1245800723@qq.com", "bncmwhiwjnvcifjg", "smtp.qq.com")
	//to := []string{"804070835@qq.com"}
	//nickname := "test"
	user := "1245800723@qq.com"
	//subject := "test mail"
	content_type := "Content-Type: text/html; charset=UTF-8"
	//body := "This is the email body."
	msg := []byte("To: " + strings.Join(to, ",") + "\r\nFrom: " +
		from + "<" + user + ">\r\nSubject: " + theme + "\r\n" +
		content_type + "\r\n\r\n" + body)
	err := smtp.SendMail("smtp.qq.com:25", auth, user, to, msg)
	if err != nil {
		fmt.Printf("send mail error: %v", err)
		return err
	}
	return nil
}
