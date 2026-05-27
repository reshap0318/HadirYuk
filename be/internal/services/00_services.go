package services

import (
	"time"

	"github.com/reshap0318/hadirYuk/internal/clients/email"
	"github.com/reshap0318/hadirYuk/internal/database"
	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/repositories"
)

// JWTConfig holds JWT configuration.
type JWTConfig struct {
	Expiration time.Duration
	RefreshExp time.Duration
}

// Services holds all service dependencies.
type Services struct {
	repo         *repositories.Repositories
	RedisClient  *database.RedisCache
	EmailClient  *email.EmailClient
	Logger       *helpers.Logger
	JWKSManager  *JWKSManager
	Access       *helpers.Access
	cfg          *JWTConfig
}

// NewServices creates and initializes all services.
func NewServices(repo *repositories.Repositories, redisClient *database.RedisCache, emailClient *email.EmailClient, logger *helpers.Logger) *Services {
	return &Services{
		repo:        repo,
		RedisClient: redisClient,
		EmailClient: emailClient,
		Logger:      logger,
		cfg: &JWTConfig{
			Expiration: time.Duration(helpers.GetEnvInt("JWT_EXPIRATION", 24)) * time.Hour,
			RefreshExp: time.Duration(helpers.GetEnvInt("JWT_REFRESH_EXPIRATION", 168)) * time.Hour,
		},
	}
}
