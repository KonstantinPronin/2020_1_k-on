package session

import "time"

const (
	CookieName     = "session_id"
	UserIdKey      = "user_id"
	CookieDuration = 10 * time.Hour
)
