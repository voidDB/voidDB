package steps

import (
	"context"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"

	"github.com/voidDB/voidDB/cursor"
)

func AddStepGetFirst(sc *godog.ScenarioContext) {
	sc.Then(`^I should get "([^"]*)", "([^"]*)" first using "([^"]*)"$`,
		getFirst,
	)

	return
}

func getFirst(ctx0 context.Context, keyExpect, valueExpect, name string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		cur *cursor.Cursor = ctx.Value(ctxKeyCursor{name}).(*cursor.Cursor)

		keyActual   []byte
		valueActual []byte
	)

	keyActual, valueActual, e = cur.GetFirst()
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
