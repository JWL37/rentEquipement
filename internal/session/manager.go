package session

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type SessionsManager struct {
	data map[string]*Session
	mu   *sync.RWMutex
}

func NewSessionsManager() *SessionsManager {
	return &SessionsManager{
		data: make(map[string]*Session, 10),
		mu:   &sync.RWMutex{},
	}
}

func (sm *SessionsManager) Check(r *http.Request) (*Session, error) {
	sessionCookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		return nil, ErrNoAuth
	}
	log.Println(sessionCookie.Value)
	sm.mu.RLock()
	sess, ok := sm.data[sessionCookie.Value]
	sm.mu.RUnlock()

	if !ok {
		return nil, ErrNoAuth
	}

	return sess, nil
}

func (sm *SessionsManager) Create(w http.ResponseWriter, userID, username, role string) (*Session, error) {
	sess := NewSession(userID, role)

	sm.mu.Lock()
	sm.data[sess.ID] = sess
	sm.mu.Unlock()

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sess.ID,
		Expires: time.Now().Add(10 * time.Minute),
		Path:    "/",
	}
	http.SetCookie(w, cookie)
	usernameCookie := &http.Cookie{
		Name:    "username",
		Value:   username,
		Expires: time.Now().Add(5 * time.Minute),
		Path:    "/",
	}
	http.SetCookie(w, usernameCookie)
	return sess, nil
}

func (sm *SessionsManager) DestroyCurrent(w http.ResponseWriter, r *http.Request) error {
	sess, err := SessionFromContext(r.Context())
	if err != nil {
		return err
	}

	sm.mu.Lock()
	delete(sm.data, sess.ID)
	sm.mu.Unlock()

	cookie := http.Cookie{
		Name:    "session_id",
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}
	http.SetCookie(w, &cookie)
	return nil
}
