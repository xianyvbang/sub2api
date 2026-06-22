<template>
  <div
    class="model-marketplace-page relative min-h-screen overflow-hidden bg-gradient-to-br from-amber-50 via-white to-indigo-50 text-[13px] text-gray-900 dark:from-dark-950 dark:via-dark-950 dark:to-dark-900 dark:text-gray-100"
  >
    <div
      class="pointer-events-none absolute inset-0 bg-[linear-gradient(rgba(15,23,42,0.04)_1px,transparent_1px),linear-gradient(90deg,rgba(15,23,42,0.04)_1px,transparent_1px)] bg-[size:30px_30px] dark:bg-[linear-gradient(rgba(148,163,184,0.08)_1px,transparent_1px),linear-gradient(90deg,rgba(148,163,184,0.08)_1px,transparent_1px)]"
    />
    <div
      class="pointer-events-none absolute inset-x-0 top-0 h-[34rem] bg-[radial-gradient(circle_at_top_left,_rgba(245,158,11,0.18),_transparent_30%),radial-gradient(circle_at_top_right,_rgba(59,130,246,0.16),_transparent_34%)] dark:bg-[radial-gradient(circle_at_top_left,_rgba(245,158,11,0.12),_transparent_32%),radial-gradient(circle_at_top_right,_rgba(96,165,250,0.14),_transparent_36%)]"
    />

    <header class="relative z-10 px-6 py-5">
      <nav class="flex w-full items-center justify-between gap-4">
        <router-link to="/home" class="flex items-center gap-3">
          <div
            class="h-11 w-11 overflow-hidden rounded-2xl bg-white shadow-lg shadow-amber-100 dark:bg-dark-800 dark:shadow-black/20"
          >
            <img :src="siteLogo || '/logo.png'" alt="Logo" class="h-full w-full object-contain" />
          </div>
          <div>
            <p class="text-[11px] font-semibold uppercase tracking-[0.24em] text-amber-700/80 dark:text-amber-300/80">
              {{ t('nav.modelMarketplace', 'Model Marketplace') }}
            </p>
            <h1 class="text-base font-semibold text-slate-900 dark:text-white">{{ siteName }}</h1>
          </div>
        </router-link>

        <div class="flex items-center gap-3">
          <router-link
            v-if="isAuthenticated"
            :to="dashboardPath"
            class="rounded-full border border-slate-200/80 bg-white/90 px-4 py-2 text-xs font-medium text-slate-700 shadow-sm backdrop-blur transition hover:border-slate-300 hover:text-slate-900 dark:border-dark-700 dark:bg-dark-800/90 dark:text-dark-200 dark:hover:border-dark-600 dark:hover:text-white"
          >
            {{ t('home.dashboard') }}
          </router-link>
          <router-link
            v-else
            to="/login"
            class="rounded-full bg-slate-900 px-4 py-2 text-xs font-medium text-white shadow-lg shadow-slate-300 transition hover:bg-slate-800 dark:bg-primary-500 dark:shadow-primary-500/20 dark:hover:bg-primary-600"
          >
            {{ t('home.login') }}
          </router-link>
        </div>
      </nav>
    </header>

    <main class="relative z-10 px-6 pb-16 pt-6">
      <div class="w-full">
        <section class="grid gap-6 xl:grid-cols-[320px_minmax(0,1fr)] xl:items-start">
          <aside
            class="h-fit rounded-[1.75rem] border border-white/70 bg-white/85 p-4 shadow-[0_20px_60px_-45px_rgba(15,23,42,0.55)] backdrop-blur dark:border-dark-700/70 dark:bg-dark-900/80 dark:shadow-[0_20px_60px_-45px_rgba(0,0,0,0.6)] xl:sticky xl:top-6"
          >
            <div class="mb-4 flex items-center justify-between gap-3 px-2">
              <div>
                <p class="text-[11px] font-semibold uppercase tracking-[0.22em] text-slate-400 dark:text-dark-400">
                  {{ t('common.filter', 'Filter') }}
                </p>
                <p class="mt-1 text-xs text-slate-600 dark:text-dark-300">{{ activeFilterSummary }}</p>
              </div>
              <button
                v-if="hasCategoryFilters"
                @click="resetCategoryFilters"
                class="text-[11px] font-semibold text-amber-700 transition hover:text-amber-600 dark:text-amber-300 dark:hover:text-amber-200"
              >
                {{ t('common.reset', 'Reset') }}
              </button>
            </div>

            <div class="space-y-3">
              <details open class="group overflow-hidden rounded-2xl border border-slate-200/80 bg-slate-50/75 dark:border-dark-700 dark:bg-dark-950/60">
                <summary
                  class="flex cursor-pointer list-none items-center justify-between gap-3 px-4 py-3 text-xs font-semibold text-slate-900 dark:text-white"
                >
                  <span>{{ t('modelMarketplace.groupLabel', '分组') }}</span>
                  <span class="rounded-full bg-white px-2 py-0.5 text-[11px] text-slate-500 dark:bg-dark-800 dark:text-dark-300">
                    {{ groupOptions.length }}
                  </span>
                </summary>
                <div class="space-y-1 border-t border-slate-200/70 p-2 dark:border-dark-700">
                  <button
                    @click="selectedGroup = ''"
                    class="flex w-full items-center justify-between rounded-xl px-3 py-2 text-left text-xs transition"
                    :class="
                      selectedGroup
                        ? 'text-slate-600 hover:bg-white hover:text-slate-900 dark:text-dark-300 dark:hover:bg-dark-800 dark:hover:text-white'
                        : 'bg-slate-900 text-white dark:bg-primary-500 dark:text-white'
                    "
                  >
                    <span>{{ t('common.all', 'All') }}</span>
                    <span class="text-[11px] opacity-75">{{ groupAllCount }}</span>
                  </button>
                  <button
                    v-for="group in groupEntries"
                    :key="group.value"
                    @click="toggleGroup(group.value)"
                    class="flex w-full items-center justify-between rounded-xl px-3 py-2 text-left text-xs transition"
                    :class="
                      selectedGroup === group.value
                        ? 'bg-slate-900 text-white dark:bg-primary-500 dark:text-white'
                        : 'text-slate-600 hover:bg-white hover:text-slate-900 dark:text-dark-300 dark:hover:bg-dark-800 dark:hover:text-white'
                    "
                  >
                    <span class="truncate pr-3">{{ group.value }}</span>
                    <span class="text-[11px] opacity-75">{{ group.count }}</span>
                  </button>
                </div>
              </details>

              <details open class="group overflow-hidden rounded-2xl border border-slate-200/80 bg-slate-50/75 dark:border-dark-700 dark:bg-dark-950/60">
                <summary
                  class="flex cursor-pointer list-none items-center justify-between gap-3 px-4 py-3 text-xs font-semibold text-slate-900 dark:text-white"
                >
                  <span>{{ t('modelMarketplace.supplierLabel', '供应商') }}</span>
                  <span class="rounded-full bg-white px-2 py-0.5 text-[11px] text-slate-500 dark:bg-dark-800 dark:text-dark-300">
                    {{ supplierOptions.length }}
                  </span>
                </summary>
                <div class="space-y-1 border-t border-slate-200/70 p-2 dark:border-dark-700">
                  <button
                    @click="selectedSupplier = ''"
                    class="flex w-full items-center justify-between rounded-xl px-3 py-2 text-left text-xs transition"
                    :class="
                      selectedSupplier
                        ? 'text-slate-600 hover:bg-white hover:text-slate-900 dark:text-dark-300 dark:hover:bg-dark-800 dark:hover:text-white'
                        : 'bg-slate-900 text-white dark:bg-primary-500 dark:text-white'
                    "
                  >
                    <span>{{ t('common.all', 'All') }}</span>
                    <span class="text-[11px] opacity-75">{{ supplierAllCount }}</span>
                  </button>
                  <button
                    v-for="supplier in supplierEntries"
                    :key="supplier.value"
                    @click="toggleSupplier(supplier.value)"
                    class="flex w-full items-center justify-between rounded-xl px-3 py-2 text-left text-xs transition"
                    :class="
                      selectedSupplier === supplier.value
                        ? 'bg-slate-900 text-white dark:bg-primary-500 dark:text-white'
                        : 'text-slate-600 hover:bg-white hover:text-slate-900 dark:text-dark-300 dark:hover:bg-dark-800 dark:hover:text-white'
                    "
                  >
                    <span class="truncate pr-3">{{ supplierLabel(supplier.value) }}</span>
                    <span class="text-[11px] opacity-75">{{ supplier.count }}</span>
                  </button>
                </div>
              </details>

              <details open class="group overflow-hidden rounded-2xl border border-slate-200/80 bg-slate-50/75 dark:border-dark-700 dark:bg-dark-950/60">
                <summary
                  class="flex cursor-pointer list-none items-center justify-between gap-3 px-4 py-3 text-xs font-semibold text-slate-900 dark:text-white"
                >
                  <span>{{ t('modelMarketplace.billingType', 'Billing Type') }}</span>
                  <span class="rounded-full bg-white px-2 py-0.5 text-[11px] text-slate-500 dark:bg-dark-800 dark:text-dark-300">
                    {{ billingTypeOptions.length }}
                  </span>
                </summary>
                <div class="space-y-1 border-t border-slate-200/70 p-2 dark:border-dark-700">
                  <button
                    @click="selectedBillingType = ''"
                    class="flex w-full items-center justify-between rounded-xl px-3 py-2 text-left text-xs transition"
                    :class="
                      selectedBillingType
                        ? 'text-slate-600 hover:bg-white hover:text-slate-900 dark:text-dark-300 dark:hover:bg-dark-800 dark:hover:text-white'
                        : 'bg-slate-900 text-white dark:bg-primary-500 dark:text-white'
                    "
                  >
                    <span>{{ t('common.all', 'All') }}</span>
                    <span class="text-[11px] opacity-75">{{ billingTypeAllCount }}</span>
                  </button>
                  <button
                    v-for="billingType in billingTypeEntries"
                    :key="billingType.value"
                    @click="toggleBillingType(billingType.value)"
                    class="flex w-full items-center justify-between rounded-xl px-3 py-2 text-left text-xs transition"
                    :class="
                      selectedBillingType === billingType.value
                        ? 'bg-slate-900 text-white dark:bg-primary-500 dark:text-white'
                        : 'text-slate-600 hover:bg-white hover:text-slate-900 dark:text-dark-300 dark:hover:bg-dark-800 dark:hover:text-white'
                    "
                  >
                    <span class="truncate pr-3">{{ billingTypeLabel(billingType.value) }}</span>
                    <span class="text-[11px] opacity-75">{{ billingType.count }}</span>
                  </button>
                </div>
              </details>
            </div>

            <div class="hidden">
              <label>
                {{ t('modelMarketplace.groupLabel', '分组') }}
                <select v-model="selectedGroup">
                  <option value="">{{ t('common.all', 'All') }}</option>
                  <option v-for="group in groupOptions" :key="group" :value="group">{{ group }}</option>
                </select>
              </label>
              <label>
                {{ t('modelMarketplace.supplierLabel', '供应商') }}
                <select v-model="selectedSupplier">
                  <option value="">{{ t('common.all', 'All') }}</option>
                  <option v-for="supplier in supplierOptions" :key="supplier" :value="supplier">{{ supplierLabel(supplier) }}</option>
                </select>
              </label>
              <label>
                Billing Type
                <select v-model="selectedBillingType">
                  <option value="">{{ t('common.all', 'All') }}</option>
                  <option v-for="billingType in billingTypeOptions" :key="billingType" :value="billingType">
                    {{ billingTypeLabel(billingType) }}
                  </option>
                </select>
              </label>
            </div>
          </aside>

          <div class="min-w-0 w-full">
            <section
              class="rounded-[2rem] border border-white/70 bg-white/80 p-6 shadow-[0_25px_80px_-40px_rgba(15,23,42,0.45)] backdrop-blur dark:border-dark-700/70 dark:bg-dark-900/75 dark:shadow-[0_25px_80px_-40px_rgba(0,0,0,0.55)]"
            >
              <div class="flex flex-col gap-6">
                <div class="max-w-3xl">
                  <p class="text-2xl font-black tracking-[0.08em] text-slate-950 dark:text-white md:text-3xl">
                    {{ t('nav.modelMarketplace', 'Model Marketplace') }}
                  </p>
                </div>

                <div class="flex flex-col gap-4 lg:flex-row lg:items-center">
                  <div class="min-w-0 flex-1">
                    <input
                      v-model="searchQuery"
                      type="text"
                      :placeholder="t('modelMarketplace.searchPlaceholder', '搜索模型名称')"
                      class="w-full rounded-2xl border border-slate-200 bg-white px-4 py-3 text-xs text-slate-900 placeholder:text-slate-400 outline-none transition focus:border-amber-400 dark:border-dark-600 dark:bg-dark-900 dark:text-white dark:placeholder:text-dark-400 dark:focus:border-amber-400"
                    />
                  </div>

                  <div class="flex flex-wrap items-center gap-3">
                    <button
                      v-if="hasCategoryFilters"
                      @click="resetCategoryFilters"
                      class="inline-flex items-center rounded-full border border-slate-200 bg-white px-4 py-2.5 text-xs font-semibold text-slate-700 transition hover:border-slate-300 hover:text-slate-900 dark:border-dark-600 dark:bg-dark-900 dark:text-dark-200 dark:hover:border-dark-500 dark:hover:text-white"
                    >
                      {{ t('common.reset', 'Reset') }}
                    </button>
                    <button
                      @click="loadMarketplace"
                      :disabled="loading"
                      class="inline-flex items-center rounded-full bg-amber-400 px-5 py-2.5 text-xs font-semibold text-slate-950 transition hover:bg-amber-300 disabled:cursor-not-allowed disabled:opacity-60"
                    >
                      {{ loading ? t('common.loading', 'Loading...') : t('common.refresh', 'Refresh') }}
                    </button>
                  </div>
                </div>
              </div>
            </section>

            <section v-if="loading" class="mt-8 grid gap-5 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4">
              <div
                v-for="index in 6"
                :key="index"
                class="h-72 animate-pulse rounded-[1.75rem] border border-white/70 bg-white/70 dark:border-dark-700/70 dark:bg-dark-800/70"
              />
            </section>

            <section
              v-else-if="filteredCards.length === 0"
              class="mt-8 rounded-[2rem] border border-dashed border-slate-300 bg-white/70 px-8 py-20 text-center shadow-sm dark:border-dark-700 dark:bg-dark-900/60"
            >
              <h3 class="text-xl font-bold text-slate-900 dark:text-white">
                {{ t('modelMarketplace.emptyTitle', 'No matching models') }}
              </h3>
              <p class="mt-3 text-xs text-slate-500 dark:text-dark-400">
                {{ t('modelMarketplace.emptyDescription', 'Adjust your search or filters and try again.') }}
              </p>
            </section>

            <section v-else class="mt-8 grid gap-5 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4">
              <article
                v-for="card in filteredCards"
                :key="cardKey(card)"
                :data-testid="`marketplace-card-${card.supplier}-${card.model_name}`"
                class="group cursor-pointer rounded-[1.75rem] border border-white/70 bg-white/88 p-6 shadow-[0_20px_70px_-45px_rgba(15,23,42,0.5)] backdrop-blur transition hover:-translate-y-1 hover:shadow-[0_30px_80px_-45px_rgba(15,23,42,0.55)] dark:border-dark-700/70 dark:bg-dark-900/82 dark:shadow-[0_20px_70px_-45px_rgba(0,0,0,0.6)] dark:hover:shadow-[0_30px_80px_-45px_rgba(0,0,0,0.7)]"
                @click="openCard(card)"
              >
                <div class="flex items-start justify-between gap-4">
                  <div class="min-w-0">
                    <div class="truncate text-[11px] font-semibold uppercase tracking-[0.2em] text-slate-400 dark:text-dark-400">
                      <span
                        class="inline-flex items-center gap-1.5 rounded-full px-2.5 py-1"
                        :class="supplierBadgeClass(card.supplier)"
                        :data-testid="`marketplace-supplier-badge-${card.supplier}-${card.model_name}`"
                      >
                        <span>{{ supplierLabel(card.supplier) }}</span>
                      </span>
                    </div>
                    <h3 class="mt-2 min-w-0 break-words text-base font-bold leading-snug text-slate-950 dark:text-white">
                      {{ card.model_name }}
                    </h3>
                  </div>
                  <span
                    class="shrink-0 rounded-full px-3 py-1 text-[11px] font-semibold"
                    :class="pricingSourceBadgeClass(cardDisplayOffer(card).pricing_source)"
                  >
                    {{ pricingSourceLabel(cardDisplayOffer(card).pricing_source) }}
                  </span>
                </div>

                <div class="mt-5 flex flex-wrap gap-2">
                  <span
                    class="inline-flex items-center gap-1.5 rounded-full bg-amber-100 px-3 py-1 text-[11px] font-medium text-amber-800 dark:bg-amber-500/15 dark:text-amber-300"
                  >
                    {{ cardDisplayOffer(card).group_name }}
                  </span>
                  <span
                    class="rounded-full bg-sky-100 px-3 py-1 text-[11px] font-medium text-sky-800 dark:bg-sky-500/15 dark:text-sky-300"
                  >
                    {{ billingTypeLabel(cardDisplayOffer(card).billing_type) }}
                  </span>
                  <span
                    class="rounded-full bg-violet-100 px-3 py-1 text-[11px] font-medium text-violet-800 dark:bg-violet-500/15 dark:text-violet-300"
                  >
                    {{ t('modelMarketplace.groupCount', { count: card.groups.length }) }}
                  </span>
                </div>

                <div class="mt-6 rounded-[1.5rem] border border-slate-200/80 bg-slate-50/90 p-4 dark:border-dark-700 dark:bg-dark-950/70">
                  <div
                    v-if="buildPricingComparisonLines(cardDisplayOffer(card).current_pricing, cardDisplayOffer(card).original_pricing, cardDisplayOffer(card).billing_type).length > 0"
                    class="space-y-3"
                  >
                    <div
                      v-for="line in buildPricingComparisonLines(cardDisplayOffer(card).current_pricing, cardDisplayOffer(card).original_pricing, cardDisplayOffer(card).billing_type)"
                      :key="`${cardKey(card)}-comparison-${line.label}`"
                      class="rounded-2xl bg-white/80 px-4 py-3 dark:bg-dark-900/70"
                    >
                      <div class="flex items-center justify-between gap-4 text-xs">
                        <span class="font-medium text-amber-700 dark:text-amber-300">{{ formatCardPriceLabel(line.label) }}:</span>
                        <span class="text-right font-semibold text-slate-950 dark:text-white">
                          {{ line.currentValue ?? t('modelMarketplace.priceUnavailable', '暂无价格') }}
                        </span>
                      </div>
                      <div class="mt-2 flex items-center justify-between gap-4 text-[11px]">
                        <span class="text-slate-500 dark:text-dark-400">{{ t('modelMarketplace.originalPrice', '原价') }}:</span>
                        <span class="text-right text-slate-500 line-through dark:text-dark-300">
                          {{ line.originalValue ?? t('availableChannels.noPricing', 'No pricing configured') }}
                        </span>
                      </div>
                    </div>
                  </div>
                  <p v-else class="text-xs text-slate-500 dark:text-dark-400">
                    {{ t('availableChannels.noPricing', 'No pricing configured') }}
                  </p>
                </div>

                <dl class="mt-6 space-y-3 text-xs">
                  <div class="flex items-center justify-between gap-4">
                    <dt class="text-slate-500 dark:text-dark-400">
                      {{
                        selectedGroup
                          ? t('modelMarketplace.displayPriceGroup', '当前展示分组')
                          : t('modelMarketplace.lowestPriceGroup', '最低展示价分组')
                      }}
                    </dt>
                    <dd class="font-medium text-slate-900 dark:text-white">{{ cardDisplayOffer(card).group_name }}</dd>
                  </div>
                  <div class="flex items-center justify-between gap-4">
                    <dt class="text-slate-500 dark:text-dark-400">{{ t('modelMarketplace.groupRate', 'Rate') }}</dt>
                    <dd class="font-medium text-slate-900 dark:text-white">{{ formatRate(cardDisplayOffer(card).group_rate) }}</dd>
                  </div>
                </dl>

                <div
                  class="mt-6 flex items-center justify-between rounded-2xl border border-slate-200/80 bg-white/90 px-4 py-3 text-xs font-medium text-slate-700 transition group-hover:border-slate-300 group-hover:text-slate-900 dark:border-dark-700 dark:bg-dark-900/90 dark:text-dark-200 dark:group-hover:border-dark-500 dark:group-hover:text-white"
                >
                  <span>{{ t('modelMarketplace.viewDetails', '查看分组明细') }}</span>
                  <Icon name="chevronRight" size="sm" />
                </div>
              </article>
            </section>
          </div>
        </section>
      </div>
    </main>

    <transition
      enter-active-class="transition duration-300 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition duration-200 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div v-if="activeCard" class="fixed inset-0 z-[120]" @keydown.esc="closeCard">
        <div class="absolute inset-0 bg-slate-950/50 backdrop-blur-sm" @click="closeCard" />
        <transition
          enter-active-class="transition duration-300 ease-out"
          enter-from-class="translate-x-full"
          enter-to-class="translate-x-0"
          leave-active-class="transition duration-200 ease-in"
          leave-from-class="translate-x-0"
          leave-to-class="translate-x-full"
        >
          <aside
            v-if="activeCard"
            data-testid="marketplace-detail-drawer"
            class="absolute inset-y-0 right-0 flex h-full w-full max-w-2xl flex-col border-l border-white/70 bg-white/96 shadow-[0_24px_80px_rgba(15,23,42,0.28)] backdrop-blur-xl dark:border-dark-700 dark:bg-dark-900/96 dark:shadow-[0_24px_80px_rgba(0,0,0,0.55)]"
            @click.stop
          >
            <div class="flex items-start justify-between gap-4 border-b border-slate-200/80 px-6 py-5 dark:border-dark-700">
              <div class="min-w-0">
                <p class="text-[11px] font-semibold uppercase tracking-[0.22em] text-slate-400 dark:text-dark-400">
                  {{ supplierLabel(activeCard.supplier) }}
                </p>
                <div class="mt-2 flex items-start gap-2">
                  <h3 class="break-words text-xl font-bold text-slate-950 dark:text-white">
                    {{ activeCard.model_name }}
                  </h3>
                  <button
                    type="button"
                    :data-testid="`copy-model-name-${activeCard.model_name}`"
                    class="mt-1 shrink-0 rounded-full border border-slate-200/80 bg-white/90 p-1.5 text-slate-500 shadow-sm transition hover:border-slate-300 hover:text-slate-900 dark:border-dark-700 dark:bg-dark-800/90 dark:text-dark-300 dark:hover:border-dark-600 dark:hover:text-white"
                    :title="copiedModelName === activeCard.model_name ? t('common.copied') : t('common.copy')"
                    :aria-label="copiedModelName === activeCard.model_name ? t('common.copied') : t('common.copy')"
                    @click.stop="copyModelName(activeCard.model_name)"
                  >
                    <Icon v-if="copiedModelName === activeCard.model_name" name="check" size="xs" :stroke-width="2" />
                    <Icon v-else name="copy" size="xs" :stroke-width="2" />
                  </button>
                </div>
                <div class="mt-3 flex flex-wrap gap-2">
                  <span
                    class="rounded-full bg-amber-100 px-3 py-1 text-[11px] font-medium text-amber-800 dark:bg-amber-500/15 dark:text-amber-300"
                  >
                    {{ t('modelMarketplace.groupCount', { count: activeCard.groups.length }) }}
                  </span>
                  <span
                    class="rounded-full px-3 py-1 text-[11px] font-medium"
                    :class="billingTypeBadgeClass(cardDisplayOffer(activeCard).billing_type)"
                  >
                    {{ billingTypeLabel(cardDisplayOffer(activeCard).billing_type) }}
                  </span>
                </div>
              </div>

              <button
                type="button"
                class="rounded-full border border-slate-200/80 bg-white/90 p-2 text-slate-500 shadow-sm transition hover:border-slate-300 hover:text-slate-900 dark:border-dark-700 dark:bg-dark-800/90 dark:text-dark-300 dark:hover:border-dark-600 dark:hover:text-white"
                :aria-label="t('common.close', 'Close')"
                @click="closeCard"
              >
                <Icon name="x" size="sm" />
              </button>
            </div>

            <div class="flex-1 overflow-y-auto px-6 py-5">
              <div class="space-y-4">
                <article
                  v-for="group in activeCard.groups"
                  :key="`${group.group_id}-${group.model_name}`"
                  class="rounded-[1.5rem] border border-slate-200/80 bg-white/92 p-5 shadow-sm dark:border-dark-700 dark:bg-dark-950/70"
                >
                  <div class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
                    <div>
                      <h4 class="text-base font-semibold text-slate-950 dark:text-white">{{ group.group_name }}</h4>
                      <div class="mt-3 flex flex-wrap gap-2">
                        <span
                          class="rounded-full px-3 py-1 text-[11px] font-medium"
                          :class="billingTypeBadgeClass(group.billing_type)"
                        >
                          {{ billingTypeLabel(group.billing_type) }}
                        </span>
                        <span
                          class="rounded-full px-3 py-1 text-[11px] font-medium"
                          :class="pricingSourceBadgeClass(group.pricing_source)"
                        >
                          {{ pricingSourceLabel(group.pricing_source) }}
                        </span>
                        <span
                          class="rounded-full bg-slate-100 px-3 py-1 text-[11px] font-medium text-slate-700 dark:bg-dark-800 dark:text-dark-200"
                        >
                          {{ group.subscription_type || '-' }}
                        </span>
                      </div>
                    </div>
                    <div class="text-xs font-semibold text-slate-950 dark:text-white">
                      {{ t('modelMarketplace.groupRate', 'Rate') }}: {{ formatRate(group.group_rate) }}
                    </div>
                  </div>

                  <div class="mt-5 rounded-2xl border border-amber-200/80 bg-amber-50/70 p-4 dark:border-amber-500/20 dark:bg-amber-500/10">
                    <div
                      v-if="buildPricingComparisonLines(group.current_pricing, group.original_pricing, group.billing_type).length > 0"
                      class="space-y-3"
                    >
                      <div
                        v-for="line in buildPricingComparisonLines(group.current_pricing, group.original_pricing, group.billing_type)"
                        :key="`${group.group_id}-comparison-${line.label}`"
                        class="rounded-2xl bg-white/85 px-4 py-3 dark:bg-dark-900/70"
                      >
                        <div class="flex items-center justify-between gap-4 text-xs">
                          <span class="font-medium text-amber-700 dark:text-amber-300">{{ formatCardPriceLabel(line.label) }}:</span>
                          <span class="text-right font-semibold text-slate-950 dark:text-white">
                            {{ line.currentValue ?? t('modelMarketplace.priceUnavailable', '暂无价格') }}
                          </span>
                        </div>
                        <div class="mt-2 flex items-center justify-between gap-4 text-[11px]">
                          <span class="text-slate-500 dark:text-dark-400">{{ t('modelMarketplace.originalPrice', '原价') }}:</span>
                          <span class="text-right text-slate-500 line-through dark:text-dark-300">
                            {{ line.originalValue ?? t('availableChannels.noPricing', 'No pricing configured') }}
                          </span>
                        </div>
                      </div>
                    </div>
                    <p v-else class="text-xs text-slate-500 dark:text-dark-400">
                      {{ t('availableChannels.noPricing', 'No pricing configured') }}
                    </p>
                  </div>
                </article>
              </div>
            </div>
          </aside>
        </transition>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  BILLING_MODE_IMAGE,
  BILLING_MODE_PER_REQUEST,
  BILLING_MODE_TOKEN,
} from '@/constants/channel'
import Icon from '@/components/icons/Icon.vue'
import { useClipboard } from '@/composables/useClipboard'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import { formatScaled } from '@/utils/pricing'
import modelMarketplaceAPI, {
  type ModelMarketplaceCard,
  type ModelMarketplaceGroupOffer,
  type ModelMarketplacePricing,
  type ModelMarketplacePricingInterval,
} from '@/api/modelMarketplace'
import { extractApiErrorMessage } from '@/utils/apiError'

