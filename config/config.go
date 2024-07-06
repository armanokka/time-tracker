package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type PostgresConfig struct {
	User     string `env:"POSTGRES_USER" env-required:"true"`
	DB       string `env:"POSTGRES_DB" env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	Driver   string `env:"POSTGRES_DRIVER" env-default:"pgx"`
	Host     string `env:"POSTGRES_HOST" env-required:"true"`
	Port     int    `env:"POSTGRES_PORT" env-required:"true"`
}

type LoggerConfig struct {
	EnableCaller     bool   `env:"LOGGER_ENABLE_CALLER" env-default:"false"`
	EnableStacktrace bool   `env:"LOGGER_ENABLE_STACKTRACE" env-default:"false"`
	Encoding         string `env:"LOGGER_ENCODING" env-default:"console"` // json/console
	Level            string `env:"LOGGER_LEVEL" env-default:"true"`       // debug/info/warn/error/dpanic/panic/fatal
}

type ServerConfig struct {
	Mode         string `env:"SERVER_MODE" env-default:"development"`
	Port         int    `env:"SERVER_PORT" env-default:"80"`
	PprofPort    int    `env:"SERVER_PPROF_PORT" env-default:"6053"`
	JWTSecretKey string `env:"SERVER_JWT_SECRET_KEY" env-required:"true"`
}

type RedisConfig struct {
	Addr         string `env:"REDIS_ADDR" env-default:"redis:6379"`
	DB           int    `env:"REDIS_DB" env-default:"0"`
	Password     string `env:"REDIS_PASSWORD" env-required:"true"`
	MinIdleConns int    `env:"REDIS_MIN_IDLE_CONNS" env-default:"200"`
	PoolTimeout  int    `env:"REDIS_POOL_TIMEOUT" env-default:"240"`
	PoolSize     int    `env:"REDIS_POOL_SIZE" env-default:"12000"`
}

type CookieConfig struct {
	Name     string `env:"COOKIE_NAME" env-default:"jwt-token"`
	MaxAge   int    `env:"COOKIE_MAX_AGE" env-default:"86400"`
	Secure   bool   `env:"COOKIE_SECURE" env-default:"false"`
	HttpOnly bool   `env:"COOKIE_HTTP_ONLY" env-default:"true"`
}

type Config struct {
	Postgres PostgresConfig
	Redis    RedisConfig
	Cookie   CookieConfig
	Logger   LoggerConfig
	Server   ServerConfig
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadConfig(".env", cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
