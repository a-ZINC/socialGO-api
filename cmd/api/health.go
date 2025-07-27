package main

import (
	"net/http"
	"social-api/cmd/utils"
	"social-api/internal/env"
)

// healthHandler returns the health status of the API
// @Summary Health check
// @Description Check the health status of the API
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} error
// @Router /v1/health [get]
// @Security ApiKeyAuth	
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
