package steps

import (
	"context"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"

	"github.com/voidDB/voidDB/tree"
)

func AddStepGetLast(sc *godog.ScenarioContext) {
	sc.Then(
		`^I should get "([^"]*)", "([^"]*)" first using "([^"]*)" in reverse$`,
		getLast,
	)

	return
}

func getLast(ctx0 context.Context, keyExpect, valueExpect, name string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		cursor *tree.Cursor = ctx.Value(ctxKeyCursor{name}).(*tree.Cursor)

		keyActual   []byte
		valueActual []byte
	)

	keyActual, valueActual, e = cursor.GetLast()
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