package sessions

import "context"

type ContextKey string

const SessionKey ContextKey = ContextKey("session")

func getSessionFromContext(ctx context.Context) *Session {
	session, ok := ctx.Value(SessionKey).(*Session)
	if !ok {
		return nil
	}
	return session
}
