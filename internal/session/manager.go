package session

import (
	"context"
	"log"
	"net/http"
	"rentEquipement/internal/config"
	redisClient "rentEquipement/internal/redis"
	"time"
)

type SessionsManager struct {
	redis      *redisClient.RedisClient
	cfg        *config.Config
	sessionTTL int
}

func NewSessionsManager(redis *redisClient.RedisClient, cfg *config.Config) *SessionsManager {
	return &SessionsManager{
		redis:      redis,
		cfg:        cfg,
		sessionTTL: cfg.SessionTTL,
	}
}

func (sm *SessionsManager) Check(r *http.Request) (*Session, error) {
	ctx := r.Context()
	sessionCookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		return nil, ErrNoAuth
	}

	sessionID := sessionCookie.Value
	log.Println("Checking session:", sessionID)

	// Получаем сессию из Redis
	data, err := sm.redis.Client.Get(ctx, "session:"+sessionID).Bytes()
	if err != nil {
		log.Println("Failed to get session from Redis:", err)
		return nil, ErrNoAuth
	}

	// Десериализуем сессию
	sess, err := SessionFromJSON(data)
	if err != nil {
		log.Println("Failed to deserialize session:", err)
		return nil, ErrNoAuth
	}

	return sess, nil
}

func (sm *SessionsManager) Create(w http.ResponseWriter, userID, username, role string) (*Session, error) {
	ctx := context.Background()
	sess := NewSession(userID, role, sm.sessionTTL)

	log.Println("TIIIIIIME", time.Until(sess.ExpiresAt))
	data, err := sess.ToJSON()
	if err != nil {
		log.Println("Failed to serialize session:", err)
		return nil, err
	}
	err = sm.redis.Client.Set(
		ctx,
		"session:"+sess.ID,
		data,
		time.Duration(sm.sessionTTL)*time.Second,
	).Err()

	if err != nil {
		log.Println("Failed to save session to Redis:", err)
		return nil, err
	}

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sess.ID,
		Expires: sess.ExpiresAt.Add(time.Until(sess.ExpiresAt)),
		Path:    "/",
	}
	http.SetCookie(w, cookie)

	usernameCookie := &http.Cookie{
		Name:    "username",
		Value:   username,
		Expires: sess.ExpiresAt.Add(time.Until(sess.ExpiresAt)),
		Path:    "/",
	}
	http.SetCookie(w, usernameCookie)

	return sess, nil
}

func (sm *SessionsManager) DestroyCurrent(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	sess, err := SessionFromContext(r.Context())
	if err != nil {
		return err
	}

	err = sm.redis.Client.Del(ctx, "session:"+sess.ID).Err()
	if err != nil {
		log.Println("Failed to delete session from Redis:", err)
	}

	cookie := http.Cookie{
		Name:    "session_id",
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}
	http.SetCookie(w, &cookie)

	usernameCookie := http.Cookie{
		Name:    "username",
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}
	http.SetCookie(w, &usernameCookie)

	return nil
}
