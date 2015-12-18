package ec2

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/margic/aws/util"
)

// create a tag with key and value provided
// returns pointer to new tag
func CreateTag(k string, v string) *ec2.Tag {
	return &ec2.Tag{
		Key:   aws.String(k),
		Value: aws.String(v),
	}
}

func TagResource(ctx *util.AwsContext, resourceId *string, tags ...*ec2.Tag) (*ec2.CreateTagsOutput, error) {
	svc := ec2.New(ctx.AwsSession)
	ids := []*string{
		resourceId,
	}

	// everything will be tagged with CreateBy tag
	newTags := make([]*ec2.Tag, len(tags)+1)
	copy(newTags, tags)
	newTags[len(newTags)-1] = CreateTag("CreatedBy", "goCli")
	createTagsInput := ec2.CreateTagsInput{
		Resources: ids,
		Tags:      newTags,
	}
	return svc.CreateTags(&createTagsInput)
}
