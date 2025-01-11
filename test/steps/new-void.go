package teststeps

import (
	"context"
	"path/filepath"

	"github.com/cucumber/godog"

	"github.com/voidDB/voidDB"
)

func AddStepNewVoid(sc *godog.ScenarioContext) {
	sc.Given(`^there is a new Void "([^"]*)" of size (\d+)$`, newVoid)

	return
}

func newVoid(ctx0 context.Context, name string, size int) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		void *voidDB.Void

		path string = filepath.Join(
			ctx.Value(ctxKeyTempDir{}).(string),
			name,
		)
	)

	void, e = voidDB.NewVoid(path, size)
	if e != nil {
		return
	}

	ctx = context.WithValue(ctx, ctxKeyVoid{name}, void)

	return
}

type ctxKeyVoid struct {
	Name string
}
