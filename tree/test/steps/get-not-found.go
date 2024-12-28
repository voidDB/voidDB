package steps

import (
	"context"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"

	"github.com/voidDB/voidDB/tree"
)

func AddStepGetNotFound(sc *godog.ScenarioContext) {
	sc.Then(`^Get-ting "([^"]*)" from "([^"]*)" should not find$`, getNotFound)

	return
}

func getNotFound(ctx0 context.Context, key, name string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		r root = ctx.Value(ctxKeyRoot{name}).(root)
	)

	_, e = tree.Get(&r.medium, r.offset,
		[]byte(key),
	)

	assert.Equal(
		godog.T(ctx),
		tree.ErrorNotFound,
		e,
	)

	e = nil

	return
}
