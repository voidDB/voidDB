package test

import (
	"testing"

	"github.com/cucumber/godog"

	"github.com/voidDB/voidDB/tree/test/steps"
)

func TestNode(t *testing.T) {
	var (
		scenarioInitializer = func(sc *godog.ScenarioContext) {
			steps.AddStepNewRootNode(sc)
			steps.AddStepPut(sc)
			steps.AddStepGet(sc)
			steps.AddStepGetNotFound(sc)
			steps.AddStepDel(sc)
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
