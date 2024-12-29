package steps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/voidDB/voidDB/tree"
)

func AddStepNewCursor(sc *godog.ScenarioContext) {
	sc.When(`^I open a new cursor "([^"]*)" at "([^"]*)"$`, newCursor)

	return
}

func newCursor(ctx0 context.Context, name, rootName string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		r root = ctx.Value(ctxKeyRoot{rootName}).(root)
	)

	ctx = context.WithValue(ctx, ctxKeyCursor{name},
		tree.NewCursor(&r.medium, r.offset),
	)

	return
}

type ctxKeyCursor struct {
	Name string
}
