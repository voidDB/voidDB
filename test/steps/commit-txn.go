package steps

import (
	"context"

	"github.com/cucumber/godog"

	"github.com/voidDB/voidDB"
)

func AddStepCommitTxn(sc *godog.ScenarioContext) {
	sc.When(`^I commit "([^"]*)"$`, commitTxn)

	return
}

func commitTxn(ctx0 context.Context, name string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		txn *voidDB.Txn = ctx.Value(ctxKeyTxn{name}).(*voidDB.Txn)
	)

	e = txn.Commit()
	if e != nil {
		return
	}

	ctx = context.WithValue(ctx, ctxKeyTxn{name}, txn)

	return
}
