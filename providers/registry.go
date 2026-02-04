package providers

import (
	"context"
	"fmt"
	"sync"

	"github.com/libdns/libdns"
)

// DNSProvider combines the libdns interfaces needed for ACME DNS-01 challenges
type DNSProvider interface {
	libdns.RecordAppender
	libdns.RecordDeleter
	libdns.RecordGetter
	libdns.RecordSetter
}

// ProviderConfig holds the configuration needed to instantiate a provider
type ProviderConfig struct {
	// Provider-specific configuration as key-value pairs
	// Populated from Kubernetes Secret data
	Credentials map[string]string
}

// ProviderFactory creates a DNSProvider from configuration
type ProviderFactory func(config ProviderConfig) (DNSProvider, error)

// Registry holds all registered provider factories
type Registry struct {
	mu        sync.RWMutex
	factories map[string]ProviderFactory
}

// globalRegistry is the default registry instance
var globalRegistry = &Registry{
	factories: make(map[string]ProviderFactory),
}

// Register adds a provider factory to the global registry
func Register(name string, factory ProviderFactory) {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	globalRegistry.factories[name] = factory
}

// Get retrieves a provider factory by name from the global registry
func Get(name string) (ProviderFactory, error) {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()

	factory, ok := globalRegistry.factories[name]
	if !ok {
		return nil, fmt.Errorf("unknown DNS provider: %s (available: %v)", name, ListProviders())
	}
	return factory, nil
}

// ListProviders returns all registered provider names
func ListProviders() []string {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()

	names := make([]string, 0, len(globalRegistry.factories))
	for name := range globalRegistry.factories {
		names = append(names, name)
	}
	return names
}

// CreateProvider is a convenience function that gets and invokes a provider factory
func CreateProvider(name string, config ProviderConfig) (DNSProvider, error) {
	factory, err := Get(name)
	if err != nil {
		return nil, err
	}
	return factory(config)
}

// Ensure DNSProvider interface is compatible with libdns at compile time
var _ DNSProvider = (*wrappedProvider)(nil)

// wrappedProvider is a compile-time check helper
type wrappedProvider struct {
	libdns.RecordAppender
	libdns.RecordDeleter
	libdns.RecordGetter
	libdns.RecordSetter
}

func (w *wrappedProvider) AppendRecords(ctx context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	return w.RecordAppender.AppendRecords(ctx, zone, recs)
}

func (w *wrappedProvider) DeleteRecords(ctx context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	return w.RecordDeleter.DeleteRecords(ctx, zone, recs)
}

func (w *wrappedProvider) GetRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
	return w.RecordGetter.GetRecords(ctx, zone)
}

func (w *wrappedProvider) SetRecords(ctx context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	return w.RecordSetter.SetRecords(ctx, zone, recs)
}
