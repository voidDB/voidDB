package steps

import (
	"context"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"

	"github.com/voidDB/voidDB/common"
	"github.com/voidDB/voidDB/tree"
)

func AddStepGetNextNotFound(sc *godog.ScenarioContext) {
	sc.Then(`^getting next using "([^"]*)" should not find$`, getNextNotFound)

	return
}

func getNextNotFound(ctx0 context.Context, cursorName string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		cursor *tree.Cursor = ctx.Value(ctxKeyCursor{cursorName}).(*tree.Cursor)
	)

	_, _, e = cursor.GetNext()

	assert.ErrorIs(
		godog.T(ctx),
		e,
		common.ErrorNotFound,
	)

	e = nil

	return
}
