package steps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/voidDB/voidDB"
	"github.com/voidDB/voidDB/cursor"
)

func AddStepDel(sc *godog.ScenarioContext) {
	sc.When(`^I delete "([^"]*)" from "([^"]*)"$`, del)

	sc.When(`^I delete "([^"]*)" using "([^"]*)"$`, delUsingCursor)

	return
}

func del(ctx0 context.Context, key, txnName string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		txn *voidDB.Txn = ctx.Value(ctxKeyTxn{txnName}).(*voidDB.Txn)
	)

	_, e = txn.Get(
		[]byte(key),
	)
	if e != nil {
		return
	}

	e = txn.Del()
	if e != nil {
		return
	}

	return
}

func delUsingCursor(ctx0 context.Context, key, cursorName string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		cur = ctx.Value(ctxKeyCursor{cursorName}).(*cursor.Cursor)
	)

	_, e = cur.Get(
		[]byte(key),
	)
	if e != nil {
		return
	}

	e = cur.Del()
	if e != nil {
		return
	}

	return
}
