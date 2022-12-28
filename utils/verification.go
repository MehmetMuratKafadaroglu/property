package utils

import (
	"fmt"
	"go-rest/models"
	"go-rest/settings"
	"math/rand"
	"net/smtp"
	"time"
)

const (
	from     = ""
	password = ""
	smtpHost = "smtp.gmail.com"
	smtpPort = "587"
)

var EmailQueue = models.NewQueue()

func MailSender() {
	for { //infinite loop
		if EmailQueue.IsEmpty() {
			time.Sleep(10000000000) // 10 seconds
		} else {
			time.Sleep(2000000000) // 2 seconds
			sendMail(EmailQueue.Pop())
		}

	}

}
func sendMail(address string, url string) {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(`
	<html><body>
		<h1>This is mail a from %s</h1>
		<p>Do not send this email to anyone else </p></br>
		<a href=%q> Please click here to verify your email</a>
	</body></html>`,
		settings.CompanyName, url)
	msg := []byte("From: " + from + "\r\n" +
		"To: " + address + "\r\n" +
		"Subject: Email Verification!\r\n" +
		mime +
		"\r\n" + body)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	temp := []string{address}
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, temp, []byte(msg))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Init() {
	rand.Seed(time.Now().UnixNano())
}

func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
