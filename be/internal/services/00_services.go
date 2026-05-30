package services

import (
	"fmt"
	"strconv"
	"time"

	"github.com/reshap0318/hadirYuk/internal/clients/email"
	"github.com/reshap0318/hadirYuk/internal/clients/face"
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
	FaceService  *FaceService
	cfg          *JWTConfig
}

// NewServices creates and initializes all services.
func NewServices(repo *repositories.Repositories, redisClient *database.RedisCache, emailClient *email.EmailClient, logger *helpers.Logger) *Services {
	services := &Services{
		repo:        repo,
		RedisClient: redisClient,
		EmailClient: emailClient,
		Logger:      logger,
		cfg: &JWTConfig{
			Expiration: time.Duration(helpers.GetEnvInt("JWT_EXPIRATION", 24)) * time.Hour,
			RefreshExp: time.Duration(helpers.GetEnvInt("JWT_REFRESH_EXPIRATION", 168)) * time.Hour,
		},
	}

	// Initialize FaceClient and FaceService
	faceClient := face.NewFaceClient()

	// Panic if face recognition models fail to load (only applies when built with -tags gocv)
	if err := faceClient.InitError(); err != nil {
		panic(fmt.Sprintf("Face recognition initialization failed: %v. Pastikan file model tersedia: haarcascade_frontalface_default.xml, nn4.small2.v1.t7", err))
	}

	// Override threshold from env if set
	if threshold := helpers.GetEnv("FACE_THRESHOLD", ""); threshold != "" {
		if t, err := strconv.ParseFloat(threshold, 64); err == nil {
			faceClient.SetThreshold(t)
		}
	}

	services.FaceService = NewFaceService(faceClient, logger)

	return services
}
