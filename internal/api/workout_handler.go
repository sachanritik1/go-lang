package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
		fmt.Println(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	createdWorkout, err := h.store.CreateWorkout(&workout)
	if err != nil {
		http.Error(w, "Failed to create workout"+strings.Split(err.Error(), "\n")[0], http.StatusInternalServerError)
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
		http.Error(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}

	workout, err := h.store.GetWorkoutByID(int(workoutID))
	if err != nil {
		http.Error(w, "Failed to get workout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workout)

}

func (h *WorkoutHandler) HandlerGetAllWorkouts(w http.ResponseWriter, r *http.Request) {
	workouts, err := h.store.ListWorkouts()
	if err != nil {
		http.Error(w, "Failed to get workouts"+strings.Split(err.Error(), "\n")[0], http.StatusInternalServerError)
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
		http.Error(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}

	err = h.store.DeleteWorkout(int(workoutID))
	if err != nil {
		http.Error(w, "Failed to delete workout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *WorkoutHandler) HandlerUpdateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)

	if err != nil {
		http.Error(w, "Invalid request payload"+strings.Split(err.Error(), "\n")[0], http.StatusBadRequest)
		return
	}

	paramWorkoutID := chi.URLParam(r, "id")

	if paramWorkoutID == "" {
		http.Error(w, "Workout ID is required", http.StatusBadRequest)
		return
	}

	workoutID, err := strconv.ParseInt(paramWorkoutID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}

	workout.ID = int(workoutID)

	updatedWorkout, err := h.store.UpdateWorkout(&workout)
	if err != nil {
		http.Error(w, "Failed to update workout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedWorkout)
}
