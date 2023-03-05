package sendmailsms

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

// send sms With Veryfication Code to a phone number using Aws SNS
func SendSMSWithVeryficationCode(code string, number string) bool {

	// fmt.Println("creating session")
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWSSNSCREDID"), os.Getenv("AWSSNSCREDSECRET"), ""),
	}))
	// fmt.Println("session created")

	svc := sns.New(sess)
	// fmt.Println("service created")
	sededmsg := "dancee Account Veryfication Code | " + code
	params := &sns.PublishInput{
		Subject:     aws.String("Veryfiy Your Account"),
		Message:     aws.String(sededmsg),
		PhoneNumber: aws.String(number),
	}
	resp, err := svc.Publish(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return false
	}
	// fmt.Println(resp)

	// Pretty-print the response data.
	if resp.MessageId != nil {
		return true
	} else {
		return false
	} //
}
