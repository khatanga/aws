package util

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Session is a shared session for svc clients to use
var Session = session.New(&aws.Config{Region: aws.String("us-west-2")})
