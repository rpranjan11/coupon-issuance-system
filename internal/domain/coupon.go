// internal/domain/coupon.go
package domain

import (
	"time"
)

type Coupon struct {
	Code       string    `json:"code"`
	CampaignID string    `json:"campaign_id"`
	IssuedAt   time.Time `json:"issued_at"`
}
