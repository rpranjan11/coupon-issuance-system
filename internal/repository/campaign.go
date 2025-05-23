// internal/repository/campaign.go
package repository

import (
	"context"

	"github.com/rpranjan11/coupon-issuance-system/internal/domain"
)

// CampaignRepository defines the interface for campaign persistence
type CampaignRepository interface {
	// Create saves a new campaign
	Create(ctx context.Context, campaign *domain.Campaign) error

	// Get retrieves a campaign by ID
	Get(ctx context.Context, id string) (*domain.Campaign, error)

	// Update updates an existing campaign
	Update(ctx context.Context, campaign *domain.Campaign) error

	// AtomicIncrementIssued atomically increments the issued_coupons counter
	// Returns true if increment was successful, false if total was reached
	AtomicIncrementIssued(ctx context.Context, campaignID string) (bool, error)

	// FindByName finds a campaign by its name
	FindByName(ctx context.Context, name string) (*domain.Campaign, error)

	// DeleteByID deletes a campaign by ID
	DeleteByID(ctx context.Context, id string) (bool, error)

	// DeleteByName deletes a campaign by name
	DeleteByName(ctx context.Context, name string) (bool, error)
}
