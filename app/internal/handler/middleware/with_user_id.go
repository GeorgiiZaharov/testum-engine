package middleware

import (
	"context"
	"net/http"
)

func WithUserID(r *http.Request, userID int) *http.Request {
	ctx := context.WithValue(r.Context(), userIDKey, userID)
	return r.WithContext(ctx)
}
