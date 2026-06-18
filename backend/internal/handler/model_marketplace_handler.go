package handler

import (
	"context"
	"sort"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

type modelMarketplaceChannelLister interface {
	ListModelMarketplace(ctx context.Context, visibleGroups []service.Group) ([]service.ModelMarketplaceCard, error)
}

type modelMarketplaceAPIKeyGroupProvider interface {
	GetAvailableGroups(ctx context.Context, userID int64) ([]service.Group, error)
}

type modelMarketplaceGroupProvider interface {
	ListPublicMarketplaceGroups(ctx context.Context) ([]service.Group, error)
}

type modelMarketplaceSettingRuntimeProvider interface {
	GetModelMarketplaceRuntime(ctx context.Context) service.ModelMarketplaceRuntime
}

type modelMarketplacePricingDTO struct {
	BillingMode      string                   `json:"billing_mode"`
	InputPrice       *float64                 `json:"input_price"`
	OutputPrice      *float64                 `json:"output_price"`
	CacheWritePrice  *float64                 `json:"cache_write_price"`
	CacheReadPrice   *float64                 `json:"cache_read_price"`
	ImageOutputPrice *float64                 `json:"image_output_price"`
	PerRequestPrice  *float64                 `json:"per_request_price"`
	Intervals        []userPricingIntervalDTO `json:"intervals"`
}

type modelMarketplaceCardDTO struct {
	GroupID          int64                       `json:"group_id"`
	GroupName        string                      `json:"group_name"`
	GroupPlatform    string                      `json:"group_platform"`
	GroupRate        float64                     `json:"group_rate"`
	GroupIsExclusive bool                        `json:"group_is_exclusive"`
	SubscriptionType string                      `json:"subscription_type"`
	ModelName        string                      `json:"model_name"`
	Platform         string                      `json:"platform"`
	BillingType      string                      `json:"billing_type"`
	Pricing          *modelMarketplacePricingDTO `json:"pricing"`
}

// ModelMarketplaceHandler serves the public model marketplace API.
type ModelMarketplaceHandler struct {
	channelService modelMarketplaceChannelLister
	apiKeyService  modelMarketplaceAPIKeyGroupProvider
	groupService   modelMarketplaceGroupProvider
	settingService modelMarketplaceSettingRuntimeProvider
}

func NewModelMarketplaceHandler(
	channelService *service.ChannelService,
	apiKeyService *service.APIKeyService,
	groupService *service.GroupService,
	settingService *service.SettingService,
) *ModelMarketplaceHandler {
	return &ModelMarketplaceHandler{
		channelService: channelService,
		apiKeyService:  apiKeyService,
		groupService:   groupService,
		settingService: settingService,
	}
}

// List returns a flattened card wall payload for the public model marketplace.
// GET /api/v1/model-marketplace
func (h *ModelMarketplaceHandler) List(c *gin.Context) {
	if h.settingService == nil {
		response.Success(c, []modelMarketplaceCardDTO{})
		return
	}

	runtime := h.settingService.GetModelMarketplaceRuntime(c.Request.Context())
	if !runtime.Enabled {
		response.Success(c, []modelMarketplaceCardDTO{})
		return
	}

	var visibleGroups []service.Group
	var err error

	if subject, ok := middleware.GetAuthSubjectFromContext(c); ok {
		visibleGroups, err = h.apiKeyService.GetAvailableGroups(c.Request.Context(), subject.UserID)
	} else {
		visibleGroups, err = h.groupService.ListPublicMarketplaceGroups(c.Request.Context())
	}
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	cards, err := h.channelService.ListModelMarketplace(c.Request.Context(), visibleGroups)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	out := make([]modelMarketplaceCardDTO, 0, len(cards))
	for _, card := range cards {
		out = append(out, modelMarketplaceCardDTO{
			GroupID:          card.GroupID,
			GroupName:        card.GroupName,
			GroupPlatform:    card.GroupPlatform,
			GroupRate:        card.GroupRate,
			GroupIsExclusive: card.GroupIsExclusive,
			SubscriptionType: card.SubscriptionType,
			ModelName:        card.ModelName,
			Platform:         card.Platform,
			BillingType:      card.BillingType,
			Pricing:          toMarketplacePricing(card.Pricing),
		})
	}

	sort.SliceStable(out, func(i, j int) bool {
		if out[i].GroupName != out[j].GroupName {
			return out[i].GroupName < out[j].GroupName
		}
		if out[i].Platform != out[j].Platform {
			return out[i].Platform < out[j].Platform
		}
		return out[i].ModelName < out[j].ModelName
	})

	response.Success(c, out)
}

func toMarketplacePricing(p *service.ChannelModelPricing) *modelMarketplacePricingDTO {
	if p == nil {
		return nil
	}
	return &modelMarketplacePricingDTO{
		BillingMode:      string(p.BillingMode),
		InputPrice:       p.InputPrice,
		OutputPrice:      p.OutputPrice,
		CacheWritePrice:  p.CacheWritePrice,
		CacheReadPrice:   p.CacheReadPrice,
		ImageOutputPrice: p.ImageOutputPrice,
		PerRequestPrice:  p.PerRequestPrice,
		Intervals:        toUserPricing(p).Intervals,
	}
}