type FilterEntry = {
  value: string
  count: number
}

type PriceSummaryLine = {
  label: string
  value: string
}

type PriceComparisonLine = {
  label: string
  currentValue: string | null
  originalValue: string | null
}

type PricingDisplayMode = 'default' | 'per_request' | 'image_output'

type PriceDescriptor = {
  field: 'input_price' | 'output_price' | 'cache_write_price' | 'cache_read_price' | 'image_output_price' | 'per_request_price'
  value: number
}

type OfferFilter = {
  groupName?: string
  billingType?: string
}

const { t } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()
const { copyToClipboard } = useClipboard()

const loading = ref(false)
const cards = ref<ModelMarketplaceCard[]>([])
const searchQuery = ref('')
const selectedGroup = ref('')
const selectedSupplier = ref('')
const selectedBillingType = ref('')
const copiedModelName = ref('')
const activeCardKey = ref('')

const siteName = computed(() => appStore.siteName || 'Sub2API')
const siteLogo = computed(() => appStore.siteLogo || '')
const isAuthenticated = computed(() => authStore.isAuthenticated)
const dashboardPath = computed(() => (authStore.isAdmin ? '/admin/dashboard' : '/dashboard'))

const groupOptions = computed(() =>
  Array.from(
    new Set(cards.value.flatMap((card) => card.groups.map((group) => group.group_name))),
  ).sort((a, b) => a.localeCompare(b)),
)
const supplierOptions = computed(() =>
  Array.from(new Set(cards.value.map((card) => card.supplier))).sort((a, b) => a.localeCompare(b)),
)
const billingTypeOptions = computed(() =>
  Array.from(
    new Set(cards.value.flatMap((card) => card.groups.map((group) => normalizeBillingType(group.billing_type))).filter(Boolean)),
  ).sort((a, b) => a.localeCompare(b)),
)

