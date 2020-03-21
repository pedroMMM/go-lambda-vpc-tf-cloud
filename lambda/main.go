package main

import (
	"context"
	"fmt"
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

		mySession := session.Must(session.NewSession())

		static, err = newAssetSet(mySession, os.Getenv("bucket_name"), os.Getenv("state_s3_key"))
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

	err = setState(ctx, static, current)
	if err != nil {
		return "", fmt.Errorf("setState failed with:\n%w", err)
	}

	return diff.ToString(), nil
}

func main() {
	lambda.Start(vpcCounter)
}
