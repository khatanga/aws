package unit

import (
	"github.com/aws/aws-sdk-go/awstesting/integration/smoke"
	"github.com/margic/aws/util"
)

var Ctx = util.NewAwsContext(smoke.Session)
