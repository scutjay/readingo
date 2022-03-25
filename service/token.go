package service

import (
	"readingo/conf"
	"readingo/constant"
	"time"
)

var allowTokenCache = make(map[string]allowToken, 0)

type allowToken struct {
	lastActiveTime time.Time
	role           string
}

func CheckIfValidToken(token string) (isValid bool) {
	if conf.Auth.Anonymous {
		isValid = true
	} else if token == "" {
		isValid = false
	} else if target, ok := allowTokenCache[token]; ok {
		now := time.Now()
		lastActiveTime := target.lastActiveTime
		if lastActiveTime.Add(constant.TokenDurationInServer).After(now) {
			if constant.AllowToRecountAliveDuration {
				target.lastActiveTime = now
				allowTokenCache[token] = target
			}
			isValid = true
		} else {
			delete(allowTokenCache, token)
		}
	}
	return isValid
}

func getRoleByToken(token string) string {
	if tar, ok := allowTokenCache[token]; ok {
		return tar.role
	}
	return ""
}

func addTokenToCache(token, role string) {
	allowTokenCache[token] = allowToken{
		lastActiveTime: time.Now(),
		role:           role,
	}
}

func removeTokenFromCache(token string) {
	delete(allowTokenCache, token)
}
