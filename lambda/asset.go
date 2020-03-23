package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/ses"
)

type assetSet struct {
	bucket    string
	stateKey  string
	emailFrom string
	emailTo   string
	ec2       []*ec2.EC2
	s3        *s3.S3
	ses       *ses.SES
}

type newAssetSetInput struct {
	session    *session.Session
	bucket     string
	states3Key string
	emailFrom  string
	emailTo    string
}

func newAssetSet(in newAssetSetInput) (assetSet, error) {
	asset := assetSet{
		bucket:    in.bucket,
		stateKey:  in.states3Key,
		emailFrom: in.emailFrom,
		emailTo:   in.emailTo,
		ec2:       make([]*ec2.EC2, 0),
		s3:        s3.New(in.session),
		ses:       ses.New(in.session),
	}

	regions, err := getAllAvailableRegions(ec2.New(in.session))
	if err != nil {
		return asset, err
	}

	for _, e := range regions {
		e := e
		asset.ec2 = append(asset.ec2, ec2.New(in.session, &aws.Config{Region: &e}))
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
