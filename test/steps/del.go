package steps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/voidDB/voidDB"
)

func AddStepDel(sc *godog.ScenarioContext) {
	sc.When(`^I delete "([^"]*)" from "([^"]*)"$`, del)

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
