package steps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/voidDB/voidDB/cursor"
)

func AddStepDel(sc *godog.ScenarioContext) {
	sc.When(`^I delete using "([^"]*)"$`, del)

	return
}

func del(ctx0 context.Context, name string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		cur *cursor.Cursor = ctx.Value(ctxKeyCursor{name}).(*cursor.Cursor)
	)

	e = cur.Del()
	if e != nil {
		return
	}

	return
}
