package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gorilla/feeds"
	ec "github.com/khatanga/aws/ec2" // imported as ec as to avoid name collision with ec2
	"github.com/khatanga/aws/util"
)

// Creates aws environment for Khatanga project
func main() {
	log.Info("Starting AWS Client")
	// create an awsContext that stores the session and results from function calls
	// for example ids of created vpc
	session := session.New(&aws.Config{Region: aws.String("us-west-2")})
	token := aws.String(feeds.NewUUID().String())

	ctx := util.NewAwsContext(session)

	// token for retry logic
	ctx.IdempotentToken = token

	err := ec.CreateVPC(&ctx)
	if err != nil {
		panic(err)
	}
	err = ec.CreateSubnets(&ctx)
	if err != nil {
		panic(err)
	}
	err = ec.CreateInternetGateway(&ctx)
	if err != nil {
		panic(err)
	}

	err = ec.CreateNatEip(&ctx)
	if err != nil {
		panic(err)
	}

	err = ec.CreateNATGateway(&ctx)
	if err != nil {
		panic(err)
	}

	log.Info(ctx.Results)
}

/*
	deleting the vpc will delete subnets and internet gateway
*/
