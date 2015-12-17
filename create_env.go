package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/margic/aws/svc"
	"github.com/margic/aws/util"
)

// Creates aws environment for Khatanga project
func main() {
	fmt.Println("Starting AWS Client")

	// create a map to store outputs and ids from func calls
	// for example ids of created vpc
	results := make(map[string]util.Resource)

	s := session.New(&aws.Config{Region: aws.String("us-west-2")})

	err := svc.CreateVPC(s, results)
	if err != nil {
		panic(err)
	}
	err = svc.CreateSubnets(s, results)
	if err != nil {
		panic(err)
	}
	err = svc.CreateInternetGateway(s, results)
	if err != nil {
		panic(err)
	}

	fmt.Println(results)
}

/*
	deleting the vpc will delete subnets and internet gateway
*/
