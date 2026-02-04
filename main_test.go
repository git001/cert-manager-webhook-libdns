package main

import (
	"os"
	"testing"

	acmetest "github.com/cert-manager/cert-manager/test/acme"
)

var (
	zone = os.Getenv("TEST_ZONE_NAME")
)

func TestRunsSuite(t *testing.T) {
	// Skip if zone is not configured
	if zone == "" {
		t.Skip("TEST_ZONE_NAME environment variable not set, skipping conformance tests")
	}

	// The manifest path should contain a file named config.json that holds
	// the webhook configuration for testing
	fixture := acmetest.NewFixture(&libdnsSolver{},
		acmetest.SetResolvedZone(zone),
		acmetest.SetAllowAmbientCredentials(false),
		acmetest.SetManifestPath("testdata/libdns-solver"),
		acmetest.SetStrict(true),
	)

	fixture.RunConformance(t)
}
