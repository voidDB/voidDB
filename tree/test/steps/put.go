package steps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/voidDB/voidDB/tree"
)

func AddStepPut(sc *godog.ScenarioContext) {
	sc.When(`^I Put "([^"]*)", "([^"]*)" into "([^"]*)"$`, _put)

	return
}

func _put(ctx0 context.Context, key, value, name string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		r root = ctx.Value(ctxKeyRoot{name}).(root)
	)

	r.offset, e = tree.Put(&r.medium, r.offset,
		[]byte(key),
		[]byte(value),
	)
	if e != nil {
		return
	}

	ctx = context.WithValue(ctx, ctxKeyRoot{name}, r)

	return
}
