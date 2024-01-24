package gql_middleware

import (
	"context"
	"fmt"
	"go-graphql-api/database"
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
			updated_request, err := ProcessAuthFromRequestHeader(r)
			if err != nil {
				logger.Err("Error processing auth from request header: %#v", err)
			}
			next.ServeHTTP(w, updated_request)
		})
	}
}

func ProcessAuthFromRequestHeader(r *http.Request) (*http.Request, error) {
	auth_bearer := r.Header.Get("Authorization")
	if len(auth_bearer) <= 7 || auth_bearer[:7] != "Bearer " {
		// No auth header, not an error, just continue normal request
		logger.Info("Serving request without auth token.")
		return r, nil
	}

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
		return r, err
	}
	// Parse the token into a user payload
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return r, fmt.Errorf("Failed to get claims from jwt auth token")
	}

	user, err := UserFromToken(&claims)
	if err != nil {
		return r, err
	}
	// Successfully parsed the user payload, store it in the request's context.
	logger.Info("User auth token translated to a valid user payload.")
	ctx := context.WithValue(r.Context(), util.ContextKey_User, user)
	return r.WithContext(ctx), nil

}

func UserFromToken(claims *jwt.MapClaims) (*dbmodel.User, error) {
	id_opaq := (*claims)["id"]
	id, ok := id_opaq.(float64)
	if !ok {
		return nil, fmt.Errorf("invalid id type in payload")
	}

	db, err := database.GetDbInstance()
	if err != nil {
		return nil, err
	}

	var user dbmodel.User
	db.Find(&user, []int{int(id)})

	if float64(user.ID) != id {
		return nil, fmt.Errorf("no user found with id %f", id)
	}
	return &user, nil
}
