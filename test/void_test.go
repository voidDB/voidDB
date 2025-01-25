package test

import (
	"testing"

	"github.com/cucumber/godog"

	"github.com/voidDB/voidDB/test/steps"
)

func TestVoid(t *testing.T) {
	var (
		scenarioInitializer = func(sc *godog.ScenarioContext) {
			teststeps.AddStepSetUp(sc)
			teststeps.AddStepWriteZeroFile(sc)
			teststeps.AddStepReadZeroFile(sc)
			teststeps.AddStepTearDown(sc)

			teststeps.AddStepNewVoid(sc)
			teststeps.AddStepBeginTxn(sc)
			teststeps.AddStepOpenCursor(sc)
			teststeps.AddStepGet(sc)
			teststeps.AddStepPut(sc)
			teststeps.AddStepDel(sc)
			teststeps.AddStepCommitTxn(sc)
			teststeps.AddStepAbortTxn(sc)
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
