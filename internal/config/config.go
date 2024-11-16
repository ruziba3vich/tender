package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		Server   ServerConfig
		MongoDb  MongoDbConfig
		JWT      JWTConfig
		Email    EmailConfig
		RedisURI string
	}
	JWTConfig struct {
		SecretKey string
	}

	ServerConfig struct {
		Port string
	}
	MongoDbConfig struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
	}

	EmailConfig struct {
		SmtpHost string
		SmtpPort int
		SmtpUser string
		SmtpPass string
	}
)

func (c *Config) Load() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return err
	}

	c.Server.Port = os.Getenv("SERVER_PORT")
	c.MongoDb.Host = os.Getenv("DB_HOST")
	c.MongoDb.Port = os.Getenv("DB_PORT")
	c.MongoDb.User = os.Getenv("DB_USER")
	c.MongoDb.Password = os.Getenv("DB_PASSWORD")
	c.MongoDb.DBName = os.Getenv("DB_NAME")
	c.RedisURI = os.Getenv("REDIS_URI")
	c.JWT.SecretKey = os.Getenv("JWT_SECRET_KEY")

	c.Email.SmtpHost = os.Getenv("SMTP_HOST")
	c.Email.SmtpPort = smtpPort
	c.Email.SmtpUser = os.Getenv("SMTP_USER")
	c.Email.SmtpPass = os.Getenv("SMTP_PASS")

	return nil
}

func New() (*Config, error) {
	var config Config
	if err := config.Load(); err != nil {
		return nil, err
	}
	return &config, nil
}
