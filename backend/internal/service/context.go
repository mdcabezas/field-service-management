package service

import (
	"context"

	"localis-backend/internal/types"
)

func UserIDFromCtx(ctx context.Context) string {
	if v, ok := ctx.Value(types.UserIDKey).(string); ok {
		return v
	}
	return ""
}

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, types.UserIDKey, userID)
}
