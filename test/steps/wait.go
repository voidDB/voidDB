package steps

import (
	"context"
	"time"

	"github.com/cucumber/godog"
)

func AddStepWait(sc *godog.ScenarioContext) {
	sc.When(`^I wait for "([^"]*)"$`, wait)

	return
}

func wait(ctx0 context.Context, waitFor string) (ctx context.Context, e error) {
	ctx = ctx0

	var (
		duration time.Duration
	)

	duration, e = time.ParseDuration(waitFor)
	if e != nil {
		return
	}

	time.Sleep(duration)

	return
}
