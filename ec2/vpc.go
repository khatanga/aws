package ec2

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/margic/aws/util"
)

// Create VPC
func CreateVPC(r map[string]util.Resource) (err error) {
	svc := ec2.New(util.Session)

	vpcInput := ec2.CreateVpcInput{
		DryRun:          aws.Bool(false),
		CidrBlock:       aws.String("10.2.0.0/16"),
		InstanceTenancy: aws.String("default"),
	}

	vpcOutput, err := svc.CreateVpc(&vpcInput)
	if err != nil {
		return err
	}

	// tag the new vpc
	vpcId := vpcOutput.Vpc.VpcId

	_, err = util.TagResource(vpcId, util.CreateTag("Name", "KhatangaVPC"))

	// save the vpcId in the results map
	r["vpcId"] = util.Resource{
		ResourceType: "vpc",
		ResourceID:   vpcId,
	}

	return err
}

func CreateSubnets(r map[string]util.Resource) error {
	svc := ec2.New(util.Session)

	publicSubnetInput := ec2.CreateSubnetInput{
		AvailabilityZone: aws.String("us-west-2b"),
		CidrBlock:        aws.String("10.2.0.0/20"),
		DryRun:           aws.Bool(false),
		VpcId:            r["vpcId"].ResourceID,
	}

	privateSubnetInput := ec2.CreateSubnetInput{
		AvailabilityZone: aws.String("us-west-2b"),
		CidrBlock:        aws.String("10.2.16.0/20"),
		DryRun:           aws.Bool(false),
		VpcId:            r["vpcId"].ResourceID,
	}

	// create public subnet
	createSubnetOutput, err := svc.CreateSubnet(&publicSubnetInput)
	if err != nil {
		return err
	}

	subnetId := createSubnetOutput.Subnet.SubnetId
	// save the subnet id in the results
	r["publicSubnet"] = util.Resource{
		ResourceType: "subnet",
		ResourceID:   subnetId,
	}
	// set the tags on the new public subnet
	_, err = util.TagResource(subnetId, util.CreateTag("Name", "KhatangaPublic"))

	if err != nil {
		return err
	}

	// create private subnet
	createSubnetOutput, err = svc.CreateSubnet(&privateSubnetInput)
	if err != nil {
		return err
	}

	subnetId = createSubnetOutput.Subnet.SubnetId
	// save the subnet id in the results
	r["privateSubnet"] = util.Resource{
		ResourceType: "subnet",
		ResourceID:   subnetId,
	}
	// set the tags on the new public subnet
	_, err = util.TagResource(subnetId, util.CreateTag("Name", "KhatangaPrivate"))

	return err
}

func CreateInternetGateway(r map[string]util.Resource) error {
	svc := ec2.New(util.Session)

	createGWInput := ec2.CreateInternetGatewayInput{
		DryRun: aws.Bool(false),
	}

	createGWOutput, err := svc.CreateInternetGateway(&createGWInput)
	gwId := createGWOutput.InternetGateway.InternetGatewayId
	// save the internet gateway id
	r["internetGateway"] = util.Resource{
		ResourceType: "internetGateway",
		ResourceID:   gwId,
	}
	// set the tags on the internet gateway
	_, err = util.TagResource(gwId, util.CreateTag("Name", "KhatangaGateway"))

	if err != nil {
		return err
	}

	attachGWInput := ec2.AttachInternetGatewayInput{
		DryRun:            aws.Bool(false),
		InternetGatewayId: gwId,
		VpcId:             r["vpcId"].ResourceID,
	}
	_, err = svc.AttachInternetGateway(&attachGWInput)
	return err
}
