package scheduler

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents"
)

func TestCloudWatchEventsScheduler_UpdateTrigger(t *testing.T) {
	tests := []struct {
		name                        string
		createRuleReturnError       error
		setTargetForRuleReturnError error
		wantErr                     bool
	}{
		{
			name:    "happy path",
			wantErr: false,
		},
		{
			name:                  "CreateRule returns an error",
			createRuleReturnError: fmt.Errorf("fake err"),
			wantErr:               true,
		},
		{
			name:                        "SetTargetForRule returns an error",
			setTargetForRuleReturnError: fmt.Errorf("fake err"),
			wantErr:                     true,
		},
	}
	for _, tt := range tests {

		cwc := MockCloudWatchClient{}

		cwc.On("PutRule").Return(&cloudwatchevents.PutRuleOutput{}, tt.createRuleReturnError)
		cwc.On("PutTargets").Return(&cloudwatchevents.PutTargetsOutput{}, tt.setTargetForRuleReturnError)

		cwes := CloudWatchEventsScheduler{
			TargetARN:       "mock-arn",
			RuleName:        "mock-rule",
			TargetId:        "mock-target-id",
			RuleDescription: "Mock description",
			Api:             cwc,
		}

		mockTime := time.Now()

		t.Run(tt.name, func(t *testing.T) {
			if err := cwes.UpdateTrigger(mockTime); (err != nil) != tt.wantErr {
				t.Errorf("CloudWatchEventsScheduler.UpdateTrigger() error = %v, wantErr %v", err, tt.wantErr)
			}
		})

		if !tt.wantErr {
			cwc.AssertExpectations(t)
		}
	}
}
