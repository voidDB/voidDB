package steps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/voidDB/voidDB"
	"github.com/voidDB/voidDB/cursor"
)

func AddStepOpenCursor(sc *godog.ScenarioContext) {
	sc.When(`^I open a cursor "([^"]*)" associated with keyspace "([^"]*)" `+
		`and "([^"]*)"$`,
		openCursor,
	)

	return
}

func openCursor(ctx0 context.Context, cursorName, keyspace, txnName string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		txn *voidDB.Txn = ctx.Value(ctxKeyTxn{txnName}).(*voidDB.Txn)

		cur *cursor.Cursor
	)

	cur, e = txn.OpenCursor(
		[]byte(keyspace),
	)
	if e != nil {
		return
	}

	ctx = context.WithValue(ctx, ctxKeyCursor{cursorName}, cur)

	return
}

type ctxKeyCursor struct {
	Name string
}
