// internal/domain/campaign.go
package domain

import (
	"time"
)

type Campaign struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	TotalCoupons  int       `json:"total_coupons"`
	IssuedCoupons int       `json:"issued_coupons"`
	StartTime     time.Time `json:"start_time"`
	CreatedAt     time.Time `json:"created_at"`
}

// CanIssue checks if a coupon can be issued for this campaign
func (c *Campaign) CanIssue() bool {
	return c.IssuedCoupons < c.TotalCoupons && time.Now().After(c.StartTime)
}

// HasStarted checks if the campaign has started
func (c *Campaign) HasStarted() bool {
	return time.Now().After(c.StartTime)
}

// RemainingCoupons returns the number of remaining coupons
func (c *Campaign) RemainingCoupons() int {
	return c.TotalCoupons - c.IssuedCoupons
}
