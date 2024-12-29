package steps

import (
	"context"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"

	"github.com/voidDB/voidDB/tree"
)

func AddStepGetNext(sc *godog.ScenarioContext) {
	sc.Then(`^I should get "([^"]*)", "([^"]*)" next from "([^"]*)"$`, getNext)

	return
}

func getNext(ctx0 context.Context, keyExpect, valueExpect, cursorName string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		cursor *tree.Cursor = ctx.Value(ctxKeyCursor{cursorName}).(*tree.Cursor)

		keyActual   []byte
		valueActual []byte
	)

	keyActual, valueActual, e = cursor.GetNext()
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