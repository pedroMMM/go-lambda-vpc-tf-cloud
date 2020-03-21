package main

import (
	"context"
	"sync"

	"github.com/hashicorp/go-multierror"

	"github.com/aws/aws-sdk-go/service/ec2"
)

func getVpcs(ctx context.Context, static assetSet) ([]string, error) {
	vpcs := make([]string, 0)
	var mVpcs sync.Mutex
	errs := &multierror.Error{}
	var mErrs sync.Mutex
	var wg sync.WaitGroup

	perPage := func(page *ec2.DescribeVpcsOutput, lastPage bool) bool {
		for _, vpc := range page.Vpcs {
			vpc := vpc
			mVpcs.Lock()
			vpcs = append(vpcs, *vpc.VpcId)
			mVpcs.Unlock()
		}

		return !lastPage
	}

	for _, ec2Client := range static.ec2 {
		ec2Client := ec2Client
		wg.Add(1)

		go func() {
			defer wg.Done()

			err := ec2Client.DescribeVpcsPagesWithContext(ctx, nil, perPage)
			if err != nil {
				mErrs.Lock()
				errs = multierror.Append(errs, err)
				mErrs.Unlock()
			}
		}()
	}

	wg.Wait()
	return vpcs, errs.ErrorOrNil()
}
