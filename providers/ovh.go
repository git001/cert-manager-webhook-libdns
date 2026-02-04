package providers

import (
	"fmt"

	"github.com/libdns/ovh"
)

func init() {
	Register("ovh", NewOVHProvider)
}

// NewOVHProvider creates an OVH DNS provider
//
// Required credentials:
//   - endpoint: OVH API endpoint (e.g., ovh-eu, ovh-ca, ovh-us)
//   - application_key: OVH application key
//   - application_secret: OVH application secret
//   - consumer_key: OVH consumer key
func NewOVHProvider(config ProviderConfig) (DNSProvider, error) {
	endpoint := config.Credentials["endpoint"]
	applicationKey := config.Credentials["application_key"]
	applicationSecret := config.Credentials["application_secret"]
	consumerKey := config.Credentials["consumer_key"]

	if endpoint == "" {
		return nil, fmt.Errorf("ovh: endpoint is required")
	}
	if applicationKey == "" {
		return nil, fmt.Errorf("ovh: application_key is required")
	}
	if applicationSecret == "" {
		return nil, fmt.Errorf("ovh: application_secret is required")
	}
	if consumerKey == "" {
		return nil, fmt.Errorf("ovh: consumer_key is required")
	}

	return &ovh.Provider{
		Endpoint:          endpoint,
		ApplicationKey:    applicationKey,
		ApplicationSecret: applicationSecret,
		ConsumerKey:       consumerKey,
	}, nil
}