const hasCategoryFilters = computed(
  () => Boolean(selectedGroup.value || selectedSupplier.value || selectedBillingType.value),
)

const activeFilterSummary = computed(() => {
  const count = [selectedGroup.value, selectedSupplier.value, selectedBillingType.value].filter(Boolean).length
  if (count === 0) return t('common.all', 'All')
  return `${count} ${t('common.filter', 'Filter')}`
})

const normalizedSearchQuery = computed(() => searchQuery.value.trim().toLowerCase())

const searchedCards = computed(() => {
  const query = normalizedSearchQuery.value
  if (!query) return cards.value
  return cards.value.filter((card) => marketplaceSearchText(card).includes(query))
})

const filteredCards = computed(() =>
  searchedCards.value
    .filter((card) => cardMatchesSelectedFilters(card))
    .slice()
    .sort((left, right) => compareCardsByDisplayPrice(left, right)),
)

const groupAllCount = computed(
  () =>
    searchedCards.value.filter((card) => {
      if (selectedSupplier.value && card.supplier !== selectedSupplier.value) return false
      if (selectedBillingType.value && !cardHasMatchingOffer(card, { billingType: selectedBillingType.value })) return false
      return true
    }).length,
)

const supplierAllCount = computed(
  () =>
    searchedCards.value.filter((card) => {
      if (
        (selectedGroup.value || selectedBillingType.value) &&
        !cardHasMatchingOffer(card, {
          groupName: selectedGroup.value,
          billingType: selectedBillingType.value,
        })
      ) {
        return false
      }
      return true
    }).length,
)

