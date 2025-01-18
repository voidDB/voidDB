package steps

import (
	"context"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"

	"github.com/voidDB/voidDB/common"
	"github.com/voidDB/voidDB/cursor"
)

func AddStepGet(sc *godog.ScenarioContext) {
	sc.Then(`^I should get "([^"]*)", "([^"]*)" using "([^"]*)"$`, get)

	sc.Then(`^getting "([^"]*)" using "([^"]*)" should not find$`, getNotFound)

	return
}

func get(ctx0 context.Context, key, valueExpect, name string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		cur *cursor.Cursor = ctx.Value(ctxKeyCursor{name}).(*cursor.Cursor)

		valueActual []byte
	)

	valueActual, e = cur.Get(
		[]byte(key),
	)
	if e != nil {
		return
	}

	assert.Equal(
		godog.T(ctx),
		[]byte(valueExpect),
		valueActual,
	)

	return
}

func getNotFound(ctx0 context.Context, key, name string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		cur *cursor.Cursor = ctx.Value(ctxKeyCursor{name}).(*cursor.Cursor)
	)

	_, e = cur.Get(
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
