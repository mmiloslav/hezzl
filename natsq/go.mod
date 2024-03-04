module natsq

go 1.21

require common v0.0.0-00010101000000-000000000000

require (
	github.com/nats-io/nats.go v1.33.1
	github.com/sirupsen/logrus v1.9.3
)

require (
	github.com/klauspost/compress v1.17.2 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/stretchr/testify v1.8.1 // indirect
	golang.org/x/crypto v0.18.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
)

replace common => ../common

replace redisdb => ../redisdb
