package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host         string
	Password     string
	TLSEnable    bool
	DialTimeout  time.Duration
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
}

type Product struct {
	product *redis.Client 
}

func NewProduct(ctx context.Context, cfg Config) (*Product, error) {

	opts := &redis.Options{
		Addr:         cfg.Host,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		Password:     cfg.Password,
	}

	if cfg.TLSEnable {
		opts.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	product := redis.NewClient(opts)

	err := product.Ping(ctx).Err()
	if err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	return &Product{product: product}, nil
}

func (p *Product) Close() error {
	err := p.product.Close()
	if err != nil {
		return fmt.Errorf("close: %w", err)
	}

	return nil
}

func (p *Product) Ping(ctx context.Context) error {
	err := p.product.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("redis: %w", err)
	}

	return nil
}

func (p *Product) Unwrap() *redis.Client {
	return p.product
}