const billingTypeAllCount = computed(
  () =>
    searchedCards.value.filter((card) => {
      if (selectedGroup.value && !cardHasMatchingOffer(card, { groupName: selectedGroup.value })) return false
      if (selectedSupplier.value && card.supplier !== selectedSupplier.value) return false
      return true
    }).length,
)

const groupEntries = computed<FilterEntry[]>(() =>
  groupOptions.value.map((value) => ({
    value,
    count: searchedCards.value.filter((card) => {
      if (selectedSupplier.value && card.supplier !== selectedSupplier.value) return false
      if (!cardHasMatchingOffer(card, { groupName: value, billingType: selectedBillingType.value })) return false
      return true
    }).length,
  })),
)

const supplierEntries = computed<FilterEntry[]>(() =>
  supplierOptions.value.map((value) => ({
    value,
    count: searchedCards.value.filter((card) => {
      if (card.supplier !== value) return false
      if (
        (selectedGroup.value || selectedBillingType.value) &&
        !cardHasMatchingOffer(card, {
          groupName: selectedGroup.value,
          billingType: selectedBillingType.value,
        })
      ) {
        return false
      }
      return true
    }).length,
  })),
)

const billingTypeEntries = computed<FilterEntry[]>(() =>
  billingTypeOptions.value.map((value) => ({
    value,
    count: searchedCards.value.filter((card) => {
      if (selectedSupplier.value && card.supplier !== selectedSupplier.value) return false
      if (!cardHasMatchingOffer(card, { groupName: selectedGroup.value, billingType: value })) return false
      return true
    }).length,
  })),
)

