package steps

import (
	"context"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"

	"github.com/voidDB/voidDB"
)

func AddStepGetNext(sc *godog.ScenarioContext) {
	sc.Then(`^I should get "([^"]*)", "([^"]*)" next from "([^"]*)"$`,
		getNext,
	)
}

func getNext(
	ctx0 context.Context, keyExpect, valueExpect, txnName string,
) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		txn *voidDB.Txn = ctx.Value(ctxKeyTxn{txnName}).(*voidDB.Txn)

		keyActual   []byte
		valueActual []byte
	)

	keyActual, valueActual, e = txn.GetNext()
	if e != nil {
		return
	}

	assert.Equal(
		godog.T(ctx),
		[]byte(keyExpect),
		keyActual,
	)

	assert.Equal(
		godog.T(ctx),
		[]byte(valueExpect),
		valueActual,
	)

	return
}
