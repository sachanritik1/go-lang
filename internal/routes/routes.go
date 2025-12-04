package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/sachanritik1/go-lang/internal/app"
)

func SetupRoutes(app *app.App) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/health", app.HealthCheckHandler)

	// Workout routes
	r.Post("/workouts", app.WorkoutHandler.HandlerCreateWorkout)
	r.Get("/workouts/{id}", app.WorkoutHandler.HandlerGetWorkoutByID)
	r.Get("/workouts", app.WorkoutHandler.HandlerGetAllWorkouts)
	r.Put("/workouts/{id}", app.WorkoutHandler.HandlerUpdateWorkout)
	r.Delete("/workouts/{id}", app.WorkoutHandler.HandlerDeleteWorkout)

	// User routes
	r.Post("/users", app.UserHandler.HandlerRegisterUser)
	// r.Post("/users/login", app.UserHandler.HandlerRegisterUser)
	r.Get("/users/{id}", app.UserHandler.HandleGetUserByID)
	// r.Put("/users/{id}", app.UserHandler.HandlerUpdateUser)
	// r.Delete("/users/{id}", app.UserHandler.HandlerDeleteUser)

	return r
}
