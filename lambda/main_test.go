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
	myBucket := ""
	myStateS3Key := "state"

	// Setup
	ctx := context.Background()
	once.Do(func() {
		mySession := session.Must(session.NewSession(&aws.Config{Region: &myRegion}))
		var err error
		static, err = newAssetSet(mySession, myBucket, myStateS3Key)
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
