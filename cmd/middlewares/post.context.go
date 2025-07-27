package middlewares

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PostIdType string

const postCtx PostIdType = "postId"

func (m *Middleware) PostContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postIdParam := chi.URLParam(r, "postId")
		postId, err := strconv.ParseInt(postIdParam, 10, 64)
		if err != nil {
			m.Err.InternalServerError(w, r, err)
		}
		ctxWithValue := context.WithValue(r.Context(), postCtx, postId)
		next.ServeHTTP(w, r.WithContext(ctxWithValue))
	})
}

func (m *Middleware) GetPostIdFromContext(r *http.Request) int64 {
	postId := r.Context().Value(postCtx).(int64)

	return postId
}
