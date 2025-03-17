package steps

import (
	"context"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/cucumber/godog"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	minioServerAddr = "127.237.93.83:9383"
	minioServerCred = "minioadmin"
	minioServerPath = "minio-data"
)

func AddStepNewMinioServer(sc *godog.ScenarioContext) {
	sc.Given(`^there is a new Minio server$`, newMinioServer)

	return
}

func newMinioServer(ctx0 context.Context) (ctx context.Context, e error) {
	ctx = ctx0

	var (
		server *exec.Cmd = exec.Command("minio", "server",
			"--address", minioServerAddr,
			filepath.Join(
				ctx.Value(ctxKeyTempDir{}).(string),
				minioServerPath,
			),
		)

		timeout context.Context
	)

	e = server.Start()
	if e != nil {
		return
	}

	ctx = context.WithValue(ctx, ctxKeyProcesses{},
		append(
			ctx.Value(ctxKeyProcesses{}).([]*os.Process),
			server.Process,
		),
	)

	timeout, _ = context.WithTimeout(ctx, time.Second)

	for {
		if timeout.Err() != nil {
			return
		}

		_, e = net.Dial("tcp", minioServerAddr)
		if e == nil {
			return
		}
	}

	return
}

var minioClient *minio.Client

func init() {
	var (
		e error
	)

	minioClient, e = minio.New(minioServerAddr,
		&minio.Options{
			Creds: credentials.NewStaticV4(minioServerCred, minioServerCred,
				"",
			),
		},
	)
	if e != nil {
		panic(e)
	}

	return
}
