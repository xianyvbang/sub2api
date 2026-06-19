//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListModelMarketplace_FlattensPerGroupPlatformModel(t *testing.T) {
	repo := &mockChannelRepository{
		listAllFn: func(ctx context.Context) ([]Channel, error) {
			price := 0.01
			return []Channel{{
				ID:     1,
				Name:   "channel-a",
				Status: StatusActive,
				GroupIDs: []int64{
					1, 2,
				},
				ModelPricing: []ChannelModelPricing{{
					Platform:        "openai",
					Models:          []string{"gpt-4o"},
					BillingMode:     BillingModePerRequest,
					PerRequestPrice: &price,
				}},
			}}, nil
		},
	}
	groupRepo := &stubGroupRepoForAvailable{
		activeGroups: []Group{
			{ID: 1, Name: "public", Platform: "openai", IsExclusive: false, RateMultiplier: 1.2, SubscriptionType: SubscriptionTypeStandard},
			{ID: 2, Name: "pro", Platform: "openai", IsExclusive: true, RateMultiplier: 1.8, SubscriptionType: SubscriptionTypeSubscription},
		},
	}
	svc := NewChannelService(repo, groupRepo, nil, nil)

	cards, err := svc.ListModelMarketplace(context.Background(), []Group{
		{ID: 1, Name: "public", Platform: "openai", IsExclusive: false, RateMultiplier: 1.2, SubscriptionType: SubscriptionTypeStandard},
		{ID: 2, Name: "pro", Platform: "openai", IsExclusive: true, RateMultiplier: 1.8, SubscriptionType: SubscriptionTypeSubscription},
	})
	require.NoError(t, err)
	require.Len(t, cards, 2)

	gotGroups := []string{cards[0].GroupName, cards[1].GroupName}
	require.ElementsMatch(t, []string{"public", "pro"}, gotGroups)

	cardsByGroup := make(map[string]ModelMarketplaceCard, len(cards))
	for _, card := range cards {
		cardsByGroup[card.GroupName] = card
	}

	require.Equal(t, "gpt-4o", cardsByGroup["public"].ModelName)
	require.Equal(t, "openai", cardsByGroup["public"].Platform)
	require.Equal(t, string(BillingModePerRequest), cardsByGroup["public"].BillingType)
	require.Equal(t, "gpt-4o", cardsByGroup["pro"].ModelName)
}

func TestListModelMarketplace_DedupsSameGroupPlatformModelAcrossChannels(t *testing.T) {
	price := 0.01
	repo := &mockChannelRepository{
		listAllFn: func(ctx context.Context) ([]Channel, error) {
			return []Channel{
				{
					ID:       1,
					Name:     "a",
					Status:   StatusActive,
					GroupIDs: []int64{1},
					ModelPricing: []ChannelModelPricing{{
						Platform:        "openai",
						Models:          []string{"gpt-4o"},
						BillingMode:     BillingModePerRequest,
						PerRequestPrice: &price,
					}},
				},
				{
					ID:       2,
					Name:     "b",
					Status:   StatusActive,
					GroupIDs: []int64{1},
					ModelPricing: []ChannelModelPricing{{
						Platform:        "openai",
						Models:          []string{"GPT-4O"},
						BillingMode:     BillingModePerRequest,
						PerRequestPrice: &price,
					}},
				},
			}, nil
		},
	}
	groupRepo := &stubGroupRepoForAvailable{
		activeGroups: []Group{{ID: 1, Name: "public", Platform: "openai"}},
	}
	svc := NewChannelService(repo, groupRepo, nil, nil)

	cards, err := svc.ListModelMarketplace(context.Background(), []Group{{ID: 1, Name: "public", Platform: "openai"}})
	require.NoError(t, err)
	require.Len(t, cards, 1)
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
