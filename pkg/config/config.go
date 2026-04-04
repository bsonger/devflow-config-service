package config

import (
	"context"
	"github.com/bsonger/devflow-common/client/logging"
	"github.com/bsonger/devflow-common/client/mongo"
	commonModel "github.com/bsonger/devflow-common/model"
	"github.com/bsonger/devflow-config-service/pkg/model"
	"github.com/bsonger/devflow-config-service/pkg/store"
	"github.com/bsonger/devflow-service-common/observability"
	"github.com/spf13/viper"
)

type Config struct {
	Server    *model.ServerConfig `mapstructure:"server" json:"server" yaml:"server"`
	Mongo     *model.MongoConfig  `mapstructure:"mongo"  json:"mongo"  yaml:"mongo"`
	Log       *model.LogConfig    `mapstructure:"log"    json:"log"    yaml:"log"`
	Otel      *model.OtelConfig   `mapstructure:"otel"   json:"otel"   yaml:"otel"`
	Pyroscope string              `mapstructure:"pyroscope" json:"pyroscope" yaml:"pyroscope"`
}

func Load() (*Config, error) {
	v := viper.New()
	//v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config/")
	v.AddConfigPath("/etc/devflow/config/")

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var config *Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func InitConfig(ctx context.Context, config *Config) error {
	_, err := InitRuntime(ctx, config, "")
	return err
}

func InitRuntime(ctx context.Context, config *Config, serviceName string) (func(context.Context) error, error) {
	shutdown, err := observability.Init(ctx, observability.RuntimeOptions{
		LogLevel:        safeLogLevel(config),
		LogFormat:       safeLogFormat(config),
		OtelEndpoint:    safeOtelEndpoint(config),
		OtelService:     safeOtelServiceName(config),
		PyroscopeAddr:   config.Pyroscope,
		ServiceOverride: serviceName,
	})
	if err != nil {
		return nil, err
	}

	client, err := mongo.InitMongo(ctx, toCommonMongoConfig(config.Mongo), logging.Logger)
	if err != nil {
		return shutdown, err
	}
	store.InitMongo(client, config.Mongo.DBName)
	return shutdown, nil
}

func toCommonMongoConfig(cfg *model.MongoConfig) *commonModel.MongoConfig {
	if cfg == nil {
		return nil
	}
	return &commonModel.MongoConfig{
		URI:    cfg.URI,
		DBName: cfg.DBName,
	}
}

func ResolveConfigPort(cfg *Config) int {
	if cfg == nil || cfg.Server == nil {
		return 0
	}
	return cfg.Server.Port
}

func safeLogLevel(cfg *Config) string {
	if cfg != nil && cfg.Log != nil {
		return cfg.Log.Level
	}
	return ""
}

func safeLogFormat(cfg *Config) string {
	if cfg != nil && cfg.Log != nil {
		return cfg.Log.Format
	}
	return ""
}

func safeOtelEndpoint(cfg *Config) string {
	if cfg != nil && cfg.Otel != nil {
		return cfg.Otel.Endpoint
	}
	return ""
}

func safeOtelServiceName(cfg *Config) string {
	if cfg != nil && cfg.Otel != nil {
		return cfg.Otel.ServiceName
	}
	return ""
}
