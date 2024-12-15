package session

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type Session struct {
	ID     string
	UserID string
	Role   string
}

func NewSession(userID, role string) *Session {
	randID := uuid.New()

	return &Session{
		ID:     randID.String(),
		UserID: userID,
		Role:   role,
	}
}

var (
	ErrNoAuth = errors.New("No session found")
)

type sessKey string

var SessionKey sessKey = "sessionKey"

func SessionFromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(SessionKey).(*Session)
	if !ok || sess == nil {
		return nil, ErrNoAuth
	}
	return sess, nil
}

func ContextWithSession(ctx context.Context, sess *Session) context.Context {
	return context.WithValue(ctx, SessionKey, sess)
}
