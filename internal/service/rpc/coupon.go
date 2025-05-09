// internal/service/rpc/coupon.go
package rpc

import (
	"context"
	"errors"
	"time"

	"github.com/bufbuild/connect-go"
	"google.golang.org/protobuf/types/known/timestamppb"

	coupon "github.com/rpranjan11/coupon-issuance-system/api/coupon"
	"github.com/rpranjan11/coupon-issuance-system/internal/domain"
	"github.com/rpranjan11/coupon-issuance-system/internal/service"
)

// CouponServiceServer implements the CouponService Connect API
type CouponServiceServer struct {
	campaignService *service.CampaignService
}

// NewCouponServiceServer creates a new CouponServiceServer
func NewCouponServiceServer(campaignService *service.CampaignService) *CouponServiceServer {
	return &CouponServiceServer{
		campaignService: campaignService,
	}
}

// CreateCampaign creates a new coupon campaign
func (s *CouponServiceServer) CreateCampaign(
	ctx context.Context,
	req *connect.Request[coupon.CreateCampaignRequest],
) (*connect.Response[coupon.CreateCampaignResponse], error) {
	// Validate request
	if req.Msg.Name == "" || req.Msg.TotalCoupons <= 0 || req.Msg.StartTime == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("invalid request parameters"))
	}

	// Extract start time
	startTime := req.Msg.StartTime.AsTime()

	// Create campaign
	campaign, err := s.campaignService.CreateCampaign(ctx, req.Msg.Name, int(req.Msg.TotalCoupons), startTime)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Convert domain model to proto model
	campaignProto := &coupon.Campaign{
		Id:            campaign.ID,
		Name:          campaign.Name,
		TotalCoupons:  int32(campaign.TotalCoupons),
		IssuedCoupons: int32(campaign.IssuedCoupons),
		StartTime:     timestamppb.New(campaign.StartTime),
		CreatedAt:     timestamppb.New(campaign.CreatedAt),
	}

	return connect.NewResponse(&coupon.CreateCampaignResponse{
		Campaign: campaignProto,
	}), nil
}

// GetCampaign gets campaign information including all issued coupon codes
func (s *CouponServiceServer) GetCampaign(
	ctx context.Context,
	req *connect.Request[coupon.GetCampaignRequest],
) (*connect.Response[coupon.GetCampaignResponse], error) {
	// Validate request
	if req.Msg.CampaignId == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("campaign ID is required"))
	}

	// Get campaign and coupons
	campaign, coupons, err := s.campaignService.GetCampaign(ctx, req.Msg.CampaignId)
	if err != nil {
		if errors.Is(err, service.ErrCampaignNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Convert domain models to proto models
	campaignProto := &coupon.Campaign{
		Id:            campaign.ID,
		Name:          campaign.Name,
		TotalCoupons:  int32(campaign.TotalCoupons),
		IssuedCoupons: int32(campaign.IssuedCoupons),
		StartTime:     timestamppb.New(campaign.StartTime),
		CreatedAt:     timestamppb.New(campaign.CreatedAt),
	}

	couponProtos := make([]*coupon.Coupon, len(coupons))
	for i, coupon := range coupons {
		couponProtos[i] = &coupon.Coupon{
			Code:       coupon.Code,
			CampaignId: coupon.CampaignID,
			IssuedAt:   timestamppb.New(coupon.IssuedAt),
		}
	}

	return connect.NewResponse(&coupon.GetCampaignResponse{
		Campaign: campaignProto,
		Coupons:  couponProtos,
	}), nil
}

// IssueCoupon requests coupon issuance on a specific campaign
func (s *CouponServiceServer) IssueCoupon(
	ctx context.Context,
	req *connect.Request[coupon.IssueCouponRequest],
) (*connect.Response[coupon.IssueCouponResponse], error) {
	// Validate request
	if req.Msg.CampaignId == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("campaign ID is required"))
	}

	// Issue coupon
	coupon, err := s.campaignService.IssueCoupon(ctx, req.Msg.CampaignId)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrCampaignNotFound):
			return nil, connect.NewError(connect.CodeNotFound, err)
		case errors.Is(err, service.ErrCampaignNotStarted):
			return connect.NewResponse(&coupon.IssueCouponResponse{
				Success: false,
				Error:   "campaign has not started yet",
			}), nil
		case errors.Is(err, service.ErrNoMoreCoupons):
			return connect.NewResponse(&coupon.IssueCouponResponse{
				Success: false,
				Error:   "no more coupons available",
			}), nil
		default:
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	// Convert domain model to proto model
	couponProto := &coupon.Coupon{
		Code:       coupon.Code,
		CampaignId: coupon.CampaignID,
		IssuedAt:   timestamppb.New(coupon.IssuedAt),
	}

	return connect.NewResponse(&coupon.IssueCouponResponse{
		Success: true,
		Coupon:  couponProto,
	}), nil
}
