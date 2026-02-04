package providers

import (
	"fmt"

	"github.com/libdns/alidns"
)

func init() {
	Register("alidns", NewAlidnsProvider)
}

// NewAlidnsProvider creates an Alibaba Cloud DNS provider
//
// Required credentials:
//   - access_key_id: Alibaba Cloud access key ID
//   - access_key_secret: Alibaba Cloud access key secret
//
// Optional credentials:
//   - region_id: Alibaba Cloud region (default: cn-hangzhou)
//   - security_token: STS security token (for temporary credentials)
func NewAlidnsProvider(config ProviderConfig) (DNSProvider, error) {
	accessKeyID := config.Credentials["access_key_id"]
	accessKeySecret := config.Credentials["access_key_secret"]

	if accessKeyID == "" {
		return nil, fmt.Errorf("alidns: access_key_id is required")
	}
	if accessKeySecret == "" {
		return nil, fmt.Errorf("alidns: access_key_secret is required")
	}

	provider := &alidns.Provider{
		CredentialInfo: alidns.CredentialInfo{
			AccessKeyID:     accessKeyID,
			AccessKeySecret: accessKeySecret,
		},
	}

	if regionID := config.Credentials["region_id"]; regionID != "" {
		provider.CredentialInfo.RegionID = regionID
	}

	if securityToken := config.Credentials["security_token"]; securityToken != "" {
		provider.CredentialInfo.SecurityToken = securityToken
	}

	return provider, nil
}
