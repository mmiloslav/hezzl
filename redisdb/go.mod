module redisdb

go 1.21

require common v0.0.0-00010101000000-000000000000

require (
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/sirupsen/logrus v1.9.3
)

require (
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.31.1 // indirect
	golang.org/x/sys v0.15.0 // indirect
)

replace common => ../common
