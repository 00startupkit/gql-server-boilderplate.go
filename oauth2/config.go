package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"go-graphql-api/database"
	"go-graphql-api/dbmodel"
	"go-graphql-api/util"
	"go-graphql-api/util/logger"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthConfig struct {
	ProviderName string
	ProviderId   string
	// The oauth version being defined.
	Version int
	// This function defines the conversion from an access token given from the
	// oauth server, to a user payload.
	UserFromToken func(token string) (*dbmodel.User, error)
	Oauth2        *oauth2.Config
	// Executed when the auth process is complete.
	OnAuthComplete func(http.ResponseWriter, *http.Request)
}

var (
	_oauth_registers = []*AuthConfig{

		// Google Oauth Config
		{
			ProviderName:   "Google",
			ProviderId:     "google",
			Version:        2,
			UserFromToken:  google_access_token_to_user_payload,
			OnAuthComplete: redirect_home,

			Oauth2: &oauth2.Config{
				ClientID:     util.EnvOrDefault("GOOGLE_CLIENT_ID", ""),
				ClientSecret: util.EnvOrDefault("GOOGLE_CLIENT_SECRET", ""),
				Scopes: []string{
					"https://www.googleapis.com/auth/userinfo.profile",
					"https://www.googleapis.com/auth/userinfo.email"},
				Endpoint: google.Endpoint,
			},
		},
	}
	_oauth_state_key = "placeholderstatekey"
)

func redirect_home(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, util.ServerUri(), http.StatusOK)
}

func RegisterOauthRoutes(router *chi.Mux) {

	handlefn_wrap := func(pattern string, h func(string, http.HandlerFunc), handlerfn func(http.ResponseWriter, *http.Request)) {
		logger.Info("Registering oauth route handler: %s", pattern)
		h(pattern, handlerfn)
	}

	for _, cfg := range _oauth_registers {
		basepath := path.Join("/oauth", fmt.Sprintf("%d", cfg.Version), cfg.ProviderId)
		logger.Info("Registering Oauth: %s/v%d", cfg.ProviderName, cfg.Version)
		cfg.Oauth2.RedirectURL = util.ServerUri() + path.Join(basepath, "callback")

		handlefn_wrap(path.Join(basepath, "login"), router.Get, oauth2_login_initiator(cfg.Oauth2))
		handlefn_wrap(path.Join(basepath, "callback"), router.Get, placeholder_oauth_callback_handler)
	}
}

func placeholder_oauth_callback_handler(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state != _oauth_state_key {
		send_json(w, r,
			http.StatusBadRequest,
			map[string]interface{}{
				"error": "Corrupted state",
			})
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		send_json(w, r,
			http.StatusBadRequest,
			map[string]interface{}{
				"error": "Invalid cllback",
			})
		return
	}
	version := parts[2]
	provider := parts[3]
	logger.Info("Processing callback from provider: [%s]", provider)

	var config *AuthConfig = find_provider_oauth2_config(provider)
	if config == nil {
		send_json(w, r,
			http.StatusBadRequest,
			map[string]interface{}{
				"error": fmt.Sprintf(`auth not configured for "%s"`, provider),
			})
		return
	}

	code := r.URL.Query().Get("code")
	token, err := config.Oauth2.Exchange(context.Background(), code)
	if err != nil {
		logger.Err("Failed to exchange code for token in oauth callback: provider=%s, code=%s, error=%v", provider, code, err)
		send_json(w, r,
			http.StatusBadRequest,
			map[string]interface{}{
				"error": "Token exchange failed",
			})
		return
	}

	user, err := config.UserFromToken(token.AccessToken)
	if err != nil {
		logger.Err("Error exchanging access token to user payload: %#v", err)
		send_json(w, r,
			http.StatusBadRequest,
			map[string]interface{}{
				"error": "User token translation failed",
			})
		return
	}
	logger.Info("User from token: %#v", user)

	if !valid_user_payload(user) {
		logger.Err("Invalid user payload from token conversion.")
		send_json(w, r,
			http.StatusBadRequest,
			map[string]interface{}{
				"error": "User token translation failed",
			})
		return
	}

	// Check if the user exists
	db, err := database.GetDbInstance()
	if err != nil {
		logger.Err("Failed to get database instance: %#v", err)
		send_json(w, r,
			http.StatusBadRequest,
			map[string]interface{}{
				"error": "Internal error",
			})
		return
	}
	var existing_user dbmodel.User
	db.Model(dbmodel.User{}).Where("email = ?", user.Email).First(&existing_user)
	if len(existing_user.Email) == 0 {
		// Non-existing user; create a new user
		logger.Info("OAuth user not yet registered; registering them.")
		result := db.Create(user)
		if result.Error != nil {
			logger.Err("Error creating new user from oauth instance: %#v", result.Error)
			send_json(w, r,
				http.StatusBadRequest,
				map[string]interface{}{
					"error": "Internal error",
				})
			return
		}
		existing_user = *user
	}
	logger.Info("Registered user payload: %#v", existing_user)

	// Remove all prior auth tokens for this user for this given provider and version
	// db.
	db.Delete(
		&dbmodel.OAuthToken{},
		"version = ? AND provider = ? AND user_id = ?", version, provider, existing_user.ID)

	new_auth_token_record := dbmodel.OAuthToken{
		Version:  version,
		Provider: provider,
		// TODO: Consider encrypting the token information
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
		LastRefresh:  time.Now(),
		UserId:       existing_user.ID,
	}
	db.Create(&new_auth_token_record)
	// Save the user information in the reqeust context.
	r.WithContext(context.WithValue(r.Context(), util.ContextKey_User, &existing_user))
	config.OnAuthComplete(w, r)
}

func find_provider_oauth2_config(providerid string) *AuthConfig {
	for _, cfg := range _oauth_registers {
		if cfg.ProviderId == providerid {
			return cfg
		}
	}
	return nil
}

func oauth2_login_initiator(c *oauth2.Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		url := c.AuthCodeURL(_oauth_state_key) // TODO: generate random state.
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}

func send_json(w http.ResponseWriter, r *http.Request, statuscode int, json_data map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statuscode)
	json.NewEncoder(w).Encode(json_data)
}

func google_access_token_to_user_payload(accesstoken string) (*dbmodel.User, error) {
	google_exchange_url := fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", accesstoken)
	res, err := http.Get(google_exchange_url)
	if err != nil {
		return nil, err
	}

	var userdat map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&userdat)
	if err != nil {
		return nil, err
	}

	email, ok := userdat["email"].(string)
	if !ok {
		return nil, fmt.Errorf("failed to parse email from google user payload")
	}
	user := dbmodel.User{}
	user.Email = email
	user.Type = dbmodel.UserType_Normal

	return &user, nil
}

func valid_user_payload(user *dbmodel.User) bool {
	return user != nil && len(user.Email) > 0 && user.Type == dbmodel.UserType_Normal
}
