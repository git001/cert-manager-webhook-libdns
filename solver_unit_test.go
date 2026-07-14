package main

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/cert-manager-webhook-libdns/libdnsregistry"
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

func registerMockProvider(t *testing.T, name string, provider libdnsregistry.Provider) {
	t.Helper()
	err := libdnsregistry.Register(name, &libdnsregistry.RegistryProvider{
		Init: func(conf [][]byte) (libdnsregistry.Provider, error) {
			return provider, nil
		},
	})
	if err != nil {
		t.Fatalf("failed to register mock provider %q: %v", name, err)
	}
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

// syncPoint lets a test observe exactly when a provider call enters its
// critical section, and hold it there until the test says to proceed.
type syncPoint struct {
	entered chan struct{}
	release chan struct{}
}

// blockingProvider wraps mockProvider and pauses inside GetRecords until
// released, so a test can deterministically force two Present()/CleanUp()
// calls to overlap the way cert-manager does for wildcard + apex challenges
// that share one DNS-01 record name.
type blockingProvider struct {
	*mockProvider
	sync *syncPoint
}

func (p *blockingProvider) GetRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
	p.sync.entered <- struct{}{}
	<-p.sync.release
	return p.mockProvider.GetRecords(ctx, zone)
}

// TestPresentSerializesConcurrentChallengesForSameRecord reproduces the
// wildcard+apex certificate scenario: cert-manager creates one Challenge per
// SAN entry, both resolving to the same DNS-01 record name, and calls
// Present() for them concurrently. Present() does a non-atomic
// get-existing-values -> merge -> set-all-values cycle; without serializing
// that cycle, the second call's SetRecords can clobber the first call's
// still-pending value. This test fails if Present() is not holding a lock
// across the whole read-modify-write cycle.
func TestPresentSerializesConcurrentChallengesForSameRecord(t *testing.T) {
	mp := &mockProvider{}
	sp := &syncPoint{entered: make(chan struct{}), release: make(chan struct{})}
	bp := &blockingProvider{mockProvider: mp, sync: sp}

	providerName := testProviderName(t, "race")
	registerMockProvider(t, providerName, bp)

	solver := newTestSolver("cert-manager", "dns-creds")
	makeChallenge := func(key string) *v1alpha1.ChallengeRequest {
		return &v1alpha1.ChallengeRequest{
			ResolvedFQDN:      "_acme-challenge.staging.example.com.",
			ResolvedZone:      "example.com.",
			Key:               key,
			ResourceNamespace: "cert-manager",
			Config:            challengeConfigJSON(t, providerName, "dns-creds", "", 120),
		}
	}

	done := make(chan error, 2)
	go func() { done <- solver.Present(makeChallenge("apex-value")) }()

	// Wait for the first call to be inside its critical section (blocked in GetRecords).
	<-sp.entered

	// Start the second, concurrent Present() call for the SAME record name.
	go func() { done <- solver.Present(makeChallenge("wildcard-value")) }()

	// It must NOT be able to enter GetRecords while the first call still
	// holds the lock - if it does, the two calls are unserialized.
	select {
	case <-sp.entered:
		t.Fatal("second Present() entered its critical section while the first was still in progress - Present() is not serialized across concurrent challenges")
	case <-time.After(50 * time.Millisecond):
	}

	// Let the first call finish, then the second proceeds through the same path.
	sp.release <- struct{}{}
	if err := <-done; err != nil {
		t.Fatalf("first Present() failed: %v", err)
	}

	<-sp.entered
	sp.release <- struct{}{}
	if err := <-done; err != nil {
		t.Fatalf("second Present() failed: %v", err)
	}

	values := txtValuesForName(mp.records, "_acme-challenge.staging")
	for _, key := range []string{"apex-value", "wildcard-value"} {
		if !slices.Contains(values, key) {
			t.Fatalf("expected both concurrent challenge values to survive, got %v (missing %q)", values, key)
		}
	}
}

// TestCleanUpSerializesConcurrentChallengesForSameRecord mirrors the Present
// test for CleanUp(): removing one challenge's value must not race with
// another challenge's concurrent CleanUp() for the same record name.
func TestCleanUpSerializesConcurrentChallengesForSameRecord(t *testing.T) {
	mp := &mockProvider{
		records: []libdns.Record{
			libdns.TXT{Name: "_acme-challenge.staging", Text: "apex-value", TTL: 120 * time.Second},
			libdns.TXT{Name: "_acme-challenge.staging", Text: "wildcard-value", TTL: 120 * time.Second},
		},
	}
	sp := &syncPoint{entered: make(chan struct{}), release: make(chan struct{})}
	bp := &blockingProvider{mockProvider: mp, sync: sp}

	providerName := testProviderName(t, "race-cleanup")
	registerMockProvider(t, providerName, bp)

	solver := newTestSolver("cert-manager", "dns-creds")
	makeChallenge := func(key string) *v1alpha1.ChallengeRequest {
		return &v1alpha1.ChallengeRequest{
			ResolvedFQDN:      "_acme-challenge.staging.example.com.",
			ResolvedZone:      "example.com.",
			Key:               key,
			ResourceNamespace: "cert-manager",
			Config:            challengeConfigJSON(t, providerName, "dns-creds", "", 120),
		}
	}

	done := make(chan error, 2)
	go func() { done <- solver.CleanUp(makeChallenge("apex-value")) }()
	<-sp.entered

	go func() { done <- solver.CleanUp(makeChallenge("wildcard-value")) }()
	select {
	case <-sp.entered:
		t.Fatal("second CleanUp() entered its critical section while the first was still in progress - CleanUp() is not serialized across concurrent challenges")
	case <-time.After(50 * time.Millisecond):
	}

	sp.release <- struct{}{}
	if err := <-done; err != nil {
		t.Fatalf("first CleanUp() failed: %v", err)
	}

	<-sp.entered
	sp.release <- struct{}{}
	if err := <-done; err != nil {
		t.Fatalf("second CleanUp() failed: %v", err)
	}

	if len(mp.records) != 0 {
		t.Fatalf("expected both values to be removed, got %v", mp.records)
	}
}

// TestConcurrentPresentsDoNotDataRace exercises Present() from many
// goroutines under `go test -race` to catch any remaining unsynchronized
// access to the provider/mockProvider state, beyond the deterministic
// interleaving already covered above.
func TestConcurrentPresentsDoNotDataRace(t *testing.T) {
	mp := &mockProvider{}
	providerName := testProviderName(t, "race-detector")
	registerMockProvider(t, providerName, mp)

	solver := newTestSolver("cert-manager", "dns-creds")

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ch := &v1alpha1.ChallengeRequest{
				ResolvedFQDN:      "_acme-challenge.staging.example.com.",
				ResolvedZone:      "example.com.",
				Key:               fmt.Sprintf("value-%d", i),
				ResourceNamespace: "cert-manager",
				Config:            challengeConfigJSON(t, providerName, "dns-creds", "", 120),
			}
			if err := solver.Present(ch); err != nil {
				t.Errorf("Present failed: %v", err)
			}
		}(i)
	}
	wg.Wait()
}
