module github.com/maxnrm/teleflood

go 1.22.3

require (
	github.com/nats-io/nats.go v1.36.0
	go.uber.org/ratelimit v0.3.1
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	gopkg.in/telebot.v3 v3.3.6
)

replace gopkg.in/telebot.v3 v3.3.6 => github.com/maxnrm/telebot v0.0.0

require (
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/klauspost/compress v1.17.2 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/stretchr/testify v1.8.1 // indirect
	golang.org/x/crypto v0.18.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)
