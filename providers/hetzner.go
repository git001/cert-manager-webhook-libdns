package providers

import (
	"fmt"

	hetzner "github.com/libdns/hetzner/v2" // Note: v2 module path
)

func init() {
	Register("hetzner", NewHetznerProvider)
}

// NewHetznerProvider creates a Hetzner DNS provider
//
// Required credentials:
//   - api_token: Hetzner DNS API token
func NewHetznerProvider(config ProviderConfig) (DNSProvider, error) {
	apiToken := config.Credentials["api_token"]
	if apiToken == "" {
		return nil, fmt.Errorf("hetzner: api_token is required")
	}

	return &hetzner.Provider{
		APIToken: apiToken,
	}, nil
}
