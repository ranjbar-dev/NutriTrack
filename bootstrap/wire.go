package bootstrap

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ranjbar-dev/nutritrack/configs"
	appAdmin "github.com/ranjbar-dev/nutritrack/internal/application/admin"
	appAuth "github.com/ranjbar-dev/nutritrack/internal/application/auth"
	appDietPlan "github.com/ranjbar-dev/nutritrack/internal/application/dietplan"
	appFood "github.com/ranjbar-dev/nutritrack/internal/application/food"
	appFoodRequest "github.com/ranjbar-dev/nutritrack/internal/application/foodrequest"
	appLabResult "github.com/ranjbar-dev/nutritrack/internal/application/labresult"
	appMed "github.com/ranjbar-dev/nutritrack/internal/application/medication"
	appMessage "github.com/ranjbar-dev/nutritrack/internal/application/message"
	appNotif "github.com/ranjbar-dev/nutritrack/internal/application/notification"
	appPush "github.com/ranjbar-dev/nutritrack/internal/application/push"
	appTracking "github.com/ranjbar-dev/nutritrack/internal/application/tracking"
	appUser "github.com/ranjbar-dev/nutritrack/internal/application/user"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	infraAdmin "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/admin"
	pgDietPlan "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/dietplan"
	foodRepo "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/food"
	frInfra "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/foodrequest"
	labResultRepo "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/labresult"
	medRepo "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/medication"
	msgInfra "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/message"
	notifInfra "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/notification"
	pushInfra "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/push"
	dbsqlc "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/sqlc"
	trackInfra "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/tracking"
	"github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/user"
	infraPush "github.com/ranjbar-dev/nutritrack/internal/infrastructure/push"
	redisInfra "github.com/ranjbar-dev/nutritrack/internal/infrastructure/redis"
	"github.com/ranjbar-dev/nutritrack/internal/infrastructure/sms"
	"github.com/ranjbar-dev/nutritrack/internal/infrastructure/storage"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// Container holds all application-level singletons wired together.
type Container struct {
	Cfg                 *configs.Config
	RateLimiter         *redisInfra.RedisRateLimiter
	AuthService         *appAuth.AuthService
	JWTService          *appAuth.JWTService
	OTPStore            *redisInfra.OTPStore
	TokenBlacklist      *redisInfra.TokenBlacklist
	NutritionistService *appUser.NutritionistService
	ClientService       *appUser.ClientService
	LocalStorage        *storage.LocalStorage
	AvatarService       *appUser.AvatarService
	FoodService         *appFood.FoodService
	FoodCategoryService *appFood.FoodCategoryService
	MedicationService   *appMed.MedicationService
	DietPlanService     *appDietPlan.DietPlanService
	TrackingService     *appTracking.TrackingService
	LabResultService    *appLabResult.LabResultService
	MessageService      *appMessage.MessageService
	FoodRequestService  *appFoodRequest.FoodRequestService
	PushService         *appPush.PushService
	NotificationService *appNotif.NotificationService
	AdminService        *appAdmin.AdminService
}

// NewContainer wires all dependencies manually (no code generation needed).
func NewContainer(db *pgxpool.Pool, rdb *redis.Client, cfg *configs.Config) *Container {
	// Infrastructure layer
	userRepo := user.NewPgUserRepository(db)
	pgFoodRepo := foodRepo.NewPgFoodRepository(db)
	cachedFoodRepo := foodRepo.NewCachedFoodRepository(pgFoodRepo, rdb)
	pgCategoryRepo := foodRepo.NewPgFoodCategoryRepository(db)
	otpStore := redisInfra.NewOTPStore(rdb)
	tokenBlacklist := redisInfra.NewTokenBlacklist(rdb)
	rateLimiter := redisInfra.NewRedisRateLimiter(rdb)
	jwtService := appAuth.NewJWTService(&cfg.JWT)

	// SMS provider selection: Kavenegar in production, mock otherwise
	var smsProvider shared.SMSProvider
	if cfg.App.Env == "production" && cfg.SMS.KavenegarAPIKey != "" {
		kavAdapter, err := sms.NewKavenegarAdapter(cfg.SMS.KavenegarAPIKey, cfg.SMS.OTPTemplate)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to initialize Kavenegar SMS adapter")
		}
		smsProvider = kavAdapter
	} else {
		smsProvider = sms.NewMockSMSProvider()
	}

	// Application layer
	authService := appAuth.NewAuthService(userRepo, otpStore, tokenBlacklist, jwtService, smsProvider)
	nutSvc := appUser.NewNutritionistService(userRepo)
	clientSvc := appUser.NewClientService(userRepo)
	localStorage := storage.NewLocalStorage("uploads", "/uploads")
	avatarSvc := appUser.NewAvatarService(userRepo, localStorage)
	foodSvc := appFood.NewFoodService(cachedFoodRepo, pgCategoryRepo)
	catSvc := appFood.NewFoodCategoryService(pgCategoryRepo)
	pgMedRepo := medRepo.NewPgMedicationRepository(db)
	medSvc := appMed.NewMedicationService(pgMedRepo)
	pgPlanRepo := pgDietPlan.NewPgDietPlanRepository(db)
	cachedPlanRepo := pgDietPlan.NewCachedDietPlanRepository(pgPlanRepo, rdb)
	planSvc := appDietPlan.NewDietPlanService(cachedPlanRepo, userRepo)
	pgTrackingRepo := trackInfra.NewPgTrackingRepository(db)
	trackingSvc := appTracking.NewTrackingService(pgTrackingRepo, userRepo)
	pgLabResultRepo := labResultRepo.NewPgLabResultRepository(db)
	labResultSvc := appLabResult.NewLabResultService(pgLabResultRepo, userRepo, localStorage)
	pgMessageRepo := msgInfra.NewPgMessageRepository(db)
	messageSvc := appMessage.NewMessageService(pgMessageRepo, userRepo, localStorage)
	pgFoodRequestRepo := frInfra.NewPgFoodRequestRepository(db)
	foodRequestSvc := appFoodRequest.NewFoodRequestService(pgFoodRequestRepo, userRepo, foodSvc)
	pgPushRepo := pushInfra.NewPgPushSubscriptionRepository(db)
	webpushSender := infraPush.NewWebpushSender(cfg.VAPID.PublicKey, cfg.VAPID.PrivateKey, "mailto:info@nutritrack.ir")
	pushSvc := appPush.NewPushService(pgPushRepo, webpushSender)
	pgNotifPrefRepo := notifInfra.NewPgNotificationPreferenceRepository(db)
	notifSvc := appNotif.NewNotificationService(pgNotifPrefRepo)
	adminSvc := appAdmin.NewAdminService(infraAdmin.NewPgAdminRepository(dbsqlc.New(db)))

	return &Container{
		Cfg:                 cfg,
		RateLimiter:         rateLimiter,
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
		LabResultService:    labResultSvc,
		MessageService:      messageSvc,
		FoodRequestService:  foodRequestSvc,
		PushService:         pushSvc,
		NotificationService: notifSvc,
		AdminService:        adminSvc,
	}
}
