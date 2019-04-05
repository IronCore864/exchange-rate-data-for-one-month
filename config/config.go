package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

func getEnvIntValue(key string, fallback int) int {
	value, exists := os.LookupEnv(key)
	if exists {
		i, _ := strconv.Atoi(value)
		return i
	}
	return fallback
}

func getEnvStrValue(key string, fallback string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return fallback
}

// Config is the struct contains the configuration
type Config struct {
	RedisHost     string
	RedisPort     int
	RedisPassword string
}

// LoadConfiguration loads config from config/config.json as default values
// if ENV vars are set, values are overwritten by ENV var values.
// Possible ENV vars are: Port, ReqPerSec, ReqPerMin, ReqPerHour, RedisHost, RedisPort
func loadConfiguration() Config {
	var conf Config
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	filename := fmt.Sprintf("%s/config.json", basepath)
	configFile, err := os.Open(filename)
	defer configFile.Close()
	if err != nil {
		log.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&conf)
	// Use ENV var if it's set
	conf.RedisHost = getEnvStrValue("RedisHost", conf.RedisHost)
	conf.RedisPort = getEnvIntValue("RedisPort", conf.RedisPort)
	conf.RedisPassword = getEnvStrValue("RedisPassword", conf.RedisPassword)
	return conf
}

// Conf is the configutation struct object
var Conf = loadConfiguration()
