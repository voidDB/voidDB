package teststeps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/voidDB/voidDB"
	"github.com/voidDB/voidDB/cursor"
)

func AddStepPut(sc *godog.ScenarioContext) {
	sc.When(`^I put "([^"]*)", "([^"]*)" in "([^"]*)"$`, put)

	sc.When(`^I put "([^"]*)", "([^"]*)" using "([^"]*)"$`, putUsingCursor)

	return
}

func put(ctx0 context.Context, key, value, txnName string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		txn *voidDB.Txn = ctx.Value(ctxKeyTxn{txnName}).(*voidDB.Txn)
	)

	e = txn.Put(
		[]byte(key),
		[]byte(value),
	)
	if e != nil {
		return
	}

	return
}

func putUsingCursor(ctx0 context.Context, key, value, cursorName string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		cur = ctx.Value(ctxKeyCursor{cursorName}).(*cursor.Cursor)
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
