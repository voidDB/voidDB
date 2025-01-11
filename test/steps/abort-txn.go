package teststeps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/voidDB/voidDB"
)

func AddStepAbortTxn(sc *godog.ScenarioContext) {
	sc.When(`^I abort "([^"]*)"$`, abortTxn)

	return
}

func abortTxn(ctx0 context.Context, name string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		txn *voidDB.Txn = ctx.Value(ctxKeyTxn{name}).(*voidDB.Txn)
	)

	e = txn.Abort()
	if e != nil {
		return
	}

	ctx = context.WithValue(ctx, ctxKeyTxn{name}, txn)

	return
}
