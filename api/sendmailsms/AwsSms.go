package sendmailsms

import (
	"fmt"
	"os"

	twilio "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

// send sms With Veryfication Code to a phone number using Aws SNS
func SendSMSWithVeryficationCode(code string, number string) bool {

	// // fmt.Println("creating session")
	// sess := session.Must(session.NewSession(&aws.Config{
	// 	Region:      aws.String("us-east-1"),
	// 	Credentials: credentials.NewStaticCredentials(os.Getenv("AWSSNSCREDID"), os.Getenv("AWSSNSCREDSECRET"), ""),
	// }))
	// // fmt.Println("session created")

	// svc := sns.New(sess)
	// // fmt.Println("service created")
	// sededmsg := "dancee Account Veryfication Code | " + code
	// params := &sns.PublishInput{
	// 	Subject:     aws.String("Veryfiy Your Account"),
	// 	Message:     aws.String(sededmsg),
	// 	PhoneNumber: aws.String(number),
	// }
	// resp, err := svc.Publish(params)

	// if err != nil {
	// 	// Print the error, cast err to awserr.Error to get the Code and
	// 	// Message from an error.
	// 	fmt.Println(err.Error())
	// 	return false
	// }
	// // fmt.Println(resp)

	// // Pretty-print the response data.
	// if resp.MessageId != nil {
	// 	return true
	// } else {
	// 	return false
	// } //

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: os.Getenv("TWILLO_USERNAME"),
		Password: os.Getenv("TWILLO_PASSWORD"),
	})

	CodeMessage := "your code is : " + code

	params := &openapi.CreateMessageParams{}
	params.SetTo(number)
	params.SetFrom(os.Getenv("TWILLO_PHONE"))
	params.SetBody(CodeMessage)
	_, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else {
		fmt.Println("SMS sent successfully!")
		return true
	}

}
