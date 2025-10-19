package marketplace

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ArtisanCloud/PowerXPlugin/internal/config"
	marketobs "github.com/ArtisanCloud/PowerXPlugin/internal/observability/marketplace"
	"github.com/sirupsen/logrus"
)

var (
	// ErrProviderNotConfigured indicates the client was created without a backing provider adapter.
	ErrProviderNotConfigured = errors.New("tax provider not configured")
	// ErrNotImplemented indicates the concrete adapter has not yet implemented an operation.
	ErrNotImplemented = errors.New("tax provider operation not implemented")
	// ErrReplayNotSupported indicates the provider does not support replay semantics.
	ErrReplayNotSupported = errors.New("tax provider replay not supported")
)

// TaxChargeRequest captures the minimal context required for tax calculations.
type TaxChargeRequest struct {
	TenantID          string
	ListingID         string
	Currency          string
	AmountCents       int64
	Jurisdiction      string
	ExternalReference string
	Items             []TaxLineItem
	Metadata          map[string]string
}

// TaxLineItem describes per-item tax metadata (SKU/plan level granularity).
type TaxLineItem struct {
	SKU         string
	Quantity    int
	AmountCents int64
	TaxCode     string
}

// TaxChargeResult captures the provider response metadata.
type TaxChargeResult struct {
	ExternalTransactionID string
	TaxAmountCents        int64
	Currency              string
	Jurisdiction          string
	RawPayload            []byte
}

// providerAdapter abstracts vendor specific integrations.
type providerAdapter interface {
	Name() string
	Dispatch(ctx context.Context, req *TaxChargeRequest) (*TaxChargeResult, error)
	Replay(ctx context.Context, externalTransactionID string) error
}

// TaxProviderClient orchestrates tax provider dispatch and replay with retry semantics.
type TaxProviderClient struct {
	provider string
	retries  []time.Duration
	adapter  providerAdapter
	logger   *logrus.Entry
}

// NewTaxProviderClient constructs a retry-aware client for the configured tax provider.
func NewTaxProviderClient(cfg *config.Config, httpClient *http.Client, logger *logrus.Entry) (*TaxProviderClient, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config required to initialize tax provider client")
	}

	provider := cfg.IntegrationTaxProvider()
	if provider == "" {
		return nil, fmt.Errorf("integration.billing.tax_provider is not configured")
	}

	schedule := cfg.BillingRetrySchedule()
	if httpClient == nil {
		httpClient = &http.Client{Timeout: cfg.BillingHTTPTimeout()}
	} else if httpClient.Timeout == 0 {
		httpClient.Timeout = cfg.BillingHTTPTimeout()
	}

	client := &TaxProviderClient{
		provider: provider,
		retries:  schedule,
		logger:   logger,
	}

	switch provider {
	case "stripe_tax":
		client.adapter = newStripeAdapter(httpClient, cfg.StripeTaxConfig(), logger)
	case "avalara":
		client.adapter = newAvalaraAdapter(httpClient, cfg.AvalaraConfig(), logger)
	default:
		return nil, fmt.Errorf("unsupported tax provider %q", provider)
	}

	return client, nil
}

// Provider returns the currently configured provider identifier.
func (c *TaxProviderClient) Provider() string {
	if c == nil {
		return ""
	}
	return c.provider
}

// CreateTransaction performs a best-effort dispatch with retries according to configuration.
func (c *TaxProviderClient) CreateTransaction(ctx context.Context, req *TaxChargeRequest) (*TaxChargeResult, error) {
	if c == nil || c.adapter == nil {
		return nil, ErrProviderNotConfigured
	}
	if ctx == nil {
		ctx = context.Background()
	}

	attempts := len(c.retries) + 1
	var lastErr error
	for i := 0; i < attempts; i++ {
		if i > 0 {
			if err := waitWithContext(ctx, c.retries[i-1]); err != nil {
				return nil, err
			}
		}

		result, err := c.adapter.Dispatch(ctx, req)
		if err == nil {
			return result, nil
		}
		lastErr = err
		c.recordFailure("dispatch", err)

		if !c.shouldRetry(err) || i == attempts-1 || ctx.Err() != nil {
			break
		}

		if c.logger != nil {
			c.logger.WithError(err).WithFields(logrus.Fields{
				"provider": c.provider,
				"attempt":  i + 1,
			}).Warn("tax provider dispatch failed, retrying")
		}
	}

	if lastErr == nil {
		lastErr = ErrNotImplemented
	}
	return nil, lastErr
}

