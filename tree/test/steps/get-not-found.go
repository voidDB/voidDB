package steps

import (
	"context"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"

	"github.com/voidDB/voidDB/common"
	"github.com/voidDB/voidDB/tree"
)

func AddStepGetNotFound(sc *godog.ScenarioContext) {
	sc.Then(`^getting "([^"]*)" using "([^"]*)" should not find$`,
		getNotFoundUsingCursor,
	)

	return
}

func getNotFoundUsingCursor(ctx0 context.Context, key, name string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		cursor *tree.Cursor = ctx.Value(ctxKeyCursor{name}).(*tree.Cursor)
	)

	_, e = cursor.Get(
		[]byte(key),
	)

	assert.ErrorIs(
		godog.T(ctx),
		e,
		common.ErrorNotFound,
	)

	e = nil

	return
}
