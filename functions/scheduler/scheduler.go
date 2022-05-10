package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents/types"
)

type IScheduler interface {
	UpdateTrigger(time.Time) error
}

type CloudWatchEventsScheduler struct {
	Api             ICloudWatchClient
	TargetARN       string
	RuleName        string
	RuleDescription string
	TargetId        string
}

type ICloudWatchClient interface {
	PutRule(ctx context.Context,
		params *cloudwatchevents.PutRuleInput,
		optFns ...func(*cloudwatchevents.Options)) (*cloudwatchevents.PutRuleOutput, error)

	PutTargets(ctx context.Context,
		params *cloudwatchevents.PutTargetsInput,
		optFns ...func(*cloudwatchevents.Options)) (*cloudwatchevents.PutTargetsOutput, error)
}

func (cw CloudWatchEventsScheduler) UpdateTrigger(t time.Time) error {
	err := cw.ClearPreviousRules()
	if err != nil {
		return err
	}

	err = cw.CreateRule(t)
	if err != nil {
		return err
	}

	err = cw.SetTargetForRule()
	if err != nil {
		return err
	}

	return nil
}

func (cw CloudWatchEventsScheduler) ClearPreviousRules() error {
	return nil
}

func (cw CloudWatchEventsScheduler) CreateRule(t time.Time) error {
	scheduleExpression := fmt.Sprintf("cron(0 12 %v %v ? %v)", t.Day(), t.Format("Jan"), t.Year())

	params := cloudwatchevents.PutRuleInput{
		Name:               &cw.RuleName,
		Description:        &cw.RuleDescription,
		ScheduleExpression: &scheduleExpression,
	}

	_, err := cw.Api.PutRule(context.TODO(), &params)

	if err != nil {
		return err
	}

	fmt.Printf("Created EventBridge rule with name=%s, which will be triggered on cron=%s\n", cw.RuleName, scheduleExpression)

	return nil
}

func (cw CloudWatchEventsScheduler) SetTargetForRule() error {

	target := types.Target{Id: &cw.TargetId, Arn: &cw.TargetARN}

	params := cloudwatchevents.PutTargetsInput{Rule: &cw.RuleName, Targets: []types.Target{target}}

	_, err := cw.Api.PutTargets(context.TODO(), &params)

	if err != nil {
		return err
	}

	fmt.Printf("Attached target with ARN=%s to EventBridge rule name=%s\n", cw.TargetARN, cw.RuleName)

	return nil
}
