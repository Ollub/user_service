package session

import "context"

const sessionKey = "session"

type Session struct {
	UserID uint32
	ID     string
}

func ToContext(ctx context.Context, sess *Session) context.Context {
	return context.WithValue(ctx, sessionKey, sess)
}

func FromContext(ctx context.Context) *Session {
	sess, ok := ctx.Value(sessionKey).(*Session)
	if !ok {
		return nil
	}
	return sess
}
