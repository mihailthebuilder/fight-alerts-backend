package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents/types"
)

type EventBridgeClient struct {
	host   string
	port   int
	handle *cloudwatchevents.Client
}

func (e *EventBridgeClient) Connect() error {

	endpoint := fmt.Sprintf("http://%s:%d", e.host, e.port)

	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion("us-east-1"),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: endpoint, SigningRegion: "us-east-1"}, nil
			},
		)),
	)
	if err != nil {
		return err
	}

	client := cloudwatchevents.NewFromConfig(cfg)
	e.handle = client

	return nil
}

func (e *EventBridgeClient) GetAllRuleNamesByNamePrefix(name string) ([]types.Rule, error) {

	lri := cloudwatchevents.ListRulesInput{NamePrefix: aws.String(name)}

	lro, err := e.handle.ListRules(context.Background(), &lri)
	if err != nil {
		return nil, err
	}

	return lro.Rules, nil
}

func (e *EventBridgeClient) GetAllTargetIdsByRuleName(name string) ([]types.Target, error) {

	lti := cloudwatchevents.ListTargetsByRuleInput{Rule: aws.String(name)}

	lto, err := e.handle.ListTargetsByRule(context.Background(), &lti)
	if err != nil {
		return nil, err
	}

	return lto.Targets, nil
}

func (e *EventBridgeClient) InsertRuleWithTarget(ruleName, targetId, targetARN string) error {
	pri := cloudwatchevents.PutRuleInput{Name: &ruleName, ScheduleExpression: aws.String("cron(0 12 10 10 10 2023)")}

	_, err := e.handle.PutRule(context.Background(), &pri)
	if err != nil {
		return err
	}

	target := types.Target{Arn: &targetARN, Id: &targetId}

	pti := cloudwatchevents.PutTargetsInput{Rule: &ruleName, Targets: []types.Target{target}}

	_, err = e.handle.PutTargets(context.Background(), &pti)
	if err != nil {
		return err
	}

	return nil
}
