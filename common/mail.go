package common

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/mail"
	"net/smtp"
	"strings"

	"blog/conf"
)

var (
	auth      smtp.Auth
	config    conf.MailConfig
	state     bool
	tlsConfig *tls.Config
	addr      string
)

func encodeRFC2047(String string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{Name: String, Address: ""}
	return strings.Trim(addr.String(), " <@>")
}

func SendMail(ctx *gin.Context, receiver string, title string, msg string) {

	if !state {
		log.Println("connection state is invalid, skipping the mail send operation")
	}
	// TODO filter mail list

	header := make(map[string]string)
	header["From"] = config.From
	header["To"] = receiver
	header["Subject"] = encodeRFC2047(title)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(msg))

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		_ = ctx.Error(errors.New("failed to connect to mail server skipping the send")).SetType(gin.ErrorTypePrivate)
	}

	c, err := smtp.NewClient(conn, config.Host)
	if err != nil {
		log.Panic(err)
	}
	// Auth
	if err = c.Auth(auth); err != nil {
		log.Panic(err)
	}
	// To && From
	if err = c.Mail(config.From); err != nil {
		log.Panic(err)
	}

	if err = c.Rcpt(receiver); err != nil {
		log.Panic(err)
	}
	// Data
	w, err := c.Data()
	if err != nil {
		log.Panic(err)
	}
	_, err = w.Write([]byte(message))
	if err != nil {
		log.Panic(err)
	}
	err = w.Close()
	if err != nil {
		log.Panic(err)
	}
	err = c.Quit()
	if err != nil {
		log.Printf("[%v] failed to send mail to %s", ctx.Value("traceId").(string), receiver)
	}
}

func init() {
	config = conf.Config.Mail
	auth = smtp.PlainAuth("", config.Username, config.Password, config.Host)
	// TLS config
	tlsConfig = &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         config.Host,
	}

	addr = fmt.Sprintf("%s:%d", config.Host, config.Port)
}
