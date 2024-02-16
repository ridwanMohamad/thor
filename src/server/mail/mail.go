package mail

import (
	"crypto/tls"
	"fmt"
	"github.com/domodwyer/mailyak/v3"
	"gopkg.in/gomail.v2"
	"net/smtp"
	"strconv"
	"thor/src/server/config"
)

type Mail struct {
	*mailyak.MailYak
}

type GoMailType struct {
	GMail    *gomail.Dialer
	MailFrom *string
}

func MailInitialize(conf config.Mail) *Mail {

	mail := mailyak.New(fmt.Sprintf("%s:%s", conf.Host, conf.Port), smtp.PlainAuth("", conf.Sender, conf.Password, conf.Host))

	mail.From(conf.Sender)
	mail.FromName(conf.SenderName)

	return &Mail{
		mail,
	}
}

func GoMailInitialize(conf config.Mail) *GoMailType {
	mailPort, _ := strconv.Atoi(conf.Port)
	mail := gomail.NewDialer(conf.Host, mailPort, conf.Sender, conf.Password)

	mail.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	return &GoMailType{GMail: mail, MailFrom: &conf.Sender}
}
