package util

import (
	"github.com/aws/aws-sdk-go/aws/session"
)

type AwsContext struct {
	AwsSession      *session.Session
	Results         map[string]resource
	IdempotentToken *string
}

func NewAwsContext(session *session.Session) AwsContext {
	return AwsContext{
		AwsSession: session,
		Results:    make(map[string]resource),
	}
}

// After calling an aws sevice add the result to the context's result map
// This can then be used to retrieve information about created resources etc
// key is the key of the result in the map. e.g. if creating a vpc the vpc resource id
// will be required by later calls. Store the result with the key: "vpc" in the map will
// make it easy to retrieve the new resouce id later
func (a AwsContext) AddResult(key string, resourceType string, resourceID *string) {
	a.Results[key] = resource{
		ResourceType: resourceType,
		ResourceID:   resourceID,
	}
}

// Session is a shared session for svc clients to use
//var Session = session.New(&aws.Config{Region: aws.String("us-west-2")})

type resource struct {
	// don't forget to update add resource if adding fields to this resource type
	ResourceType string
	ResourceID   *string
}
