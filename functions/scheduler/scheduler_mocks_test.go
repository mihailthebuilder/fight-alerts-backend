package scheduler

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents"
	"github.com/stretchr/testify/mock"
)

type MockCloudWatchClient struct {
	mock.Mock
}

func (cwc MockCloudWatchClient) PutRule(ctx context.Context,
	params *cloudwatchevents.PutRuleInput,
	optFns ...func(*cloudwatchevents.Options)) (*cloudwatchevents.PutRuleOutput, error) {

	args := cwc.Called()
	return args.Get(0).(*cloudwatchevents.PutRuleOutput), args.Error(1)
}

func (cwc MockCloudWatchClient) PutTargets(ctx context.Context,
	params *cloudwatchevents.PutTargetsInput,
	optFns ...func(*cloudwatchevents.Options)) (*cloudwatchevents.PutTargetsOutput, error) {

	args := cwc.Called()
	return args.Get(0).(*cloudwatchevents.PutTargetsOutput), args.Error(1)
}
