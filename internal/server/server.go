package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"gorm.io/gorm"

	"github.com/ivan-ca97/life/pkg/api/http_errors"

	appUpdateApp "github.com/ivan-ca97/life/internal/applications/app_update"
	authenticationApp "github.com/ivan-ca97/life/internal/applications/authentication"
	authorizationApp "github.com/ivan-ca97/life/internal/applications/authorization"
	dataExportApp "github.com/ivan-ca97/life/internal/applications/data_export"
	fitnessAdvisorApp "github.com/ivan-ca97/life/internal/applications/fitness_advisor"
	healthConnectImportApp "github.com/ivan-ca97/life/internal/applications/health_connect_import"
	hevyImportApp "github.com/ivan-ca97/life/internal/applications/hevy_import"
	watchdogApp "github.com/ivan-ca97/life/internal/applications/integrity_watchdog"
	mealAIApp "github.com/ivan-ca97/life/internal/applications/meal_ai"
	aiUsage "github.com/ivan-ca97/life/internal/features/ai_usage"
	"github.com/ivan-ca97/life/internal/features/authentication"
	"github.com/ivan-ca97/life/internal/features/authorization"
	"github.com/ivan-ca97/life/internal/features/daily"
	"github.com/ivan-ca97/life/internal/features/exercise"
	"github.com/ivan-ca97/life/internal/features/food"
	"github.com/ivan-ca97/life/internal/features/goal"
	"github.com/ivan-ca97/life/internal/features/meal"
	"github.com/ivan-ca97/life/internal/features/measurements"
	"github.com/ivan-ca97/life/internal/features/media"
	"github.com/ivan-ca97/life/internal/features/user"
	"github.com/ivan-ca97/life/internal/features/weight"

	"github.com/ivan-ca97/life/internal/infrastructure/llm/openai"
)

type server struct {
	router   chi.Router
	port     int
	watchdog *watchdogApp.WatchdogApplication
}

