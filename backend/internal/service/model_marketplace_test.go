//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListModelMarketplace_AggregatesModelsAndChoosesLowestCurrentPrice(t *testing.T) {
	svc := NewChannelService(&mockChannelRepository{}, &stubGroupRepoForAvailable{}, nil, &PricingService{
		pricingData: map[string]*LiteLLMModelPricing{
			"gpt-4o": {
				InputCostPerToken:  1e-6,
				OutputCostPerToken: 2e-6,
				LiteLLMProvider:    "openai",
				Mode:               "chat",
			},
		},
	})

	cards, err := svc.ListModelMarketplace(context.Background(), []Group{
		{ID: 1, Name: "public", Platform: "openai", IsExclusive: false, RateMultiplier: 1.2, SubscriptionType: SubscriptionTypeStandard, Status: StatusActive},
		{ID: 2, Name: "promo", Platform: "openai", IsExclusive: false, RateMultiplier: 0.8, SubscriptionType: SubscriptionTypeStandard, Status: StatusActive},
	})
	require.NoError(t, err)
	require.Len(t, cards, 1)

	card := cards[0]
	require.Equal(t, "gpt-4o", card.ModelName)
	require.Equal(t, "promo", card.GroupName)
	require.Equal(t, ModelMarketplacePricingSourceGroup, card.PricingSource)
	require.Len(t, card.Groups, 2)
	require.NotNil(t, card.OriginalPricing)
	require.NotNil(t, card.CurrentPricing)
	require.NotNil(t, card.OriginalPricing.InputPrice)
	require.NotNil(t, card.CurrentPricing.InputPrice)
	require.InDelta(t, 1e-6, *card.OriginalPricing.InputPrice, 1e-12)
	require.InDelta(t, 0.8e-6, *card.CurrentPricing.InputPrice, 1e-12)
}

func TestListModelMarketplace_ChannelPricingWinsWithoutApplyingGroupMultiplier(t *testing.T) {
	channel := Channel{
		ID:       1,
		Status:   StatusActive,
		GroupIDs: []int64{10},
		ModelPricing: []ChannelModelPricing{
			{
				ID:          100,
				Platform:    "openai",
				Models:      []string{"gpt-4o"},
				BillingMode: BillingModeToken,
				InputPrice:  testPtrFloat64(9.9e-5),
				OutputPrice: testPtrFloat64(1.99e-4),
			},
		},
	}

	svc := NewChannelService(
		makeStandardRepo(channel, map[int64]string{10: "openai"}),
		&stubGroupRepoForAvailable{},
		nil,
		&PricingService{
			pricingData: map[string]*LiteLLMModelPricing{
				"gpt-4o": {
					InputCostPerToken:  1e-6,
					OutputCostPerToken: 2e-6,
					LiteLLMProvider:    "openai",
					Mode:               "chat",
				},
			},
		},
	)

	cards, err := svc.ListModelMarketplace(context.Background(), []Group{
		{ID: 10, Name: "enterprise", Platform: "openai", RateMultiplier: 3.5, Status: StatusActive},
	})
	require.NoError(t, err)
	require.Len(t, cards, 1)

	card := cards[0]
	require.Equal(t, ModelMarketplacePricingSourceChannel, card.PricingSource)
	require.NotNil(t, card.OriginalPricing)
	require.NotNil(t, card.CurrentPricing)
	require.NotNil(t, card.OriginalPricing.InputPrice)
	require.NotNil(t, card.CurrentPricing.InputPrice)
	require.InDelta(t, 9.9e-5, *card.OriginalPricing.InputPrice, 1e-12)
	require.InDelta(t, 9.9e-5, *card.CurrentPricing.InputPrice, 1e-12)
}

func TestListModelMarketplace_ModelsWithoutPricingRemainVisibleAndSortLast(t *testing.T) {
	channel := Channel{
		ID:       1,
		Status:   StatusActive,
		GroupIDs: []int64{10},
		ModelMapping: map[string]map[string]string{
			"openai": {
				"mystery-model": "mystery-model",
			},
		},
	}

	svc := NewChannelService(
		makeStandardRepo(channel, map[int64]string{10: "openai"}),
		&stubGroupRepoForAvailable{},
		nil,
		&PricingService{
			pricingData: map[string]*LiteLLMModelPricing{
				"gpt-4o": {
					InputCostPerToken:  1e-6,
					OutputCostPerToken: 2e-6,
					LiteLLMProvider:    "openai",
					Mode:               "chat",
				},
			},
		},
	)

	cards, err := svc.ListModelMarketplace(context.Background(), []Group{
		{ID: 10, Name: "public", Platform: "openai", RateMultiplier: 1.0, Status: StatusActive},
	})
	require.NoError(t, err)
	require.Len(t, cards, 2)
	require.Equal(t, "gpt-4o", cards[0].ModelName)
	require.Equal(t, "mystery-model", cards[1].ModelName)
	require.Nil(t, cards[1].OriginalPricing)
	require.Nil(t, cards[1].CurrentPricing)
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
