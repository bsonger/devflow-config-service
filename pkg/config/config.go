package config

import (
	"context"
	"database/sql"

	"github.com/bsonger/devflow-config-service/pkg/model"
	"github.com/bsonger/devflow-config-service/pkg/store"
	"github.com/bsonger/devflow-service-common/observability"
	"github.com/spf13/viper"
)

type Config struct {
	Server    *model.ServerConfig   `mapstructure:"server" json:"server" yaml:"server"`
	Postgres  *model.PostgresConfig `mapstructure:"postgres" json:"postgres" yaml:"postgres"`
	Log       *model.LogConfig      `mapstructure:"log" json:"log" yaml:"log"`
	Otel      *model.OtelConfig     `mapstructure:"otel" json:"otel" yaml:"otel"`
	Pyroscope string                `mapstructure:"pyroscope" json:"pyroscope" yaml:"pyroscope"`
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

	db, err := sql.Open("pgx", stringValue(config.Postgres, func(v *model.PostgresConfig) string { return v.DSN }))
	if err != nil {
		return shutdown, err
	}
	store.ApplyPool(db,
		intValue(config.Postgres, func(v *model.PostgresConfig) int { return v.MaxOpenConns }),
		intValue(config.Postgres, func(v *model.PostgresConfig) int { return v.MaxIdleConns }),
		intValue(config.Postgres, func(v *model.PostgresConfig) int { return v.ConnMaxLifetimeMinutes }),
	)
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return shutdown, err
	}

	store.InitPostgres(db)
	return func(shutdownCtx context.Context) error {
		closeErr := db.Close()
		shutdownErr := shutdown(shutdownCtx)
		if shutdownErr != nil {
			return shutdownErr
		}
		return closeErr
	}, nil
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

func stringValue[T any](value *T, getter func(*T) string) string {
	if value == nil {
		return ""
	}
	return getter(value)
}

func intValue[T any](value *T, getter func(*T) int) int {
	if value == nil {
		return 0
	}
	return getter(value)
}
