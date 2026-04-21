package router

import (
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
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.ErrorHandler())

	// Health check (public, no auth)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
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
		authGroup.POST("/otp/send",   authHandler.SendOTP)
		authGroup.POST("/otp/verify", authHandler.VerifyOTP)
		authGroup.POST("/refresh",    authHandler.Refresh)
	}

	// --- Protected routes (require JWT) ---
	protected := v1.Group("")
	protected.Use(middleware.RequireAuth(container.JWTService))
	{
		protected.POST("/auth/logout", authHandler.Logout)
		protected.GET("/auth/me",      authHandler.Me)

		// Avatar upload (access control is service-level)
		avatarHandler := handler.NewAvatarHandler(container.AvatarService)
		protected.PUT("/users/:id/avatar", avatarHandler.Upload)

		// Super admin: nutritionist management
		nutHandler := handler.NewNutritionistHandler(container.NutritionistService)
		adminGroup := protected.Group("/admin")
		adminGroup.Use(middleware.RequireRole(middleware.RoleSuperAdmin))
		{
			adminGroup.POST("/nutritionists",             nutHandler.Create)
			adminGroup.GET("/nutritionists",              nutHandler.List)
			adminGroup.GET("/nutritionists/:id",          nutHandler.Get)
			adminGroup.PATCH("/nutritionists/:id",        nutHandler.Update)
			adminGroup.PATCH("/nutritionists/:id/status", nutHandler.SetStatus)
		}

		// Nutritionist: client management
		clientHandler := handler.NewClientHandler(container.ClientService)
		clientGroup := protected.Group("/clients")
		clientGroup.Use(middleware.RequireRole(middleware.RoleNutritionist))
		{
			clientGroup.POST("",      clientHandler.RegisterClient)
			clientGroup.GET("",       clientHandler.ListClients)
			clientGroup.GET("/:id",   clientHandler.GetClientProfile)
			clientGroup.PATCH("/:id", clientHandler.UpdateClient)
		}

		// Foods: available to any authenticated user; role-based logic is in the service
		foodHandler := handler.NewFoodHandler(container.FoodService)
		protected.POST("/foods",       foodHandler.Create)
		protected.GET("/foods",        foodHandler.Search)
		protected.GET("/foods/:id",    foodHandler.GetOne)
		protected.PATCH("/foods/:id",  foodHandler.Update)
		protected.DELETE("/foods/:id", foodHandler.Delete)

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

		// Diet plans: nutritionist creates plans for clients; clients/nutritionists can view
		planHandler := handler.NewDietPlanHandler(container.DietPlanService)
		protected.POST("/clients/:id/plans",                                                               planHandler.CreatePlan)
		protected.GET("/clients/:id/plans",                                                                planHandler.ListClientPlans)
		protected.GET("/plans/:id",                                                                        planHandler.GetPlan)
		protected.POST("/plans/:id/days",                                                                  planHandler.AddDay)
		protected.POST("/plans/:id/days/:day_id/meals",                                                    planHandler.AddMeal)
		protected.POST("/plans/:id/days/:day_id/meals/:meal_id/options",                                   planHandler.AddOption)
		protected.POST("/plans/:id/days/:day_id/meals/:meal_id/options/:option_id/items",                  planHandler.AddItem)
		protected.DELETE("/plans/:id/days/:day_id/meals/:meal_id/options/:option_id/items/:item_id",       planHandler.RemoveItem)
		protected.DELETE("/plans/:id",                                                                     planHandler.DeletePlan)
	}

	// 404 handler
	r.NoRoute(middleware.NotFound())

	return r
}
