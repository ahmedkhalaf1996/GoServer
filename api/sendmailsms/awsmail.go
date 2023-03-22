package sendmailsms

import (
	"fmt"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	// "embed"
	// "fmt"
	// "time"
)

// //go:embed "templates"
// var templateFS embed.FS

// type MailerConfig struct {
// 	Timeout      time.Duration
// 	Host         string
// 	Port         int
// 	Username     string
// 	Password     string
// 	Sender       string
// 	TemplatePath string
// }

// type Mailer struct {
// 	dailer *mail.Dialer
// 	config MailerConfig
// 	sender string
// }

// //
// var Config = MailerConfig{
// 	Host:     os.Getenv("AWSMAILHOST"),
// 	Port:     587,
// 	Username: os.Getenv("AWSMAILUSERNAME"),
// 	Password: os.Getenv("AWSMAILPASSWORD"),
// 	Timeout:  30 * time.Second,
// 	Sender:   "danccee@danccee.com",
// 	// Sender: "mockacademy666@gmail.com",
// }

func SendEmailWithVeryficationCodeToEmail(IncomeCode string, Email string) bool {

	// dailer := mail.NewDialer(Config.Host, Config.Port, Config.Username, Config.Password)
	// dailer.Timeout = Config.Timeout

	// m := Mailer{
	// 	dailer: dailer,
	// 	sender: Config.Sender,
	// 	config: Config,
	// }

	// var body bytes.Buffer
	// t, _ := template.ParseFiles("./api/sendmailsms/mail.html")
	// // api\sendmailsms\mail.html
	// t.Execute(&body, struct{ Name string }{Name: IncomeCode})

	// Mess := body.String()

	// msg := mail.NewMessage()
	// msg.SetHeader("To", Email)
	// msg.SetHeader("Subject", "VeryiFication Code")
	// msg.SetHeader("From", m.sender)
	// msg.SetBody("text/plain", IncomeCode)
	// msg.AddAlternative("text/html", Mess)

	// res := m.dailer.DialAndSend(msg)

	// if res != nil {
	// 	fmt.Println("err", res)
	// 	return false
	// } else {
	// 	fmt.Println("email send")
	// 	return true

	// }

	// send mail ------------

	EditedCode := "Your Verfication Code is " + "<h1><strong>" + IncomeCode + "</strong></h1>"
	from := mail.NewEmail(os.Getenv("MAIL_SEND_FROM_NAME"), os.Getenv("MAIL_SEND_FROM_ADDRESS"))
	subject := "Veryfication Code From GolfScore"
	to := mail.NewEmail("", Email)
	plainTextContent := "VeryFiCation Code"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, EditedCode)
	client := sendgrid.NewSendClient(os.Getenv("MAIL_SEND_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		fmt.Println("Unable to send your email")
		log.Fatal(err)
		return false
	}
	// Check if it was sent
	statusCode := response.StatusCode
	if statusCode == 200 || statusCode == 201 || statusCode == 202 {
		fmt.Println("Email sent!")
		return true
	}

	return false
	// ----------------------
}
