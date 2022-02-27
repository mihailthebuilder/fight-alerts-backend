package main

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_run(t *testing.T) {
	type args struct {
		s IScraper
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	successMockScraper := NewMockIScraper(ctrl)
	successMockScraper.EXPECT().getResultsFromUrl().Return([]FightRecord{}, nil)

	failMockScraper := NewMockIScraper(ctrl)
	failMockScraper.EXPECT().getResultsFromUrl().Return(nil, fmt.Errorf("Panic example"))

	tests := []struct {
		name      string
		args      args
		wantPanic bool
	}{
		{name: "Test run function successful", args: args{s: successMockScraper}},
		{name: "Test run function failed", args: args{s: failMockScraper}, wantPanic: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic == true {
				assert.Panics(t, func() { run(tt.args.s) }, "The code did not panic")
				return
			}

			run(tt.args.s)
		})
	}
}
