package session

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Role      string    `json:"role"`
	ExpiresAt time.Time `json:"expires_at"`
}

func NewSession(userID, role string, ttl int) *Session {
	randID := uuid.New()

	return &Session{
		ID:        randID.String(),
		UserID:    userID,
		Role:      role,
		ExpiresAt: time.Now().Add(time.Duration(ttl) * time.Second),
	}
}

// ToJSON конвертирует сессию в JSON для хранения в Redis
func (s *Session) ToJSON() ([]byte, error) {
	return json.Marshal(s)
}

// FromJSON конвертирует JSON из Redis в сессию
func SessionFromJSON(data []byte) (*Session, error) {
	var s Session
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

var (
	ErrNoAuth = errors.New("no session found")
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