const activeCard = computed(() => {
  if (!activeCardKey.value) return null
  return cards.value.find((card) => cardKey(card) === activeCardKey.value) ?? null
})

const rateFormatter = new Intl.NumberFormat('en-US', {
  minimumFractionDigits: 0,
  maximumFractionDigits: 10,
  useGrouping: false,
})

function cardKey(card: Pick<ModelMarketplaceCard, 'supplier' | 'model_name'>): string {
  return `${card.supplier}::${card.model_name}`.toLowerCase()
}

function marketplaceSearchText(card: ModelMarketplaceCard): string {
  return [
    card.model_name,
    card.supplier,
    card.group_name,
    ...card.groups.map((group) => group.group_name),
  ]
    .join(' ')
    .toLowerCase()
}

function cardMatchesSelectedFilters(card: ModelMarketplaceCard): boolean {
  if (selectedSupplier.value && card.supplier !== selectedSupplier.value) return false
  if (selectedGroup.value || selectedBillingType.value) {
    return cardHasMatchingOffer(card, {
      groupName: selectedGroup.value,
      billingType: selectedBillingType.value,
    })
  }
  return true
}

function cardDisplayOffer(card: ModelMarketplaceCard): ModelMarketplaceCard | ModelMarketplaceGroupOffer {
  if (!selectedGroup.value && !selectedBillingType.value) return card

  const candidates = card.groups.filter((group) => groupOfferMatchesSelectedFilters(group))
  return pickBestDisplayOffer(candidates) ?? card
}

