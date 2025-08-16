package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"social-api/cmd/utils"
	"social-api/internal/store"
	"social-api/model"

	"github.com/google/uuid"
)
type AuthenticationPayload struct {
	Name     string `json:"name" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

// @Summary Register a new user
// @Description Registers a new user with the provided email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param input body AuthenticationPayload true "User registration payload"
// @Success 201 {object} model.User
// @Failure 400 {object} error
// @Router /authentication/user [post]
func (app *Application) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	payload := &AuthenticationPayload{}
	ctx := r.Context()
	if err := utils.ReadJson(w, r, payload); err != nil {
		app.Logger.Errorf("Failed to read JSON", "error", err)
		app.Err.BadRequestError(w, r, err)
		return
	}
	if err := utils.Validator.Struct(payload); err != nil {
		app.Logger.Errorf("Validation error", "error", err)
		app.Err.BadRequestError(w, r, err)
		return
	}

	user := &model.User{
		Name:  payload.Name,
		Email: payload.Email,
	}

	if err := user.Password.Set(payload.Password); err != nil {
		app.Logger.Errorf("Failed to set password", "error", err)
		app.Err.BadRequestError(w, r, err)
		return
	}

	plainToken := uuid.New().String()
	hash := sha256.Sum256([]byte(plainToken))
	hashedToken := hex.EncodeToString(hash[:])

	err := app.Store.Users.CreateAndInvitation(ctx, user, hashedToken, app.Config.Email.ExpiryTime)
	if err != nil {
		switch err {
		case store.ErrEmailTaken:
			app.Logger.Errorf("Email is already taken", "email", payload.Email)
			app.Err.BadRequestError(w, r, err)
		case store.ErrNameTaken:
			app.Logger.Errorf("Name is already taken", "name", payload.Name)
			app.Err.BadRequestError(w, r, err)
		default:
			app.Logger.Errorf("Failed to register user", "error", err)
			app.Err.InternalServerError(w, r, err)
		}
		return
	}
	tokenRes := TokenResponse{Token: plainToken}
	app.Logger.Infof("Registering user", "email", user.Email, "name", user.Name)
	utils.WriteJson(w, http.StatusCreated, tokenRes)
}
