package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/sachanritik1/go-lang/internal/app"
)

func SetupRoutes(app *app.App) *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(app.Middleware.Authenticate)

		// Workout routes
		// Protected route but AnonymousUser can create workouts
		// r.Get("/workouts", app.WorkoutHandler.HandlerGetAllWorkouts)

		// The following routes require an authenticated user
		r.Get("/workouts", app.Middleware.RequireUser(app.WorkoutHandler.HandlerGetAllWorkouts))
		r.Get("/workouts/{id}", app.Middleware.RequireUser(app.WorkoutHandler.HandlerGetWorkoutByID))

		r.Post("/workouts", app.Middleware.RequireUser(app.WorkoutHandler.HandlerCreateWorkout))
		r.Put("/workouts/{id}", app.Middleware.RequireUser(app.WorkoutHandler.HandlerUpdateWorkout))
		r.Delete("/workouts/{id}", app.Middleware.RequireUser(app.WorkoutHandler.HandlerDeleteWorkout))

		r.Get("/users/{id}", app.Middleware.RequireUser(app.UserHandler.HandleGetUserByID))
		// r.Put("/users/{id}", app.Middleware.RequireUser(app.UserHandler.HandlerUpdateUser))
		// r.Delete("/users/{id}", app.Middleware.RequireUser(app.UserHandler.HandlerDeleteUser))

	})

	// Public routes
	r.Get("/health", app.HealthCheckHandler)

	r.Post("/users", app.UserHandler.HandlerRegisterUser)
	r.Post("/tokens/authentication", app.TokenHandler.HandleCreateToken)

	return r
}
