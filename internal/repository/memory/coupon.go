// internal/repository/memory/coupon.go
package memory

import (
	"context"
	"sync"

	"github.com/rpranjan11/coupon-issuance-system/internal/domain"
	"github.com/rpranjan11/coupon-issuance-system/internal/repository"
)

// CouponRepository is an in-memory implementation of repository.CouponRepository
type CouponRepository struct {
	coupons map[string][]*domain.Coupon
	mutex   sync.RWMutex
}

// NewCouponRepository creates a new in-memory coupon repository
func NewCouponRepository() repository.CouponRepository {
	return &CouponRepository{
		coupons: make(map[string][]*domain.Coupon),
	}
}

// Create saves a new coupon
func (r *CouponRepository) Create(ctx context.Context, coupon *domain.Coupon) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.coupons[coupon.CampaignID]; !exists {
		r.coupons[coupon.CampaignID] = make([]*domain.Coupon, 0)
	}

	r.coupons[coupon.CampaignID] = append(r.coupons[coupon.CampaignID], coupon)
	return nil
}

// GetByCampaign retrieves all coupons for a campaign
func (r *CouponRepository) GetByCampaign(ctx context.Context, campaignID string) ([]*domain.Coupon, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	coupons, exists := r.coupons[campaignID]
	if !exists {
		return make([]*domain.Coupon, 0), nil
	}

	// Return a copy to prevent concurrent modification
	result := make([]*domain.Coupon, len(coupons))
	copy(result, coupons)

	return result, nil
}
