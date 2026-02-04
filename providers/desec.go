package providers

import (
	"fmt"

	"github.com/libdns/desec"
)

func init() {
	Register("desec", NewDesecProvider)
}

// NewDesecProvider creates a deSEC DNS provider
//
// Required credentials:
//   - api_token: deSEC API token
func NewDesecProvider(config ProviderConfig) (DNSProvider, error) {
	apiToken := config.Credentials["api_token"]
	if apiToken == "" {
		return nil, fmt.Errorf("desec: api_token is required")
	}

	return &desec.Provider{
		Token: apiToken,
	}, nil
}
