module main

go 1.21

require (
	natsq v0.0.0-00010101000000-000000000000
	postgres v0.0.0-00010101000000-000000000000
	redisdb v0.0.0-00010101000000-000000000000
)

require (
	github.com/go-playground/validator v9.31.0+incompatible
	github.com/gorilla/mux v1.8.1
	github.com/sirupsen/logrus v1.9.3
)

require (
	common v0.0.0-00010101000000-000000000000 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-redis/redis v6.15.9+incompatible // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.17.2 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/nats-io/nats.go v1.33.1 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.19.0 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gorm.io/driver/postgres v1.5.6 // indirect
	gorm.io/gorm v1.25.7 // indirect
)

replace common => ../common

replace postgres => ../postgres

replace redisdb => ../redisdb

replace natsq => ../natsq
