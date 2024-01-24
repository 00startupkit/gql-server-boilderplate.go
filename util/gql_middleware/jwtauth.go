package gql_middleware

import (
	"context"
	"fmt"
	"go-graphql-api/dbmodel"
	"go-graphql-api/util"
	"go-graphql-api/util/logger"
	"net/http"

	"github.com/golang-jwt/jwt"
)

// Auth example modified from: https://gqlgen.com/recipes/authentication/

// Check if there is user auth information attached to the request
// and if so, store it in the context that will be sent to the graphql
// request.
//
// The auth token is expected to be in the request header in this format:
//
//	Authorization: Bearer <jwt token>
func JwtAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth_bearer := r.Header.Get("Authorization")
			if len(auth_bearer) > 7 &&
				auth_bearer[:7] == "Bearer " {
				tokenstr := auth_bearer[7:]

				logger.Info("Attempting to validate auth token: %s", tokenstr)
				token, err := jwt.Parse(tokenstr, func(t *jwt.Token) (interface{}, error) {
					// Must validate that the token is using the expected algo.
					if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
					}
					secret := util.EnvOrDefault("JWT_SECRET", "")
					if len(secret) == 0 {
						return nil, fmt.Errorf("No jwt secret found.")
					}
					return []byte(secret), nil
				})
				if err != nil {
					logger.Err("Failed to validate auth token [%s]: %v", tokenstr, err)
				} else {
					// Parse the token into a user payload
					if claims, ok := token.Claims.(jwt.Claims); ok {

						user, err := UserFromToken(&claims)
						if err != nil {
							logger.Err("Failed to parse user payload from jwt claims: %#v", err)
						} else {
							// Successfully parsed the user payload, store it in the request's context.
							ctx := context.WithValue(r.Context(), "user", user)
							r = r.WithContext(ctx)
						}
					} else {
						logger.Err("Failed to get claims from jwt auth token.")
					}
				}
			} else {
				logger.Info("Serving request without auth token.")
			}
			next.ServeHTTP(w, r)
		})
	}
}

func UserFromToken(claims *jwt.Claims) (*dbmodel.User, error) {
	return nil, fmt.Errorf("UserFromToken not implemented")
}