function cardHasMatchingOffer(card: ModelMarketplaceCard, filters: OfferFilter): boolean {
  return card.groups.some((group) => groupOfferMatchesFilters(group, filters))
}

function groupOfferMatchesSelectedFilters(group: ModelMarketplaceGroupOffer): boolean {
  return groupOfferMatchesFilters(group, {
    groupName: selectedGroup.value,
    billingType: selectedBillingType.value,
  })
}

function groupOfferMatchesFilters(group: ModelMarketplaceGroupOffer, filters: OfferFilter): boolean {
  if (filters.groupName && group.group_name !== filters.groupName) return false
  if (filters.billingType && normalizeBillingType(group.billing_type) !== filters.billingType) return false
  return true
}

function pickBestDisplayOffer(offers: ModelMarketplaceGroupOffer[]): ModelMarketplaceGroupOffer | null {
  if (offers.length === 0) return null
  return offers.slice().sort(compareGroupOffersByDisplayPrice)[0]
}

function compareGroupOffersByDisplayPrice(left: ModelMarketplaceGroupOffer, right: ModelMarketplaceGroupOffer): number {
  const leftPrice = primaryPriceDescriptor(left.current_pricing, left.billing_type)?.value
  const rightPrice = primaryPriceDescriptor(right.current_pricing, right.billing_type)?.value

  if (leftPrice == null && rightPrice != null) return 1
  if (leftPrice != null && rightPrice == null) return -1
  if (leftPrice != null && rightPrice != null && leftPrice !== rightPrice) {
    return leftPrice - rightPrice
  }
  if (left.group_name !== right.group_name) return left.group_name.localeCompare(right.group_name)
  return left.model_name.localeCompare(right.model_name)
}

function compareCardsByDisplayPrice(left: ModelMarketplaceCard, right: ModelMarketplaceCard): number {
  const leftOffer = cardDisplayOffer(left)
  const rightOffer = cardDisplayOffer(right)
  const leftPrice = primaryPriceDescriptor(leftOffer.current_pricing, leftOffer.billing_type)?.value
  const rightPrice = primaryPriceDescriptor(rightOffer.current_pricing, rightOffer.billing_type)?.value

  if (leftPrice == null && rightPrice != null) return 1
  if (leftPrice != null && rightPrice == null) return -1
  if (leftPrice != null && rightPrice != null && leftPrice !== rightPrice) {
    return leftPrice - rightPrice
  }
  if (left.supplier !== right.supplier) return left.supplier.localeCompare(right.supplier)
  return left.model_name.localeCompare(right.model_name)
}

function normalizeBillingType(mode: string | null | undefined): string {
  switch (mode) {
    case BILLING_MODE_PER_REQUEST:
      return BILLING_MODE_PER_REQUEST
    case BILLING_MODE_IMAGE:
      return BILLING_MODE_IMAGE
    default:
      return BILLING_MODE_TOKEN
  }
}

