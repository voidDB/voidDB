package test

import (
	"testing"

	"github.com/cucumber/godog"

	"test/steps"
)

func TestVoidlink(t *testing.T) {
	var (
		scenarioInitializer = func(sc *godog.ScenarioContext) {
			steps.AddStepSetUp(sc)
			steps.AddStepTearDown(sc)

			steps.AddStepNewVoid(sc)
			steps.AddStepBeginTxn(sc)
			steps.AddStepOpenCursor(sc)
			steps.AddStepGet(sc)
			steps.AddStepPut(sc)
			steps.AddStepDel(sc)
			steps.AddStepCommitTxn(sc)

			steps.AddStepNewMinioServer(sc)
			steps.AddStepNewVoidlink(sc)
			steps.AddStepWait(sc)
		}

		options = &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features/voidlink.feature"},
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
