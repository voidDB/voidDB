package test

import (
	"testing"

	"github.com/cucumber/godog"

	"test/steps"
)

func TestCursor(t *testing.T) {
	var (
		scenarioInitializer = func(sc *godog.ScenarioContext) {
			steps.AddStepNewTree(sc)
			steps.AddStepPut(sc)
			steps.AddStepGet(sc)
			steps.AddStepDel(sc)
			steps.AddStepNewCursor(sc)
			steps.AddStepGetNext(sc)
			steps.AddStepGetPrev(sc)
			steps.AddStepGetFirst(sc)
			steps.AddStepGetLast(sc)
		}

		options = &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		}

		suite = godog.TestSuite{
			ScenarioInitializer: scenarioInitializer,
			Options:             options,
		}
	)

	if suite.Run() != 0 {
		t.Fatal()
	}

	return
}
