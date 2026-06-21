package service

import (
	"context"
	"log/slog"
	"sort"
	"strings"
)

const (
	ModelMarketplacePricingSourceChannel = "channel"
	ModelMarketplacePricingSourceGroup   = "group"
)

// ModelMarketplaceGroupOffer is one concrete group offer for a model.
// It carries both the original/base price and the displayed current price.
type ModelMarketplaceGroupOffer struct {
	GroupID          int64
	GroupName        string
	GroupPlatform    string
	GroupRate        float64
	GroupIsExclusive bool
	SubscriptionType string
	ModelName        string
	Platform         string
	BillingType      string
	PricingSource    string
	OriginalPricing  *ChannelModelPricing
	CurrentPricing   *ChannelModelPricing
}

// ModelMarketplaceCard is the aggregated marketplace card for one
// platform + model pair.
//
// The top-level group fields always point at the representative offer,
// which is selected by the lowest displayed current price.
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
	PricingSource    string
	OriginalPricing  *ChannelModelPricing
	CurrentPricing   *ChannelModelPricing
	Groups           []ModelMarketplaceGroupOffer
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

// ListModelMarketplace returns one aggregated card per platform + model.
//
// Each card keeps the cheapest offer as its representative row and includes
// every matched group in Groups for the detail drawer.
func (s *ChannelService) ListModelMarketplace(
	ctx context.Context,
	visibleGroups []Group,
) ([]ModelMarketplaceCard, error) {
	if s.pricingService == nil {
		return []ModelMarketplaceCard{}, nil
	}

	globalModelNamesByPlatform := make(map[string][]string)

	type rowKey struct {
		groupID   int64
		modelName string
	}
	seenRows := make(map[rowKey]struct{})
	rows := make([]ModelMarketplaceGroupOffer, 0)

	for _, group := range visibleGroups {
		if group.Status != "" && group.Status != StatusActive {
			continue
		}
		platform := strings.ToLower(strings.TrimSpace(group.Platform))
		if platform == "" || strings.TrimSpace(group.Name) == "" {
			continue
		}

		modelNames := s.listMarketplaceModelNamesForGroup(ctx, group, globalModelNamesByPlatform)
		for _, modelName := range modelNames {
			modelKey := strings.ToLower(strings.TrimSpace(modelName))
			if modelKey == "" {
				continue
			}
			key := rowKey{groupID: group.ID, modelName: modelKey}
			if _, exists := seenRows[key]; exists {
				continue
			}
			seenRows[key] = struct{}{}
			rows = append(rows, s.buildMarketplaceGroupOffer(ctx, group, modelName))
		}
	}

	if len(rows) == 0 {
		return []ModelMarketplaceCard{}, nil
	}

	return aggregateMarketplaceRows(rows), nil
}

func (s *ChannelService) listMarketplaceModelNamesForGroup(
	ctx context.Context,
	group Group,
	globalModelNamesByPlatform map[string][]string,
) []string {
	platform := strings.ToLower(strings.TrimSpace(group.Platform))
	if platform == "" {
		return nil
	}

	globalNames, ok := globalModelNamesByPlatform[platform]
	if !ok {
		provider := marketplacePricingProviderForPlatform(platform)
		globalNames = s.pricingService.ListModelNamesByProvider(provider)
		globalModelNamesByPlatform[platform] = globalNames
	}

	ch, err := s.GetChannelForGroup(ctx, group.ID)
	if err != nil {
		slog.Warn("list marketplace group models: failed to load group channel",
			"group_id", group.ID,
			"group_name", group.Name,
			"error", err)
	}

	channelNames := make([]string, 0)
	if ch != nil {
		supported := ch.SupportedModels()
		s.fillGlobalPricingFallback(supported)
		for _, model := range supported {
			if !strings.EqualFold(strings.TrimSpace(model.Platform), platform) {
				continue
			}
			channelNames = append(channelNames, model.Name)
		}
	}

	if ch != nil && ch.RestrictModels && len(channelNames) > 0 {
		return dedupeAndSortMarketplaceNames(channelNames)
	}

	combined := make([]string, 0, len(globalNames)+len(channelNames))
	combined = append(combined, globalNames...)
	combined = append(combined, channelNames...)
	return dedupeAndSortMarketplaceNames(combined)
}

