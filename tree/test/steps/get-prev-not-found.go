package steps

import (
	"context"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"

	"github.com/voidDB/voidDB/tree"
)

func AddStepGetPrevNotFound(sc *godog.ScenarioContext) {
	sc.Then(`^getting next using "([^"]*)" in reverse should not find$`,
		getPrevNotFound,
	)

	return
}

func getPrevNotFound(ctx0 context.Context, cursorName string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		cursor *tree.Cursor = ctx.Value(ctxKeyCursor{cursorName}).(*tree.Cursor)
	)

	_, _, e = cursor.GetPrev()

	assert.Equal(
		godog.T(ctx),
		tree.ErrorNotFound,
		e,
	)

	e = nil

	return
}
