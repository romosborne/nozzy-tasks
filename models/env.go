package models

type Env struct {
	Db                Datastore
	SessionKey        []byte
	OauthClientID     string `json:"cid"`
	OauthClientSecret string `json:"csecret"`
	OauthRedirectUrl  string `json:"oauth_redirect_url"`
}
