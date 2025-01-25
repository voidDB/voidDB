package teststeps

import (
	"context"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"

	"github.com/voidDB/voidDB"
	"github.com/voidDB/voidDB/common"
	"github.com/voidDB/voidDB/cursor"
)

func AddStepGet(sc *godog.ScenarioContext) {
	sc.Then(`^I should get "([^"]*)", "([^"]*)" from "([^"]*)"$`, get)

	sc.Then(`^I should get "([^"]*)", "([^"]*)" using "([^"]*)"$`,
		getUsingCursor,
	)

	sc.Then(`^getting "([^"]*)" from "([^"]*)" should not find$`, getNotFound)

	return
}

func get(ctx0 context.Context, key, valueExpect, txnName string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		txn *voidDB.Txn = ctx.Value(ctxKeyTxn{txnName}).(*voidDB.Txn)

		valueActual []byte
	)

	valueActual, e = txn.Get(
		[]byte(key),
	)
	if e != nil {
		return
	}

	assert.Equal(
		godog.T(ctx),
		[]byte(valueExpect),
		valueActual,
	)

	return
}

func getUsingCursor(ctx0 context.Context, key, valueExpect, cursorName string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		cur = ctx.Value(ctxKeyCursor{cursorName}).(*cursor.Cursor)

		valueActual []byte
	)

	valueActual, e = cur.Get(
		[]byte(key),
	)
	if e != nil {
		return
	}

	assert.Equal(
		godog.T(ctx),
		[]byte(valueExpect),
		valueActual,
	)

	return
}

func getNotFound(ctx0 context.Context, key, txnName string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		txn *voidDB.Txn = ctx.Value(ctxKeyTxn{txnName}).(*voidDB.Txn)
	)

	_, e = txn.Get(
		[]byte(key),
	)

	assert.Equal(
		godog.T(ctx),
		common.ErrorNotFound,
		e,
	)

	e = nil

	return
}
