//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type marketplaceAccountRepoStub struct {
	accountsByGroup map[int64][]Account
}

func (s *marketplaceAccountRepoStub) ListSchedulableByGroupID(_ context.Context, groupID int64) ([]Account, error) {
	return append([]Account(nil), s.accountsByGroup[groupID]...), nil
}

func TestListModelMarketplace_AggregatesModelsAndChoosesLowestCurrentPrice(t *testing.T) {
	channelSvc := NewChannelService(&mockChannelRepository{}, &stubGroupRepoForAvailable{}, nil, &PricingService{
		pricingData: map[string]*LiteLLMModelPricing{
			"gpt-4o": {
				InputCostPerToken:  1e-6,
				OutputCostPerToken: 2e-6,
				LiteLLMProvider:    "openai",
				Mode:               "chat",
			},
		},
	})
	svc := NewModelMarketplaceService(&marketplaceAccountRepoStub{
		accountsByGroup: map[int64][]Account{
			1: {{Platform: "openai", Credentials: map[string]any{"model_mapping": map[string]any{"gpt-4o": "gpt-4o"}}}},
			2: {{Platform: "openai", Credentials: map[string]any{"model_mapping": map[string]any{"gpt-4o": "gpt-4o"}}}},
		},
	}, channelSvc, channelSvc.pricingService)

	cards, err := svc.ListModelMarketplace(context.Background(), []Group{
		{ID: 1, Name: "public", Platform: "openai", IsExclusive: false, RateMultiplier: 1.2, SubscriptionType: SubscriptionTypeStandard, Status: StatusActive},
		{ID: 2, Name: "promo", Platform: "openai", IsExclusive: false, RateMultiplier: 0.8, SubscriptionType: SubscriptionTypeStandard, Status: StatusActive},
	})
	require.NoError(t, err)
	require.Len(t, cards, 1)

	card := cards[0]
	require.Equal(t, "gpt-4o", card.ModelName)
	require.Equal(t, "openai", card.Supplier)
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

	channelSvc := NewChannelService(
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
	svc := NewModelMarketplaceService(&marketplaceAccountRepoStub{
		accountsByGroup: map[int64][]Account{
			10: {{Platform: "openai", Credentials: map[string]any{"model_mapping": map[string]any{"gpt-4o": "gpt-4o"}}}},
		},
	}, channelSvc, channelSvc.pricingService)

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
	channelSvc := NewChannelService(
		makeStandardRepo(Channel{ID: 1, Status: StatusActive, GroupIDs: []int64{10}}, map[int64]string{10: "openai"}),
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
	svc := NewModelMarketplaceService(&marketplaceAccountRepoStub{
		accountsByGroup: map[int64][]Account{
			10: {{
				Platform: "openai",
				Credentials: map[string]any{
					"model_mapping": map[string]any{
						"gpt-4o":        "gpt-4o",
						"mystery-model": "mystery-model",
					},
				},
			}},
		},
	}, channelSvc, channelSvc.pricingService)

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

func TestListModelMarketplace_UsesOnlyExplicitModelMappings(t *testing.T) {
	channelSvc := NewChannelService(&mockChannelRepository{}, &stubGroupRepoForAvailable{}, nil, &PricingService{
		pricingData: map[string]*LiteLLMModelPricing{
			"claude-sonnet-4-5": {LiteLLMProvider: "anthropic", InputCostPerToken: 1e-6, OutputCostPerToken: 2e-6, Mode: "chat"},
		},
	})
	svc := NewModelMarketplaceService(&marketplaceAccountRepoStub{
		accountsByGroup: map[int64][]Account{
			11: {
				{Platform: "antigravity", Credentials: map[string]any{}},
				{Platform: "antigravity", Credentials: map[string]any{"model_mapping": map[string]any{"claude-sonnet-4-5": "deepseek-v4-pro"}}},
			},
		},
	}, channelSvc, channelSvc.pricingService)

	cards, err := svc.ListModelMarketplace(context.Background(), []Group{
		{ID: 11, Name: "ag", Platform: "antigravity", RateMultiplier: 1.0, Status: StatusActive},
	})
	require.NoError(t, err)
	require.Len(t, cards, 1)
	require.Equal(t, "claude-sonnet-4-5", cards[0].ModelName)
	require.Equal(t, "anthropic", cards[0].Supplier)
}

func TestListModelMarketplace_ExpandsWildcardMappingsAndDedupes(t *testing.T) {
	channelSvc := NewChannelService(&mockChannelRepository{}, &stubGroupRepoForAvailable{}, nil, &PricingService{
		pricingData: map[string]*LiteLLMModelPricing{
			"claude-sonnet-4-5":      {LiteLLMProvider: "anthropic", InputCostPerToken: 1e-6, OutputCostPerToken: 2e-6, Mode: "chat"},
			"claude-opus-4-5":        {LiteLLMProvider: "anthropic", InputCostPerToken: 3e-6, OutputCostPerToken: 6e-6, Mode: "chat"},
			"claude-opus-4-5-sonnet": {LiteLLMProvider: "anthropic", InputCostPerToken: 4e-6, OutputCostPerToken: 8e-6, Mode: "chat"},
		},
	})
	svc := NewModelMarketplaceService(&marketplaceAccountRepoStub{
		accountsByGroup: map[int64][]Account{
			12: {
				{Platform: "anthropic", Credentials: map[string]any{"model_mapping": map[string]any{"claude-*": "x"}}},
				{Platform: "anthropic", Credentials: map[string]any{"model_mapping": map[string]any{"claude-opus-4-5": "y"}}},
			},
		},
	}, channelSvc, channelSvc.pricingService)

	cards, err := svc.ListModelMarketplace(context.Background(), []Group{
		{ID: 12, Name: "wild", Platform: "anthropic", RateMultiplier: 1.0, Status: StatusActive},
	})
	require.NoError(t, err)
	require.Len(t, cards, 3)
	require.ElementsMatch(t, []string{
		"claude-opus-4-5",
		"claude-opus-4-5-sonnet",
		"claude-sonnet-4-5",
	}, []string{
		cards[0].ModelName,
		cards[1].ModelName,
		cards[2].ModelName,
	})
}

func TestListModelMarketplace_FallsBackSupplierToPricingProvider(t *testing.T) {
	channelSvc := NewChannelService(&mockChannelRepository{}, &stubGroupRepoForAvailable{}, nil, &PricingService{
		pricingData: map[string]*LiteLLMModelPricing{
			"custom-model": {LiteLLMProvider: "openai", InputCostPerToken: 1e-6, OutputCostPerToken: 2e-6, Mode: "chat"},
		},
	})
	svc := NewModelMarketplaceService(&marketplaceAccountRepoStub{
		accountsByGroup: map[int64][]Account{
			13: {
				{Platform: "openai", Credentials: map[string]any{"model_mapping": map[string]any{"custom-model": "custom-model"}}},
			},
		},
	}, channelSvc, channelSvc.pricingService)

	cards, err := svc.ListModelMarketplace(context.Background(), []Group{
		{ID: 13, Name: "custom", Platform: "openai", RateMultiplier: 1.0, Status: StatusActive},
	})
	require.NoError(t, err)
	require.Len(t, cards, 1)
	require.Equal(t, "openai", cards[0].Supplier)
}

func TestListModelMarketplace_ImageModelUsesGroupImagePricing(t *testing.T) {
	groupImage1K := 0.2
	groupImage4K := 0.4

	channelSvc := NewChannelService(&mockChannelRepository{}, &stubGroupRepoForAvailable{}, nil, &PricingService{
		pricingData: map[string]*LiteLLMModelPricing{
			"gpt-image-1.5": {
				LiteLLMProvider: "openai",
				Mode:            "image_generation",
			},
		},
	})
	svc := NewModelMarketplaceService(&marketplaceAccountRepoStub{
		accountsByGroup: map[int64][]Account{
			21: {
				{Platform: "openai", Credentials: map[string]any{"model_mapping": map[string]any{"gpt-image-1.5": "gpt-image-1.5"}}},
			},
		},
	}, channelSvc, channelSvc.pricingService)

	cards, err := svc.ListModelMarketplace(context.Background(), []Group{
		{
			ID:             21,
			Name:           "images",
			Platform:       "openai",
			RateMultiplier: 2.0,
			ImagePrice1K:   &groupImage1K,
			ImagePrice4K:   &groupImage4K,
			Status:         StatusActive,
		},
	})
	require.NoError(t, err)
	require.Len(t, cards, 1)

	card := cards[0]
	require.Equal(t, ModelMarketplacePricingSourceGroup, card.PricingSource)
	require.Equal(t, string(BillingModeImage), card.BillingType)
	require.NotNil(t, card.OriginalPricing)
	require.NotNil(t, card.CurrentPricing)
	require.Len(t, card.OriginalPricing.Intervals, 3)
	require.Len(t, card.CurrentPricing.Intervals, 3)

	originalByTier := make(map[string]float64, len(card.OriginalPricing.Intervals))
	for _, interval := range card.OriginalPricing.Intervals {
		require.NotNil(t, interval.PerRequestPrice)
		originalByTier[interval.TierLabel] = *interval.PerRequestPrice
	}
	require.InDelta(t, 0.2, originalByTier["1K"], 1e-12)
	require.InDelta(t, 0.201, originalByTier["2K"], 1e-12)
	require.InDelta(t, 0.4, originalByTier["4K"], 1e-12)

	currentByTier := make(map[string]float64, len(card.CurrentPricing.Intervals))
	for _, interval := range card.CurrentPricing.Intervals {
		require.NotNil(t, interval.PerRequestPrice)
		currentByTier[interval.TierLabel] = *interval.PerRequestPrice
	}
	require.InDelta(t, 0.4, currentByTier["1K"], 1e-12)
	require.InDelta(t, 0.402, currentByTier["2K"], 1e-12)
	require.InDelta(t, 0.8, currentByTier["4K"], 1e-12)
	require.Nil(t, card.CurrentPricing.InputPrice)
	require.Nil(t, card.CurrentPricing.OutputPrice)
}

func TestListModelMarketplace_ImageModelChannelImagePriceWins(t *testing.T) {
	channel := Channel{
		ID:       1,
		Status:   StatusActive,
		GroupIDs: []int64{22},
		ModelPricing: []ChannelModelPricing{
			{
				ID:               100,
				Platform:         "openai",
				Models:           []string{"gpt-image-1.5"},
				BillingMode:      BillingModeToken,
				InputPrice:       testPtrFloat64(5e-6),
				OutputPrice:      testPtrFloat64(1e-5),
				ImageOutputPrice: testPtrFloat64(3.2e-5),
			},
		},
	}

	channelSvc := NewChannelService(
		makeStandardRepo(channel, map[int64]string{22: "openai"}),
		&stubGroupRepoForAvailable{},
		nil,
		&PricingService{
			pricingData: map[string]*LiteLLMModelPricing{
				"gpt-image-1.5": {
					LiteLLMProvider:         "openai",
					Mode:                    "image_generation",
					OutputCostPerImage:      0.25,
					OutputCostPerToken:      1e-5,
					InputCostPerToken:       5e-6,
					OutputCostPerImageToken: 3.2e-5,
				},
			},
		},
	)
	svc := NewModelMarketplaceService(&marketplaceAccountRepoStub{
		accountsByGroup: map[int64][]Account{
			22: {
				{Platform: "openai", Credentials: map[string]any{"model_mapping": map[string]any{"gpt-image-1.5": "gpt-image-1.5"}}},
			},
		},
	}, channelSvc, channelSvc.pricingService)

	cards, err := svc.ListModelMarketplace(context.Background(), []Group{
		{ID: 22, Name: "images", Platform: "openai", RateMultiplier: 3.0, Status: StatusActive},
	})
	require.NoError(t, err)
	require.Len(t, cards, 1)

	card := cards[0]
	require.Equal(t, ModelMarketplacePricingSourceChannel, card.PricingSource)
	require.Equal(t, string(BillingModeImage), card.BillingType)
	require.NotNil(t, card.CurrentPricing)
	require.NotNil(t, card.CurrentPricing.ImageOutputPrice)
	require.InDelta(t, 3.2e-5, *card.CurrentPricing.ImageOutputPrice, 1e-12)
	require.Nil(t, card.CurrentPricing.InputPrice)
	require.Nil(t, card.CurrentPricing.OutputPrice)
	require.Nil(t, card.CurrentPricing.PerRequestPrice)
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
