package steps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/voidDB/voidDB/tree"
)

func AddStepNewTree(sc *godog.ScenarioContext) {
	sc.Given(`^there is a new tree "([^"]*)"$`, newTree)

	return
}

func newTree(ctx0 context.Context, name string) (
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

	ctx = context.WithValue(ctx, ctxKeyTree{name},
		root{medium, offset},
	)

	return
}

type ctxKeyTree struct {
	Name string
}

type root struct {
	medium Medium
	offset int
}
