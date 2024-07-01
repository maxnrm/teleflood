package config

import (
	"fmt"
	"os"
	"strconv"
)

// COMMON
var NATS_URL = getEnvStr("NATS_URL")
var BOT_TOKEN = getEnvStr("BOT_TOKEN")
var TELEGRAM_BASE_URL = getEnvStrWithDefault("TELEGRAM_BASE_URL", "https://api.telegram.org")

// MESSAGES
var NATS_MESSAGES_STREAM = getEnvStrWithDefault("NATS_MESSAGES_STREAM", "tg-messages")
var NATS_MESSAGES_CONSUMER = getEnvStrWithDefault("NATS_MESSAGES_CONSUMER", "tg-messages-user-consumer")
var NATS_MESSAGES_SUBJECT = getEnvStrWithDefault("NATS_MESSAGES_SUBJECT", "tg.messages.user")

var RATE_LIMIT_GLOBAL = getEnvIntWithDefault("RATE_LIMIT_GLOBAL", 30) // per sec
var RATE_LIMIT_USER = getEnvIntWithDefault("RATE_LIMIT_USER", 1)      // per sec

func getEnvStrWithDefault(key string, dv string) string {
	v, exists := os.LookupEnv(key)
	if !exists || v == "" {
		fmt.Printf("Using default value for env %s: %s\n", key, dv)
		return dv
	}
	return v
}

func getEnvIntWithDefault(key string, dv int) int {
	v, exists := os.LookupEnv(key)
	if !exists || v == "" {
		fmt.Printf("env var %s not found but is required. using default value: %d\n", key, dv)
		return dv
	}

	intVal, err := strconv.Atoi(v)
	if err != nil {
		fmt.Printf("cannot parse env var %s as int. using default value: %d\n", key, dv)
		return dv
	}

	return intVal
}

func getEnvStr(key string) string {
	v, exists := os.LookupEnv(key)
	if !exists || v == "" {
		fmt.Fprintf(os.Stderr, "error: %v\n", ErrNoEnv{key})
		os.Exit(1)
	}
	return v
}

// func getEnvInt(key string) int {
// 	strVal := getEnvStr(key)

// 	intVal, err := strconv.Atoi(strVal)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "error: %v\n", ErrParseIntFailed{key})
// 		os.Exit(1)
// 	}

// 	return intVal
// }

type ErrNoEnv struct {
	key string
}

func (e ErrNoEnv) Error() string {
	return fmt.Sprintf("env var %s is required", e.key)
}

type ErrParseIntFailed struct {
	key string
}

func (e ErrParseIntFailed) Error() string {
	return fmt.Sprintf("env var %s can not be parsed as int and is required", e.key)
}
