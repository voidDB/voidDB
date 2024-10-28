package teststeps

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/cucumber/godog"

	_ "github.com/voidDB/voidDB/lib"
)

// #include "../../lib/include/void.h"
import "C"

func AddStepWriteZeroFile(sc *godog.ScenarioContext) {
	sc.When(`^I invoke a C function that writes (\d+) zero bytes to file `+
		`"([^"]*)"$`,
		writeZeroFile,
	)

	return
}

func writeZeroFile(ctx0 context.Context, n int, fileName string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		filePath string = filepath.Join(
			ctx.Value(ctxKeyTempDir{}).(string),
			fileName,
		)

		status C.int
	)

	status = C.write_zero_file(
		C.size_t(n),
		C.CString(filePath),
	)

	if status != 0 {
		e = fmt.Errorf("C function returned non-zero value %d", status)

		return
	}

	return
}
