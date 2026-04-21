package bootstrap

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	appAuth "github.com/ranjbar-dev/nutritrack/internal/application/auth"
	appFood "github.com/ranjbar-dev/nutritrack/internal/application/food"
	appUser "github.com/ranjbar-dev/nutritrack/internal/application/user"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	foodRepo "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/food"
	"github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/user"
	redisInfra "github.com/ranjbar-dev/nutritrack/internal/infrastructure/redis"
	"github.com/ranjbar-dev/nutritrack/internal/infrastructure/sms"
	"github.com/ranjbar-dev/nutritrack/internal/infrastructure/storage"
	"github.com/ranjbar-dev/nutritrack/configs"
)

// Container holds all application-level singletons wired together.
type Container struct {
	AuthService          *appAuth.AuthService
	JWTService           *appAuth.JWTService
	OTPStore             *redisInfra.OTPStore
	TokenBlacklist       *redisInfra.TokenBlacklist
	NutritionistService  *appUser.NutritionistService
	ClientService        *appUser.ClientService
	LocalStorage         *storage.LocalStorage
	AvatarService        *appUser.AvatarService
	FoodService          *appFood.FoodService
	FoodCategoryService  *appFood.FoodCategoryService
}

// NewContainer wires all dependencies manually (no code generation needed).
func NewContainer(db *pgxpool.Pool, rdb *redis.Client, cfg *configs.Config) *Container {
	// Infrastructure layer
	userRepo       := user.NewPgUserRepository(db)
	pgFoodRepo     := foodRepo.NewPgFoodRepository(db)
	pgCategoryRepo := foodRepo.NewPgFoodCategoryRepository(db)
	otpStore       := redisInfra.NewOTPStore(rdb)
	tokenBlacklist := redisInfra.NewTokenBlacklist(rdb)
	jwtService     := appAuth.NewJWTService(&cfg.JWT)

	// SMS provider selection: Kavenegar in production, mock otherwise
	var smsProvider shared.SMSProvider
	if cfg.App.Env == "production" && cfg.SMS.KavenegarAPIKey != "" {
		smsProvider = sms.NewKavenegarAdapter(cfg.SMS.KavenegarAPIKey, cfg.SMS.OTPTemplate)
	} else {
		smsProvider = sms.NewMockSMSProvider()
	}

	// Application layer
	authService := appAuth.NewAuthService(userRepo, otpStore, tokenBlacklist, jwtService, smsProvider)
	nutSvc := appUser.NewNutritionistService(userRepo)
	clientSvc := appUser.NewClientService(userRepo)
	localStorage := storage.NewLocalStorage("uploads", "/uploads")
	avatarSvc := appUser.NewAvatarService(userRepo, localStorage)
	foodSvc := appFood.NewFoodService(pgFoodRepo, pgCategoryRepo)
	catSvc := appFood.NewFoodCategoryService(pgCategoryRepo)

	return &Container{
		AuthService:         authService,
		JWTService:          jwtService,
		OTPStore:            otpStore,
		TokenBlacklist:      tokenBlacklist,
		NutritionistService: nutSvc,
		ClientService:       clientSvc,
		LocalStorage:        localStorage,
		AvatarService:       avatarSvc,
		FoodService:         foodSvc,
		FoodCategoryService: catSvc,
	}
}
