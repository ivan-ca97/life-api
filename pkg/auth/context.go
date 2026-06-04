package auth

import (
	"context"

	"github.com/google/uuid"
)

type contextKey string

const (
	actorKey   contextKey = "actor_id"
	sessionKey contextKey = "session_id"
)

func WithActor(ctx context.Context, actorId uuid.UUID) context.Context {
	return context.WithValue(ctx, actorKey, actorId)
}

func ActorFromContext(ctx context.Context) (uuid.UUID, error) {
	id, ok := ctx.Value(actorKey).(uuid.UUID)
	if !ok || id == uuid.Nil {
		return uuid.UUID{}, ErrNoActor
	}
	return id, nil
}

func WithSession(ctx context.Context, sessionId uuid.UUID) context.Context {
	return context.WithValue(ctx, sessionKey, sessionId)
}

func SessionFromContext(ctx context.Context) (uuid.UUID, error) {
	id, ok := ctx.Value(sessionKey).(uuid.UUID)
	if !ok || id == uuid.Nil {
		return uuid.UUID{}, ErrNoActor
	}
	return id, nil
}
