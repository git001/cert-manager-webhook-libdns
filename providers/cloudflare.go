package providers

import (
	"fmt"

	"github.com/libdns/cloudflare"
)

func init() {
	Register("cloudflare", NewCloudflareProvider)
}

// NewCloudflareProvider creates a Cloudflare DNS provider
//
// Required credentials:
//   - api_token: Cloudflare API token with Zone:DNS:Edit permissions
//
// Note: Use scoped API tokens, NOT global API keys
func NewCloudflareProvider(config ProviderConfig) (DNSProvider, error) {
	apiToken := config.Credentials["api_token"]
	if apiToken == "" {
		return nil, fmt.Errorf("cloudflare: api_token is required")
	}

	return &cloudflare.Provider{
		APIToken: apiToken,
	}, nil
}
