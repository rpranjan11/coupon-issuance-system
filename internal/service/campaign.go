// internal/service/campaign.go
package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rpranjan11/coupon-issuance-system/internal/domain"
	"github.com/rpranjan11/coupon-issuance-system/internal/repository"
	"github.com/rpranjan11/coupon-issuance-system/pkg/coupongen"
)

var (
	ErrInvalidRequest     = errors.New("invalid request parameters")
	ErrCampaignNotFound   = errors.New("campaign not found")
	ErrCampaignNotStarted = errors.New("campaign has not started yet")
	ErrNoMoreCoupons      = errors.New("no more coupons available")
	ErrDuplicateCampaign  = errors.New("a campaign with this name already exists")
)

// CampaignService handles campaign-related business logic
type CampaignService struct {
	campaignRepo repository.CampaignRepository
	couponRepo   repository.CouponRepository
}

// NewCampaignService creates a new campaign service
func NewCampaignService(campaignRepo repository.CampaignRepository, couponRepo repository.CouponRepository) *CampaignService {
	return &CampaignService{
		campaignRepo: campaignRepo,
		couponRepo:   couponRepo,
	}
}

// CreateCampaign creates a new coupon campaign
func (s *CampaignService) CreateCampaign(ctx context.Context, name string, totalCoupons int, startTime time.Time) (*domain.Campaign, error) {
	// Validate input
	if name == "" || totalCoupons <= 0 {
		return nil, ErrInvalidRequest
	}

	// Check if a campaign with the same name already exists
	existingCampaign, err := s.campaignRepo.FindByName(ctx, name)
	if err == nil && existingCampaign != nil {
		return nil, ErrDuplicateCampaign
	}

	// Create campaign
	campaign := &domain.Campaign{
		ID:            uuid.New().String(),
		Name:          name,
		TotalCoupons:  totalCoupons,
		IssuedCoupons: 0,
		StartTime:     startTime,
		CreatedAt:     time.Now(),
	}

	// Save campaign
	err = s.campaignRepo.Create(ctx, campaign)
	if err != nil {
		return nil, err
	}

	return campaign, nil
}

// GetCampaign retrieves a campaign by ID
func (s *CampaignService) GetCampaign(ctx context.Context, id string) (*domain.Campaign, []*domain.Coupon, error) {
	// Get campaign
	campaign, err := s.campaignRepo.Get(ctx, id)
	if err != nil {
		return nil, nil, ErrCampaignNotFound
	}

	// Get coupons for this campaign
	coupons, err := s.couponRepo.GetByCampaign(ctx, id)
	if err != nil {
		return campaign, nil, err
	}

	return campaign, coupons, nil
}

// IssueCoupon issues a coupon for a campaign
func (s *CampaignService) IssueCoupon(ctx context.Context, campaignID string) (*domain.Coupon, error) {
	// Get campaign
	campaign, err := s.campaignRepo.Get(ctx, campaignID)
	if err != nil {
		return nil, ErrCampaignNotFound
	}

	// Check if campaign has started
	if !campaign.HasStarted() {
		return nil, ErrCampaignNotStarted
	}

	// Try to atomically increment the issued count
	success, err := s.campaignRepo.AtomicIncrementIssued(ctx, campaignID)
	if err != nil {
		return nil, err
	}
	if !success {
		return nil, ErrNoMoreCoupons
	}

	// Generate unique coupon code
	couponCode := coupongen.GenerateCode(10)

	// Create coupon
	coupon := &domain.Coupon{
		Code:       couponCode,
		CampaignID: campaignID,
		IssuedAt:   time.Now(),
	}

	// Save coupon
	err = s.couponRepo.Create(ctx, coupon)
	if err != nil {
		// This is a critical error - we incremented the counter but failed to save the coupon
		// In a production system, this should be handled with a transaction or compensation logic
		return nil, err
	}

	return coupon, nil
}

// DeleteCampaign deletes a campaign by ID or name
func (s *CampaignService) DeleteCampaign(ctx context.Context, id, name string) (bool, string, error) {
	// If both ID and name are empty, return an error
	if id == "" && name == "" {
		return false, "Campaign ID or name is required", ErrInvalidRequest
	}

	var campaignID string
	var campaignDeleted bool
	var err error

	// Try to delete by ID first if provided
	if id != "" {
		// Store the ID before deleting (we'll need it to delete coupons)
		campaignID = id

		// Delete the campaign
		campaignDeleted, err = s.campaignRepo.DeleteByID(ctx, id)
		if err != nil {
			return false, "Failed to delete campaign", err
		}
	}

	// If ID wasn't provided or deletion by ID failed, try by name
	if !campaignDeleted && name != "" {
		// First, find the campaign to get its ID
		campaign, err := s.campaignRepo.FindByName(ctx, name)
		if err != nil {
			return false, "Failed to find campaign by name", err
		}

		if campaign != nil {
			// Store the ID before deleting
			campaignID = campaign.ID

			// Delete the campaign
			campaignDeleted, err = s.campaignRepo.DeleteByName(ctx, name)
			if err != nil {
				return false, "Failed to delete campaign", err
			}
		}
	}

	// If a campaign was deleted, also delete its coupons
	if campaignDeleted && campaignID != "" {
		err = s.couponRepo.DeleteByCampaignID(ctx, campaignID)
		if err != nil {
			// This is a partial failure - the campaign was deleted but coupons weren't
			// In a production system, this should be handled with transactions
			return true, "Campaign deleted but failed to delete associated coupons", err
		}
		return true, "Campaign and associated coupons deleted successfully", nil
	}

	// If we get here, the campaign wasn't found
	return false, "No such campaign exists", nil
}
