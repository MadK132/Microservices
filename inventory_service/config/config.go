package config

import (
	"github.com/recktt77/Microservices-First-/inventory_service/pkg/mongo"
	"time"

	"github.com/caarlos0/env/v6"
)

type (
	Config struct {
		Mongo  mongo.Config
		Server Server

		Version string `env:"VERSION"`

		NATSUrl string `env:"NATS_URL" env-default:"nats://localhost:4222"`
		Redis  RedisConfig         
		Cache  CacheConfig  
	}

	Server struct {
		HTTPServer HTTPServer
		GRPCServer GRPCServer
	}

	HTTPServer struct {
		Port           int           `env:"HTTP_PORT,required"`
		ReadTimeout    time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"30s"`
		WriteTimeout   time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"30s"`
		IdleTimeout    time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"60s"`
		MaxHeaderBytes int           `env:"HTTP_MAX_HEADER_BYTES" envDefault:"1048576"`
		TrustedProxies []string      `env:"HTTP_TRUSTED_PROXIES" envSeparator:","`
		Mode           string        `env:"GIN_MODE" envDefault:"release"`
	}
	

	GRPCServer struct {
		Port                  int16         `env:"GRPC_PORT,notEmpty"`
		MaxRecvMsgSizeMiB     int           `env:"GRPC_MAX_MESSAGE_SIZE_MIB" envDefault:"12"`
		MaxConnectionAge      time.Duration `env:"GRPC_MAX_CONNECTION_AGE" envDefault:"30s"`
		MaxConnectionAgeGrace time.Duration `env:"GRPC_MAX_CONNECTION_AGE_GRACE" envDefault:"10s"`
	}

	RedisConfig struct {
		Host         string        `env:"REDIS_HOST,notEmpty"`
		Password     string        `env:"REDIS_PASSWORD"`
		TLSEnable    bool          `env:"REDIS_TLS_ENABLE" envDefault:"false"`
		DialTimeout  time.Duration `env:"REDIS_DIAL_TIMEOUT" envDefault:"5s"`
		WriteTimeout time.Duration `env:"REDIS_WRITE_TIMEOUT" envDefault:"5s"`
		ReadTimeout  time.Duration `env:"REDIS_READ_TIMEOUT" envDefault:"5s"`
	}

	CacheConfig struct {
		ProductTTL time.Duration `env:"PRODUCT_CACHE_TTL" envDefault:"24h"`
	}
)

func New() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)

	return &cfg, err
}