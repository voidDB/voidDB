package test

import (
	"testing"

	"github.com/cucumber/godog"

	"test/steps"
)

func TestVoid(t *testing.T) {
	var (
		scenarioInitializer = func(sc *godog.ScenarioContext) {
			steps.AddStepSetUp(sc)
			steps.AddStepTearDown(sc)

			steps.AddStepNewVoid(sc)
			steps.AddStepBeginTxn(sc)
			steps.AddStepOpenCursor(sc)
			steps.AddStepGet(sc)
			steps.AddStepGetNext(sc)
			steps.AddStepPut(sc)
			steps.AddStepDel(sc)
			steps.AddStepCommitTxn(sc)
			steps.AddStepAbortTxn(sc)
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
