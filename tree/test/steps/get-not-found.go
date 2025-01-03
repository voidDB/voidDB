package steps

import (
	"context"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"

	"github.com/voidDB/voidDB/tree"
)

func AddStepGetNotFound(sc *godog.ScenarioContext) {
	sc.Then(`getting "([^"]*)" from "([^"]*)" should not find$`,
		getNotFound,
	)

	sc.Then(`^getting "([^"]*)" using "([^"]*)" should not find$`,
		getNotFoundUsingCursor,
	)

	return
}

func getNotFound(ctx0 context.Context, key, name string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		r root = ctx.Value(ctxKeyTree{name}).(root)
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

	assert.Equal(
		godog.T(ctx),
		tree.ErrorNotFound,
		e,
	)

	e = nil

	return
}
