package utils

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"net/smtp"
	"win-a-thon/config"
)

func Notify(toEmail string, subject string, content string) (err error) {

	appConfig := config.GetConfig()
	if _, err := toml.DecodeFile("config/env.default.toml", &appConfig); err != nil {
		fmt.Println(err)
		return err
	}

	from := appConfig.Application.Email
	password := appConfig.Application.Password

	//receiver details
	to := []string{toEmail}

	//smtp
	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port

	//message
	subject = "Subject : " + subject + "\r\n"
	body := "Dear User, \r\n \r\n" + content + "\r\n \r\nThanks & regards\r\nTeam Winathon"
	message := []byte(subject + body)

	auth := smtp.PlainAuth("", from, password, host)

	err = smtp.SendMail(address, auth, from, to, message)

	if err != nil {
		fmt.Println(err)
		return err
	} else {
		fmt.Println("mail sent successfully!")
	}
	return nil
}
