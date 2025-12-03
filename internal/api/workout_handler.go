package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sachanritik1/go-lang/internal/store"
)

type WorkoutHandler struct {
	store store.WorkoutStore
}

func NewWorkoutHandler(store store.WorkoutStore) *WorkoutHandler {
	return &WorkoutHandler{store: store}
}

func (h *WorkoutHandler) HandlerCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdWorkout, err := h.store.CreateWorkout(&workout)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdWorkout)
}

func (h *WorkoutHandler) HandlerGetWorkoutByID(w http.ResponseWriter, r *http.Request) {
	paramWorkoutID := chi.URLParam(r, "id")

	if paramWorkoutID == "" {
		http.Error(w, "Workout ID is required", http.StatusBadRequest)
		return
	}

	workoutID, err := strconv.ParseInt(paramWorkoutID, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	workout, err := h.store.GetWorkoutByID(int(workoutID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workout)

}

func (h *WorkoutHandler) HandlerGetAllWorkouts(w http.ResponseWriter, r *http.Request) {
	workouts, err := h.store.ListWorkouts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workouts)
}

func (h *WorkoutHandler) HandlerDeleteWorkout(w http.ResponseWriter, r *http.Request) {
	paramWorkoutID := chi.URLParam(r, "id")

	if paramWorkoutID == "" {
		http.Error(w, "Workout ID is required", http.StatusBadRequest)
		return
	}

	workoutID, err := strconv.ParseInt(paramWorkoutID, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.store.DeleteWorkout(int(workoutID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *WorkoutHandler) HandlerUpdateWorkout(w http.ResponseWriter, r *http.Request) {

	paramWorkoutID := chi.URLParam(r, "id")

	if paramWorkoutID == "" {
		http.Error(w, "Workout ID is required", http.StatusBadRequest)
		return
	}

	workoutID, err := strconv.ParseInt(paramWorkoutID, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// find workout by id
	workout, err := h.store.GetWorkoutByID(int(workoutID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// make a struct for update
	var UpdateWorkoutRequest struct {
		Title           *string              `json:"title"`
		Description     *string              `json:"description"`
		DurationMinutes *int                 `json:"duration_minutes"`
		CaloriesBurned  *int                 `json:"calories_burned"`
		Entries         []store.WorkoutEntry `json:"workout_entries"`
	}

	err = json.NewDecoder(r.Body).Decode(&UpdateWorkoutRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if UpdateWorkoutRequest.Title != nil {
		workout.Title = *UpdateWorkoutRequest.Title
	}
	if UpdateWorkoutRequest.Description != nil {
		workout.Description = *UpdateWorkoutRequest.Description
	}
	if UpdateWorkoutRequest.DurationMinutes != nil {
		workout.DurationMinutes = *UpdateWorkoutRequest.DurationMinutes
	}
	if UpdateWorkoutRequest.CaloriesBurned != nil {
		workout.CaloriesBurned = *UpdateWorkoutRequest.CaloriesBurned
	}
	if UpdateWorkoutRequest.Entries != nil {
		workout.Entries = UpdateWorkoutRequest.Entries
	}

	// update workout
	updatedWorkout, err := h.store.UpdateWorkout(workout)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedWorkout)
}
