package steps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/voidDB/voidDB/tree"
)

func AddStepDel(sc *godog.ScenarioContext) {
	sc.When(`^I delete using "([^"]*)"$`, delUsingCursor)

	return
}

func delUsingCursor(ctx0 context.Context, name string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		cursor *tree.Cursor = ctx.Value(ctxKeyCursor{name}).(*tree.Cursor)
	)

	e = cursor.Del()
	if e != nil {
		return
	}

	return
}
