package main

import (
	"fmt"
	//"github.com/aws/aws-sdk-go/awstesting/integration/smoke"
	ec "github.com/margic/aws/ec2" // imported as ec as to avoid name collision with ec2
	"github.com/margic/aws/util"
)

// Creates aws environment for Khatanga project
func main() {
	fmt.Println("Starting AWS Client")

	// create a map to store outputs and ids from func calls
	// for example ids of created vpc
	results := make(map[string]util.Resource)

	err := ec.CreateVPC(results)
	if err != nil {
		panic(err)
	}
	err = ec.CreateSubnets(results)
	if err != nil {
		panic(err)
	}
	err = ec.CreateInternetGateway(results)
	if err != nil {
		panic(err)
	}

	fmt.Println(results)
}

/*
	deleting the vpc will delete subnets and internet gateway
*/
