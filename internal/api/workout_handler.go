package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/sachanritik1/go-lang/internal/store"
	"github.com/sachanritik1/go-lang/internal/utils"
)

type WorkoutHandler struct {
	store  store.WorkoutStore
	logger *log.Logger
}

func NewWorkoutHandler(store store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{store: store, logger: logger}
}

func (h *WorkoutHandler) HandlerCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		h.logger.Printf("ERROR: decoding create workout request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	createdWorkout, err := h.store.CreateWorkout(&workout)
	if err != nil {
		h.logger.Printf("ERROR: creating workout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "could not create workout"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"workout": createdWorkout})
}

func (h *WorkoutHandler) HandlerGetWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("ERROR: reading ID parameter: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout ID parameter"})
		return
	}

	workout, err := h.store.GetWorkoutByID(int(workoutID))
	if err != nil {
		h.logger.Printf("ERROR: getting workout by ID: %v", err)
		if err == sql.ErrNoRows {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout not found"})
		} else {
			utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "could not retrieve workout"})
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": workout})

}

func (h *WorkoutHandler) HandlerGetAllWorkouts(w http.ResponseWriter, r *http.Request) {
	workouts, err := h.store.ListWorkouts()
	if err != nil {
		h.logger.Printf("ERROR: listing workouts: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "could not retrieve workouts"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workouts": workouts})
}

func (h *WorkoutHandler) HandlerDeleteWorkout(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("ERROR: reading ID parameter: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout ID parameter"})
		return
	}

	err = h.store.DeleteWorkout(int(workoutID))
	if err != nil {
		h.logger.Printf("ERROR: deleting workout: %v", err)
		if err == sql.ErrNoRows {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout not found"})
		} else {
			utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "could not delete workout"})
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "workout deleted successfully"})
}

func (h *WorkoutHandler) HandlerUpdateWorkout(w http.ResponseWriter, r *http.Request) {

	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("ERROR: reading ID parameter: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout ID parameter"})
		return
	}

	// find workout by id
	workout, err := h.store.GetWorkoutByID(int(workoutID))
	if err != nil {
		h.logger.Printf("ERROR: getting workout by ID: %v", err)
		if err == sql.ErrNoRows {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout not found"})
		} else {
			utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "could not retrieve workout"})
		}
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
		h.logger.Printf("ERROR: decoding update workout request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
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
		h.logger.Printf("ERROR: updating workout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "could not update workout"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": updatedWorkout})
}
