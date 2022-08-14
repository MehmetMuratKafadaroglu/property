package utils

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"time"
)

const (
	from     = "mailsender@adreact.com"
	password = "gheb4gsiPPWC"
	smtpHost = "smtp.gmail.com"
	smtpPort = "587"
)

func SendMail(address string, message string) {
	auth := smtp.PlainAuth("", from, password, smtpHost)
	temp := []string{address}
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, temp, []byte(message))
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