func dedupeAndSortMarketplaceNames(names []string) []string {
	seen := make(map[string]string)
	for _, name := range names {
		trimmed := strings.TrimSpace(name)
		if trimmed == "" {
			continue
		}
		lower := strings.ToLower(trimmed)
		if _, exists := seen[lower]; exists {
			continue
		}
		seen[lower] = trimmed
	}

	out := make([]string, 0, len(seen))
	for _, original := range seen {
		out = append(out, original)
	}
	sort.SliceStable(out, func(i, j int) bool {
		return strings.ToLower(out[i]) < strings.ToLower(out[j])
	})
	return out
}

func (s *ChannelService) buildMarketplaceGroupOffer(
	ctx context.Context,
	group Group,
	modelName string,
) ModelMarketplaceGroupOffer {
	channelPricing := s.GetChannelModelPricing(ctx, group.ID, modelName)

	pricingSource := ModelMarketplacePricingSourceGroup
	var originalPricing *ChannelModelPricing
	var currentPricing *ChannelModelPricing

	switch {
	case marketplaceHasExplicitChannelPricing(channelPricing):
		pricingSource = ModelMarketplacePricingSourceChannel
		originalPricing = normalizeMarketplacePricing(channelPricing)
		currentPricing = cloneMarketplacePricing(originalPricing)
	default:
		originalPricing = s.marketplaceBasePricing(modelName, channelPricing)
		currentPricing = scaleMarketplacePricing(originalPricing, group.RateMultiplier)
	}

	billingType := marketplaceBillingType(channelPricing, originalPricing, currentPricing)

	return ModelMarketplaceGroupOffer{
		GroupID:          group.ID,
		GroupName:        group.Name,
		GroupPlatform:    group.Platform,
		GroupRate:        group.RateMultiplier,
		GroupIsExclusive: group.IsExclusive,
		SubscriptionType: group.SubscriptionType,
		ModelName:        modelName,
		Platform:         group.Platform,
		BillingType:      billingType,
		PricingSource:    pricingSource,
		OriginalPricing:  originalPricing,
		CurrentPricing:   currentPricing,
	}
}

func marketplaceHasExplicitChannelPricing(pricing *ChannelModelPricing) bool {
	return pricing != nil && !pricingNeedsFallback(pricing)
}

func (s *ChannelService) marketplaceBasePricing(modelName string, existing *ChannelModelPricing) *ChannelModelPricing {
	var lp *LiteLLMModelPricing
	if s.pricingService != nil {
		lp = s.pricingService.GetModelPricing(modelName)
	}
	return normalizeMarketplacePricing(synthesizePricingFromLiteLLM(lp, existing))
}

func normalizeMarketplacePricing(pricing *ChannelModelPricing) *ChannelModelPricing {
	if pricing == nil {
		return nil
	}
	cp := pricing.Clone()
	if cp.BillingMode == "" {
		cp.BillingMode = BillingModeToken
	}
	return &cp
}

func cloneMarketplacePricing(pricing *ChannelModelPricing) *ChannelModelPricing {
	if pricing == nil {
		return nil
	}
	cp := pricing.Clone()
	return &cp
}

