package steps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/voidDB/voidDB/tree"
)

func AddStepDel(sc *godog.ScenarioContext) {
	sc.When(`^I delete "([^"]*)" from "([^"]*)"$`, del)

	return
}

func del(ctx0 context.Context, key, name string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		r root = ctx.Value(ctxKeyTree{name}).(root)
	)

	r.offset, e = tree.Del(&r.medium, r.offset,
		[]byte(key),
	)
	if e != nil {
		return
	}

	ctx = context.WithValue(ctx, ctxKeyTree{name}, r)

	return
}
