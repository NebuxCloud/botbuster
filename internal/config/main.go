package config

import "time"

type Config struct {
	Debug bool `env:"DEBUG" envDefault:"false"`

	ListenPort     string   `env:"LISTEN_PORT" envDefault:"8000"`
	AllowedOrigins []string `env:"ALLOWED_ORIGINS"`

	HmacKey    string        `env:"HMAC_KEY"`
	Expiration time.Duration `env:"EXPIRATION" envDefault:"10m"`

	ValkeyURL    string `env:"VALKEY_URL"`
	ValkeyPrefix string `env:"VALKEY_PREFIX" envDefault:"botbuster:"`
}
