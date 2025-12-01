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

	return r
}
