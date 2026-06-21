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
	out := make([]service.ModelMarketplaceCard, 0, len(visibleGroups))
	allowed := make(map[int64]struct{}, len(visibleGroups))
	for _, group := range visibleGroups {
		allowed[group.ID] = struct{}{}
	}
	for _, card := range s.cards {
		if _, ok := allowed[card.GroupID]; ok {
			out = append(out, card)
		}
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
	price := 0.02
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
					Platform:         "openai",
					BillingType:      string(service.BillingModePerRequest),
					Pricing: &service.ChannelModelPricing{
						BillingMode:     service.BillingModePerRequest,
						PerRequestPrice: &price,
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
					Platform:         "openai",
					BillingType:      string(service.BillingModePerRequest),
					Pricing: &service.ChannelModelPricing{
						BillingMode:     service.BillingModePerRequest,
						PerRequestPrice: &price,
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
					Platform:         "openai",
					BillingType:      string(service.BillingModePerRequest),
					Pricing: &service.ChannelModelPricing{
						BillingMode:     service.BillingModePerRequest,
						PerRequestPrice: &price,
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

func TestModelMarketplace_AnonymousSeesAllActiveGroups(t *testing.T) {
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
}

func TestModelMarketplace_AuthenticatedSeesAllActiveGroups(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := newMarketplaceHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/model-marketplace", nil)
	c.Set(string(servermiddleware.ContextKeyUser), servermiddleware.AuthSubject{UserID: 7})

	h.List(c)

	require.Equal(t, http.StatusOK, w.Code)
	rows := decodeMarketplaceRows(t, w.Body.Bytes())
	require.Len(t, rows, 3)
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
