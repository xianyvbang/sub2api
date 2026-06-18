package service

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

// ModelMarketplaceCard is a flattened, user-facing record for one
// group + platform + model combination.
type ModelMarketplaceCard struct {
	GroupID           int64
	GroupName         string
	GroupPlatform     string
	GroupRate         float64
	GroupIsExclusive  bool
	SubscriptionType  string
	ModelName         string
	Platform          string
	BillingType       string
	Pricing           *ChannelModelPricing
}

// ModelMarketplaceRuntime is the lightweight view of the public model marketplace switches.
type ModelMarketplaceRuntime struct {
	Enabled       bool
	RequiresLogin bool
}

// GetModelMarketplaceRuntime reads the marketplace switches directly from the settings store.
// Fail-closed for Enabled, fail-safe-default for RequiresLogin.
func (s *SettingService) GetModelMarketplaceRuntime(ctx context.Context) ModelMarketplaceRuntime {
	vals, err := s.settingRepo.GetMultiple(ctx, []string{
		SettingKeyModelMarketplaceEnabled,
		SettingKeyModelMarketplaceRequiresLogin,
	})
	if err != nil {
		return ModelMarketplaceRuntime{
			Enabled:       false,
			RequiresLogin: true,
		}
	}
	return ModelMarketplaceRuntime{
		Enabled:       vals[SettingKeyModelMarketplaceEnabled] == "true",
		RequiresLogin: vals[SettingKeyModelMarketplaceRequiresLogin] != "false",
	}
}

// ListModelMarketplace returns one record per group + platform + model.
// Anonymous users only see public groups; authenticated users see the groups they may access.
func (s *ChannelService) ListModelMarketplace(
	ctx context.Context,
	visibleGroups []Group,
) ([]ModelMarketplaceCard, error) {
	channels, err := s.ListAvailable(ctx)
	if err != nil {
		return nil, err
	}

	groupByID := make(map[int64]Group, len(visibleGroups))
	for i := range visibleGroups {
		groupByID[visibleGroups[i].ID] = visibleGroups[i]
	}

	type dedupKey struct {
		groupID   int64
		platform  string
		modelName string
	}
	seen := make(map[dedupKey]struct{})
	out := make([]ModelMarketplaceCard, 0)

	for _, ch := range channels {
		if ch.Status != StatusActive {
			continue
		}

		modelsByPlatform := make(map[string][]SupportedModel, 4)
		for _, model := range ch.SupportedModels {
			if strings.TrimSpace(model.Platform) == "" || strings.TrimSpace(model.Name) == "" {
				continue
			}
			modelsByPlatform[model.Platform] = append(modelsByPlatform[model.Platform], model)
		}

		for _, ref := range ch.Groups {
			group, ok := groupByID[ref.ID]
			if !ok {
				continue
			}
			models := modelsByPlatform[group.Platform]
			if len(models) == 0 {
				continue
			}
			for _, model := range models {
				key := dedupKey{
					groupID:   group.ID,
					platform:  strings.ToLower(strings.TrimSpace(model.Platform)),
					modelName: strings.ToLower(strings.TrimSpace(model.Name)),
				}
				if _, exists := seen[key]; exists {
					continue
				}
				seen[key] = struct{}{}

				billingType := ""
				if model.Pricing != nil {
					billingType = string(model.Pricing.BillingMode)
				}

				out = append(out, ModelMarketplaceCard{
					GroupID:          group.ID,
					GroupName:        group.Name,
					GroupPlatform:    group.Platform,
					GroupRate:        group.RateMultiplier,
					GroupIsExclusive: group.IsExclusive,
					SubscriptionType: group.SubscriptionType,
					ModelName:        model.Name,
					Platform:         model.Platform,
					BillingType:      billingType,
					Pricing:          model.Pricing,
				})
			}
		}
	}

	sort.SliceStable(out, func(i, j int) bool {
		if out[i].GroupName != out[j].GroupName {
			return strings.ToLower(out[i].GroupName) < strings.ToLower(out[j].GroupName)
		}
		if out[i].Platform != out[j].Platform {
			return out[i].Platform < out[j].Platform
		}
		return strings.ToLower(out[i].ModelName) < strings.ToLower(out[j].ModelName)
	})

	return out, nil
}

// ListPublicMarketplaceGroups returns the groups visible to anonymous marketplace visitors.
func (s *GroupService) ListPublicMarketplaceGroups(ctx context.Context) ([]Group, error) {
	groups, err := s.ListActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("list public marketplace groups: %w", err)
	}
	out := make([]Group, 0, len(groups))
	for _, group := range groups {
		if group.IsExclusive {
			continue
		}
		out = append(out, group)
	}
	return out, nil
}