// ReplayTransaction attempts to replay a previously failed tax transaction.
func (c *TaxProviderClient) ReplayTransaction(ctx context.Context, externalTransactionID string) error {
	if c == nil || c.adapter == nil {
		return ErrProviderNotConfigured
	}
	if ctx == nil {
		ctx = context.Background()
	}

	attempts := len(c.retries) + 1
	var lastErr error
	for i := 0; i < attempts; i++ {
		if i > 0 {
			if err := waitWithContext(ctx, c.retries[i-1]); err != nil {
				return err
			}
		}

		err := c.adapter.Replay(ctx, externalTransactionID)
		if err == nil {
			return nil
		}
		lastErr = err
		c.recordFailure("replay", err)

		if !c.shouldRetry(err) || i == attempts-1 || ctx.Err() != nil {
			break
		}

		if c.logger != nil {
			c.logger.WithError(err).WithFields(logrus.Fields{
				"provider": c.provider,
				"attempt":  i + 1,
			}).Warn("tax provider replay failed, retrying")
		}
	}

	if lastErr == nil {
		lastErr = ErrNotImplemented
	}
	return lastErr
}

func (c *TaxProviderClient) recordFailure(stage string, err error) {
	if err == nil {
		return
	}
	code := errorCode(err)
	marketobs.IncrementTaxProviderError(c.provider, fmt.Sprintf("%s_%s", stage, code))
}

func (c *TaxProviderClient) shouldRetry(err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, ErrNotImplemented),
		errors.Is(err, ErrReplayNotSupported),
		errors.Is(err, context.Canceled),
		errors.Is(err, context.DeadlineExceeded):
		return false
	default:
		return true
	}
}

func waitWithContext(ctx context.Context, delay time.Duration) error {
	if delay <= 0 {
		return nil
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func errorCode(err error) string {
	switch {
	case err == nil:
		return "ok"
	case errors.Is(err, ErrNotImplemented):
		return "not_implemented"
	case errors.Is(err, ErrReplayNotSupported):
		return "replay_unsup"
	case errors.Is(err, context.Canceled):
		return "canceled"
	case errors.Is(err, context.DeadlineExceeded):
		return "timeout"
	default:
		return "unknown"
	}
}

type stripeAdapter struct {
	httpClient *http.Client
	cfg        config.IntegrationStripeTaxConfig
	logger     *logrus.Entry
}

func newStripeAdapter(httpClient *http.Client, cfg config.IntegrationStripeTaxConfig, logger *logrus.Entry) providerAdapter {
	return &stripeAdapter{
		httpClient: httpClient,
		cfg:        cfg,
		logger:     logger,
	}
}

func (a *stripeAdapter) Name() string {
	return "stripe_tax"
}

func (a *stripeAdapter) Dispatch(ctx context.Context, req *TaxChargeRequest) (*TaxChargeResult, error) {
	return nil, ErrNotImplemented
}

func (a *stripeAdapter) Replay(ctx context.Context, externalTransactionID string) error {
	return ErrNotImplemented
}

type avalaraAdapter struct {
	httpClient *http.Client
	cfg        config.IntegrationAvalaraConfig
	logger     *logrus.Entry
}

func newAvalaraAdapter(httpClient *http.Client, cfg config.IntegrationAvalaraConfig, logger *logrus.Entry) providerAdapter {
	return &avalaraAdapter{
		httpClient: httpClient,
		cfg:        cfg,
		logger:     logger,
	}
}

func (a *avalaraAdapter) Name() string {
	return "avalara"
}

func (a *avalaraAdapter) Dispatch(ctx context.Context, req *TaxChargeRequest) (*TaxChargeResult, error) {
	return nil, ErrNotImplemented
}

func (a *avalaraAdapter) Replay(ctx context.Context, externalTransactionID string) error {
	return ErrNotImplemented
}
