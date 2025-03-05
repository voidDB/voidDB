module voidlink

go 1.24.0

replace github.com/voidDB/voidDB => ../..

require (
	github.com/minio/minio-go v6.0.14+incompatible
	github.com/voidDB/voidDB v0.0.0-00010101000000-000000000000
	golang.org/x/sync v0.11.0
)

require (
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	golang.org/x/crypto v0.35.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
)
