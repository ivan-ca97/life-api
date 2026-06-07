package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"gorm.io/gorm"

	"github.com/ivan-ca97/life/pkg/api/http_errors"

	authorizationApp "github.com/ivan-ca97/life/internal/applications/authorization"
	"github.com/ivan-ca97/life/internal/features/authorization"
	featureAuth "github.com/ivan-ca97/life/internal/features/auth"
	"github.com/ivan-ca97/life/internal/features/daily"
	"github.com/ivan-ca97/life/internal/features/exercise"
	"github.com/ivan-ca97/life/internal/features/food"
	"github.com/ivan-ca97/life/internal/features/goal"
	"github.com/ivan-ca97/life/internal/features/meal"
	"github.com/ivan-ca97/life/internal/features/user"
	"github.com/ivan-ca97/life/internal/features/weight"
)

type server struct {
	router chi.Router
	port   int
}

func NewServer(database *gorm.DB, port int, version, corsOrigins, seedEmail, seedPassword, googleClientId string) (*server, error) {
	logger := slog.Default()
	errorHandler := http_errors.NewErrorContextBagHandler(logger)

	authorizationFeature := authorization.NewAuthorizationFeature(database)
	authorizer := authorizationFeature.AuthorizationService()

	userFeature := user.NewUserFeature(database, authorizer, errorHandler)
	authFeature := featureAuth.NewAuthFeature(database, userFeature.Service(), authorizationFeature.RoleRepository(), errorHandler, googleClientId)
	authorizationApplication := authorizationApp.NewAuthorizationApplication(authorizationFeature.ShareRepository(), authorizer, userFeature.Service(), errorHandler)
	foodFeature := food.NewFoodFeature(database, authorizer, errorHandler)
	mealFeature := meal.NewMealFeature(database, authorizer, errorHandler)
	exerciseFeature := exercise.NewExerciseFeature(database, authorizer, errorHandler)
	weightFeature := weight.NewWeightFeature(database, authorizer, errorHandler)
	goalFeature := goal.NewGoalFeature(database, authorizer, errorHandler)
	dailyFeature := daily.NewDailyFeature(database, authorizer, errorHandler)

	router := chi.NewRouter()
	origins := []string{"http://localhost:3000"}
	if corsOrigins != "" {
		origins = strings.Split(corsOrigins, ",")
	}
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"version": version,
		})
	})

	router.Route("/api/v1", func(r chi.Router) {
		r.Use(http_errors.Middleware)

		authFeature.PublicRoutes(r)

		r.Group(func(r chi.Router) {
			r.Use(authFeature.Middleware().Handle)

			authFeature.ProtectedRoutes(r)
			foodFeature.GlobalRoutes(r)
			userFeature.AdminRoutes(r)

			r.Route("/users/{userId}", func(r chi.Router) {
				userFeature.ProtectedRoutes(r)
				foodFeature.ProtectedRoutes(r)
				mealFeature.ProtectedRoutes(r)
				exerciseFeature.ProtectedRoutes(r)
				weightFeature.ProtectedRoutes(r)
				goalFeature.ProtectedRoutes(r)
				dailyFeature.ProtectedRoutes(r)
				authorizationApplication.ProtectedRoutes(r)
			})
		})
	})

	err := seed(userFeature.Service(), seedEmail, seedPassword)
	if err != nil {
		return nil, fmt.Errorf("seeding: %w", err)
	}

	s := &server{
		router: router,
		port:   port,
	}
	return s, nil
}

func (s *server) Start() error {
	address := fmt.Sprintf(":%d", s.port)
	slog.Info("server listening", "address", address)
	err := http.ListenAndServe(address, s.router)
	if err != nil {
		return err
	}
	return nil
}
