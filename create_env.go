package main

// package for automating the build of the khatanga aws environment
// follows vpc scenario 2 with public and private subnets
import (
	"flag"
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

	// set up all the command line flag variables
	var fAll, fVpc, fSubnets, fInternetGw, fNatEip, fNatGw bool

	flag.BoolVar(&fAll, "createAll", false, "Special flag to execute all create statements")
	flag.BoolVar(&fVpc, "createVpc", false, "Create Virtual Private Cloud (VPC) for environment")
	flag.BoolVar(&fSubnets, "createSubnets", false, "Create private and public Subnets in VPC")
	flag.BoolVar(&fInternetGw, "createInternetGw", false, "Creates an internet gateway")
	flag.BoolVar(&fNatEip, "createNatEip", false, "Creates an EIP for the Nat Gateway")
	flag.BoolVar(&fNatGw, "createNatGw", false, "Creates the Nat Gateway. Will require a EIP")

	flag.Parse()

	// create an awsContext that stores the session and results from function calls
	// for example ids of created vpc
	session := session.New(&aws.Config{Region: aws.String("us-west-2")})
	token := aws.String(feeds.NewUUID().String())
	ctx := util.NewAwsContext(session)
	// token for retry logic
	ctx.IdempotentToken = token

	// determine work
	if !(fAll || fVpc || fSubnets || fInternetGw || fNatEip || fNatGw) {
		log.Error("Nothing to do no flags set")
		flag.Usage()
	}

	if fAll || fVpc {
		err := ec.CreateVPC(&ctx)
		if err != nil {
			panic(err)
		}
	}

	if fAll || fSubnets {
		err := ec.CreateSubnets(&ctx)
		if err != nil {
			panic(err)
		}
	}

	if fAll || fInternetGw {
		err := ec.CreateInternetGateway(&ctx)
		if err != nil {
			panic(err)
		}
	}

	if fAll || fNatEip {
		err := ec.CreateNatEip(&ctx)
		if err != nil {
			panic(err)
		}
	}
	if fAll || fNatGw {
		err := ec.CreateNATGateway(&ctx)
		if err != nil {
			panic(err)
		}
	}
	log.Info(ctx.Results)
}

/*
	deleting the vpc will delete subnets and internet gateway
*/
