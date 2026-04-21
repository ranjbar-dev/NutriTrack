package bootstrap

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	appAuth "github.com/ranjbar-dev/nutritrack/internal/application/auth"
	appDietPlan "github.com/ranjbar-dev/nutritrack/internal/application/dietplan"
	appFood "github.com/ranjbar-dev/nutritrack/internal/application/food"
	appMed "github.com/ranjbar-dev/nutritrack/internal/application/medication"
	appTracking "github.com/ranjbar-dev/nutritrack/internal/application/tracking"
	appUser "github.com/ranjbar-dev/nutritrack/internal/application/user"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	pgDietPlan "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/dietplan"
	foodRepo "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/food"
	medRepo "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/medication"
	trackInfra "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/tracking"
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
	MedicationService    *appMed.MedicationService
	DietPlanService      *appDietPlan.DietPlanService
	TrackingService      *appTracking.TrackingService
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
	pgMedRepo := medRepo.NewPgMedicationRepository(db)
	medSvc := appMed.NewMedicationService(pgMedRepo)
	pgPlanRepo := pgDietPlan.NewPgDietPlanRepository(db)
	planSvc := appDietPlan.NewDietPlanService(pgPlanRepo, userRepo)
	pgTrackingRepo := trackInfra.NewPgTrackingRepository(db)
	trackingSvc := appTracking.NewTrackingService(pgTrackingRepo, userRepo)

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
		MedicationService:   medSvc,
		DietPlanService:     planSvc,
		TrackingService:     trackingSvc,
	}
}
