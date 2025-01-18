package steps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/voidDB/voidDB/cursor"
)

func AddStepPut(sc *godog.ScenarioContext) {
	sc.When(`^I put "([^"]*)", "([^"]*)" using "([^"]*)"$`, put)

	return
}

func put(ctx0 context.Context, key, value, name string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		cur *cursor.Cursor = ctx.Value(ctxKeyCursor{name}).(*cursor.Cursor)
	)

	e = cur.Put(
		[]byte(key),
		[]byte(value),
	)
	if e != nil {
		return
	}

	return
}
