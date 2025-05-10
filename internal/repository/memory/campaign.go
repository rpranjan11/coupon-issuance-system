// internal/repository/memory/campaign.go
package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/rpranjan11/coupon-issuance-system/internal/domain"
	"github.com/rpranjan11/coupon-issuance-system/internal/repository"
)

var (
	ErrCampaignNotFound = errors.New("campaign not found")
	ErrLimitReached     = errors.New("coupon limit reached")
)

// CampaignRepository is an in-memory implementation of repository.CampaignRepository
type CampaignRepository struct {
	campaigns map[string]*domain.Campaign
	mutex     sync.RWMutex
}

// NewCampaignRepository creates a new in-memory campaign repository
func NewCampaignRepository() repository.CampaignRepository {
	return &CampaignRepository{
		campaigns: make(map[string]*domain.Campaign),
	}
}

// Create saves a new campaign
func (r *CampaignRepository) Create(ctx context.Context, campaign *domain.Campaign) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.campaigns[campaign.ID] = campaign
	return nil
}

// Get retrieves a campaign by ID
func (r *CampaignRepository) Get(ctx context.Context, id string) (*domain.Campaign, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	campaign, exists := r.campaigns[id]
	if !exists {
		return nil, ErrCampaignNotFound
	}

	return campaign, nil
}

// Update updates an existing campaign
func (r *CampaignRepository) Update(ctx context.Context, campaign *domain.Campaign) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.campaigns[campaign.ID]
	if !exists {
		return ErrCampaignNotFound
	}

	r.campaigns[campaign.ID] = campaign
	return nil
}

// AtomicIncrementIssued atomically increments the issued_coupons counter
func (r *CampaignRepository) AtomicIncrementIssued(ctx context.Context, campaignID string) (bool, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	campaign, exists := r.campaigns[campaignID]
	if !exists {
		return false, ErrCampaignNotFound
	}

	if campaign.IssuedCoupons >= campaign.TotalCoupons {
		return false, ErrLimitReached
	}

	// Check if the campaign has started
	if time.Now().Before(campaign.StartTime) {
		return false, errors.New("campaign has not started yet")
	}

	campaign.IssuedCoupons++
	return true, nil
}

// FindByName finds a campaign by its name
func (r *CampaignRepository) FindByName(ctx context.Context, name string) (*domain.Campaign, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, campaign := range r.campaigns {
		if campaign.Name == name {
			return campaign, nil
		}
	}

	return nil, ErrCampaignNotFound
}

// DeleteByID deletes a campaign by ID
func (r *CampaignRepository) DeleteByID(ctx context.Context, id string) (bool, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.campaigns[id]
	if !exists {
		return false, nil
	}

	delete(r.campaigns, id)
	return true, nil
}

// DeleteByName deletes a campaign by name
func (r *CampaignRepository) DeleteByName(ctx context.Context, name string) (bool, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for id, campaign := range r.campaigns {
		if campaign.Name == name {
			delete(r.campaigns, id)
			return true, nil
		}
	}

	return false, nil
}