function billingTypeLabel(mode: string | null | undefined): string {
  switch (normalizeBillingType(mode)) {
    case BILLING_MODE_PER_REQUEST:
      return t('modelMarketplace.billingModePerRequest', '按次计费')
    case BILLING_MODE_IMAGE:
      return t('modelMarketplace.billingModeImage', '按次计费(图片)')
    default:
      return t('modelMarketplace.billingModeToken', '按量计费(tokens)')
  }
}

function billingTypeBadgeClass(mode: string | null | undefined): string {
  switch (normalizeBillingType(mode)) {
    case BILLING_MODE_PER_REQUEST:
      return 'bg-sky-100 text-sky-800 dark:bg-sky-500/15 dark:text-sky-300'
    case BILLING_MODE_IMAGE:
      return 'bg-rose-100 text-rose-800 dark:bg-rose-500/15 dark:text-rose-300'
    default:
      return 'bg-emerald-100 text-emerald-800 dark:bg-emerald-500/15 dark:text-emerald-300'
  }
}

function pricingSourceLabel(source: string): string {
  return source === 'channel'
    ? t('modelMarketplace.pricingSourceChannel', '渠道定价')
    : t('modelMarketplace.pricingSourceGroup', '分组现价')
}

function pricingSourceBadgeClass(source: string): string {
  return source === 'channel'
    ? 'bg-slate-900 text-white dark:bg-primary-500 dark:text-white'
    : 'bg-white text-slate-700 ring-1 ring-inset ring-slate-200 dark:bg-dark-800 dark:text-dark-200 dark:ring-dark-600'
}

function toggleGroup(value: string) {
  selectedGroup.value = selectedGroup.value === value ? '' : value
}

function toggleSupplier(value: string) {
  selectedSupplier.value = selectedSupplier.value === value ? '' : value
}

function toggleBillingType(value: string) {
  selectedBillingType.value = selectedBillingType.value === value ? '' : value
}

function resetCategoryFilters() {
  selectedGroup.value = ''
  selectedSupplier.value = ''
  selectedBillingType.value = ''
}

function supplierLabel(value: string | null | undefined): string {
  const normalized = (value || '').trim().toLowerCase()
  switch (normalized) {
    case 'openai':
      return 'OpenAI'
    case 'anthropic':
      return 'Anthropic'
    case 'google':
      return 'Google'
    case 'deepseek':
      return 'DeepSeek'
    case 'kimi':
      return 'Kimi'
    case 'moonshot':
      return 'Moonshot'
    case 'glm':
      return 'GLM'
    case 'qwen':
      return 'Qwen'
    case 'minimax':
      return 'MiniMax'
    case 'doubao':
      return 'Doubao'
    case 'unknown':
      return t('modelMarketplace.supplierUnknown', '未知供应商')
    default:
      return value || t('modelMarketplace.supplierUnknown', '未知供应商')
  }
}

function supplierBadgeClass(value: string | null | undefined): string {
  switch ((value || '').trim().toLowerCase()) {
    case 'openai':
      return 'bg-emerald-100 text-emerald-800 dark:bg-emerald-500/15 dark:text-emerald-300'
    case 'anthropic':
      return 'bg-orange-100 text-orange-800 dark:bg-orange-500/15 dark:text-orange-300'
    case 'google':
      return 'bg-blue-100 text-blue-800 dark:bg-blue-500/15 dark:text-blue-300'
    case 'deepseek':
      return 'bg-cyan-100 text-cyan-800 dark:bg-cyan-500/15 dark:text-cyan-300'
    case 'kimi':
      return 'bg-lime-100 text-lime-800 dark:bg-lime-500/15 dark:text-lime-300'
    case 'moonshot':
      return 'bg-indigo-100 text-indigo-800 dark:bg-indigo-500/15 dark:text-indigo-300'
    case 'glm':
      return 'bg-fuchsia-100 text-fuchsia-800 dark:bg-fuchsia-500/15 dark:text-fuchsia-300'
    case 'qwen':
      return 'bg-sky-100 text-sky-800 dark:bg-sky-500/15 dark:text-sky-300'
    case 'minimax':
      return 'bg-rose-100 text-rose-800 dark:bg-rose-500/15 dark:text-rose-300'
    case 'doubao':
      return 'bg-amber-100 text-amber-800 dark:bg-amber-500/15 dark:text-amber-300'
    default:
      return 'bg-slate-100 text-slate-700 dark:bg-dark-800 dark:text-dark-200'
  }
}

function formatTokenPrice(value: number | null | undefined): string {
  if (value == null) return '-'
  return `${formatScaled(value, 1_000_000)} / 1M Tokens`
}

function formatRequestPrice(value: number | null | undefined): string {
  if (value == null) return '-'
  return `${formatScaled(value, 1)} ${t('modelMarketplace.unitPerRequest', '/ 次')}`
}

function formatFieldPrice(field: PriceDescriptor['field'], value: number): string {
  if (field === 'per_request_price') {
    return formatRequestPrice(value)
  }
  return formatTokenPrice(value)
}

function pricingHasPerRequestPrice(pricing: ModelMarketplacePricing | null): boolean {
  if (!pricing) return false
  if (pricing.per_request_price != null) return true
  return pricing.intervals.some((interval) => interval.per_request_price != null)
}

function pricingHasImageOutputPrice(pricing: ModelMarketplacePricing | null): boolean {
  return pricing?.image_output_price != null
}

function shouldPreferPerRequestPricing(
  billingType: string | null | undefined,
  ...pricings: Array<ModelMarketplacePricing | null>
): boolean {
  if (normalizeBillingType(billingType) === BILLING_MODE_PER_REQUEST) return true
  return pricings.some((pricing) => pricingHasPerRequestPrice(pricing))
}

function shouldPreferImageOutputPricing(
  billingType: string | null | undefined,
  ...pricings: Array<ModelMarketplacePricing | null>
): boolean {
  if (normalizeBillingType(billingType) !== BILLING_MODE_IMAGE) return false
  return !shouldPreferPerRequestPricing(billingType, ...pricings) && pricings.some((pricing) => pricingHasImageOutputPrice(pricing))
}

