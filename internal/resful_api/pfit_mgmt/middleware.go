package pfit_mgmt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"surf_be/internal/database/redis"
)

const (
	HeaderField      = "email"
	ContextEmailKey  = "email"
	ContextAPIKey    = "api_key"
	ContextSecretKey = "secret_key"
)

func Authorize(redis redis.Redis) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case PathToLogout, PathToPing, PathToLogin:
				next.ServeHTTP(w, r)
				return
			}

			ctx := r.Context()

			emailWithEnv := r.Header.Get(HeaderField)
			if len(emailWithEnv) == 0 {
				WriteJSON(w)(HandleError(errors.New("missing email header field"), 9998))
				return
			}

			redisValue, err := redis.Get(ctx, emailWithEnv)
			if err != nil {
				WriteJSON(w)(HandleError(err, 9999))
				return
			}

			if redisValue != nil {
				var redisValueMap map[string]interface{}
				if err = json.Unmarshal([]byte(fmt.Sprintf("%v", redisValue)), &redisValueMap); err != nil {
					WriteJSON(w)(HandleError(err, 9999))
					return
				}

				ctx = context.WithValue(ctx, ContextEmailKey, redisValueMap[ContextEmailKey])
				ctx = context.WithValue(ctx, ContextAPIKey, redisValueMap[ContextAPIKey])
				ctx = context.WithValue(ctx, ContextSecretKey, redisValueMap[ContextSecretKey])
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
