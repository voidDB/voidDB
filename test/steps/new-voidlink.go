package steps

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cucumber/godog"
)

const (
	voidlinkDefaultPeriod = "100ms"
	voidlinkDisablePeriod = "0"
)

func AddStepNewVoidlink(sc *godog.ScenarioContext) {
	sc.Given(`^there is a new VoidLink between "([^"]*)" and "([^"]*)"$`,
		newVoidlink,
	)

	sc.Given(`^there is a new uplink-only VoidLink between "([^"]*)" and `+
		`"([^"]*)"$`,
		newVoidlinkUpOnly,
	)

	sc.Given(`^there is a new downlink-only VoidLink between "([^"]*)" and `+
		`"([^"]*)"$`,
		newVoidlinkDownOnly,
	)

	return
}

func newVoidlink(ctx context.Context, voidName, bucketName string) (
	context.Context, error,
) {

	return newVoidlinkSetPeriods(ctx, voidName, bucketName,
		voidlinkDefaultPeriod, voidlinkDefaultPeriod,
	)
}

func newVoidlinkUpOnly(ctx context.Context, voidName, bucketName string) (
	context.Context, error,
) {

	return newVoidlinkSetPeriods(ctx, voidName, bucketName,
		voidlinkDefaultPeriod, voidlinkDisablePeriod,
	)
}

func newVoidlinkDownOnly(ctx context.Context, voidName, bucketName string) (
	context.Context, error,
) {

	return newVoidlinkSetPeriods(ctx, voidName, bucketName,
		voidlinkDisablePeriod, voidlinkDefaultPeriod,
	)
}

func newVoidlinkSetPeriods(ctx0 context.Context, voidName, bucketName,
	uplinkPeriod, downlinkPeriod string,
) (
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
			"--uplink-period", uplinkPeriod,
			"--downlink-period", downlinkPeriod,
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
