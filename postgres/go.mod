module postgres

go 1.21

require (
	common v0.0.0-00010101000000-000000000000
	natsq v0.0.0-00010101000000-000000000000
	redisdb v0.0.0-00010101000000-000000000000
)

require (
	github.com/lib/pq v1.10.9
	github.com/sirupsen/logrus v1.9.3
	gorm.io/driver/postgres v1.5.6
	gorm.io/gorm v1.25.7
)

require (
	github.com/go-redis/redis v6.15.9+incompatible // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.17.2 // indirect
	github.com/nats-io/nats.go v1.33.1 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/nxadm/tail v1.4.11 // indirect
	golang.org/x/crypto v0.18.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)

replace common => ../common

replace redisdb => ../redisdb

replace natsq => ../natsq
