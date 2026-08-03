package di

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_trans "github.com/go-playground/validator/v10/translations/en"
	"gorm.io/gorm"

	clientEmail "github.com/reshap0318/hadirYuk/internal/clients/email"
	clientOpenai "github.com/reshap0318/hadirYuk/internal/clients/openai"
	"github.com/reshap0318/hadirYuk/internal/database"
	"github.com/reshap0318/hadirYuk/internal/handlers"
	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/repositories"
	"github.com/reshap0318/hadirYuk/internal/services"
)

// Container holds all dependencies.
type Container struct {
	DB           *gorm.DB
	AiReadOnlyDB *gorm.DB
	Redis        *database.RedisCache
	Access       *helpers.Access
	EmailClient  *clientEmail.EmailClient
	Logger       *helpers.Logger
	RateLimiter  *helpers.RateLimiter
	Repositories *repositories.Repositories
	Services     *services.Services
	Handlers     *handlers.Handlers
}

// Close closes all connections.
func (c *Container) Close() error {
	if c.DB != nil {
		sqlDB, err := c.DB.DB()
		if err != nil {
			return fmt.Errorf("error getting database connection: %w", err)
		}
		if err := sqlDB.Close(); err != nil {
			return fmt.Errorf("error closing database connection: %w", err)
		}
		log.Println("Database connection closed")
	}

	if c.AiReadOnlyDB != nil {
		if sqlDB, err := c.AiReadOnlyDB.DB(); err == nil {
			sqlDB.Close()
			log.Println("AI readonly database connection closed")
		}
	}

	if c.Redis != nil {
		log.Println("Redis connection closed")
	}

	if c.Logger != nil {
		c.Logger.Close()
		log.Println("Logger closed")
	}

	return nil
}

