package constant

import "time"

var (
	RoleReadOnly  = "readonly"
	RoleReadWrite = "readwrite"

	DBTreeAutoRefreshNormalDuration = time.Minute
	DBTreeAutoRefreshLongDuration   = 10 * time.Minute
	DBTreeRefreshMinInterval        = 10 * time.Second
	DBTreeRefreshFailTimes          = 10

	TokenDurationInCookie       = 3600 * 24
	TokenDurationInServer       = 24 * time.Hour
	AllowToRecountAliveDuration = false
)
