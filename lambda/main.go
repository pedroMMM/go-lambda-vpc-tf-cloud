package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
)

var once sync.Once
var static assetSet

func vpcCounter(ctx context.Context) (string, error) {
	once.Do(func() {
		var err error
		static, err = newAssetSet(newAssetSetInput{
			session:    session.Must(session.NewSession()),
			bucket:     os.Getenv("bucket_name"),
			states3Key: os.Getenv("state_s3_key"),
			emailFrom:  os.Getenv("email_from"),
			emailTo:    os.Getenv("email_to"),
		})
		if err != nil {
			panic(fmt.Errorf("newAssetSet failed with:\n%w", err))
		}
	})

	last, err := getState(ctx, static)
	if err != nil {
		return "", fmt.Errorf("getState failed with:\n%w", err)
	}

	current, err := getVpcs(ctx, static)
	if err != nil {
		return "", fmt.Errorf("getVpcs failed with:\n%w", err)
	}

	diff := calculateDiff(last, current)
	result := diff.ToString()

	err = sendEmail(ctx, static, diff)
	if err != nil {
		return "", fmt.Errorf("sendEmail failed with:\n%w", err)
	}

	err = setState(ctx, static, current)
	if err != nil {
		return "", fmt.Errorf("setState failed with:\n%w", err)
	}

	log.Printf("vpcCounter completed with: \n%s\n", result)
	return result, nil
}

func main() {
	lambda.Start(vpcCounter)
}
