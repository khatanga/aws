package util

import (
	"github.com/aws/aws-sdk-go/private/waiter"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func WaitForNatGateway(input *ec2.DescribeNatGatewaysInput, svc *ec2.EC2) error {
	waiterCfg := waiter.Config{
		Operation:   "DescribeNatGateways",
		Delay:       15,
		MaxAttempts: 40,
		Acceptors: []waiter.WaitAcceptor{
			{
				State:    "success",
				Matcher:  "pathAll",
				Argument: "NatGateways[].State",
				Expected: "available",
			},
		},
	}

	w := waiter.Waiter{
		Client: svc,
		Input:  input,
		Config: waiterCfg,
	}
	return w.Wait()

}
