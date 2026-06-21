//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListModelMarketplace_GroupsByPlatformAndUsesCatalogModels(t *testing.T) {
	svc := NewChannelService(&mockChannelRepository{}, &stubGroupRepoForAvailable{}, nil, &PricingService{
		pricingData: map[string]*LiteLLMModelPricing{
			"gpt-4o": {
				InputCostPerToken:           1,
				OutputCostPerToken:          2,
				LiteLLMProvider:             "openai",
				Mode:                        "chat",
				SupportsPromptCaching:       true,
				OutputCostPerImageToken:     0.03,
				CacheReadInputTokenCost:     0.04,
				CacheCreationInputTokenCost: 0.05,
			},
			"claude-sonnet-4": {
				InputCostPerToken:     3,
				OutputCostPerToken:    4,
				LiteLLMProvider:       "anthropic",
				Mode:                  "chat",
				SupportsPromptCaching: true,
			},
		},
	})

	cards, err := svc.ListModelMarketplace(context.Background(), []Group{
		{ID: 1, Name: "openai-public", Platform: "openai", IsExclusive: false, RateMultiplier: 1.1, SubscriptionType: SubscriptionTypeStandard, Status: StatusActive},
		{ID: 2, Name: "anthropic-pro", Platform: "anthropic", IsExclusive: true, RateMultiplier: 1.8, SubscriptionType: SubscriptionTypeSubscription, Status: StatusActive},
	})
	require.NoError(t, err)
	require.Len(t, cards, 2)

	cardsByGroup := make(map[string]ModelMarketplaceCard, len(cards))
	for _, card := range cards {
		cardsByGroup[card.GroupName] = card
	}

	require.Contains(t, cardsByGroup, "openai-public")
	require.Contains(t, cardsByGroup, "anthropic-pro")
	require.Equal(t, "gpt-4o", cardsByGroup["openai-public"].ModelName)
	require.Equal(t, "openai", cardsByGroup["openai-public"].Platform)
	require.Equal(t, string(BillingModeToken), cardsByGroup["openai-public"].BillingType)
	require.NotNil(t, cardsByGroup["openai-public"].Pricing)
	require.Equal(t, "claude-sonnet-4", cardsByGroup["anthropic-pro"].ModelName)
	require.Equal(t, "anthropic", cardsByGroup["anthropic-pro"].Platform)
	require.NotNil(t, cardsByGroup["anthropic-pro"].Pricing)
}

func TestListModelMarketplace_SkipsPlatformsWithoutCatalogModels(t *testing.T) {
	svc := NewChannelService(&mockChannelRepository{}, &stubGroupRepoForAvailable{}, nil, &PricingService{
		pricingData: map[string]*LiteLLMModelPricing{
			"gpt-4o": {
				InputCostPerToken:     1,
				OutputCostPerToken:    2,
				LiteLLMProvider:       "openai",
				Mode:                  "chat",
				SupportsPromptCaching: true,
			},
		},
	})

	cards, err := svc.ListModelMarketplace(context.Background(), []Group{
		{ID: 1, Name: "openai-public", Platform: "openai", Status: StatusActive},
		{ID: 2, Name: "unknown-platform", Platform: "mystery", Status: StatusActive},
	})
	require.NoError(t, err)
	require.Len(t, cards, 1)
	require.Equal(t, "openai-public", cards[0].GroupName)
}

func TestListPublicMarketplaceGroups_AnonymousOnlySeesNonExclusive(t *testing.T) {
	groupSvc := NewGroupService(&stubGroupRepoForAvailable{
		activeGroups: []Group{
			{ID: 1, Name: "public", IsExclusive: false},
			{ID: 2, Name: "private", IsExclusive: true},
		},
	}, nil)

	groups, err := groupSvc.ListPublicMarketplaceGroups(context.Background())
	require.NoError(t, err)
	require.Len(t, groups, 1)
	require.Equal(t, int64(1), groups[0].ID)
}
