package main

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/cert-manager-webhook-libdns/providers"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/libdns/libdns"
	corev1 "k8s.io/api/core/v1"
	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

type mockProvider struct {
	records      []libdns.Record
	getErr       error
	appendCalls  int
	setCalls     int
	deleteCalls  int
	lastZoneSeen string
}

func (m *mockProvider) AppendRecords(_ context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	m.appendCalls++
	m.lastZoneSeen = zone
	m.records = append(m.records, recs...)
	return recs, nil
}

func (m *mockProvider) DeleteRecords(_ context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	m.deleteCalls++
	m.lastZoneSeen = zone

	for _, del := range recs {
		delRR := del.RR()
		filtered := m.records[:0]
		for _, cur := range m.records {
			curRR := cur.RR()
			if curRR.Type == delRR.Type && curRR.Name == delRR.Name && curRR.Data == delRR.Data {
				continue
			}
			filtered = append(filtered, cur)
		}
		m.records = filtered
	}
	return recs, nil
}

func (m *mockProvider) GetRecords(_ context.Context, zone string) ([]libdns.Record, error) {
	m.lastZoneSeen = zone
	if m.getErr != nil {
		return nil, m.getErr
	}
	out := make([]libdns.Record, len(m.records))
	copy(out, m.records)
	return out, nil
}

func (m *mockProvider) SetRecords(_ context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	m.setCalls++
	m.lastZoneSeen = zone

	if len(recs) == 0 {
		return recs, nil
	}

	// Mimic libdns "set" behavior for the same RR type+name by replacing those tuples.
	replacedKeys := make(map[string]struct{}, len(recs))
	for _, rec := range recs {
		rr := rec.RR()
		replacedKeys[rr.Type+"|"+rr.Name] = struct{}{}
	}

	filtered := m.records[:0]
	for _, cur := range m.records {
		rr := cur.RR()
		if _, ok := replacedKeys[rr.Type+"|"+rr.Name]; ok {
			continue
		}
		filtered = append(filtered, cur)
	}
	m.records = append(filtered, recs...)
	return recs, nil
}

func testProviderName(t *testing.T, suffix string) string {
	t.Helper()
	return fmt.Sprintf("test-%s-%s", strings.ReplaceAll(strings.ToLower(t.Name()), "/", "-"), suffix)
}

func registerMockProvider(t *testing.T, name string, mp *mockProvider) {
	t.Helper()
	providers.Register(name, func(config providers.ProviderConfig) (providers.DNSProvider, error) {
		if len(config.Credentials) == 0 {
			return nil, fmt.Errorf("expected credentials")
		}
		return mp, nil
	})
}

func challengeConfigJSON(t *testing.T, providerName, secretName, secretNamespace string, ttl int) *extapi.JSON {
	t.Helper()
	cfg := LibdnsConfig{
		Provider: providerName,
		SecretRef: SecretReference{
			Name:      secretName,
			Namespace: secretNamespace,
		},
		TTL: ttl,
	}
	raw, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("failed to marshal config: %v", err)
	}
	return &extapi.JSON{Raw: raw}
}

func newTestSolver(namespace, secretName string) *libdnsSolver {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		},
		Data: map[string][]byte{
			"api_token": []byte("dummy"),
		},
	}
	return &libdnsSolver{
		client: fake.NewSimpleClientset(secret),
	}
}

func txtValuesForName(records []libdns.Record, name string) []string {
	var out []string
	for _, rec := range records {
		rr := rec.RR()
		if rr.Type == "TXT" && rr.Name == name {
			out = append(out, rr.Data)
		}
	}
	return out
}

func TestPresentMergesTXTValues(t *testing.T) {
	mp := &mockProvider{
		records: []libdns.Record{
			libdns.TXT{Name: "_acme-challenge", Text: "existing", TTL: 120 * time.Second},
		},
	}
	providerName := testProviderName(t, "merge")
	registerMockProvider(t, providerName, mp)

	solver := newTestSolver("cert-manager", "dns-creds")
	ch := &v1alpha1.ChallengeRequest{
		ResolvedFQDN:      "_acme-challenge.example.com.",
		ResolvedZone:      "example.com.",
		Key:               "new-value",
		ResourceNamespace: "cert-manager",
		Config:            challengeConfigJSON(t, providerName, "dns-creds", "", 120),
	}

	if err := solver.Present(ch); err != nil {
		t.Fatalf("Present failed: %v", err)
	}

	if mp.appendCalls != 0 {
		t.Fatalf("AppendRecords should not be called, got %d calls", mp.appendCalls)
	}
	if mp.setCalls != 1 {
		t.Fatalf("SetRecords should be called once, got %d calls", mp.setCalls)
	}

	values := txtValuesForName(mp.records, "_acme-challenge")
	if !slices.Contains(values, "existing") || !slices.Contains(values, "new-value") {
		t.Fatalf("expected merged TXT values [existing,new-value], got %v", values)
	}
}

