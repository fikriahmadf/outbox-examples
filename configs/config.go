package configs

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"sync"
	"time"
)

// Config is a struct that will receive configuration options via environment
// variables.
type Config struct {
	App struct {
		CORS struct {
			AllowCredentials bool     `mapstructure:"ALLOW_CREDENTIALS"`
			AllowedHeaders   []string `mapstructure:"ALLOWED_HEADERS"`
			AllowedMethods   []string `mapstructure:"ALLOWED_METHODS"`
			AllowedOrigins   []string `mapstructure:"ALLOWED_ORIGINS"`
			Enable           bool     `mapstructure:"ENABLE"`
			MaxAgeSeconds    int      `mapstructure:"MAX_AGE_SECONDS"`
		}
		Name     string `mapstructure:"NAME"`
		URL      string `mapstructure:"URL"`
		Revision string `mapstructure:"REVISION"`
	}

	DB struct {
		Postgres struct {
			Read struct {
				Host            string        `mapstructure:"HOST"`
				Port            string        `mapstructure:"PORT"`
				User            string        `mapstructure:"USER"`
				Password        string        `mapstructure:"PASSWORD"`
				Name            string        `mapstructure:"NAME"`
				SSLMode         string        `mapstructure:"SSLMODE"`
				MaxConnLifetime time.Duration `mapstructure:"MAX_CONNECTION_LIFETIME"`
				MaxIdleConn     int           `mapstructure:"MAX_IDLE_CONNECTION"`
				MaxOpenConn     int           `mapstructure:"MAX_OPEN_CONNECTION"`
			}
			Write struct {
				Host            string        `mapstructure:"HOST"`
				Port            string        `mapstructure:"PORT"`
				User            string        `mapstructure:"USER"`
				Password        string        `mapstructure:"PASSWORD"`
				Name            string        `mapstructure:"NAME"`
				SSLMode         string        `mapstructure:"SSLMODE"`
				MaxConnLifetime time.Duration `mapstructure:"MAX_CONNECTION_LIFETIME"`
				MaxIdleConn     int           `mapstructure:"MAX_IDLE_CONNECTION"`
				MaxOpenConn     int           `mapstructure:"MAX_OPEN_CONNECTION"`
			}
		} `mapstructure:"PG"`
	}

	Server struct {
		Env string `mapstructure:"ENV"`
		Log struct {
			Level  string `mapstructure:"LEVEL"`
			Pretty bool   `mapstructure:"PRETTY"`
			Color  bool   `mapstructure:"COLOR"`
		} `mapstructure:"LOG"`
		Port     string `mapstructure:"PORT"`
		Shutdown struct {
			CleanupPeriodSeconds int64 `mapstructure:"CLEANUP_PERIOD_SECONDS"`
			GracePeriodSeconds   int64 `mapstructure:"GRACE_PERIOD_SECONDS"`
		}
		LogLevel string `mapstructure:"LOG_LEVEL"`
	}
	Email struct {
		Memo struct {
			Recipient string `mapstructure:"RECIPIENT"`
		}
	}
	External struct {
		N8N struct {
			BaseURL       string        `mapstructure:"BASE_URL"`
			RetryCount    int           `mapstructure:"RETRY_COUNT"`
			RetryWaitTime time.Duration `mapstructure:"RETRY_WAIT_TIME"`
			Endpoints     struct {
				SendMemoNotifPath string `mapstructure:"SEND_MEMO_NOTIF_PATH"`
			}
		}
	}
}

func (c Config) IsServerEnvDevelopment() bool { return c.Server.Env == "development" }

var (
	conf Config
	once sync.Once
)

// Get are responsible to load env and get data an return the struct
func Get() *Config {
	once.Do(func() {
		viper.SetConfigFile(".env")
		err := viper.ReadInConfig()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed reading config file")
		}
		log.Info().Msg("Service configuration initialized.")
		err = viper.Unmarshal(&conf)
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}
	})

	return &conf
}
