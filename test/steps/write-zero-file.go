package teststeps

import (
	"context"
	"fmt"
	"path/filepath"
	"unsafe"

	"github.com/cucumber/godog"

	_ "github.com/voidDB/voidDB/libvoid"
)

// #include <stdlib.h>
// #include "../../libvoid/include/void.h"
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

		filePathC *C.char = C.CString(filePath)

		status C.int
	)

	defer C.free(
		unsafe.Pointer(filePathC),
	)

	status = C.void_write_zero_file(
		C.size_t(n),
		filePathC,
	)

	if status != 0 {
		e = fmt.Errorf("C function returned non-zero value %d", status)

		return
	}

	return
}
