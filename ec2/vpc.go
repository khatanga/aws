package ec2

import (
	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/margic/aws/util"
)

// Create VPC
func CreateVPC(ctx *util.AwsContext) (err error) {
	svc := ec2.New(ctx.AwsSession)
	vpcInput := ec2.CreateVpcInput{
		DryRun:          aws.Bool(false),
		CidrBlock:       aws.String("10.2.0.0/16"),
		InstanceTenancy: aws.String("default"),
	}

	log.Info("Creating vpc")
	vpcOutput, err := svc.CreateVpc(&vpcInput)
	if err != nil {
		return err
	}

	// tag the new vpc
	vpcId := vpcOutput.Vpc.VpcId

	_, err = TagResource(ctx, vpcId, CreateTag("Name", "KhatangaVPC"))

	// save the vpcId in the results map
	ctx.AddResult("vpc", "vpc", vpcId)

	return err
}

func CreateSubnets(ctx *util.AwsContext) error {
	svc := ec2.New(ctx.AwsSession)

	publicSubnetInput := ec2.CreateSubnetInput{
		AvailabilityZone: aws.String("us-west-2b"),
		CidrBlock:        aws.String("10.2.0.0/20"),
		DryRun:           aws.Bool(false),
		VpcId:            ctx.Results["vpc"].ResourceID,
	}

	privateSubnetInput := ec2.CreateSubnetInput{
		AvailabilityZone: aws.String("us-west-2b"),
		CidrBlock:        aws.String("10.2.16.0/20"),
		DryRun:           aws.Bool(false),
		VpcId:            ctx.Results["vpc"].ResourceID,
	}

	// create public subnet
	log.Info("Creating public Subnet")
	createSubnetOutput, err := svc.CreateSubnet(&publicSubnetInput)
	if err != nil {
		return err
	}

	subnetId := createSubnetOutput.Subnet.SubnetId
	// save the subnet id in the results
	ctx.AddResult("publicSubnet", "subnet", subnetId)

	// set the tags on the new public subnet
	_, err = TagResource(ctx, subnetId, CreateTag("Name", "KhatangaPublic"))

	if err != nil {
		return err
	}

	// create private subnet
	log.Info("Creating private Subnet")
	createSubnetOutput, err = svc.CreateSubnet(&privateSubnetInput)
	if err != nil {
		return err
	}

	subnetId = createSubnetOutput.Subnet.SubnetId
	// save the subnet id in the results
	ctx.AddResult("privateSubnet", "subnet", subnetId)
	// set the tags on the new public subnet
	_, err = TagResource(ctx, subnetId, CreateTag("Name", "KhatangaPrivate"))

	return err
}

func CreateInternetGateway(ctx *util.AwsContext) error {
	svc := ec2.New(ctx.AwsSession)

	createGWInput := ec2.CreateInternetGatewayInput{
		DryRun: aws.Bool(false),
	}
	log.Info("Creating Internet Gateway")
	createGWOutput, err := svc.CreateInternetGateway(&createGWInput)
	gwId := createGWOutput.InternetGateway.InternetGatewayId
	// save the internet gateway id
	ctx.AddResult("gateway", "internetGateway", gwId)
	// set the tags on the internet gateway
	_, err = TagResource(ctx, gwId, CreateTag("Name", "KhatangaGateway"))

	if err != nil {
		return err
	}

	attachGWInput := ec2.AttachInternetGatewayInput{
		DryRun:            aws.Bool(false),
		InternetGatewayId: gwId,
		VpcId:             ctx.Results["vpc"].ResourceID,
	}
	log.Info("Attach Internet Gateway to VPC")
	_, err = svc.AttachInternetGateway(&attachGWInput)
	return err
}

func CreateNatEip(ctx *util.AwsContext) error {
	svc := ec2.New(ctx.AwsSession)

	eipInput := ec2.AllocateAddressInput{
		DryRun: aws.Bool(false),
	}
	log.Info("Creating Nat Elastic IP")
	eipOutput, err := svc.AllocateAddress(&eipInput)
	if err != nil {
		return err
	}
	ctx.AddResult("natEip", "ElasticIp", eipOutput.AllocationId)
	return err
}

func CreateNATGateway(ctx *util.AwsContext) error {
	svc := ec2.New(ctx.AwsSession)

	natGwInput := ec2.CreateNatGatewayInput{
		// elastic ip id
		AllocationId: ctx.Results["natEip"].ResourceID,
		ClientToken:  ctx.IdempotentToken,
		SubnetId:     ctx.Results["publicSubnet"].ResourceID,
	}

	log.Info("Creating NAT Gateway")
	natGwOutput, err := svc.CreateNatGateway(&natGwInput)
	if err != nil {
		return err
	}
	ctx.AddResult("natGateway", "Nat Gateway", natGwOutput.NatGateway.NatGatewayId)
	return err
}
