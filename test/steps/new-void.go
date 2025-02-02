package steps

import (
	"context"
	"path/filepath"

	"github.com/cucumber/godog"

	"github.com/voidDB/voidDB"
)

func AddStepNewVoid(sc *godog.ScenarioContext) {
	sc.Given(`^there is a new Void "([^"]*)"$`, newVoidDefault)

	sc.Given(`^there is a new Void "([^"]*)" of capacity (\d+)$`, newVoid)

	return
}

func newVoidDefault(ctx context.Context, name string) (
	context.Context, error,
) {
	const (
		capacity = 1 << 20 // 1 MiB
	)

	return newVoid(ctx, name, capacity)
}

func newVoid(ctx0 context.Context, name string, capacity int) (
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

	void, e = voidDB.NewVoid(path, capacity)
	if e != nil {
		return
	}

	ctx = context.WithValue(ctx, ctxKeyVoid{name}, void)

	return
}

type ctxKeyVoid struct {
	Name string
}
