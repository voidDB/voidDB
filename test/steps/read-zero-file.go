package teststeps

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
)

func AddStepReadZeroFile(sc *godog.ScenarioContext) {
	sc.Then(`^I should read exactly (\d+) zero bytes from file "([^"]*)"$`,
		readZeroFile,
	)

	return
}

func readZeroFile(ctx0 context.Context, n int, fileName string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		file *os.File

		path string = filepath.Join(
			ctx.Value(ctxKeyTempDir{}).(string),
			fileName,
		)

		empty []byte = make([]byte, n)
		slice []byte = make([]byte, n)
	)

	file, e = os.Open(path)
	if e != nil {
		return
	}

	_, e = io.ReadFull(file, slice)
	if e != nil {
		return
	}

	assert.Equal(
		godog.T(ctx),
		empty,
		slice,
	)

	_, e = io.ReadFull(file, slice)

	assert.Equal(
		godog.T(ctx),
		io.EOF,
		e,
	)

	e = nil

	return
}
