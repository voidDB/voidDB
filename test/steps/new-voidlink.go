package steps

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cucumber/godog"
)

func AddStepNewVoidlink(sc *godog.ScenarioContext) {
	sc.Given(`^there is a new VoidLink between "([^"]*)" and "([^"]*)"$`,
		newVoidlink,
	)

	return
}

func newVoidlink(ctx0 context.Context, voidName, bucketName string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		binPath string = filepath.Join(
			ctx.Value(ctxKeyTempDir{}).(string),
			"voidlink",
		)

		builder *exec.Cmd = exec.Command("go", "build",
			"-o", binPath,
			"voidlink",
		)

		command *exec.Cmd = exec.Command(binPath,
			filepath.Join(
				ctx.Value(ctxKeyTempDir{}).(string),
				voidName,
			),
			minioServerAddr,
			bucketName,
			"--uplink-period", "100ms",
			"--downlink-period", "100ms",
		)

		exitError   *exec.ExitError
		isExitError bool
	)

	e = builder.Run()

	if exitError, isExitError = e.(*exec.ExitError); isExitError {
		e = fmt.Errorf("%s", exitError.Stderr)
	}

	if e != nil {
		return
	}

	command.Stderr = os.Stderr

	e = command.Start()
	if e != nil {
		return
	}

	ctx = context.WithValue(ctx, ctxKeyProcesses{},
		append(
			ctx.Value(ctxKeyProcesses{}).([]*os.Process),
			command.Process,
		),
	)

	return
}
