package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/s3"
)

type assetSet struct {
	bucket   string
	stateKey string
	ec2      []*ec2.EC2
	s3       *s3.S3
}

func newAssetSet(mySession *session.Session, bucket, states3Key string) (assetSet, error) {
	asset := assetSet{
		bucket:   bucket,
		stateKey: states3Key,
		ec2:      make([]*ec2.EC2, 0),
		s3:       s3.New(mySession),
	}

	regions, err := getAllAvailableRegions(ec2.New(mySession))
	if err != nil {
		return asset, err
	}

	for _, e := range regions {
		e := e
		asset.ec2 = append(asset.ec2, ec2.New(mySession, &aws.Config{Region: &e}))
	}

	return asset, nil
}

func getAllAvailableRegions(ec2 *ec2.EC2) ([]string, error) {
	result, err := ec2.DescribeRegions(nil)
	if err != nil {
		return nil, err
	}

	out := make([]string, 0)
	for _, e := range result.Regions {
		out = append(out, *e.RegionName)
	}

	return out, nil
}
