package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/ranjbar-dev/nutritrack/bootstrap"
	"github.com/ranjbar-dev/nutritrack/configs"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/handler"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/middleware"
)

func New(db *pgxpool.Pool, rdb *redis.Client, cfg *configs.Config) *gin.Engine {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Wire all dependencies via DI container
	container := bootstrap.NewContainer(db, rdb, cfg)

	r := gin.New()

	// Global middleware chain (order matters)
	r.Use(middleware.CORS(cfg.App.CORSAllowedOrigins))
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.ErrorHandler())

	// Health check (public, no auth)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "nutritrack",
			"version": "1.0.0",
		})
	})

	// Serve uploaded files (avatars etc.) from local filesystem
	r.Static("/uploads", "./uploads")

	// API v1
	v1 := r.Group("/api/v1")

	// --- Auth routes (public) ---
	authHandler := handler.NewAuthHandler(container.AuthService)
	authGroup := v1.Group("/auth")
	{
		authGroup.POST("/login",      authHandler.Login)
		authGroup.POST("/otp/send",   middleware.RateLimitByIP(rdb, 60), authHandler.SendOTP)
		authGroup.POST("/otp/verify", authHandler.VerifyOTP)
		authGroup.POST("/refresh",    authHandler.Refresh)
	}

	// --- Protected routes (require JWT) ---
	protected := v1.Group("")
	protected.Use(middleware.RequireAuth(container.JWTService, container.TokenBlacklist))
	{
		protected.POST("/auth/logout", authHandler.Logout)
		protected.GET("/auth/me",      authHandler.Me)

		// Avatar upload (access control is service-level)
		avatarHandler := handler.NewAvatarHandler(container.AvatarService)
		protected.PUT("/users/:id/avatar", avatarHandler.Upload)

		// Super admin: nutritionist management + platform stats
		nutHandler := handler.NewNutritionistHandler(container.NutritionistService)
		adminHandler := handler.NewAdminHandler(container.AdminService)
		adminGroup := protected.Group("/admin")
		adminGroup.Use(middleware.RequireRole(middleware.RoleSuperAdmin))
		{
			adminGroup.GET("/stats",                      adminHandler.GetStats)
			adminGroup.POST("/nutritionists",             nutHandler.Create)
			adminGroup.GET("/nutritionists",              nutHandler.List)
			adminGroup.GET("/nutritionists/:id",          nutHandler.Get)
			adminGroup.PATCH("/nutritionists/:id",        nutHandler.Update)
			adminGroup.PATCH("/nutritionists/:id/status", nutHandler.SetStatus)
			adminGroup.GET("/nutritionists/:id/clients", nutHandler.GetClients)
		}

		// Nutritionist: client management
		clientHandler := handler.NewClientHandler(container.ClientService)
		clientGroup := protected.Group("/clients")
		clientGroup.Use(middleware.RequireRole(middleware.RoleNutritionist))
		{
			clientGroup.POST("",      clientHandler.RegisterClient)
			clientGroup.GET("",       clientHandler.ListClients)
			clientGroup.GET("/:id",   clientHandler.GetClientProfile)
			clientGroup.PATCH("/:id",        clientHandler.UpdateClient)
			clientGroup.PATCH("/:id/status", clientHandler.SetStatus)
		}

		// Foods: available to any authenticated user; role-based logic is in the service
		foodHandler := handler.NewFoodHandler(container.FoodService)
		protected.POST("/foods",       foodHandler.Create)
		protected.GET("/foods",        foodHandler.Search)
		protected.GET("/foods/:id",    foodHandler.GetOne)
		protected.PATCH("/foods/:id",  foodHandler.Update)
		protected.DELETE("/foods/:id", foodHandler.Delete)
		// Admin: list all foods and force-delete any food
		adminGroup.GET("/foods",        foodHandler.Search)
		adminGroup.DELETE("/foods/:id", foodHandler.Delete)

		// Food categories
		catHandler := handler.NewFoodCategoryHandler(container.FoodCategoryService)
		protected.GET("/food-categories", catHandler.ListAll)
		adminGroup.POST("/food-categories",        catHandler.Create)
		adminGroup.DELETE("/food-categories/:id",  catHandler.Delete)

		// Medications: any authenticated user; role-based logic is in the service
		medHandler := handler.NewMedicationHandler(container.MedicationService)
		protected.POST("/medications",       medHandler.Create)
		protected.GET("/medications",        medHandler.Search)
		protected.GET("/medications/:id",    medHandler.GetOne)
		protected.PATCH("/medications/:id",  medHandler.Update)
		protected.DELETE("/medications/:id", medHandler.Delete)
		// Admin: list all medications and force-delete any medication
		adminGroup.GET("/medications",        medHandler.Search)
		adminGroup.DELETE("/medications/:id", medHandler.Delete)

		// Diet plans: nutritionist creates plans for clients; clients/nutritionists can view
		planHandler := handler.NewDietPlanHandler(container.DietPlanService, container.PushService)
		protected.GET("/plans/active",                                                                     planHandler.GetActivePlan)
		protected.POST("/clients/:id/plans",                                                               planHandler.CreatePlan)
		protected.GET("/clients/:id/plans",                                                                planHandler.ListClientPlans)
		protected.GET("/plans/:id",                                                                        planHandler.GetPlan)
		protected.PATCH("/plans/:id",                                                                      planHandler.UpdatePlan)
		protected.POST("/plans/:id/days",                                                                  planHandler.AddDay)
		protected.DELETE("/plans/:id/days/:day_id",                                                        planHandler.DeleteDay)
		protected.POST("/plans/:id/days/:day_id/meals",                                                    planHandler.AddMeal)
		protected.DELETE("/plans/:id/days/:day_id/meals/:meal_id",                                         planHandler.DeleteMeal)
		protected.POST("/plans/:id/days/:day_id/meals/:meal_id/options",                                   planHandler.AddOption)
		protected.DELETE("/plans/:id/days/:day_id/meals/:meal_id/options/:option_id",                      planHandler.DeleteOption)
		protected.POST("/plans/:id/days/:day_id/meals/:meal_id/options/:option_id/items",                  planHandler.AddItem)
		protected.DELETE("/plans/:id/days/:day_id/meals/:meal_id/options/:option_id/items/:item_id",       planHandler.RemoveItem)
		protected.DELETE("/plans/:id",                                                                     planHandler.DeletePlan)
		protected.POST("/plans/:id/days/:day_id/exercises",                                                planHandler.AddExercise)
		protected.DELETE("/plans/:id/days/:day_id/exercises/:exercise_id",                                 planHandler.RemoveExercise)
		protected.POST("/plans/:id/days/:day_id/prescriptions",                                            planHandler.AddPrescription)
		protected.DELETE("/plans/:id/days/:day_id/prescriptions/:prescription_id",                         planHandler.RemovePrescription)

		// Tracking — clients log their daily activity; nutritionists + superadmin can read
		trackingHandler := handler.NewTrackingHandler(container.TrackingService)
		protected.POST("/tracking/food",       trackingHandler.LogFood)
		protected.POST("/tracking/water",      trackingHandler.LogWater)
		protected.POST("/tracking/sleep",      trackingHandler.LogSleep)
		protected.POST("/tracking/exercise",   trackingHandler.LogExercise)
		protected.POST("/tracking/medication", trackingHandler.LogMedication)
		protected.POST("/tracking/body",       trackingHandler.LogBody)
		protected.POST("/tracking/sync",       trackingHandler.BulkSync)
		protected.GET("/clients/:id/tracking", trackingHandler.GetTracking)

		// Lab results — clients and nutritionists upload/download lab result files
		labResultHandler := handler.NewLabResultHandler(container.LabResultService)
		protected.POST("/clients/:id/lab-results", labResultHandler.Upload)
		protected.GET("/clients/:id/lab-results",  labResultHandler.List)
		protected.GET("/lab-results/:id/download", labResultHandler.Download)

		// Messages — client-nutritionist chat
		messageHandler := handler.NewMessageHandler(container.MessageService, container.PushService)
		protected.GET("/messages/unread-count",     messageHandler.GetUnreadCount)
		protected.GET("/messages",                  messageHandler.GetClientMessages)
		protected.POST("/messages",                 messageHandler.SendAsClient)
		protected.GET("/clients/:id/messages",      messageHandler.GetNutritionistMessages)
		protected.POST("/clients/:id/messages",     messageHandler.SendAsNutritionist)

		// Food requests — client submits, nutritionist approves/rejects
		frHandler := handler.NewFoodRequestHandler(container.FoodRequestService, container.PushService)
		protected.POST("/food-requests",             frHandler.Submit)
		protected.GET("/food-requests",              frHandler.ListPending)
		protected.POST("/food-requests/:id/approve", frHandler.Approve)
		protected.POST("/food-requests/:id/reject",  frHandler.Reject)

		// Push subscriptions — Web Push notification management
		pushHandler := handler.NewPushHandler(container.PushService)
		protected.POST("/push/subscribe",   pushHandler.Subscribe)
		protected.DELETE("/push/subscribe", pushHandler.Unsubscribe)

		// Notification preferences — user controls which notifications they receive
		notifHandler := handler.NewNotificationHandler(container.NotificationService)
		protected.PATCH("/notifications/preferences", notifHandler.UpdatePreferences)
		protected.GET("/notifications/preferences",   notifHandler.GetPreferences)
	}

	// 404 handler
	r.NoRoute(middleware.NotFound())

	return r
}
