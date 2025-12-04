package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/sachanritik1/go-lang/internal/store"
	"github.com/sachanritik1/go-lang/internal/utils"
)

type UserHandler struct {
	store  store.UserStore
	logger *log.Logger
}

func NewUserHandler(store store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{store: store, logger: logger}
}

type RegisterUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
	Bio      *string `json:"bio,omitempty"`
}

func (h *UserHandler) validateRegisterUserRequest(req *RegisterUserRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}
	if len(req.Username) < 3 || len(req.Username) > 10 {
		return errors.New("username must be between 3 and 10 characters")
	}

	if req.Email == "" {
		return errors.New("email is required")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}

	// TODO: Uncomment and use the following regex for stronger password validation
	// passwordRegex := regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$`)
	// if !passwordRegex.MatchString(req.Password) {
	// 	return errors.New("password must be at least 8 characters long, contain at least one uppercase letter, one lowercase letter, one number, and one special character")
	// }

	return nil
}

func (h *UserHandler) HandlerRegisterUser(w http.ResponseWriter, r *http.Request) {
	var registerUserRequest RegisterUserRequest
	err := json.NewDecoder(r.Body).Decode(&registerUserRequest)
	if err != nil {
		h.logger.Printf("ERROR: decoding register user request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	err = h.validateRegisterUserRequest(&registerUserRequest)
	if err != nil {
		h.logger.Printf("ERROR: validating register user request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	user := &store.User{
		Username: registerUserRequest.Username,
		Email:    registerUserRequest.Email,
	}
	if registerUserRequest.Bio != "" {
		user.Bio = registerUserRequest.Bio
	}

	err = user.PasswordHash.Set(registerUserRequest.Password)
	if err != nil {
		h.logger.Printf("ERROR: setting password hash: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "could not process password"})
		return
	}

	err = h.store.CreateUser(user)
	if err != nil {
		h.logger.Printf("ERROR: creating user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "could not create user"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"user": user})

}

func (h *UserHandler) HandleGetUserByID(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("ERROR: reading ID parameter: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid user ID parameter"})
		return
	}

	user, err := h.store.GetUserByID(int(userID))
	if err != nil {
		h.logger.Printf("ERROR: getting user by ID: %v", err)
		if err == sql.ErrNoRows {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "user not found"})
		} else {
			utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "could not retrieve user"})
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user": user})
}