func TestPresentFallsBackToAppendWhenGetFails(t *testing.T) {
	mp := &mockProvider{getErr: fmt.Errorf("transient get error")}
	providerName := testProviderName(t, "append-fallback")
	registerMockProvider(t, providerName, mp)

	solver := newTestSolver("cert-manager", "dns-creds")
	ch := &v1alpha1.ChallengeRequest{
		ResolvedFQDN:      "_acme-challenge.example.com.",
		ResolvedZone:      "example.com.",
		Key:               "new-value",
		ResourceNamespace: "cert-manager",
		Config:            challengeConfigJSON(t, providerName, "dns-creds", "", 300),
	}

	if err := solver.Present(ch); err != nil {
		t.Fatalf("Present failed: %v", err)
	}

	if mp.appendCalls != 1 {
		t.Fatalf("AppendRecords should be called once, got %d calls", mp.appendCalls)
	}
	if mp.setCalls != 0 {
		t.Fatalf("SetRecords should not be called, got %d calls", mp.setCalls)
	}
}

func TestCleanUpRemovesOnlyMatchingTXTValue(t *testing.T) {
	mp := &mockProvider{
		records: []libdns.Record{
			libdns.TXT{Name: "_acme-challenge", Text: "keep", TTL: 120 * time.Second},
			libdns.TXT{Name: "_acme-challenge", Text: "remove", TTL: 120 * time.Second},
		},
	}
	providerName := testProviderName(t, "cleanup-selective")
	registerMockProvider(t, providerName, mp)

	solver := newTestSolver("cert-manager", "dns-creds")
	ch := &v1alpha1.ChallengeRequest{
		ResolvedFQDN:      "_acme-challenge.example.com.",
		ResolvedZone:      "example.com.",
		Key:               "remove",
		ResourceNamespace: "cert-manager",
		Config:            challengeConfigJSON(t, providerName, "dns-creds", "", 120),
	}

	if err := solver.CleanUp(ch); err != nil {
		t.Fatalf("CleanUp failed: %v", err)
	}

	values := txtValuesForName(mp.records, "_acme-challenge")
	if len(values) != 1 || values[0] != "keep" {
		t.Fatalf("expected remaining TXT value [keep], got %v", values)
	}
	if mp.setCalls != 1 {
		t.Fatalf("SetRecords should be called once, got %d calls", mp.setCalls)
	}
}

func TestCleanUpFallsBackToDeleteWhenGetFails(t *testing.T) {
	mp := &mockProvider{
		getErr: fmt.Errorf("transient get error"),
		records: []libdns.Record{
			libdns.TXT{Name: "_acme-challenge", Text: "remove"},
		},
	}
	providerName := testProviderName(t, "cleanup-delete-fallback")
	registerMockProvider(t, providerName, mp)

	solver := newTestSolver("cert-manager", "dns-creds")
	ch := &v1alpha1.ChallengeRequest{
		ResolvedFQDN:      "_acme-challenge.example.com.",
		ResolvedZone:      "example.com.",
		Key:               "remove",
		ResourceNamespace: "cert-manager",
		Config:            challengeConfigJSON(t, providerName, "dns-creds", "", 120),
	}

	if err := solver.CleanUp(ch); err != nil {
		t.Fatalf("CleanUp failed: %v", err)
	}

	if mp.deleteCalls != 1 {
		t.Fatalf("DeleteRecords should be called once, got %d calls", mp.deleteCalls)
	}
}

func TestGetProviderAppliesDesecMinTTL(t *testing.T) {
	solver := newTestSolver("cert-manager", "dns-creds")
	ch := &v1alpha1.ChallengeRequest{
		ResolvedFQDN:      "_acme-challenge.example.com.",
		ResolvedZone:      "example.com.",
		Key:               "value",
		ResourceNamespace: "cert-manager",
		Config:            challengeConfigJSON(t, "desec", "dns-creds", "", 0),
	}

	_, _, ttl, err := solver.getProvider(ch)
	if err != nil {
		t.Fatalf("getProvider failed: %v", err)
	}
	if ttl != desecMinTTL*time.Second {
		t.Fatalf("expected TTL %s for deSEC, got %s", desecMinTTL*time.Second, ttl)
	}
}
