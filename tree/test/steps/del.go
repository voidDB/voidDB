package steps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/voidDB/voidDB/tree"
)

func AddStepDel(sc *godog.ScenarioContext) {
	sc.When(`^I Del "([^"]*)" from "([^"]*)"$`, _del)

	return
}

func _del(ctx0 context.Context, key, name string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		r root = ctx.Value(ctxKeyRoot{name}).(root)
	)

	r.offset, e = tree.Del(&r.medium, r.offset,
		[]byte(key),
	)
	if e != nil {
		return
	}

	ctx = context.WithValue(ctx, ctxKeyRoot{name}, r)

	return
}
