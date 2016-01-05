package ec2

// package with operations to create khatanga environment
// using convention for variables c|d for create of describe
// followed by name if resource type
// followed by I|O for input or output
// eg. for CreateVpcInput type varialbe will be cVpcI

import (
	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/khatanga/aws/util"
)



// Create VPC
func CreateVPC(ctx *util.AwsContext) (err error) {
	svc := ec2.New(ctx.AwsSession)
	cVpcI := ec2.CreateVpcInput{
		DryRun:          aws.Bool(util.Config.Vpc.DryRun),
		CidrBlock:       aws.String(util.Config.Vpc.CidrBlock),
		InstanceTenancy: aws.String(util.Config.Vpc.InstanceTenancy),
	}

	log.Info("Creating vpc")
	cVpcO, err := svc.CreateVpc(&cVpcI)
	if err != nil {
		return err
	}

	// tag the new vpc
	vpcId := cVpcO.Vpc.VpcId
	_, err = TagResource(ctx, vpcId, CreateTag("Name", util.Config.Vpc.TagName))
	// save the vpcId in the results map
	ctx.AddResult("vpc", "vpc", vpcId)
	return err
}

func CreateSubnets(ctx *util.AwsContext) error {
	svc := ec2.New(ctx.AwsSession)

	cPubSubnetI := ec2.CreateSubnetInput{
		AvailabilityZone: aws.String(util.Config.PubSubnet.AvailabilityZone),
		CidrBlock:        aws.String(util.Config.PubSubnet.CidrBlock),
		DryRun:           aws.Bool(util.Config.PubSubnet.DryRun),
		VpcId:            ctx.Results[util.Config.PubSubnet.VpcId].ResourceID,
	}

	cPvtSubnetI := ec2.CreateSubnetInput{
		AvailabilityZone: aws.String(util.Config.PvtSubnet.AvailabilityZone),
		CidrBlock:        aws.String(util.Config.PvtSubnet.CidrBlock),
		DryRun:           aws.Bool(util.Config.PvtSubnet.DryRun),
		VpcId:            ctx.Results[util.Config.PvtSubnet.VpcId].ResourceID,
	}

	// create public subnet
	log.Info("Creating public Subnet")
	cPubSubnetO, err := svc.CreateSubnet(&cPubSubnetI)
	if err != nil {
		return err
	}

	subnetId := cPubSubnetO.Subnet.SubnetId
	// save the subnet id in the results
	ctx.AddResult("publicSubnet", "subnet", subnetId)
	// set the tags on the new public subnet
	_, err = TagResource(ctx, subnetId, CreateTag("Name", util.Config.PubSubnet.TagName))

	if err != nil {
		return err
	}

	// create private subnet
	log.Info("Creating private Subnet")
	cPvtSubnetO, err := svc.CreateSubnet(&cPvtSubnetI)
	if err != nil {
		return err
	}

	subnetId = cPvtSubnetO.Subnet.SubnetId
	// save the subnet id in the results
	ctx.AddResult("privateSubnet", "subnet", subnetId)
	// set the tags on the new public subnet
	_, err = TagResource(ctx, subnetId, CreateTag("Name", util.Config.PvtSubnet.TagName))
	return err
}

func CreateInternetGateway(ctx *util.AwsContext) error {
	svc := ec2.New(ctx.AwsSession)

	cIgwI := ec2.CreateInternetGatewayInput{
		DryRun: aws.Bool(false),
	}
	log.Info("Creating Internet Gateway")
	cIgwO, err := svc.CreateInternetGateway(&cIgwI)
	gwId := cIgwO.InternetGateway.InternetGatewayId
	// save the internet gateway id
	ctx.AddResult("gateway", "internetGateway", gwId)
	// set the tags on the internet gateway
	_, err = TagResource(ctx, gwId, CreateTag("Name", util.Config.InetGateway.Name))

	if err != nil {
		return err
	}

	aIgwI := ec2.AttachInternetGatewayInput{
		DryRun:            aws.Bool(false),
		InternetGatewayId: gwId,
		VpcId:             ctx.Results["vpc"].ResourceID,
	}
	log.Info("Attach Internet Gateway to VPC")
	_, err = svc.AttachInternetGateway(&aIgwI)
	return err
}

func CreateNatEip(ctx *util.AwsContext) error {
	svc := ec2.New(ctx.AwsSession)

	aEipI := ec2.AllocateAddressInput{
		DryRun: aws.Bool(false),
	}
	log.Info("Creating Nat Elastic IP")
	aEipO, err := svc.AllocateAddress(&aEipI)
	if err != nil {
		return err
	}
	ctx.AddResult("natEip", "ElasticIp", aEipO.AllocationId)
	// tag the ip
	_, err = TagResource(ctx, aEipO.AllocationId, CreateTag("Name", util.Config.NatEip.Name))
	return err
}

func CreateNATGateway(ctx *util.AwsContext) error {
	svc := ec2.New(ctx.AwsSession)

	cNatGwI := ec2.CreateNatGatewayInput{
		// elastic ip id
		AllocationId: ctx.Results["natEip"].ResourceID,
		ClientToken:  ctx.IdempotentToken,
		SubnetId:     ctx.Results["publicSubnet"].ResourceID,
	}

	log.Info("Creating NAT Gateway")
	cNatGwO, err := svc.CreateNatGateway(&cNatGwI)
	if err != nil {
		return err
	}
	gwId := cNatGwO.NatGateway.NatGatewayId
	ctx.AddResult("natGateway", "Nat Gateway", gwId)
	// tag the gateway
	_, err = TagResource(ctx, gwId, CreateTag("Name", util.Config.NatGateway.Name))
	// the gateway takes time to start run a function to test if this is up yet.
	gwIds := []*string{gwId}
	dNatGwI := ec2.DescribeNatGatewaysInput{
		NatGatewayIds: gwIds,
	}
	dNatGwO, err := svc.DescribeNatGateways(&dNatGwI)
	log.WithField("NatGwState", dNatGwO.NatGateways[0].State).Info("NatGateway state")
	return err
}
