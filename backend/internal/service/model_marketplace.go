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
	GroupID          int64
	GroupName        string
	GroupPlatform    string
	GroupRate        float64
	GroupIsExclusive bool
	SubscriptionType string
	ModelName        string
	Platform         string
	BillingType      string
	Pricing          *ChannelModelPricing
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
//
// The marketplace is group-driven: it enumerates the visible groups, groups them by
// platform, then expands all catalog models for each platform. Pricing is sourced from
// the global pricing catalog, not channel pricing rows.
func (s *ChannelService) ListModelMarketplace(
	ctx context.Context,
	visibleGroups []Group,
) ([]ModelMarketplaceCard, error) {
	if s.pricingService == nil {
		return []ModelMarketplaceCard{}, nil
	}

	groupsByPlatform := make(map[string][]Group)
	platforms := make([]string, 0)
	seenPlatforms := make(map[string]struct{})
	for _, group := range visibleGroups {
		if group.Status != "" && group.Status != StatusActive {
			continue
		}
		platform := strings.ToLower(strings.TrimSpace(group.Platform))
		if platform == "" || strings.TrimSpace(group.Name) == "" {
			continue
		}
		groupsByPlatform[platform] = append(groupsByPlatform[platform], group)
		if _, ok := seenPlatforms[platform]; ok {
			continue
		}
		seenPlatforms[platform] = struct{}{}
		platforms = append(platforms, platform)
	}
	if len(groupsByPlatform) == 0 {
		return []ModelMarketplaceCard{}, nil
	}

	sort.Strings(platforms)

	modelNamesByPlatform := make(map[string][]string, len(platforms))
	pricingByModelName := make(map[string]*ChannelModelPricing)
	for _, platform := range platforms {
		provider := marketplacePricingProviderForPlatform(platform)
		modelNames := s.pricingService.ListModelNamesByProvider(provider)
		if len(modelNames) == 0 {
			continue
		}
		modelNamesByPlatform[platform] = modelNames
		for _, modelName := range modelNames {
			modelKey := strings.ToLower(strings.TrimSpace(modelName))
			if modelKey == "" {
				continue
			}
			if _, exists := pricingByModelName[modelKey]; exists {
				continue
			}
			lp := s.pricingService.GetModelPricing(modelName)
			if lp == nil {
				continue
			}
			pricingByModelName[modelKey] = synthesizePricingFromLiteLLM(lp, nil)
		}
	}

	type dedupKey struct {
		groupID   int64
		modelName string
	}
	seen := make(map[dedupKey]struct{})
	out := make([]ModelMarketplaceCard, 0)

	for _, platform := range platforms {
		modelNames := modelNamesByPlatform[platform]
		if len(modelNames) == 0 {
			continue
		}
		for _, group := range groupsByPlatform[platform] {
			for _, modelName := range modelNames {
				modelKey := strings.ToLower(strings.TrimSpace(modelName))
				if modelKey == "" {
					continue
				}
				key := dedupKey{groupID: group.ID, modelName: modelKey}
				if _, exists := seen[key]; exists {
					continue
				}
				seen[key] = struct{}{}

				pricing := pricingByModelName[modelKey]
				billingType := ""
				if pricing != nil {
					billingType = string(pricing.BillingMode)
					if billingType == "" {
						billingType = string(BillingModeToken)
					}
				}

				out = append(out, ModelMarketplaceCard{
					GroupID:          group.ID,
					GroupName:        group.Name,
					GroupPlatform:    group.Platform,
					GroupRate:        group.RateMultiplier,
					GroupIsExclusive: group.IsExclusive,
					SubscriptionType: group.SubscriptionType,
					ModelName:        modelName,
					Platform:         group.Platform,
					BillingType:      billingType,
					Pricing:          pricing,
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

func marketplacePricingProviderForPlatform(platform string) string {
	switch strings.ToLower(strings.TrimSpace(platform)) {
	case PlatformAnthropic:
		return "anthropic"
	case PlatformOpenAI:
		return "openai"
	case PlatformGemini:
		return "google"
	case PlatformAntigravity:
		return "anthropic"
	default:
		return strings.ToLower(strings.TrimSpace(platform))
	}
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
