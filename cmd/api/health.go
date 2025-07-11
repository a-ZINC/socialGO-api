package main

import (
	"net/http"
	"social-api/cmd/utils"
	"social-api/internal/env"
)

func (app *Application) healthHandler(w http.ResponseWriter, r *http.Request) {
	msg := map[string]string{
		"status":      "healthy",
		"version":     Version,
		"environment": env.GetString("ENV", "dev"),
	}
	if err := utils.WriteJson(w, http.StatusOK, msg); err != nil {
		if writeErr := utils.WriteJsonError(w, http.StatusInternalServerError, "Failed to write response"); writeErr != nil {
			http.Error(w, writeErr.Error(), http.StatusInternalServerError)
		}
	}
}