function primaryPriceDescriptor(pricing: ModelMarketplacePricing | null, billingType: string): PriceDescriptor | null {
  if (!pricing) return null

  const preferPerRequest = shouldPreferPerRequestPricing(billingType, pricing)
  const preferImageOutput = shouldPreferImageOutputPricing(billingType, pricing)
  const directFields = preferPerRequest
    ? ['per_request_price']
    : preferImageOutput
      ? ['image_output_price']
      : ['input_price', 'output_price', 'cache_write_price', 'cache_read_price', 'image_output_price']

  for (const field of directFields as PriceDescriptor['field'][]) {
    const value = pricing[field]
    if (value != null) {
      return { field, value }
    }
  }

  const intervalCandidate = pricing.intervals
    .map((interval) => intervalPriceDescriptor(interval, billingType, preferPerRequest))
    .filter((entry): entry is PriceDescriptor => entry !== null)
    .sort((left, right) => left.value - right.value)[0]

  if (intervalCandidate) {
    return intervalCandidate
  }

  if (preferPerRequest || preferImageOutput) {
    return null
  }

  for (const field of ['output_price', 'cache_write_price', 'cache_read_price', 'image_output_price'] as PriceDescriptor['field'][]) {
    const value = pricing[field]
    if (value != null) {
      return { field, value }
    }
  }

  return null
}

function intervalPriceDescriptor(
  interval: ModelMarketplacePricingInterval,
  billingType: string,
  preferPerRequest = false,
): PriceDescriptor | null {
  if (preferPerRequest && interval.per_request_price != null) {
    return { field: 'per_request_price', value: interval.per_request_price }
  }

  if (!preferPerRequest && normalizeBillingType(billingType) !== BILLING_MODE_TOKEN && interval.per_request_price != null) {
    return { field: 'per_request_price', value: interval.per_request_price }
  }

  if (preferPerRequest) {
    return null
  }

  for (const field of ['input_price', 'output_price', 'cache_write_price', 'cache_read_price', 'per_request_price'] as const) {
    const value = interval[field]
    if (value != null) {
      return { field, value }
    }
  }

  return null
}

function fieldLabel(field: PriceDescriptor['field']): string {
  switch (field) {
    case 'output_price':
      return t('modelMarketplace.outputPrice', '输出')
    case 'cache_write_price':
      return t('modelMarketplace.cacheWritePrice', '缓存写入')
    case 'cache_read_price':
      return t('modelMarketplace.cacheReadPrice', '缓存读取')
    case 'image_output_price':
      return t('modelMarketplace.imagePrice', '图片输出')
    case 'per_request_price':
      return t('modelMarketplace.requestPrice', '按次')
    default:
      return t('modelMarketplace.inputPrice', '输入')
  }
}

function formatCardPriceLabel(label: string): string {
  if (label.includes(t('modelMarketplace.intervalPricing', '区间价格'))) {
    return label
  }
  return `${label}${t('modelMarketplace.pricing', '价格')}`
}

function buildPricingLines(pricing: ModelMarketplacePricing | null, displayMode: PricingDisplayMode = 'default'): PriceSummaryLine[] {
  if (!pricing) return []

  const lines: PriceSummaryLine[] = []
  const directFields: PriceDescriptor['field'][] = displayMode === 'per_request'
    ? ['per_request_price']
    : displayMode === 'image_output'
      ? ['image_output_price']
      : [
        'input_price',
        'output_price',
        'cache_write_price',
        'cache_read_price',
        'per_request_price',
        'image_output_price',
      ]

  directFields.forEach((field) => {
    const value = pricing[field]
    if (value == null) return
    lines.push({
      label: fieldLabel(field),
      value: formatFieldPrice(field, value),
    })
  })

  if (pricing.intervals.length === 0) {
    return lines
  }

  pricing.intervals.forEach((interval) => {
    const descriptor = intervalPriceDescriptor(interval, pricing.billing_mode, displayMode === 'per_request')
    if (!descriptor) return
    lines.push({
      label: `${t('modelMarketplace.intervalPricing', '区间价格')} ${intervalLabel(interval)}`,
      value: formatFieldPrice(descriptor.field, descriptor.value),
    })
  })

  return lines
}

function buildPricingComparisonLines(
  currentPricing: ModelMarketplacePricing | null,
  originalPricing: ModelMarketplacePricing | null,
  billingType: string | null | undefined,
): PriceComparisonLine[] {
  const preferPerRequest = shouldPreferPerRequestPricing(billingType, currentPricing, originalPricing)
  const preferImageOutput = shouldPreferImageOutputPricing(billingType, currentPricing, originalPricing)
  const displayMode: PricingDisplayMode = preferPerRequest
    ? 'per_request'
    : preferImageOutput
      ? 'image_output'
      : 'default'
  const currentLines = buildPricingLines(currentPricing, displayMode)
  const originalLines = buildPricingLines(originalPricing, displayMode)
  const currentByLabel = new Map(currentLines.map((line) => [line.label, line.value]))
  const originalByLabel = new Map(originalLines.map((line) => [line.label, line.value]))
  const orderedLabels = Array.from(new Set([...currentLines.map((line) => line.label), ...originalLines.map((line) => line.label)]))

  return orderedLabels.map((label) => ({
    label,
    currentValue: currentByLabel.get(label) ?? null,
    originalValue: originalByLabel.get(label) ?? null,
  }))
}

function intervalLabel(interval: ModelMarketplacePricingInterval): string {
  if (interval.tier_label) return interval.tier_label
  const upper = interval.max_tokens == null ? 'max' : String(interval.max_tokens)
  return `${interval.min_tokens}-${upper}`
}

function formatRate(value: number): string {
  if (!Number.isFinite(value)) return '-'
  return `${rateFormatter.format(value)}x`
}

function openCard(card: ModelMarketplaceCard) {
  activeCardKey.value = cardKey(card)
}

function closeCard() {
  activeCardKey.value = ''
}

async function copyModelName(modelName: string) {
  const success = await copyToClipboard(modelName, t('common.copiedToClipboard'))
  if (!success) return
  copiedModelName.value = modelName
  window.setTimeout(() => {
    if (copiedModelName.value === modelName) {
      copiedModelName.value = ''
    }
  }, 1800)
}

async function loadMarketplace() {
  loading.value = true
  try {
    cards.value = await modelMarketplaceAPI.getModelMarketplace()
    if (activeCardKey.value && !cards.value.some((card) => cardKey(card) === activeCardKey.value)) {
      activeCardKey.value = ''
    }
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, t('common.error')))
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  if (!appStore.publicSettingsLoaded) {
    await appStore.fetchPublicSettings()
  }
  await loadMarketplace()
})
</script>
