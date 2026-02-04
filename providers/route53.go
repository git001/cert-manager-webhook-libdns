package providers

import (
	"fmt"

	"github.com/libdns/route53"
)

func init() {
	Register("route53", NewRoute53Provider)
}

// NewRoute53Provider creates an AWS Route53 DNS provider
//
// Required credentials:
//   - access_key_id: AWS access key ID
//   - secret_access_key: AWS secret access key
//
// Optional credentials:
//   - region: AWS region (default: us-east-1)
//   - session_token: AWS session token (for temporary credentials)
func NewRoute53Provider(config ProviderConfig) (DNSProvider, error) {
	accessKeyID := config.Credentials["access_key_id"]
	secretAccessKey := config.Credentials["secret_access_key"]

	if accessKeyID == "" {
		return nil, fmt.Errorf("route53: access_key_id is required")
	}
	if secretAccessKey == "" {
		return nil, fmt.Errorf("route53: secret_access_key is required")
	}

	provider := &route53.Provider{
		AccessKeyId:     accessKeyID,
		SecretAccessKey: secretAccessKey,
	}

	if region := config.Credentials["region"]; region != "" {
		provider.Region = region
	}

	if sessionToken := config.Credentials["session_token"]; sessionToken != "" {
		provider.SessionToken = sessionToken
	}

	return provider, nil
}
