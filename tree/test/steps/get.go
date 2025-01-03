package steps

import (
	"context"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"

	"github.com/voidDB/voidDB/tree"
)

func AddStepGet(sc *godog.ScenarioContext) {
	sc.Then(`^I should get "([^"]*)", "([^"]*)" from "([^"]*)"$`,
		get,
	)

	sc.Then(`^I should get "([^"]*)", "([^"]*)" using "([^"]*)"$`,
		getUsingCursor,
	)

	return
}

func get(ctx0 context.Context, key, valueExpect, name string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		r root = ctx.Value(ctxKeyTree{name}).(root)

		valueActual []byte
	)

	valueActual, e = tree.Get(&r.medium, r.offset,
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

func getUsingCursor(ctx0 context.Context, key, valueExpect, name string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		cursor *tree.Cursor = ctx.Value(ctxKeyCursor{name}).(*tree.Cursor)

		valueActual []byte
	)

	valueActual, e = cursor.Get(
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
