package config

import (
	"context"
	"database/sql"

	"github.com/bsonger/devflow-config-service/pkg/app"
	"github.com/bsonger/devflow-config-service/pkg/domain"
	configrepo "github.com/bsonger/devflow-config-service/pkg/infra/config_repo"
	"github.com/bsonger/devflow-config-service/pkg/infra/store"
	"github.com/bsonger/devflow-service-common/observability"
	"github.com/spf13/viper"
)

type ConfigRepoConfig struct {
	RootDir    string `mapstructure:"root_dir" json:"root_dir" yaml:"root_dir"`
	DefaultRef string `mapstructure:"default_ref" json:"default_ref" yaml:"default_ref"`
}

type Config struct {
	Server         *domain.ServerConfig   `mapstructure:"server" json:"server" yaml:"server"`
	Postgres       *domain.PostgresConfig `mapstructure:"postgres" json:"postgres" yaml:"postgres"`
	Log            *domain.LogConfig      `mapstructure:"log" json:"log" yaml:"log"`
	Otel           *domain.OtelConfig     `mapstructure:"otel" json:"otel" yaml:"otel"`
	ConfigRepo     *ConfigRepoConfig      `mapstructure:"config_repo" json:"config_repo" yaml:"config_repo"`
	AppServiceBase string                 `mapstructure:"app_service_base_url" json:"app_service_base_url" yaml:"app_service_base_url"`
	Pyroscope      string                 `mapstructure:"pyroscope" json:"pyroscope" yaml:"pyroscope"`
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

	db, err := sql.Open("pgx", stringValue(config.Postgres, func(v *domain.PostgresConfig) string { return v.DSN }))
	if err != nil {
		return shutdown, err
	}
	store.ApplyPool(db,
		intValue(config.Postgres, func(v *domain.PostgresConfig) int { return v.MaxOpenConns }),
		intValue(config.Postgres, func(v *domain.PostgresConfig) int { return v.MaxIdleConns }),
		intValue(config.Postgres, func(v *domain.PostgresConfig) int { return v.ConnMaxLifetimeMinutes }),
	)
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return shutdown, err
	}

	store.InitPostgres(db)
	configrepo.DefaultRepository = ResolveConfigRepo(config)
	app.ConfigureAppConfigRepository(configrepo.DefaultRepository)
	app.ConfigureEnvironmentResolver(app.ResolveEnvironmentResolver(resolveAppServiceBaseURL(config)))
	return func(shutdownCtx context.Context) error {
		closeErr := db.Close()
		shutdownErr := shutdown(shutdownCtx)
		if shutdownErr != nil {
			return shutdownErr
		}
		return closeErr
	}, nil
}

func ResolveConfigRepo(cfg *Config) *configrepo.Repository {
	if cfg == nil || cfg.ConfigRepo == nil || cfg.ConfigRepo.RootDir == "" {
		return nil
	}

	defaultRef := cfg.ConfigRepo.DefaultRef
	if defaultRef == "" {
		defaultRef = "main"
	}

	return configrepo.NewRepository(configrepo.Options{
		RootDir:    cfg.ConfigRepo.RootDir,
		DefaultRef: defaultRef,
	})
}

func ResolveConfigPort(cfg *Config) int {
	if cfg == nil || cfg.Server == nil {
		return 0
	}
	return cfg.Server.Port
}

func resolveAppServiceBaseURL(cfg *Config) string {
	if cfg != nil && cfg.AppServiceBase != "" {
		return cfg.AppServiceBase
	}
	return "http://devflow-app-service.devflow-staging.svc.cluster.local:8081"
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
