package steps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/voidDB/voidDB/tree"
)

func AddStepNewRootNode(sc *godog.ScenarioContext) {
	sc.Given(`^there is a new root Node "([^"]*)"$`, newRootNode)

	return
}

func newRootNode(ctx0 context.Context, name string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		medium Medium
		offset int
	)

	offset, e = medium.Save(
		tree.NewNode(),
	)
	if e != nil {
		return
	}

	ctx = context.WithValue(ctx, ctxKeyRoot{name},
		root{medium, offset},
	)

	return
}

type ctxKeyRoot struct {
	Name string
}

type root struct {
	medium Medium
	offset int
}
