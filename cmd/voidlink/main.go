package main

import (
	"context"
	_ "embed"
	"errors"
	"os"
	"os/signal"
	"time"

	"github.com/alecthomas/kong"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/voidDB/voidDB"
)

var (
	//go:embed README
	readme string
)

type CLI struct {
	PathToVoid     string        `arg type:path`
	S3Endpoint     string        `arg name:"s3-endpoint"`
	BucketName     string        `arg`
	UploadName     string        `arg optional help:"Defaults to a UUIDv6."`
	AccessKey      string        `env:"ACCESS_KEY" default:"minioadmin"`
	SecretKey      string        `env:"SECRET_KEY" default:"minioadmin"`
	UplinkPeriod   time.Duration `default:"1s" help:"0 disables uplink."`
	DownlinkPeriod time.Duration `default:"1s" help:"0 disables downlink."`
	VoidCapacity   int           `default:"1099511627776"`
}

func (cli *CLI) Run() (e error) {
	var (
		client *minio.Client
		ctx    context.Context
		exists bool
		loop   *Loop
		uuidv6 uuid.UUID
		void   *voidDB.Void
	)

	void, e = voidDB.NewVoid(cli.PathToVoid, cli.VoidCapacity)

	if errors.Is(e, os.ErrExist) {
		void, e = voidDB.OpenVoid(cli.PathToVoid, cli.VoidCapacity)
	}

	if e != nil {
		return
	}

	client, e = minio.New(cli.S3Endpoint,
		&minio.Options{
			Creds: credentials.NewStaticV4(cli.AccessKey, cli.SecretKey, ""),
		},
	)
	if e != nil {
		return
	}

	ctx, _ = signal.NotifyContext(
		context.Background(),
		os.Interrupt,
	)

	exists, e = client.BucketExists(ctx, cli.BucketName)
	if e != nil {
		return
	}

	if !exists {
		e = client.MakeBucket(ctx, cli.BucketName, minio.MakeBucketOptions{})
		if e != nil {
			return
		}
	}

	if cli.UploadName == "" {
		uuidv6, e = uuid.NewV6()
		if e != nil {
			return
		}

		cli.UploadName = uuidv6.String()
	}

	loop = NewLoop(
		NewLink(void, client),
		cli.BucketName,
		cli.UploadName,
		cli.UplinkPeriod,
		cli.DownlinkPeriod,
	)

	go loop.Uplink(ctx)

	go loop.Downlink(ctx)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func main() {
	var (
		app *kong.Context = kong.Parse(&CLI{},
			kong.Description(readme),
		)
	)

	app.FatalIfErrorf(
		app.Run(),
	)

	return
}
