package main

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/aws/session"
)

// Test_vpcCounter runs local testing for the lambda handler
func Test_vpcCounter(t *testing.T) {
	t.Skip("used for local testing, you will need actual AWS credentials and to manually input for configuration")

	// Local configuration
	myRegion := "us-east-1"
	myBucket := "527158362817-go-lambda-vpc-tf-cloud"
	myStateS3Key := "state"
	myEmailFrom := "noreply@pmedeiros.dev"
	myEmailTo := "mmedeirospedrom@gmail.com"

	// Setup
	ctx := context.Background()
	once.Do(func() {
		var err error
		static, err = newAssetSet(newAssetSetInput{
			session:    session.Must(session.NewSession(&aws.Config{Region: &myRegion})),
			bucket:     myBucket,
			states3Key: myStateS3Key,
			emailFrom:  myEmailFrom,
			emailTo:    myEmailTo,
		})
		if err != nil {
			t.Errorf("newAssetSet() error = %v", err)
		}
	})

	got, err := vpcCounter(ctx)
	if err != nil {
		t.Errorf("vpcCounter() error = %v", err)
	}
	if got != "no changes" {
		t.Errorf("vpcCounter() got = %v, want %v", got, "no changes")
	}
}
