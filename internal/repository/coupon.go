// internal/repository/coupon.go
package repository

import (
	"context"

	"github.com/rpranjan11/coupon-issuance-system/internal/domain"
)

// CouponRepository defines the interface for coupon persistence
type CouponRepository interface {
	// Create saves a new coupon
	Create(ctx context.Context, coupon *domain.Coupon) error

	// GetByCampaign retrieves all coupons for a campaign
	GetByCampaign(ctx context.Context, campaignID string) ([]*domain.Coupon, error)

	// DeleteByCampaignID deletes all coupons for a specific campaign
	DeleteByCampaignID(ctx context.Context, campaignID string) error
}
