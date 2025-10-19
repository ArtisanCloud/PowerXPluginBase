package marketplace

import (
	"context"
	"testing"

	"github.com/ArtisanCloud/PowerXPlugin/internal/config"
)

func TestNewTaxProviderClientUnsupported(t *testing.T) {
	cfg := &config.Config{
		Integration: &config.IntegrationConfig{
			Billing: config.IntegrationBillingConfig{
				TaxProvider: "unknown",
			},
		},
	}

	_, err := NewTaxProviderClient(cfg, nil, nil)
	if err == nil {
		t.Fatal("expected error for unsupported provider, got nil")
	}
}

func TestCreateTransactionNotImplemented(t *testing.T) {
	cfg := &config.Config{
		Integration: &config.IntegrationConfig{
			Billing: config.IntegrationBillingConfig{
				TaxProvider: "stripe_tax",
			},
		},
	}

	client, err := NewTaxProviderClient(cfg, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}

	_, err = client.CreateTransaction(context.Background(), &TaxChargeRequest{Currency: "USD"})
	if err == nil {
		t.Fatal("expected ErrNotImplemented, got nil")
	}
	if err != ErrNotImplemented {
		t.Fatalf("expected ErrNotImplemented, got %v", err)
	}
}
