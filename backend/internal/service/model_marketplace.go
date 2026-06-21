package service

import (
	"context"
	"sort"
	"strings"
)

const (
	ModelMarketplacePricingSourceChannel = "channel"
	ModelMarketplacePricingSourceGroup   = "group"
	marketplaceDefaultImageBasePrice     = 0.134
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
	Supplier         string
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
	Supplier         string
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

type marketplaceAccountLister interface {
	ListSchedulableByGroupID(ctx context.Context, groupID int64) ([]Account, error)
}

// ModelMarketplaceService aggregates the public model marketplace cards.
type ModelMarketplaceService struct {
	accountLister  marketplaceAccountLister
	channelService *ChannelService
	pricingService *PricingService
}

func NewModelMarketplaceService(
	accountLister marketplaceAccountLister,
	channelService *ChannelService,
	pricingService *PricingService,
) *ModelMarketplaceService {
	return &ModelMarketplaceService{
		accountLister:  accountLister,
		channelService: channelService,
		pricingService: pricingService,
	}
}

// ListModelMarketplace returns one aggregated card per supplier + model.
//
// Each card keeps the cheapest offer as its representative row and includes
// every matched group in Groups for the detail drawer.
func (s *ModelMarketplaceService) ListModelMarketplace(
	ctx context.Context,
	visibleGroups []Group,
) ([]ModelMarketplaceCard, error) {
	if s == nil || s.channelService == nil || s.pricingService == nil {
		return []ModelMarketplaceCard{}, nil
	}

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

		modelNames, err := s.listMarketplaceModelNamesForGroup(ctx, group)
		if err != nil {
			return nil, err
		}
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

func (s *ModelMarketplaceService) listMarketplaceModelNamesForGroup(
	ctx context.Context,
	group Group,
) ([]string, error) {
	platform := strings.ToLower(strings.TrimSpace(group.Platform))
	if platform == "" {
		return nil, nil
	}

	if s.accountLister == nil {
		return nil, nil
	}

	accounts, err := s.accountLister.ListSchedulableByGroupID(ctx, group.ID)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0)
	for _, account := range accounts {
		if !strings.EqualFold(strings.TrimSpace(account.Platform), platform) {
			continue
		}
		for modelName := range explicitModelMapping(account.Credentials) {
			expanded := s.expandMarketplaceModelPattern(modelName)
			if len(expanded) == 0 {
				continue
			}
			names = append(names, expanded...)
		}
	}
	return dedupeAndSortMarketplaceNames(names), nil
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

func explicitModelMapping(credentials map[string]any) map[string]string {
	if credentials == nil {
		return nil
	}
	return stringMappingFromRaw(credentials["model_mapping"])
}

func (s *ModelMarketplaceService) expandMarketplaceModelPattern(pattern string) []string {
	trimmed := strings.TrimSpace(pattern)
	if trimmed == "" {
		return nil
	}
	if !strings.HasSuffix(trimmed, "*") {
		return []string{trimmed}
	}

	if s == nil || s.pricingService == nil {
		return nil
	}

	allNames := s.pricingService.ListModelNamesByPrefix(trimmed[:len(trimmed)-1])
	return dedupeAndSortMarketplaceNames(allNames)
}

func (s *ModelMarketplaceService) buildMarketplaceGroupOffer(
	ctx context.Context,
	group Group,
	modelName string,
) ModelMarketplaceGroupOffer {
	channelPricing := s.channelService.GetChannelModelPricing(ctx, group.ID, modelName)

	if s.marketplaceUsesImagePricing(modelName, channelPricing) {
		pricingSource, originalPricing, currentPricing := s.marketplaceImagePricing(group, modelName, channelPricing)
		return ModelMarketplaceGroupOffer{
			GroupID:          group.ID,
			GroupName:        group.Name,
			GroupPlatform:    group.Platform,
			GroupRate:        group.RateMultiplier,
			GroupIsExclusive: group.IsExclusive,
			SubscriptionType: group.SubscriptionType,
			ModelName:        modelName,
			Supplier:         s.marketplaceSupplierForModel(modelName),
			BillingType:      string(BillingModeImage),
			PricingSource:    pricingSource,
			OriginalPricing:  originalPricing,
			CurrentPricing:   currentPricing,
		}
	}

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
		Supplier:         s.marketplaceSupplierForModel(modelName),
		BillingType:      billingType,
		PricingSource:    pricingSource,
		OriginalPricing:  originalPricing,
		CurrentPricing:   currentPricing,
	}
}

func (s *ModelMarketplaceService) marketplaceUsesImagePricing(modelName string, channelPricing *ChannelModelPricing) bool {
	if channelPricing != nil && channelPricing.BillingMode == BillingModeImage {
		return true
	}
	if s == nil || s.pricingService == nil {
		return false
	}
	pricing := s.pricingService.GetModelPricing(modelName)
	return pricing != nil && pricing.Mode == "image_generation"
}

func (s *ModelMarketplaceService) marketplaceImagePricing(
	group Group,
	modelName string,
	channelPricing *ChannelModelPricing,
) (string, *ChannelModelPricing, *ChannelModelPricing) {
	if marketplaceHasExplicitChannelPricing(channelPricing) {
		original := sanitizeMarketplaceImagePricing(normalizeMarketplacePricing(channelPricing))
		return ModelMarketplacePricingSourceChannel, original, cloneMarketplacePricing(original)
	}

	original := s.marketplaceGroupImagePricing(group, modelName)
	current := scaleMarketplacePricing(original, group.RateMultiplier)
	return ModelMarketplacePricingSourceGroup, original, current
}

func sanitizeMarketplaceImagePricing(pricing *ChannelModelPricing) *ChannelModelPricing {
	if pricing == nil {
		return nil
	}

	cp := pricing.Clone()
	cp.BillingMode = BillingModeImage
	cp.InputPrice = nil
	cp.OutputPrice = nil
	cp.CacheWritePrice = nil
	cp.CacheReadPrice = nil
	return &cp
}

func (s *ModelMarketplaceService) marketplaceGroupImagePricing(group Group, modelName string) *ChannelModelPricing {
	pricing := &ChannelModelPricing{
		BillingMode: BillingModeImage,
		Intervals:   make([]PricingInterval, 0, 3),
	}

	for idx, tier := range []string{"1K", "2K", "4K"} {
		unitPrice := s.marketplaceImageUnitPrice(group, modelName, tier)
		if unitPrice <= 0 {
			continue
		}
		price := unitPrice
		pricing.Intervals = append(pricing.Intervals, PricingInterval{
			TierLabel:       tier,
			PerRequestPrice: &price,
			SortOrder:       idx,
		})
	}

	return normalizeMarketplacePricing(pricing)
}

func (s *ModelMarketplaceService) marketplaceImageUnitPrice(group Group, modelName string, tier string) float64 {
	tier = NormalizeImageBillingTierOrDefault(tier)
	if configured := group.GetImagePrice(tier); configured != nil {
		return *configured
	}
	return s.marketplaceDefaultImagePrice(modelName, tier)
}

func (s *ModelMarketplaceService) marketplaceDefaultImagePrice(modelName string, tier string) float64 {
	basePrice := 0.0
	if s != nil && s.pricingService != nil {
		if pricing := s.pricingService.GetModelPricing(modelName); pricing != nil && pricing.OutputCostPerImage > 0 {
			basePrice = pricing.OutputCostPerImage
		}
	}
	if basePrice <= 0 {
		basePrice = marketplaceDefaultImageBasePrice
	}

	switch NormalizeImageBillingTierOrDefault(tier) {
	case "2K":
		return basePrice * 1.5
	case "4K":
		return basePrice * 2
	default:
		return basePrice
	}
}

func marketplaceHasExplicitChannelPricing(pricing *ChannelModelPricing) bool {
	return pricing != nil && !pricingNeedsFallback(pricing)
}

func (s *ModelMarketplaceService) marketplaceBasePricing(modelName string, existing *ChannelModelPricing) *ChannelModelPricing {
	var lp *LiteLLMModelPricing
	if s.pricingService != nil {
		lp = s.pricingService.GetModelPricing(modelName)
	}
	return normalizeMarketplacePricing(synthesizePricingFromLiteLLM(lp, existing))
}

func (s *ModelMarketplaceService) marketplaceSupplierForModel(modelName string) string {
	if supplier := inferMarketplaceSupplier(modelName); supplier != "" {
		return supplier
	}
	if s != nil && s.pricingService != nil {
		if pricing := s.pricingService.GetModelPricing(modelName); pricing != nil {
			if supplier := normalizeMarketplaceSupplier(pricing.LiteLLMProvider); supplier != "" {
				return supplier
			}
		}
	}
	return "unknown"
}

func inferMarketplaceSupplier(modelName string) string {
	name := strings.ToLower(strings.TrimSpace(modelName))
	switch {
	case name == "":
		return ""
	case strings.HasPrefix(name, "gpt-"),
		strings.HasPrefix(name, "o1"),
		strings.HasPrefix(name, "o3"),
		strings.HasPrefix(name, "o4"):
		return "openai"
	case strings.HasPrefix(name, "claude-"),
		strings.HasPrefix(name, "opus-"),
		strings.HasPrefix(name, "sonnet-"),
		strings.HasPrefix(name, "haiku-"):
		return "anthropic"
	case strings.HasPrefix(name, "gemini-"):
		return "google"
	case strings.HasPrefix(name, "deepseek-"):
		return "deepseek"
	case strings.HasPrefix(name, "kimi-"):
		return "kimi"
	case strings.HasPrefix(name, "moonshot-"):
		return "moonshot"
	case strings.HasPrefix(name, "glm-"):
		return "glm"
	case strings.HasPrefix(name, "qwen-"),
		strings.HasPrefix(name, "qwen2-"),
		strings.HasPrefix(name, "qwen3-"),
		strings.HasPrefix(name, "qwen4-"):
		return "qwen"
	case strings.HasPrefix(name, "minimax-"):
		return "minimax"
	case strings.HasPrefix(name, "doubao-"):
		return "doubao"
	default:
		return ""
	}
}

func normalizeMarketplaceSupplier(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	switch value {
	case "anthropic", "openai", "google", "deepseek", "kimi", "moonshot", "glm", "qwen", "minimax", "doubao":
		return value
	default:
		return value
	}
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
		supplier string
		model    string
	}

	grouped := make(map[cardKey][]ModelMarketplaceGroupOffer)
	for _, row := range rows {
		key := cardKey{
			supplier: strings.ToLower(strings.TrimSpace(row.Supplier)),
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
			Supplier:         rep.Supplier,
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
		if cards[i].Supplier != cards[j].Supplier {
			return cards[i].Supplier < cards[j].Supplier
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
