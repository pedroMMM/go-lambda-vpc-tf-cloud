package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/ses"
)

func sendEmail(ctx context.Context, static assetSet, diff diff) error {
	charset := "UTF-8"
	message := diff.ToString()

	if diff.NoChanges() {
		return nil
	}

	result, err := static.ses.SendEmailWithContext(ctx, &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{&static.emailTo},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: &charset,
					Data:    &message,
				},
				Text: &ses.Content{
					Charset: &charset,
					Data:    &message,
				},
			},
			Subject: &ses.Content{
				Charset: &charset,
				Data:    aws.String("VPC Count Results"),
			},
		},
		ReturnPath: &static.emailFrom,
		Source:     &static.emailFrom,
	})
	if err != nil {
		return err
	}

	log.Printf("email sent from %s to %s via SES with Id: %s\n", static.emailFrom, static.emailTo, *result.MessageId)
	return nil
}
