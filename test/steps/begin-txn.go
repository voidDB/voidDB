package teststeps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/voidDB/voidDB"
)

func AddStepBeginTxn(sc *godog.ScenarioContext) {
	sc.When(`^I begin a read-only transaction "([^"]*)" in "([^"]*)"$`,
		beginTxnReadonly,
	)

	sc.When(`^I begin a transaction "([^"]*)" in "([^"]*)"$`,
		beginTxnWrite,
	)

	return
}

func beginTxnReadonly(ctx context.Context, txnName, voidName string) (
	context.Context, error,
) {
	return beginTxn(ctx, txnName, voidName, true)
}

func beginTxnWrite(ctx context.Context, txnName, voidName string) (
	context.Context, error,
) {
	return beginTxn(ctx, txnName, voidName, false)
}

func beginTxn(ctx0 context.Context, txnName, voidName string, readonly bool) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		void *voidDB.Void = ctx.Value(ctxKeyVoid{voidName}).(*voidDB.Void)

		txn *voidDB.Txn
	)

	txn, e = void.BeginTxn(readonly)
	if e != nil {
		return
	}

	ctx = context.WithValue(ctx, ctxKeyTxn{txnName}, txn)

	return
}

type ctxKeyTxn struct {
	Name string
}