func NewServer(database *gorm.DB, port int, version, corsOrigins, seedEmail, seedPassword, googleClientId, r2AccountId, r2AccessKeyId, r2SecretAccessKey, r2Bucket, r2PublicURL, githubWebhookSecret, githubToken, openaiApiKey, openaiModel string, watchdogInterval time.Duration) (*server, error) {
	logger := slog.Default()
	errorHandler := http_errors.NewErrorContextBagHandler(logger)

	authorizationFeature := authorization.NewAuthorizationFeature(database)
	authorizer := authorizationFeature.AuthorizationService()

	userFeature := user.NewUserFeature(database, authorizer, errorHandler)
	authenticationFeature := authentication.NewAuthenticationFeature(database, errorHandler)
	authenticationApplication := authenticationApp.NewAuthenticationApplication(
		authenticationFeature.Service(),
		userFeature.Service(),
		authorizationFeature.RoleRepository(),
		authenticationFeature.GoogleVerifier(),
		googleClientId,
		errorHandler,
	)
	authorizationApplication := authorizationApp.NewAuthorizationApplication(authorizationFeature.ShareRepository(), authorizer, userFeature.Service(), errorHandler)
	foodFeature := food.NewFoodFeature(database, authorizer, errorHandler)
	dailyFeature := daily.NewDailyFeature(database, authorizer, errorHandler)
	closureChecker := dailyFeature.DayClosureChecker()
	mealFeature := meal.NewMealFeature(database, authorizer, closureChecker, errorHandler)
	exerciseFeature := exercise.NewExerciseFeature(database, authorizer, closureChecker, errorHandler)
	hevyImportApplication := hevyImportApp.NewHevyImportApplication(
		exerciseFeature.ExerciseService(),
		exerciseFeature.Repository(),
		authorizer,
		errorHandler,
	)
	weightFeature := weight.NewWeightFeature(database, authorizer, closureChecker, errorHandler)
	goalFeature := goal.NewGoalFeature(database, authorizer, errorHandler)
	healthConnectImportApplication := healthConnectImportApp.NewHealthConnectImportApplication(
		database,
		weightFeature.WeightEntryService(),
		weightFeature.Repository(),
		exerciseFeature.ExerciseService(),
		exerciseFeature.Repository(),
		authorizer,
		errorHandler,
	)
	measurementsFeature := measurements.NewMeasurementsFeature(database, authorizer, errorHandler)
	dataExport := dataExportApp.NewDataExportApplication(database, authorizer, errorHandler)
	appUpdate := appUpdateApp.NewAppUpdateApplication(database, errorHandler, githubWebhookSecret, githubToken, r2AccountId, r2AccessKeyId, r2SecretAccessKey, r2Bucket, r2PublicURL)
	mediaFeature := media.NewMediaFeature(r2AccountId, r2AccessKeyId, r2SecretAccessKey, r2Bucket, r2PublicURL, errorHandler)
	fitnessAdvisor := fitnessAdvisorApp.NewFitnessAdvisorApplication(weightFeature.Repository(), authorizer, errorHandler)
	watchdog := watchdogApp.NewWatchdogApplication(database, watchdogInterval, r2AccountId, r2AccessKeyId, r2SecretAccessKey, r2Bucket, r2PublicURL, authorizer, errorHandler)

	// AI: usage/quota management is always available; the meal estimation
	// application is only wired when an OpenAI key is configured (graceful
	// degradation — the rest of the app boots normally without it).
	aiUsageFeature := aiUsage.NewAiUsageFeature(database, authorizer, errorHandler)
	if openaiModel == "" {
		openaiModel = "gpt-4o"
	}
	var mealAI *mealAIApp.MealAIApplication
	if openaiApiKey != "" {
		openaiClient := openai.NewClient(openaiApiKey, openaiModel)
		mealAI = mealAIApp.NewMealAIApplication(foodFeature.FoodService(), openaiClient, aiUsageFeature.QuotaGuard(), aiUsageFeature.InteractionLogger(), aiUsageFeature.Service(), authorizer, errorHandler)
	} else {
		slog.Warn("OPENAI_API_KEY not set; meal AI estimation endpoint disabled")
	}

	router := chi.NewRouter()
	origins := []string{"http://localhost:3000"}
	if corsOrigins != "" {
		origins = strings.Split(corsOrigins, ",")
	}
	corsOptions := cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}
	router.Use(cors.Handler(corsOptions))
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"version": version,
		})
	})

	router.Route("/api/v1", func(r chi.Router) {
		r.Use(http_errors.Middleware)

		authenticationApplication.PublicRoutes(r)
		appUpdate.PublicRoutes(r)

		r.Group(func(r chi.Router) {
			r.Use(authenticationFeature.Middleware().Handle)

			authenticationApplication.ProtectedRoutes(r)
			foodFeature.GlobalRoutes(r)
			userFeature.AdminRoutes(r)
			watchdog.ProtectedRoutes(r)
			aiUsageFeature.Routes(r)

			r.Route("/users/{userId}", func(r chi.Router) {
				userFeature.ProtectedRoutes(r)
				foodFeature.ProtectedRoutes(r)
				mealFeature.ProtectedRoutes(r)
				exerciseFeature.ProtectedRoutes(r)
				hevyImportApplication.ProtectedRoutes(r)
				healthConnectImportApplication.ProtectedRoutes(r)
				weightFeature.ProtectedRoutes(r)
				goalFeature.ProtectedRoutes(r)
				dailyFeature.ProtectedRoutes(r)
				authorizationApplication.ProtectedRoutes(r)
				mediaFeature.ProtectedRoutes(r)
				measurementsFeature.ProtectedRoutes(r)
				fitnessAdvisor.ProtectedRoutes(r)
				dataExport.ProtectedRoutes(r)
				if mealAI != nil {
					mealAI.ProtectedRoutes(r)
				}
			})
		})
	})

	err := seed(userFeature.Service(), seedEmail, seedPassword)
	if err != nil {
		return nil, fmt.Errorf("seeding: %w", err)
	}

	s := &server{
		router:   router,
		port:     port,
		watchdog: watchdog,
	}
	return s, nil
}

func (s *server) Start() error {
	go s.watchdog.Start(context.Background())

	address := fmt.Sprintf(":%d", s.port)
	slog.Info("server listening", "address", address)
	return http.ListenAndServe(address, s.router)
}
