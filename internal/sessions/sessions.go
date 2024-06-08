package sessions

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"sync"
	"time"
)

type Session struct {
	ID      string
	Created time.Time
	Values  map[string]any
}

type SessionManager struct {
	sessions map[string]*Session
	mu       sync.Mutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*Session),
	}
}

func (sm *SessionManager) CreateSession() *Session {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	id := GenerateSessionID()
	session := &Session{ID: id, Created: time.Now(), Values: make(map[string]any)}
	sm.sessions[id] = session

	return session
}

func (sm *SessionManager) GetSession(id string) (*Session, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, exists := sm.sessions[id]
	return session, exists
}

func (sm *SessionManager) DestroySession(id string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	delete(sm.sessions, id)
}

func (sm *SessionManager) Put(ctx context.Context, key string, value string) {
	session := getSessionFromContext(ctx)
	session.Values[key] = value
}

func (sm *SessionManager) PopString(ctx context.Context, key string) string {
	session := getSessionFromContext(ctx)
	defer delete(session.Values, key)

	if session.Values[key] == nil {
		return ""
	}

	return session.Values[key].(string)
}

func GenerateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func setSessionCookie(w http.ResponseWriter, id string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "sessid",
		Value:    id,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	})
}

func (sm *SessionManager) ServeSessions(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := r.Cookie("sessid")

		var session *Session

		// If session cookie is not present, or blank, then create a new one
		if err != nil || id.Value == "" {
			session = sm.CreateSession()
			setSessionCookie(w, session.ID)
		} else {
			var exists bool
			session, exists = sm.GetSession(id.Value)
			if !exists {
				session = sm.CreateSession()
				setSessionCookie(w, session.ID)
			}
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, SessionKey, session)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
