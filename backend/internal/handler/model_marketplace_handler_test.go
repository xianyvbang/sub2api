//go:build unit

package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	servermiddleware "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type marketplaceChannelServiceStub struct {
	cards []service.ModelMarketplaceCard
}

func (s *marketplaceChannelServiceStub) ListModelMarketplace(_ context.Context, visibleGroups []service.Group) ([]service.ModelMarketplaceCard, error) {
	allowed := make(map[int64]struct{}, len(visibleGroups))
	for _, group := range visibleGroups {
		allowed[group.ID] = struct{}{}
	}

	out := make([]service.ModelMarketplaceCard, 0, len(s.cards))
	for _, card := range s.cards {
		filteredGroups := make([]service.ModelMarketplaceGroupOffer, 0, len(card.Groups))
		for _, group := range card.Groups {
			if _, ok := allowed[group.GroupID]; ok {
				filteredGroups = append(filteredGroups, group)
			}
		}
		if len(filteredGroups) == 0 {
			continue
		}
		cp := card
		cp.GroupID = filteredGroups[0].GroupID
		cp.GroupName = filteredGroups[0].GroupName
		cp.GroupPlatform = filteredGroups[0].GroupPlatform
		cp.GroupRate = filteredGroups[0].GroupRate
		cp.GroupIsExclusive = filteredGroups[0].GroupIsExclusive
		cp.SubscriptionType = filteredGroups[0].SubscriptionType
		cp.BillingType = filteredGroups[0].BillingType
		cp.PricingSource = filteredGroups[0].PricingSource
		cp.OriginalPricing = filteredGroups[0].OriginalPricing
		cp.CurrentPricing = filteredGroups[0].CurrentPricing
		cp.Groups = filteredGroups
		out = append(out, cp)
	}
	return out, nil
}

type marketplaceAPIKeyServiceStub struct {
	groups []service.Group
}

func (s *marketplaceAPIKeyServiceStub) GetAvailableGroups(_ context.Context, _ int64) ([]service.Group, error) {
	return s.groups, nil
}

type marketplaceGroupServiceStub struct {
	groups []service.Group
}

func (s *marketplaceGroupServiceStub) ListPublicMarketplaceGroups(_ context.Context) ([]service.Group, error) {
	return s.groups, nil
}

type marketplaceSettingServiceStub struct {
	runtime service.ModelMarketplaceRuntime
}

func (s *marketplaceSettingServiceStub) GetModelMarketplaceRuntime(_ context.Context) service.ModelMarketplaceRuntime {
	return s.runtime
}

