package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/sachanritik1/go-lang/internal/store"
	"github.com/sachanritik1/go-lang/internal/tokens"
	"github.com/sachanritik1/go-lang/internal/utils"
)

type TokenHandler struct {
	store     store.TokenStore
	userStore store.UserStore
	logger    *log.Logger
}

func NewTokenHandler(store store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{store: store, userStore: userStore, logger: logger}
}

type CreateTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req CreateTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("ERROR: decoding create token request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{
			"error": "invalid request payload",
		})
		return
	}

	user, err := h.userStore.GetUserByUsername(req.Username)
	if err != nil {
		h.logger.Printf("ERROR: getting user by username: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{
			"error": "internal server error",
		})
		return
	}

	passwordDoMatch, err := user.PasswordHash.Matches(req.Password)
	if err != nil {
		h.logger.Printf("ERROR: checking password match: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{
			"error": "internal server error",
		})
		return
	}
	if !passwordDoMatch {
		h.logger.Printf("WARNING: invalid password for user %s", req.Username)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{
			"error": "invalid credentials",
		})
		return
	}

	token, err := h.store.CreateNewToken(int64(user.ID), int64(24*time.Hour), tokens.ScopeAuth)
	if err != nil {
		h.logger.Printf("ERROR: creating new token: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{
			"error": "internal server error",
		})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{
		"auth_token": token,
	})

}