func scaleMarketplacePricing(pricing *ChannelModelPricing, factor float64) *ChannelModelPricing {
	if pricing == nil {
		return nil
	}
	if factor < 0 {
		factor = 0
	}

	cp := pricing.Clone()
	if cp.BillingMode == "" {
		cp.BillingMode = BillingModeToken
	}
	cp.InputPrice = scaleMarketplacePricePtr(cp.InputPrice, factor)
	cp.OutputPrice = scaleMarketplacePricePtr(cp.OutputPrice, factor)
	cp.CacheWritePrice = scaleMarketplacePricePtr(cp.CacheWritePrice, factor)
	cp.CacheReadPrice = scaleMarketplacePricePtr(cp.CacheReadPrice, factor)
	cp.ImageOutputPrice = scaleMarketplacePricePtr(cp.ImageOutputPrice, factor)
	cp.PerRequestPrice = scaleMarketplacePricePtr(cp.PerRequestPrice, factor)
	for i := range cp.Intervals {
		cp.Intervals[i].InputPrice = scaleMarketplacePricePtr(cp.Intervals[i].InputPrice, factor)
		cp.Intervals[i].OutputPrice = scaleMarketplacePricePtr(cp.Intervals[i].OutputPrice, factor)
		cp.Intervals[i].CacheWritePrice = scaleMarketplacePricePtr(cp.Intervals[i].CacheWritePrice, factor)
		cp.Intervals[i].CacheReadPrice = scaleMarketplacePricePtr(cp.Intervals[i].CacheReadPrice, factor)
		cp.Intervals[i].PerRequestPrice = scaleMarketplacePricePtr(cp.Intervals[i].PerRequestPrice, factor)
	}
	return &cp
}

func scaleMarketplacePricePtr(value *float64, factor float64) *float64 {
	if value == nil {
		return nil
	}
	scaled := *value * factor
	return &scaled
}

func marketplaceBillingType(pricings ...*ChannelModelPricing) string {
	for _, pricing := range pricings {
		if pricing == nil {
			continue
		}
		mode := pricing.BillingMode
		if mode == "" {
			mode = BillingModeToken
		}
		return string(mode)
	}
	return string(BillingModeToken)
}

func aggregateMarketplaceRows(rows []ModelMarketplaceGroupOffer) []ModelMarketplaceCard {
	type cardKey struct {
		platform string
		model    string
	}

	grouped := make(map[cardKey][]ModelMarketplaceGroupOffer)
	for _, row := range rows {
		key := cardKey{
			platform: strings.ToLower(strings.TrimSpace(row.Platform)),
			model:    strings.ToLower(strings.TrimSpace(row.ModelName)),
		}
		grouped[key] = append(grouped[key], row)
	}

	cards := make([]ModelMarketplaceCard, 0, len(grouped))
	for _, offers := range grouped {
		sortMarketplaceOffers(offers)
		rep := pickMarketplaceRepresentative(offers)
		cards = append(cards, ModelMarketplaceCard{
			GroupID:          rep.GroupID,
			GroupName:        rep.GroupName,
			GroupPlatform:    rep.GroupPlatform,
			GroupRate:        rep.GroupRate,
			GroupIsExclusive: rep.GroupIsExclusive,
			SubscriptionType: rep.SubscriptionType,
			ModelName:        rep.ModelName,
			Platform:         rep.Platform,
			BillingType:      rep.BillingType,
			PricingSource:    rep.PricingSource,
			OriginalPricing:  cloneMarketplacePricing(rep.OriginalPricing),
			CurrentPricing:   cloneMarketplacePricing(rep.CurrentPricing),
			Groups:           cloneMarketplaceOffers(offers),
		})
	}

	sortMarketplaceCards(cards)
	return cards
}

func cloneMarketplaceOffers(offers []ModelMarketplaceGroupOffer) []ModelMarketplaceGroupOffer {
	out := make([]ModelMarketplaceGroupOffer, 0, len(offers))
	for _, offer := range offers {
		cp := offer
		cp.OriginalPricing = cloneMarketplacePricing(offer.OriginalPricing)
		cp.CurrentPricing = cloneMarketplacePricing(offer.CurrentPricing)
		out = append(out, cp)
	}
	return out
}

func sortMarketplaceOffers(offers []ModelMarketplaceGroupOffer) {
	sort.SliceStable(offers, func(i, j int) bool {
		leftPrice, leftOK := marketplaceDisplayPrice(offers[i].CurrentPricing, offers[i].BillingType)
		rightPrice, rightOK := marketplaceDisplayPrice(offers[j].CurrentPricing, offers[j].BillingType)

		if leftOK != rightOK {
			return leftOK
		}
		if leftOK && rightOK && leftPrice != rightPrice {
			return leftPrice < rightPrice
		}
		if offers[i].GroupName != offers[j].GroupName {
			return strings.ToLower(offers[i].GroupName) < strings.ToLower(offers[j].GroupName)
		}
		return strings.ToLower(offers[i].ModelName) < strings.ToLower(offers[j].ModelName)
	})
}

