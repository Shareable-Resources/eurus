package user

import (
	"eurus-backend/foundation/log"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

func bodyBuilder(userEmail, code string, template string) string {
	template = strings.ReplaceAll(template, "<!--mailTo-->", userEmail)
	template = strings.ReplaceAll(template, "<!--verificationCode-->", code)
	return template
}

func SendEmail(config *UserServerConfig, recipient string, subject string, code string, template string) {
	var verbose bool = true
	sess, err := session.NewSession(&aws.Config{
		Region:                        aws.String(config.EmailServiceZone),
		Credentials:                   credentials.NewStaticCredentials(config.AwsAccessKeyId, config.AwsAccessSecretAccessKey, ""),
		CredentialsChainVerboseErrors: &verbose,
	},
	)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to create AWS session: ", err, " receipient email: ", recipient)
		return
	}

	// Create an SES session.
	svc := ses.New(sess)

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(bodyBuilder(recipient, code, template)),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(config.EmailFrom),
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	result, err := svc.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				log.GetLogger(log.Name.Root).Error(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				log.GetLogger(log.Name.Root).Error(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				log.GetLogger(log.Name.Root).Error(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				log.GetLogger(log.Name.Root).Error(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.GetLogger(log.Name.Root).Error(err.Error())
		}

		return
	}

	fmt.Println(result)
}
