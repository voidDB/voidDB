package steps

import (
	"context"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"

	"github.com/voidDB/voidDB/common"
	"github.com/voidDB/voidDB/cursor"
)

func AddStepGetNext(sc *godog.ScenarioContext) {
	sc.Then(`^I should get "([^"]*)", "([^"]*)" next using "([^"]*)"$`, getNext)

	sc.Then(`^getting next using "([^"]*)" should not find$`, getNextNotFound)

	return
}

func getNext(ctx0 context.Context, keyExpect, valueExpect, name string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		cur *cursor.Cursor = ctx.Value(ctxKeyCursor{name}).(*cursor.Cursor)

		keyActual   []byte
		valueActual []byte
	)

	keyActual, valueActual, e = cur.GetNext()
	if e != nil {
		return
	}

	assert.Equal(
		godog.T(ctx),
		[]byte(keyExpect),
		keyActual,
	)

	assert.Equal(
		godog.T(ctx),
		[]byte(valueExpect),
		valueActual,
	)

	return
}

func getNextNotFound(ctx0 context.Context, name string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		cur *cursor.Cursor = ctx.Value(ctxKeyCursor{name}).(*cursor.Cursor)
	)

	_, _, e = cur.GetNext()

	assert.ErrorIs(
		godog.T(ctx),
		e,
		common.ErrorNotFound,
	)

	e = nil

	return
}
