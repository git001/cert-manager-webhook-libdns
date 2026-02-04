package providers

import (
	"fmt"

	"github.com/libdns/linode"
)

func init() {
	Register("linode", NewLinodeProvider)
}

// NewLinodeProvider creates a Linode DNS provider
//
// Required credentials:
//   - api_token: Linode API token with DNS access
func NewLinodeProvider(config ProviderConfig) (DNSProvider, error) {
	apiToken := config.Credentials["api_token"]
	if apiToken == "" {
		return nil, fmt.Errorf("linode: api_token is required")
	}

	return &linode.Provider{
		APIToken: apiToken,
	}, nil
}
