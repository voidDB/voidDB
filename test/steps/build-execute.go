package teststeps

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"
)

func AddStepBuildAndExecute(sc *godog.ScenarioContext) {
	sc.When(`^I build "([^"]*)" and execute the resulting binary `+
		`with argument(?:s?)$`,
		buildAndExecute,
	)

	return
}

func buildAndExecute(
	ctx0 context.Context, srcPath string, args *godog.DocString,
) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		binName string = strings.TrimSuffix(
			filepath.Base(srcPath),
			filepath.Ext(srcPath),
		)

		tempDir string = ctx.Value(ctxKeyTempDir{}).(string)

		binPath string = filepath.Join(tempDir, binName)

		builder *exec.Cmd = exec.Command("gcc",
			"-o", binPath,
			srcPath,
		)

		command *exec.Cmd = exec.Command(binPath,
			strings.Fields(args.Content)...,
		)
	)

	builder.Stderr = os.Stderr

	e = builder.Run()
	if e != nil {
		return
	}

	command.Dir = tempDir

	command.Stderr = os.Stderr

	e = command.Run()
	if e != nil {
		return
	}

	return
}
