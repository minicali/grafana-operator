package grafana

import (
	"context"
	"net/http"
	"time"

	gapi "github.com/grafana/grafana-api-golang-client"
)

type GrafanaClient struct {
	Client *gapi.Client
}

func NewClient(apiURL string, timeout time.Duration, auth Authenticator) (*GrafanaClient, error) {
	clientConfig := gapi.Config{
		HTTPHeaders: nil,
		Client: &http.Client{
			Timeout: time.Second * timeout,
		},
		OrgID:      0,
		NumRetries: 0,
	}
	auth.ApplyCredentials(&clientConfig)

	grafanaClient, err := gapi.New(apiURL, clientConfig)
	if err != nil {
		return nil, err
	}

	return &GrafanaClient{
		Client: grafanaClient,
	}, nil
}

func FromContext(ctx context.Context) *GrafanaClient {
	return ctx.Value(clientCtxKey{}).(*GrafanaClient)
}
