package internal

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"reflect"

	"github.com/spf13/viper"
)

type Env string

type Config struct {
	Environment Env
	Name        string
	Version     string
	Server      *ServerConfig
	Database    *DatabaseConfig
	Logging     *LoggingConfig
}

type ServerConfig struct {
	Host string
	Port string
}

type DatabaseConfig struct {
	BaseConnStr string
	User        string
	Password    string
}

type LoggingConfig struct {
	DefaultLevel   slog.Level
	SpecificLogger map[string]slog.Level
}

func LoadConfig() (*Config, error) {
	err := viper.BindEnv("environment", "APP_ENV")
	if err != nil {
		log.Fatalf("failed to fetch environment: %v", err)
	}

	println("Environment:", viper.GetString("environment"))

	readConfigFiles(
		"config.toml",
		fmt.Sprintf("config.%s.toml", viper.GetString("environment")))

	var config Config
	// as I'm overriding the default decode hook, keep an eye out for something that may break because of it
	err = viper.Unmarshal(&config, viper.DecodeHook(parseLogLevelHook))
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %v", err)
	}

	return &config, nil
}

func parseLogLevelHook(from, to reflect.Value) (interface{}, error) {
	if from.Kind() == reflect.String && to.Type() == reflect.TypeOf(slog.LevelInfo) {
		return parseLogLevel(from.String())
	}
	return from.Interface(), nil
}

func readConfigFiles(filePaths ...string) {
	var err error = nil
	for _, filePath := range filePaths {
		viper.SetConfigFile(filePath)
		err = viper.MergeInConfig()
		if err != nil && !isConfigFileNotFoundError(err) {
			log.Fatalf("failed to read %s config file: %v", filePath, err)
		}
	}
}

func isConfigFileNotFoundError(err error) bool {
	var pathError *fs.PathError
	return errors.As(err, &pathError)
}

func parseLogLevel(level string) (slog.Level, error) {
	switch level {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("unknown log minLevel: %s", level)
	}
}

func (e Env) IsDev() bool {
	return e == "development"
}