func newMarketplaceHandler() *ModelMarketplaceHandler {
	publicOriginal := 0.01
	publicCurrent := 0.012
	exclusiveCurrent := 0.009
	subscriptionPrice := 0.02

	return &ModelMarketplaceHandler{
		channelService: &marketplaceChannelServiceStub{
			cards: []service.ModelMarketplaceCard{
				{
					GroupID:          1,
					GroupName:        "public",
					GroupPlatform:    "openai",
					GroupRate:        1.2,
					GroupIsExclusive: false,
					SubscriptionType: service.SubscriptionTypeStandard,
					ModelName:        "gpt-4o",
					Supplier:         "openai",
					BillingType:      string(service.BillingModePerRequest),
					PricingSource:    service.ModelMarketplacePricingSourceGroup,
					OriginalPricing: &service.ChannelModelPricing{
						BillingMode:     service.BillingModePerRequest,
						PerRequestPrice: &publicOriginal,
					},
					CurrentPricing: &service.ChannelModelPricing{
						BillingMode:     service.BillingModePerRequest,
						PerRequestPrice: &publicCurrent,
					},
					Groups: []service.ModelMarketplaceGroupOffer{
						{
							GroupID:          1,
							GroupName:        "public",
							GroupPlatform:    "openai",
							GroupRate:        1.2,
							GroupIsExclusive: false,
							SubscriptionType: service.SubscriptionTypeStandard,
							ModelName:        "gpt-4o",
							Supplier:         "openai",
							BillingType:      string(service.BillingModePerRequest),
							PricingSource:    service.ModelMarketplacePricingSourceGroup,
							OriginalPricing: &service.ChannelModelPricing{
								BillingMode:     service.BillingModePerRequest,
								PerRequestPrice: &publicOriginal,
							},
							CurrentPricing: &service.ChannelModelPricing{
								BillingMode:     service.BillingModePerRequest,
								PerRequestPrice: &publicCurrent,
							},
						},
						{
							GroupID:          2,
							GroupName:        "exclusive",
							GroupPlatform:    "openai",
							GroupRate:        1.8,
							GroupIsExclusive: true,
							SubscriptionType: service.SubscriptionTypeStandard,
							ModelName:        "gpt-4o",
							Supplier:         "openai",
							BillingType:      string(service.BillingModePerRequest),
							PricingSource:    service.ModelMarketplacePricingSourceChannel,
							OriginalPricing: &service.ChannelModelPricing{
								BillingMode:     service.BillingModePerRequest,
								PerRequestPrice: &exclusiveCurrent,
							},
							CurrentPricing: &service.ChannelModelPricing{
								BillingMode:     service.BillingModePerRequest,
								PerRequestPrice: &exclusiveCurrent,
							},
						},
					},
				},
				{
					GroupID:          3,
					GroupName:        "public-subscription",
					GroupPlatform:    "openai",
					GroupRate:        1.5,
					GroupIsExclusive: false,
					SubscriptionType: service.SubscriptionTypeSubscription,
					ModelName:        "gpt-4o-mini",
					Supplier:         "openai",
					BillingType:      string(service.BillingModePerRequest),
					PricingSource:    service.ModelMarketplacePricingSourceGroup,
					OriginalPricing: &service.ChannelModelPricing{
						BillingMode:     service.BillingModePerRequest,
						PerRequestPrice: &subscriptionPrice,
					},
					CurrentPricing: &service.ChannelModelPricing{
						BillingMode:     service.BillingModePerRequest,
						PerRequestPrice: &subscriptionPrice,
					},
					Groups: []service.ModelMarketplaceGroupOffer{
						{
							GroupID:          3,
							GroupName:        "public-subscription",
							GroupPlatform:    "openai",
							GroupRate:        1.5,
							GroupIsExclusive: false,
							SubscriptionType: service.SubscriptionTypeSubscription,
							ModelName:        "gpt-4o-mini",
							Supplier:         "openai",
							BillingType:      string(service.BillingModePerRequest),
							PricingSource:    service.ModelMarketplacePricingSourceGroup,
							OriginalPricing: &service.ChannelModelPricing{
								BillingMode:     service.BillingModePerRequest,
								PerRequestPrice: &subscriptionPrice,
							},
							CurrentPricing: &service.ChannelModelPricing{
								BillingMode:     service.BillingModePerRequest,
								PerRequestPrice: &subscriptionPrice,
							},
						},
					},
				},
			},
		},
		apiKeyService: &marketplaceAPIKeyServiceStub{
			groups: []service.Group{
				{ID: 1, Name: "public", Platform: "openai", IsExclusive: false},
				{ID: 2, Name: "exclusive", Platform: "openai", IsExclusive: true},
			},
		},
		groupService: &marketplaceGroupServiceStub{
			groups: []service.Group{
				{ID: 1, Name: "public", Platform: "openai", IsExclusive: false},
				{ID: 3, Name: "public-subscription", Platform: "openai", IsExclusive: false, SubscriptionType: service.SubscriptionTypeSubscription},
			},
		},
		settingService: &marketplaceSettingServiceStub{
			runtime: service.ModelMarketplaceRuntime{
				Enabled:       true,
				RequiresLogin: true,
			},
		},
	}
}

func decodeMarketplaceRows(t *testing.T, body []byte) []map[string]any {
	t.Helper()
	var payload struct {
		Code int              `json:"code"`
		Data []map[string]any `json:"data"`
	}
	require.NoError(t, json.Unmarshal(body, &payload))
	return payload.Data
}

func TestModelMarketplace_AnonymousSeesOnlyPublicCards(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := newMarketplaceHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/model-marketplace", nil)

	h.List(c)

	require.Equal(t, http.StatusOK, w.Code)
	rows := decodeMarketplaceRows(t, w.Body.Bytes())
	require.Len(t, rows, 2)
	require.Equal(t, "public", rows[0]["group_name"])
	require.Equal(t, "public-subscription", rows[1]["group_name"])

	first := rows[0]
	require.Equal(t, "group", first["pricing_source"])
	require.Contains(t, first, "original_pricing")
	require.Contains(t, first, "current_pricing")
	groups, ok := first["groups"].([]any)
	require.True(t, ok)
	require.Len(t, groups, 1)
}

func TestModelMarketplace_AuthenticatedCanSeeExclusiveGroupOffer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := newMarketplaceHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/model-marketplace", nil)
	c.Set(string(servermiddleware.ContextKeyUser), servermiddleware.AuthSubject{UserID: 7})

	h.List(c)

	require.Equal(t, http.StatusOK, w.Code)
	rows := decodeMarketplaceRows(t, w.Body.Bytes())
	require.Len(t, rows, 2)

	firstGroups, ok := rows[0]["groups"].([]any)
	require.True(t, ok)
	require.Len(t, firstGroups, 2)
}

func TestModelMarketplace_DisabledReturnsEmptyList(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := newMarketplaceHandler()
	h.settingService = &marketplaceSettingServiceStub{
		runtime: service.ModelMarketplaceRuntime{
			Enabled:       false,
			RequiresLogin: true,
		},
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/model-marketplace", nil)

	h.List(c)

	require.Equal(t, http.StatusOK, w.Code)
	rows := decodeMarketplaceRows(t, w.Body.Bytes())
	require.Len(t, rows, 0)
}
