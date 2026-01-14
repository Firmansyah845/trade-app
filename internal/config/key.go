package config

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func mustGetString(key string) string {
	mustHave(key)
	return viper.GetString(key)
}

// func mustGetBool(key string) bool {
// 	mustHave(key)
// 	return viper.GetBool(key)
// }

func mustGetStringArray(key string) []string {
	mustHave(key)
	return optionalGetStringArray(key)
}

func mustGetInt(key string) int {
	mustHave(key)
	v, err := strconv.Atoi(viper.GetString(key))
	if err != nil {
		panic(fmt.Sprintf("key %s is not a valid Integer value", key))
	}

	return v
}

func mustGetDurationMs(key string) time.Duration {
	return time.Millisecond * time.Duration(mustGetInt(key))
}

func mustGetDurationMinute(key string) time.Duration {
	return time.Minute * time.Duration(mustGetInt(key))
}

func optionalGetStringArray(key string) []string {
	value := viper.GetString(key)
	if value == "" {
		return []string{}
	}

	strS := strings.Split(value, ",")
	for i, str := range strS {
		strS[i] = strings.TrimSpace(str)
	}

	return strS
}

func mustHave(key string) {
	if !viper.IsSet(key) {
		panic(fmt.Sprintf("key %s is not set", key))
	}
}
