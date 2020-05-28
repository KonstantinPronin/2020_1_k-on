package constants

import "time"

const (
	CookieName     = "session_id"
	UserIdKey      = "user_id"
	CookieDuration = 10 * time.Hour
)

const (
	CSRFHeader = "X-CSRF-TOKEN"
	CSRFKey    = "eE%yh?aAH_hYk*5h$DXvTddAGt2eWCt^+TT_4*$ADxz^X$5ue74jmeJT@z^+c_*v"
)

const (
	ImgDir      = "/frontend/static/img/"
	LogConfFile = "/backend/conf/log.json"
	DbConfFile  = "/backend/conf/db.json"
)
