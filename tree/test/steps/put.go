package steps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/voidDB/voidDB/tree"
)

func AddStepPut(sc *godog.ScenarioContext) {
	sc.When(`^I put "([^"]*)", "([^"]*)" into "([^"]*)"$`, put)

	return
}

func put(ctx0 context.Context, key, value, name string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		r root = ctx.Value(ctxKeyTree{name}).(root)
	)

	r.offset, e = tree.Put(&r.medium, r.offset,
		[]byte(key),
		[]byte(value),
	)
	if e != nil {
		return
	}

	ctx = context.WithValue(ctx, ctxKeyTree{name}, r)

	return
}
