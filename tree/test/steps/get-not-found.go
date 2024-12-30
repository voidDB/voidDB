package steps

import (
	"context"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"

	"github.com/voidDB/voidDB/tree"
)

func AddStepGetNotFound(sc *godog.ScenarioContext) {
	sc.Then(`^getting "([^"]*)" from "([^"]*)" should not find$`, getNotFound)

	return
}

func getNotFound(ctx0 context.Context, key, name string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		cursor *tree.Cursor = ctx.Value(ctxKeyCursor{name}).(*tree.Cursor)
	)

	_, e = cursor.Get(
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
