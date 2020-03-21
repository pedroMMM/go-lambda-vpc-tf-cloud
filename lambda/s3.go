package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

func getState(ctx context.Context, static assetSet) ([]string, error) {
	out := make([]string, 0)

	result, err := static.s3.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: &static.bucket,
		Key:    &static.stateKey,
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == s3.ErrCodeNoSuchKey {
			return out, nil
		}
		return nil, err
	}

	buf, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	return out, json.Unmarshal(buf, &out)
}

func setState(ctx context.Context, static assetSet, state []string) error {
	buf, err := json.Marshal(&state)
	if err != nil {
		return err
	}

	_, err = static.s3.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(bytes.NewReader(buf)),
		Bucket: &static.bucket,
		Key:    &static.stateKey,
	})
	return err
}
