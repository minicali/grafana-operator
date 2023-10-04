package grafana

import (
	"context"
	"net/url"

	gapi "github.com/grafana/grafana-api-golang-client"
)

type Authenticator interface {
	ApplyCredentials(*gapi.Config)
}

type APITokenAuthenticator struct {
	Token string
}

func (auth *APITokenAuthenticator) ApplyCredentials(config *gapi.Config) {
	config.APIKey = auth.Token
}

type BasicAuthenticator struct {
	Username string
	Password string
}

func (auth *BasicAuthenticator) ApplyCredentials(config *gapi.Config) {
	config.BasicAuth = url.UserPassword(auth.Username, auth.Password)
}

type (
	authenticatorCtxKey struct{}
	clientCtxKey        struct{}
)

func WithAuthenticator(ctx context.Context, auth Authenticator) context.Context {
	return context.WithValue(ctx, authenticatorCtxKey{}, auth)
}

func WithGrafanaClient(ctx context.Context, client *GrafanaClient) context.Context {
	return context.WithValue(ctx, clientCtxKey{}, client)
}
