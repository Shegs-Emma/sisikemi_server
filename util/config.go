package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBSource string `mapstructure:"DB_SOURCE"`
	Environment string `mapstructure:"ENVIRONMENT"`
	AllowedOrigins []string `mapstructure:"ALLOWED_ORIGINS"`
	HttpServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
	GrpcServerAddress string `mapstructure:"GRPC_SERVER_ADDRESS"`
	RedisAddress string `mapstructure:"REDIS_ADDRESS"`
	MigrationURL string `mapstructure:"MIGRATION_URL"`
	TokenSymmetricKey string `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	EmailSenderName string `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress string `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword string `mapstructure:"EMAIL_SENDER_PASSWORD"`
	SendGridApiKey string `mapstructure:"SENDGRID_API_KEY"`
	CloudinaryUrl string `mapstructure:"CLOUDINARY_URL"`
	MailDomain string `mapstructure:"MAIL_DOMAIN"`
	MailHost string `mapstructure:"MAIL_HOST"`
	MailPort int `mapstructure:"MAIL_PORT"`
	MailEncryption string `mapstructure:"MAIL_ENCRYPTION"`
	MailUsername string `mapstructure:"MAIL_USERNAME"`
	MailPassword string `mapstructure:"MAIL_PASSWORD"`
	FromName string `mapstructure:"FROM_NAME"`
	FromAddress string `mapstructure:"FROM_ADDRESS"`
	ClientEndpoint string `mapstructure:"CLIENT_ENDPOINT"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
