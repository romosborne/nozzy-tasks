package models

// Env contains all the shared variables
type Env struct {
	Db                Datastore
	SessionKey        []byte
	ContextKey        string
	OauthClientID     string `json:"cid"`
	OauthClientSecret string `json:"csecret"`
	OauthRedirectURL  string `json:"oauth_redirect_url"`
}
