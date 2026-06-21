package handler

import (
	"context"

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

type modelMarketplaceGroupDTO struct {
	GroupID          int64                       `json:"group_id"`
	GroupName        string                      `json:"group_name"`
	GroupPlatform    string                      `json:"group_platform"`
	GroupRate        float64                     `json:"group_rate"`
	GroupIsExclusive bool                        `json:"group_is_exclusive"`
	SubscriptionType string                      `json:"subscription_type"`
	ModelName        string                      `json:"model_name"`
	Supplier         string                      `json:"supplier"`
	BillingType      string                      `json:"billing_type"`
	PricingSource    string                      `json:"pricing_source"`
	OriginalPricing  *modelMarketplacePricingDTO `json:"original_pricing"`
	CurrentPricing   *modelMarketplacePricingDTO `json:"current_pricing"`
}

type modelMarketplaceCardDTO struct {
	GroupID          int64                       `json:"group_id"`
	GroupName        string                      `json:"group_name"`
	GroupPlatform    string                      `json:"group_platform"`
	GroupRate        float64                     `json:"group_rate"`
	GroupIsExclusive bool                        `json:"group_is_exclusive"`
	SubscriptionType string                      `json:"subscription_type"`
	ModelName        string                      `json:"model_name"`
	Supplier         string                      `json:"supplier"`
	BillingType      string                      `json:"billing_type"`
	PricingSource    string                      `json:"pricing_source"`
	OriginalPricing  *modelMarketplacePricingDTO `json:"original_pricing"`
	CurrentPricing   *modelMarketplacePricingDTO `json:"current_pricing"`
	Groups           []modelMarketplaceGroupDTO  `json:"groups"`
}

// ModelMarketplaceHandler serves the public model marketplace API.
type ModelMarketplaceHandler struct {
	channelService modelMarketplaceChannelLister
	apiKeyService  modelMarketplaceAPIKeyGroupProvider
	groupService   modelMarketplaceGroupProvider
	settingService modelMarketplaceSettingRuntimeProvider
}

func NewModelMarketplaceHandler(
	channelService *service.ModelMarketplaceService,
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

// List returns one aggregated marketplace card per supplier + model.
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

	publicGroups, err := h.groupService.ListPublicMarketplaceGroups(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	visibleGroups := publicGroups
	if subject, ok := middleware.GetAuthSubjectFromContext(c); ok {
		allowedGroups, err := h.apiKeyService.GetAvailableGroups(c.Request.Context(), subject.UserID)
		if err != nil {
			response.ErrorFrom(c, err)
			return
		}
		visibleGroups = mergeMarketplaceGroups(publicGroups, allowedGroups)
	}

	cards, err := h.channelService.ListModelMarketplace(c.Request.Context(), visibleGroups)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	out := make([]modelMarketplaceCardDTO, 0, len(cards))
	for _, card := range cards {
		groups := make([]modelMarketplaceGroupDTO, 0, len(card.Groups))
		for _, group := range card.Groups {
			groups = append(groups, modelMarketplaceGroupDTO{
				GroupID:          group.GroupID,
				GroupName:        group.GroupName,
				GroupPlatform:    group.GroupPlatform,
				GroupRate:        group.GroupRate,
				GroupIsExclusive: group.GroupIsExclusive,
				SubscriptionType: group.SubscriptionType,
				ModelName:        group.ModelName,
				Supplier:         group.Supplier,
				BillingType:      group.BillingType,
				PricingSource:    group.PricingSource,
				OriginalPricing:  toMarketplacePricing(group.OriginalPricing),
				CurrentPricing:   toMarketplacePricing(group.CurrentPricing),
			})
		}

		out = append(out, modelMarketplaceCardDTO{
			GroupID:          card.GroupID,
			GroupName:        card.GroupName,
			GroupPlatform:    card.GroupPlatform,
			GroupRate:        card.GroupRate,
			GroupIsExclusive: card.GroupIsExclusive,
			SubscriptionType: card.SubscriptionType,
			ModelName:        card.ModelName,
			Supplier:         card.Supplier,
			BillingType:      card.BillingType,
			PricingSource:    card.PricingSource,
			OriginalPricing:  toMarketplacePricing(card.OriginalPricing),
			CurrentPricing:   toMarketplacePricing(card.CurrentPricing),
			Groups:           groups,
		})
	}

	response.Success(c, out)
}

func mergeMarketplaceGroups(groupLists ...[]service.Group) []service.Group {
	seen := make(map[int64]struct{})
	merged := make([]service.Group, 0)
	for _, groups := range groupLists {
		for _, group := range groups {
			if _, ok := seen[group.ID]; ok {
				continue
			}
			seen[group.ID] = struct{}{}
			merged = append(merged, group)
		}
	}
	return merged
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
