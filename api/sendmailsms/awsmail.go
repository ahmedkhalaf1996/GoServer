package sendmailsms

import (
	"bytes"
	"fmt"
	"os"

	// "embed"
	// "fmt"
	"html/template"

	"time"

	"github.com/go-mail/mail/v2"
)

// //go:embed "templates"
// var templateFS embed.FS

type MailerConfig struct {
	Timeout      time.Duration
	Host         string
	Port         int
	Username     string
	Password     string
	Sender       string
	TemplatePath string
}

type Mailer struct {
	dailer *mail.Dialer
	config MailerConfig
	sender string
}

//
var Config = MailerConfig{
	Host:     os.Getenv("AWSMAILHOST"),
	Port:     587,
	Username: os.Getenv("AWSMAILUSERNAME"),
	Password: os.Getenv("AWSMAILPASSWORD"),
	Timeout:  30 * time.Second,
	Sender:   "danccee@danccee.com",
	// Sender: "mockacademy666@gmail.com",
}

func SendEmailWithVeryficationCodeToEmail(IncomeCode string, Email string) bool {

	dailer := mail.NewDialer(Config.Host, Config.Port, Config.Username, Config.Password)
	dailer.Timeout = Config.Timeout

	m := Mailer{
		dailer: dailer,
		sender: Config.Sender,
		config: Config,
	}

	var body bytes.Buffer
	t, _ := template.ParseFiles("./api/sendmailsms/mail.html")
	// api\sendmailsms\mail.html
	t.Execute(&body, struct{ Name string }{Name: IncomeCode})

	Mess := body.String()

	msg := mail.NewMessage()
	msg.SetHeader("To", Email)
	msg.SetHeader("Subject", "VeryiFication Code")
	msg.SetHeader("From", m.sender)
	msg.SetBody("text/plain", IncomeCode)
	msg.AddAlternative("text/html", Mess)

	res := m.dailer.DialAndSend(msg)

	if res != nil {
		fmt.Println("err", res)
		return false
	} else {
		fmt.Println("email send")
		return true

	}

}
