// internal/domain/coupon.go
package domain

import (
	"time"
)

// Coupon represents an issued coupon
type Coupon struct {
	Code       string    `json:"code"`
	CampaignID string    `json:"campaign_id"`
	IssuedAt   time.Time `json:"issued_at"`
}
