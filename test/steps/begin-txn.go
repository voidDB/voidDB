package steps

import (
	"context"
	"syscall"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"

	"github.com/voidDB/voidDB"
)

func AddStepBeginTxn(sc *godog.ScenarioContext) {
	sc.When(`^I begin a read-only transaction "([^"]*)" in "([^"]*)"$`,
		beginTxnReadonly,
	)

	sc.When(`^I begin a transaction "([^"]*)" in "([^"]*)"$`,
		beginTxnWrite,
	)

	sc.Then(`^beginning a transaction in "([^"]*)" should fail with EAGAIN$`,
		beginTxnFails,
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

func beginTxnFails(ctx0 context.Context, voidName string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	ctx, e = beginTxnWrite(ctx, "", voidName)

	assert.Equal(
		godog.T(ctx),
		syscall.EAGAIN,
		e,
	)

	e = nil

	return
}

func beginTxn(ctx0 context.Context, txnName, voidName string, readonly bool) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		void *voidDB.Void = ctx.Value(ctxKeyVoid{voidName}).(*voidDB.Void)

		txn *voidDB.Txn
	)

	txn, e = void.BeginTxn(readonly, false)
	if e != nil {
		return
	}

	ctx = context.WithValue(ctx, ctxKeyTxn{txnName}, txn)

	return
}

type ctxKeyTxn struct {
	Name string
}