// NewContainer creates and initializes all dependencies.
func NewContainer() (*Container, error) {
	container := &Container{}

	// Initialize Logger (early, before other components)
	logger, err := helpers.NewLogger("storage/logs")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}
	container.Logger = logger

	// Determine database connection type (default: mysql)
	dbConnection := helpers.GetEnv("DB_CONNECTION", "mysql")

	// Initialize database based on DB_CONNECTION
	if dbConnection == "postgres" || dbConnection == "postgresql" {
		postgres, err := database.NewPostgreSQL(database.PostgreSQLConfig{
			Host:     helpers.GetEnv("DB_HOST", "127.0.0.1"),
			Port:     helpers.GetEnv("DB_PORT", "5432"),
			User:     helpers.GetEnv("DB_USERNAME", "postgres"),
			Password: helpers.GetEnv("DB_PASSWORD", ""),
			DBName:   helpers.GetEnv("DB_DATABASE", "boilerplate"),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to initialize PostgreSQL: %w", err)
		}
		container.DB = postgres
	} else {
		// Default to MySQL
		mysql, err := database.NewMySQL(database.MySQLConfig{
			Host:     helpers.GetEnv("DB_HOST", "127.0.0.1"),
			Port:     helpers.GetEnv("DB_PORT", "3306"),
			User:     helpers.GetEnv("DB_USERNAME", "root"),
			Password: helpers.GetEnv("DB_PASSWORD", ""),
			DBName:   helpers.GetEnv("DB_DATABASE", "boilerplate"),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to initialize MySQL: %w", err)
		}
		container.DB = mysql
	}

	// Initialize Redis (optional)
	if helpers.GetEnv("REDIS_ENABLED", "false") == "true" {
		redisClient, err := database.NewRedis(database.RedisConfig{
			Host:     helpers.GetEnv("REDIS_HOST", "localhost"),
			Port:     helpers.GetEnv("REDIS_PORT", "6379"),
			User:     helpers.GetEnv("REDIS_USER", "root"),
			Password: helpers.GetEnv("REDIS_PASSWORD", ""),
			DB:       helpers.GetEnvInt("REDIS_DB", 0),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Redis: %w", err)
		}
		container.Redis = database.NewRedisCache(redisClient)
	}

	// Initialize Email Client (optional)
	container.EmailClient = clientEmail.NewEmailClient()

	// Initialize Rate Limiter
	rateLimitRequests := helpers.GetEnvInt("RATE_LIMIT_REQUESTS", 100)
	rateLimitWindow := helpers.GetEnvInt("RATE_LIMIT_WINDOW", 60)
	container.RateLimiter = helpers.NewRateLimiter(rateLimitRequests, rateLimitWindow)
	log.Printf("Rate limiting enabled: %d requests per %d seconds", rateLimitRequests, rateLimitWindow)

	// DB is required for repositories and services
	if container.DB == nil {
		panic("database connection is required but not initialized. Set DB_CONNECTION=mysql or DB_CONNECTION=postgres in your .env file")
	}

	// Always initialize repositories
	repos, err := repositories.NewRepositories(container.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize repositories: %w", err)
	}
	container.Repositories = repos

	// Always initialize services (Redis can be nil)
	container.Services = services.NewServices(container.Repositories, container.Redis, container.EmailClient, container.Logger)

	// Initialize JWKS Manager
	passphrase, err := helpers.LoadPassphrase(
		helpers.GetEnv("JWT_PASSPHRASE_PATH", "storage/keys/passphrase"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load JWT passphrase: %w", err)
	}

	jwksManager := &services.JWKSManager{}
	if err := jwksManager.Initialize(
		helpers.GetEnv("JWT_PRIVATE_KEY_PATH", "storage/keys/private.pem"),
		helpers.GetEnv("JWT_PUBLIC_KEY_PATH", "storage/keys/public.pem"),
		passphrase,
	); err != nil {
		return nil, fmt.Errorf("failed to initialize JWKS Manager: %w", err)
	}
	container.Services.JWKSManager = jwksManager

	// Initialize Access
	acc := helpers.NewAccess(container.Redis, container.DB)
	container.Access = acc
	container.Services.Access = acc

	// Initialize Validator with translator
	validate := validator.New(validator.WithRequiredStructEnabled())

	// Use JSON tag as field name for validation errors
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	uni := ut.New(en.New(), en.New())
	trans, _ := uni.GetTranslator("en")
	en_trans.RegisterDefaultTranslations(validate, trans)

	// Always initialize handlers
	container.Handlers = handlers.NewHandlers(container.Services, validate, trans)

	// Initialize AI Chat feature (readonly DB, OpenAI client, session store, rate limiters)
	if aiUser := helpers.GetEnv("AI_CHAT_DB_USERNAME", ""); aiUser != "" {
		aiReadOnlyDB, err := database.NewMySQLReadOnly(database.MySQLConfig{
			Host:     helpers.GetEnv("DB_HOST", "127.0.0.1"),
			Port:     helpers.GetEnv("DB_PORT", "3306"),
			User:     aiUser,
			Password: helpers.GetEnv("AI_CHAT_DB_PASSWORD", ""),
			DBName:   helpers.GetEnv("DB_DATABASE", "hadir_yuk"),
		})
		if err != nil {
			log.Printf("Warning: AI chat readonly DB not initialized: %v", err)
		} else {
			container.AiReadOnlyDB = aiReadOnlyDB
			container.Services.AiReadOnlyDB = aiReadOnlyDB
		}
	} else {
		log.Println("Warning: AI_CHAT_DB_USERNAME not set, AI chat query execution will fail")
	}

	container.Services.AiChatClient = clientOpenai.NewClient(
		helpers.GetEnv("AI_CHAT_API_KEY", ""),
		helpers.GetEnv("AI_CHAT_MODEL", "gpt-4o-mini"),
		helpers.GetEnv("AI_CHAT_BASE_URL", ""),
	)
	container.Services.AiChatStore = clientOpenai.NewStore(
		time.Duration(helpers.GetEnvInt("AI_CHAT_SESSION_TTL_MINUTES", 1440))*time.Minute,
		helpers.GetEnvInt("AI_CHAT_MAX_STORED_MESSAGES", 100),
	)
	container.Services.AiChatCfg = &services.AiChatConfig{
		QueryTimeout:       time.Duration(helpers.GetEnvInt("AI_CHAT_QUERY_TIMEOUT_SECONDS", 10)) * time.Second,
		LLMTimeout:         120 * time.Second, // hardcoded, not from .env - see AiChatConfig doc comment
		MaxRows:            helpers.GetEnvInt("AI_CHAT_MAX_ROWS", 100),
		MaxContextMessages: helpers.GetEnvInt("AI_CHAT_MAX_CONTEXT_MESSAGES", 10),
	}
	// Windows are hardcoded per FSD §2 point 8 (4h global / 1h per-user), not from .env.
	container.Services.AiChatGlobalLimiter = helpers.NewRateLimiter(helpers.GetEnvInt("AI_CHAT_RATE_LIMIT_GLOBAL_MAX", 200), 4*3600)
	container.Services.AiChatUserLimiter = helpers.NewRateLimiter(helpers.GetEnvInt("AI_CHAT_RATE_LIMIT_USER_MAX", 20), 3600)

	return container, nil
}
