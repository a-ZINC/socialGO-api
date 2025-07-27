package utils

import (
	"log"
	"net/http"
)

type ErrorHandler struct{}

func (e *ErrorHandler) InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Internal Server Error method: %s path: %s error: %v", r.Method, r.URL.Path, err)
	WriteJsonError(w, http.StatusInternalServerError, "Internal Server Error")
}

func (e *ErrorHandler) BadRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Bad Request method: %s path: %s error: %v", r.Method, r.URL.Path, err)
	WriteJsonError(w, http.StatusBadRequest, err.Error())
}

func (e *ErrorHandler) NotFoundError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Not Found method: %s path: %s error: %v", r.Method, r.URL.Path, err)
	WriteJsonError(w, http.StatusNotFound, "Not Found")
}
