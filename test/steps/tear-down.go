package steps

import (
	"context"
	"os"

	"github.com/cucumber/godog"
)

func AddStepTearDown(sc *godog.ScenarioContext) {
	sc.After(tearDown)

	return
}

func tearDown(ctx0 context.Context, scenario *godog.Scenario, e0 error) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	e = os.RemoveAll(
		ctx.Value(ctxKeyTempDir{}).(string),
	)
	if e != nil {
		return
	}

	return
}
