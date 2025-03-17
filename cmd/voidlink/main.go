package main

import (
	"context"
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

type CLI struct {
	PathToVoid      string        `arg`
	S3Endpoint      string        `arg name:"s3-endpoint"`
	BucketName      string        `arg`
	ObjectName      string        `arg optional`
	AccessKeyID     string        `default:"minioadmin"`
	SecretAccessKey string        `default:"minioadmin"`
	UplinkPeriod    time.Duration `default:"1s"`
	DownlinkPeriod  time.Duration `default:"1s"`
	VoidCapacity    int           `default:"1099511627776"`
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
			Creds: credentials.NewStaticV4(cli.AccessKeyID, cli.SecretAccessKey,
				"",
			),
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

	if cli.ObjectName == "" {
		uuidv6, e = uuid.NewV6()
		if e != nil {
			return
		}

		cli.ObjectName = uuidv6.String()
	}

	loop = NewLoop(
		NewLink(void, client),
		cli.BucketName,
		cli.ObjectName,
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
		app *kong.Context = kong.Parse(&CLI{})
	)

	app.FatalIfErrorf(
		app.Run(),
	)

	return
}