func pickMarketplaceRepresentative(offers []ModelMarketplaceGroupOffer) ModelMarketplaceGroupOffer {
	best := offers[0]
	for i := 1; i < len(offers); i++ {
		if marketplaceOfferIsBetter(offers[i], best) {
			best = offers[i]
		}
	}
	return best
}

func marketplaceOfferIsBetter(candidate, current ModelMarketplaceGroupOffer) bool {
	candidatePrice, candidateOK := marketplaceDisplayPrice(candidate.CurrentPricing, candidate.BillingType)
	currentPrice, currentOK := marketplaceDisplayPrice(current.CurrentPricing, current.BillingType)

	switch {
	case candidateOK && !currentOK:
		return true
	case !candidateOK && currentOK:
		return false
	case candidateOK && currentOK && candidatePrice != currentPrice:
		return candidatePrice < currentPrice
	}

	if candidate.GroupName != current.GroupName {
		return strings.ToLower(candidate.GroupName) < strings.ToLower(current.GroupName)
	}
	return strings.ToLower(candidate.ModelName) < strings.ToLower(current.ModelName)
}

func sortMarketplaceCards(cards []ModelMarketplaceCard) {
	sort.SliceStable(cards, func(i, j int) bool {
		leftPrice, leftOK := marketplaceDisplayPrice(cards[i].CurrentPricing, cards[i].BillingType)
		rightPrice, rightOK := marketplaceDisplayPrice(cards[j].CurrentPricing, cards[j].BillingType)

		if leftOK != rightOK {
			return leftOK
		}
		if leftOK && rightOK && leftPrice != rightPrice {
			return leftPrice < rightPrice
		}
		if cards[i].Platform != cards[j].Platform {
			return cards[i].Platform < cards[j].Platform
		}
		return strings.ToLower(cards[i].ModelName) < strings.ToLower(cards[j].ModelName)
	})
}

func marketplaceDisplayPrice(pricing *ChannelModelPricing, billingType string) (float64, bool) {
	if pricing == nil {
		return 0, false
	}

	switch billingType {
	case string(BillingModePerRequest), string(BillingModeImage):
		if pricing.PerRequestPrice != nil {
			return *pricing.PerRequestPrice, true
		}
		if price, ok := minMarketplaceIntervalPrice(pricing.Intervals, func(iv PricingInterval) *float64 {
			return iv.PerRequestPrice
		}); ok {
			return price, true
		}
	default:
		if pricing.InputPrice != nil {
			return *pricing.InputPrice, true
		}
		if price, ok := minMarketplaceIntervalPrice(pricing.Intervals, func(iv PricingInterval) *float64 {
			return iv.InputPrice
		}); ok {
			return price, true
		}
	}

	for _, ptr := range []*float64{
		pricing.OutputPrice,
		pricing.CacheWritePrice,
		pricing.CacheReadPrice,
		pricing.ImageOutputPrice,
	} {
		if ptr != nil {
			return *ptr, true
		}
	}

	for _, getter := range []func(PricingInterval) *float64{
		func(iv PricingInterval) *float64 { return iv.OutputPrice },
		func(iv PricingInterval) *float64 { return iv.CacheWritePrice },
		func(iv PricingInterval) *float64 { return iv.CacheReadPrice },
		func(iv PricingInterval) *float64 { return iv.PerRequestPrice },
	} {
		if price, ok := minMarketplaceIntervalPrice(pricing.Intervals, getter); ok {
			return price, true
		}
	}

	return 0, false
}

func minMarketplaceIntervalPrice(
	intervals []PricingInterval,
	getter func(PricingInterval) *float64,
) (float64, bool) {
	var best float64
	found := false
	for _, iv := range intervals {
		value := getter(iv)
		if value == nil {
			continue
		}
		if !found || *value < best {
			best = *value
			found = true
		}
	}
	return best, found
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
		return nil, err
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
