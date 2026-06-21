<template>
  <div
    class="relative min-h-screen overflow-hidden bg-gradient-to-br from-amber-50 via-white to-indigo-50 text-gray-900 dark:from-dark-950 dark:via-dark-950 dark:to-dark-900 dark:text-gray-100"
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
            <p class="text-xs font-semibold uppercase tracking-[0.24em] text-amber-700/80 dark:text-amber-300/80">
              {{ t('nav.modelMarketplace', 'Model Marketplace') }}
            </p>
            <h1 class="text-lg font-semibold text-slate-900 dark:text-white">{{ siteName }}</h1>
          </div>
        </router-link>

        <div class="flex items-center gap-3">
          <router-link
            v-if="isAuthenticated"
            :to="dashboardPath"
            class="rounded-full border border-slate-200/80 bg-white/90 px-4 py-2 text-sm font-medium text-slate-700 shadow-sm backdrop-blur transition hover:border-slate-300 hover:text-slate-900 dark:border-dark-700 dark:bg-dark-800/90 dark:text-dark-200 dark:hover:border-dark-600 dark:hover:text-white"
          >
            {{ t('home.dashboard') }}
          </router-link>
          <router-link
            v-else
            to="/login"
            class="rounded-full bg-slate-900 px-4 py-2 text-sm font-medium text-white shadow-lg shadow-slate-300 transition hover:bg-slate-800 dark:bg-primary-500 dark:shadow-primary-500/20 dark:hover:bg-primary-600"
          >
            {{ t('home.login') }}
          </router-link>
        </div>
      </nav>
    </header>

    <main class="relative z-10 px-6 pb-16 pt-6">
      <div class="w-full">
        <section
          class="w-full rounded-[2rem] border border-white/70 bg-white/80 p-8 shadow-[0_25px_80px_-40px_rgba(15,23,42,0.45)] backdrop-blur dark:border-dark-700/70 dark:bg-dark-900/75 dark:shadow-[0_25px_80px_-40px_rgba(0,0,0,0.55)]"
        >
          <div class="flex flex-col gap-8 xl:flex-row xl:items-end xl:justify-between">
            <div class="max-w-3xl">
              <p class="mb-3 text-xs font-semibold uppercase tracking-[0.3em] text-amber-700/80 dark:text-amber-300/80">
                {{ t('nav.modelMarketplace', 'Model Marketplace') }}
              </p>
              <h2 class="text-4xl font-black leading-tight text-slate-950 dark:text-white md:text-5xl">
                {{ t('modelMarketplace.heroTitle', '按模型聚合浏览可用能力与最低展示价') }}
              </h2>
              <p class="mt-4 max-w-2xl text-base leading-7 text-slate-600 dark:text-dark-300">
                {{
                  t(
                    'modelMarketplace.heroDescription',
                    '每张卡片只保留一个模型，并展示当前最低价、原价对比，以及该模型所在分组的完整明细。',
                  )
                }}
              </p>
            </div>

            <div class="grid gap-4 sm:grid-cols-3 xl:min-w-[24rem]">
              <div class="rounded-2xl border border-amber-100 bg-amber-50/90 p-4 dark:border-amber-500/20 dark:bg-amber-500/10">
                <p class="text-xs font-semibold uppercase tracking-[0.2em] text-amber-700 dark:text-amber-300">
                  {{ t('modelMarketplace.cardCount', 'Cards') }}
                </p>
                <p class="mt-2 text-3xl font-bold text-slate-900 dark:text-white">{{ filteredCards.length }}</p>
              </div>
              <div class="rounded-2xl border border-sky-100 bg-sky-50/90 p-4 dark:border-sky-500/20 dark:bg-sky-500/10">
                <p class="text-xs font-semibold uppercase tracking-[0.2em] text-sky-700 dark:text-sky-300">
                  {{ t('modelMarketplace.groupLabel', '分组') }}
                </p>
                <p class="mt-2 text-3xl font-bold text-slate-900 dark:text-white">{{ groupOptions.length }}</p>
              </div>
              <div
                class="rounded-2xl border border-violet-100 bg-violet-50/90 p-4 dark:border-violet-500/20 dark:bg-violet-500/10"
              >
                <p class="text-xs font-semibold uppercase tracking-[0.2em] text-violet-700 dark:text-violet-300">
                  {{ t('availableChannels.columns.platform', 'Platforms') }}
                </p>
                <p class="mt-2 text-3xl font-bold text-slate-900 dark:text-white">{{ platformOptions.length }}</p>
              </div>
            </div>
          </div>

          <div
            class="mt-8 flex flex-col gap-4 rounded-[1.75rem] border border-slate-200/80 bg-white/85 p-5 shadow-sm dark:border-dark-700 dark:bg-dark-950/50 lg:flex-row lg:items-center"
          >
            <div class="min-w-0 flex-1">
              <label class="mb-2 block text-sm font-medium text-slate-700 dark:text-dark-200">
                {{ t('common.search', 'Search') }}
              </label>
              <input
                v-model="searchQuery"
                type="text"
                :placeholder="t('modelMarketplace.searchPlaceholder', '搜索模型名称')"
                class="w-full rounded-2xl border border-slate-200 bg-white px-4 py-3 text-sm text-slate-900 placeholder:text-slate-400 outline-none transition focus:border-amber-400 dark:border-dark-600 dark:bg-dark-900 dark:text-white dark:placeholder:text-dark-400 dark:focus:border-amber-400"
              />
            </div>

            <div class="flex flex-wrap items-center gap-3">
              <button
                v-if="hasCategoryFilters"
                @click="resetCategoryFilters"
                class="inline-flex items-center rounded-full border border-slate-200 bg-white px-4 py-2.5 text-sm font-semibold text-slate-700 transition hover:border-slate-300 hover:text-slate-900 dark:border-dark-600 dark:bg-dark-900 dark:text-dark-200 dark:hover:border-dark-500 dark:hover:text-white"
              >
                {{ t('common.reset', 'Reset') }}
              </button>
              <button
                @click="loadMarketplace"
                :disabled="loading"
                class="inline-flex items-center rounded-full bg-amber-400 px-5 py-2.5 text-sm font-semibold text-slate-950 transition hover:bg-amber-300 disabled:cursor-not-allowed disabled:opacity-60"
              >
                {{ loading ? t('common.loading', 'Loading...') : t('common.refresh', 'Refresh') }}
              </button>
            </div>
          </div>
        </section>

        <section class="mt-8 grid gap-6 xl:grid-cols-[320px_minmax(0,1fr)]">
          <aside
            class="h-fit rounded-[1.75rem] border border-white/70 bg-white/85 p-4 shadow-[0_20px_60px_-45px_rgba(15,23,42,0.55)] backdrop-blur dark:border-dark-700/70 dark:bg-dark-900/80 dark:shadow-[0_20px_60px_-45px_rgba(0,0,0,0.6)] xl:sticky xl:top-6"
          >
            <div class="mb-4 flex items-center justify-between gap-3 px-2">
              <div>
                <p class="text-xs font-semibold uppercase tracking-[0.22em] text-slate-400 dark:text-dark-400">
                  {{ t('common.filter', 'Filter') }}
                </p>
                <p class="mt-1 text-sm text-slate-600 dark:text-dark-300">{{ activeFilterSummary }}</p>
              </div>
              <button
                v-if="hasCategoryFilters"
                @click="resetCategoryFilters"
                class="text-xs font-semibold text-amber-700 transition hover:text-amber-600 dark:text-amber-300 dark:hover:text-amber-200"
              >
                {{ t('common.reset', 'Reset') }}
              </button>
            </div>

            <div class="space-y-3">
              <details open class="group overflow-hidden rounded-2xl border border-slate-200/80 bg-slate-50/75 dark:border-dark-700 dark:bg-dark-950/60">
                <summary
                  class="flex cursor-pointer list-none items-center justify-between gap-3 px-4 py-3 text-sm font-semibold text-slate-900 dark:text-white"
                >
                  <span>{{ t('modelMarketplace.groupLabel', '分组') }}</span>
                  <span class="rounded-full bg-white px-2 py-0.5 text-xs text-slate-500 dark:bg-dark-800 dark:text-dark-300">
                    {{ groupOptions.length }}
                  </span>
                </summary>
                <div class="space-y-1 border-t border-slate-200/70 p-2 dark:border-dark-700">
                  <button
                    @click="selectedGroup = ''"
                    class="flex w-full items-center justify-between rounded-xl px-3 py-2 text-left text-sm transition"
                    :class="
                      selectedGroup
                        ? 'text-slate-600 hover:bg-white hover:text-slate-900 dark:text-dark-300 dark:hover:bg-dark-800 dark:hover:text-white'
                        : 'bg-slate-900 text-white dark:bg-primary-500 dark:text-white'
                    "
                  >
                    <span>{{ t('common.all', 'All') }}</span>
                    <span class="text-xs opacity-75">{{ groupAllCount }}</span>
                  </button>
                  <button
                    v-for="group in groupEntries"
                    :key="group.value"
                    @click="toggleGroup(group.value)"
                    class="flex w-full items-center justify-between rounded-xl px-3 py-2 text-left text-sm transition"
                    :class="
                      selectedGroup === group.value
                        ? 'bg-slate-900 text-white dark:bg-primary-500 dark:text-white'
                        : 'text-slate-600 hover:bg-white hover:text-slate-900 dark:text-dark-300 dark:hover:bg-dark-800 dark:hover:text-white'
                    "
                  >
                    <span class="truncate pr-3">{{ group.value }}</span>
                    <span class="text-xs opacity-75">{{ group.count }}</span>
                  </button>
                </div>
              </details>

              <details open class="group overflow-hidden rounded-2xl border border-slate-200/80 bg-slate-50/75 dark:border-dark-700 dark:bg-dark-950/60">
                <summary
                  class="flex cursor-pointer list-none items-center justify-between gap-3 px-4 py-3 text-sm font-semibold text-slate-900 dark:text-white"
                >
                  <span>{{ t('availableChannels.columns.platform', 'Platforms') }}</span>
                  <span class="rounded-full bg-white px-2 py-0.5 text-xs text-slate-500 dark:bg-dark-800 dark:text-dark-300">
                    {{ platformOptions.length }}
                  </span>
                </summary>
                <div class="space-y-1 border-t border-slate-200/70 p-2 dark:border-dark-700">
                  <button
                    @click="selectedPlatform = ''"
                    class="flex w-full items-center justify-between rounded-xl px-3 py-2 text-left text-sm transition"
                    :class="
                      selectedPlatform
                        ? 'text-slate-600 hover:bg-white hover:text-slate-900 dark:text-dark-300 dark:hover:bg-dark-800 dark:hover:text-white'
                        : 'bg-slate-900 text-white dark:bg-primary-500 dark:text-white'
                    "
                  >
                    <span>{{ t('common.all', 'All') }}</span>
                    <span class="text-xs opacity-75">{{ platformAllCount }}</span>
                  </button>
                  <button
                    v-for="platform in platformEntries"
                    :key="platform.value"
                    @click="togglePlatform(platform.value)"
                    class="flex w-full items-center justify-between rounded-xl px-3 py-2 text-left text-sm transition"
                    :class="
                      selectedPlatform === platform.value
                        ? 'bg-slate-900 text-white dark:bg-primary-500 dark:text-white'
                        : 'text-slate-600 hover:bg-white hover:text-slate-900 dark:text-dark-300 dark:hover:bg-dark-800 dark:hover:text-white'
                    "
                  >
                    <span class="truncate pr-3">{{ platform.value }}</span>
                    <span class="text-xs opacity-75">{{ platform.count }}</span>
                  </button>
                </div>
              </details>

              <details open class="group overflow-hidden rounded-2xl border border-slate-200/80 bg-slate-50/75 dark:border-dark-700 dark:bg-dark-950/60">
                <summary
                  class="flex cursor-pointer list-none items-center justify-between gap-3 px-4 py-3 text-sm font-semibold text-slate-900 dark:text-white"
                >
                  <span>{{ t('modelMarketplace.billingType', 'Billing Type') }}</span>
                  <span class="rounded-full bg-white px-2 py-0.5 text-xs text-slate-500 dark:bg-dark-800 dark:text-dark-300">
                    {{ billingTypeOptions.length }}
                  </span>
                </summary>
                <div class="space-y-1 border-t border-slate-200/70 p-2 dark:border-dark-700">
                  <button
                    @click="selectedBillingType = ''"
                    class="flex w-full items-center justify-between rounded-xl px-3 py-2 text-left text-sm transition"
                    :class="
                      selectedBillingType
                        ? 'text-slate-600 hover:bg-white hover:text-slate-900 dark:text-dark-300 dark:hover:bg-dark-800 dark:hover:text-white'
                        : 'bg-slate-900 text-white dark:bg-primary-500 dark:text-white'
                    "
                  >
                    <span>{{ t('common.all', 'All') }}</span>
                    <span class="text-xs opacity-75">{{ billingTypeAllCount }}</span>
                  </button>
                  <button
                    v-for="billingType in billingTypeEntries"
                    :key="billingType.value"
                    @click="toggleBillingType(billingType.value)"
                    class="flex w-full items-center justify-between rounded-xl px-3 py-2 text-left text-sm transition"
                    :class="
                      selectedBillingType === billingType.value
                        ? 'bg-slate-900 text-white dark:bg-primary-500 dark:text-white'
                        : 'text-slate-600 hover:bg-white hover:text-slate-900 dark:text-dark-300 dark:hover:bg-dark-800 dark:hover:text-white'
                    "
                  >
                    <span class="truncate pr-3">{{ billingTypeLabel(billingType.value) }}</span>
                    <span class="text-xs opacity-75">{{ billingType.count }}</span>
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
                Platform
                <select v-model="selectedPlatform">
                  <option value="">{{ t('common.all', 'All') }}</option>
                  <option v-for="platform in platformOptions" :key="platform" :value="platform">{{ platform }}</option>
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
            <section v-if="loading" class="grid gap-5 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4">
              <div
                v-for="index in 6"
                :key="index"
                class="h-72 animate-pulse rounded-[1.75rem] border border-white/70 bg-white/70 dark:border-dark-700/70 dark:bg-dark-800/70"
              />
            </section>

            <section
              v-else-if="filteredCards.length === 0"
              class="rounded-[2rem] border border-dashed border-slate-300 bg-white/70 px-8 py-20 text-center shadow-sm dark:border-dark-700 dark:bg-dark-900/60"
            >
              <h3 class="text-2xl font-bold text-slate-900 dark:text-white">
                {{ t('modelMarketplace.emptyTitle', 'No matching models') }}
              </h3>
              <p class="mt-3 text-sm text-slate-500 dark:text-dark-400">
                {{ t('modelMarketplace.emptyDescription', 'Adjust your search or filters and try again.') }}
              </p>
            </section>

            <section v-else class="grid gap-5 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4">
              <article
                v-for="card in filteredCards"
                :key="cardKey(card)"
                :data-testid="`marketplace-card-${card.platform}-${card.model_name}`"
                class="group cursor-pointer rounded-[1.75rem] border border-white/70 bg-white/88 p-6 shadow-[0_20px_70px_-45px_rgba(15,23,42,0.5)] backdrop-blur transition hover:-translate-y-1 hover:shadow-[0_30px_80px_-45px_rgba(15,23,42,0.55)] dark:border-dark-700/70 dark:bg-dark-900/82 dark:shadow-[0_20px_70px_-45px_rgba(0,0,0,0.6)] dark:hover:shadow-[0_30px_80px_-45px_rgba(0,0,0,0.7)]"
                @click="openCard(card)"
              >
                <div class="flex items-start justify-between gap-4">
                  <div class="min-w-0">
                    <p class="truncate text-xs font-semibold uppercase tracking-[0.2em] text-slate-400 dark:text-dark-400">
                      {{ card.platform }}
                    </p>
                    <h3 class="mt-2 min-w-0 break-words text-[1.125rem] font-bold leading-snug text-slate-950 dark:text-white">
                      {{ card.model_name }}
                    </h3>
                  </div>
                  <span
                    class="shrink-0 rounded-full px-3 py-1 text-xs font-semibold"
                    :class="pricingSourceBadgeClass(card.pricing_source)"
                  >
                    {{ pricingSourceLabel(card.pricing_source) }}
                  </span>
                </div>

                <div class="mt-5 flex flex-wrap gap-2">
                  <span
                    class="inline-flex items-center gap-1.5 rounded-full bg-amber-100 px-3 py-1 text-xs font-medium text-amber-800 dark:bg-amber-500/15 dark:text-amber-300"
                  >
                    {{ card.group_name }}
                  </span>
                  <span
                    class="rounded-full bg-sky-100 px-3 py-1 text-xs font-medium text-sky-800 dark:bg-sky-500/15 dark:text-sky-300"
                  >
                    {{ billingTypeLabel(card.billing_type) }}
                  </span>
                  <span
                    class="rounded-full bg-violet-100 px-3 py-1 text-xs font-medium text-violet-800 dark:bg-violet-500/15 dark:text-violet-300"
                  >
                    {{ t('modelMarketplace.groupCount', '{count} 个分组').replace('{count}', String(card.groups.length)) }}
                  </span>
                </div>

                <div class="mt-6 rounded-[1.5rem] border border-slate-200/80 bg-slate-50/90 p-4 dark:border-dark-700 dark:bg-dark-950/70">
                  <div class="grid grid-cols-2 gap-3">
                    <div>
                      <p class="text-[11px] font-semibold uppercase tracking-[0.18em] text-slate-400 dark:text-dark-400">
                        {{ t('modelMarketplace.currentPrice', '现价') }}
                      </p>
                      <p class="mt-2 text-lg font-bold text-slate-950 dark:text-white">
                        {{ formatPrimaryPrice(card.current_pricing, card.billing_type) }}
                      </p>
                      <p class="mt-1 text-xs text-slate-500 dark:text-dark-400">
                        {{ primaryPriceLabel(card.current_pricing, card.billing_type) }}
                      </p>
                    </div>
                    <div>
                      <p class="text-[11px] font-semibold uppercase tracking-[0.18em] text-slate-400 dark:text-dark-400">
                        {{ t('modelMarketplace.originalPrice', '原价') }}
                      </p>
                      <p class="mt-2 text-lg font-semibold text-slate-500 line-through dark:text-dark-400">
                        {{ formatPrimaryPrice(card.original_pricing, card.billing_type) }}
                      </p>
                      <p class="mt-1 text-xs text-slate-500 dark:text-dark-400">
                        {{ primaryPriceLabel(card.original_pricing, card.billing_type) }}
                      </p>
                    </div>
                  </div>
                </div>

                <dl class="mt-6 space-y-3 text-sm">
                  <div class="flex items-center justify-between gap-4">
                    <dt class="text-slate-500 dark:text-dark-400">{{ t('modelMarketplace.lowestPriceGroup', '最低展示价分组') }}</dt>
                    <dd class="font-medium text-slate-900 dark:text-white">{{ card.group_name }}</dd>
                  </div>
                  <div class="flex items-center justify-between gap-4">
                    <dt class="text-slate-500 dark:text-dark-400">{{ t('modelMarketplace.groupRate', 'Rate') }}</dt>
                    <dd class="font-medium text-slate-900 dark:text-white">{{ formatRate(card.group_rate) }}</dd>
                  </div>
                </dl>

                <div
                  class="mt-6 flex items-center justify-between rounded-2xl border border-slate-200/80 bg-white/90 px-4 py-3 text-sm font-medium text-slate-700 transition group-hover:border-slate-300 group-hover:text-slate-900 dark:border-dark-700 dark:bg-dark-900/90 dark:text-dark-200 dark:group-hover:border-dark-500 dark:group-hover:text-white"
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
                <p class="text-xs font-semibold uppercase tracking-[0.22em] text-slate-400 dark:text-dark-400">
                  {{ activeCard.platform }}
                </p>
                <div class="mt-2 flex items-start gap-2">
                  <h3 class="break-words text-2xl font-bold text-slate-950 dark:text-white">
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
                    class="rounded-full bg-amber-100 px-3 py-1 text-xs font-medium text-amber-800 dark:bg-amber-500/15 dark:text-amber-300"
                  >
                    {{ t('modelMarketplace.groupCount', '{count} 个分组').replace('{count}', String(activeCard.groups.length)) }}
                  </span>
                  <span
                    class="rounded-full px-3 py-1 text-xs font-medium"
                    :class="billingTypeBadgeClass(activeCard.billing_type)"
                  >
                    {{ billingTypeLabel(activeCard.billing_type) }}
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
                      <h4 class="text-lg font-semibold text-slate-950 dark:text-white">{{ group.group_name }}</h4>
                      <div class="mt-3 flex flex-wrap gap-2">
                        <span
                          class="rounded-full px-3 py-1 text-xs font-medium"
                          :class="billingTypeBadgeClass(group.billing_type)"
                        >
                          {{ billingTypeLabel(group.billing_type) }}
                        </span>
                        <span
                          class="rounded-full px-3 py-1 text-xs font-medium"
                          :class="pricingSourceBadgeClass(group.pricing_source)"
                        >
                          {{ pricingSourceLabel(group.pricing_source) }}
                        </span>
                        <span
                          class="rounded-full bg-slate-100 px-3 py-1 text-xs font-medium text-slate-700 dark:bg-dark-800 dark:text-dark-200"
                        >
                          {{ group.subscription_type || '-' }}
                        </span>
                      </div>
                    </div>
                    <div class="text-sm text-slate-500 dark:text-dark-400">
                      {{ t('modelMarketplace.groupRate', 'Rate') }}: {{ formatRate(group.group_rate) }}
                    </div>
                  </div>

                  <div class="mt-5 grid gap-3 md:grid-cols-2">
                    <section class="rounded-2xl border border-slate-200/80 bg-slate-50/90 p-4 dark:border-dark-700 dark:bg-dark-900/70">
                      <p class="text-[11px] font-semibold uppercase tracking-[0.18em] text-slate-400 dark:text-dark-400">
                        {{ t('modelMarketplace.originalPrice', '原价') }}
                      </p>
                      <div v-if="buildPricingLines(group.original_pricing).length > 0" class="mt-3 space-y-2 text-sm">
                        <div
                          v-for="line in buildPricingLines(group.original_pricing)"
                          :key="`${group.group_id}-original-${line.label}`"
                          class="flex items-center justify-between gap-4"
                        >
                          <span class="text-slate-500 dark:text-dark-400">{{ line.label }}</span>
                          <span class="text-right font-medium text-slate-900 dark:text-white">{{ line.value }}</span>
                        </div>
                      </div>
                      <p v-else class="mt-3 text-sm text-slate-500 dark:text-dark-400">
                        {{ t('availableChannels.noPricing', 'No pricing configured') }}
                      </p>
                    </section>

                    <section class="rounded-2xl border border-amber-200/80 bg-amber-50/80 p-4 dark:border-amber-500/20 dark:bg-amber-500/10">
                      <p class="text-[11px] font-semibold uppercase tracking-[0.18em] text-amber-700 dark:text-amber-300">
                        {{ t('modelMarketplace.currentPrice', '现价') }}
                      </p>
                      <div v-if="buildPricingLines(group.current_pricing).length > 0" class="mt-3 space-y-2 text-sm">
                        <div
                          v-for="line in buildPricingLines(group.current_pricing)"
                          :key="`${group.group_id}-current-${line.label}`"
                          class="flex items-center justify-between gap-4"
                        >
                          <span class="text-amber-700/80 dark:text-amber-200/80">{{ line.label }}</span>
                          <span class="text-right font-semibold text-slate-950 dark:text-white">{{ line.value }}</span>
                        </div>
                      </div>
                      <p v-else class="mt-3 text-sm text-amber-700/80 dark:text-amber-200/80">
                        {{ t('modelMarketplace.priceUnavailable', '暂无价格') }}
                      </p>
                    </section>
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

type PriceDescriptor = {
  field: 'input_price' | 'output_price' | 'cache_write_price' | 'cache_read_price' | 'image_output_price' | 'per_request_price'
  value: number
}

const { t } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()
const { copyToClipboard } = useClipboard()

const loading = ref(false)
const cards = ref<ModelMarketplaceCard[]>([])
const searchQuery = ref('')
const selectedGroup = ref('')
const selectedPlatform = ref('')
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
const platformOptions = computed(() =>
  Array.from(new Set(cards.value.map((card) => card.platform))).sort((a, b) => a.localeCompare(b)),
)
const billingTypeOptions = computed(() =>
  Array.from(
    new Set(cards.value.flatMap((card) => card.groups.map((group) => normalizeBillingType(group.billing_type))).filter(Boolean)),
  ).sort((a, b) => a.localeCompare(b)),
)

const hasCategoryFilters = computed(
  () => Boolean(selectedGroup.value || selectedPlatform.value || selectedBillingType.value),
)

const activeFilterSummary = computed(() => {
  const count = [selectedGroup.value, selectedPlatform.value, selectedBillingType.value].filter(Boolean).length
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
      if (selectedPlatform.value && card.platform !== selectedPlatform.value) return false
      if (selectedBillingType.value && !cardHasBillingType(card, selectedBillingType.value)) return false
      return true
    }).length,
)

const platformAllCount = computed(
  () =>
    searchedCards.value.filter((card) => {
      if (selectedGroup.value && !cardBelongsToGroup(card, selectedGroup.value)) return false
      if (selectedBillingType.value && !cardHasBillingType(card, selectedBillingType.value)) return false
      return true
    }).length,
)

const billingTypeAllCount = computed(
  () =>
    searchedCards.value.filter((card) => {
      if (selectedGroup.value && !cardBelongsToGroup(card, selectedGroup.value)) return false
      if (selectedPlatform.value && card.platform !== selectedPlatform.value) return false
      return true
    }).length,
)

const groupEntries = computed<FilterEntry[]>(() =>
  groupOptions.value.map((value) => ({
    value,
    count: searchedCards.value.filter((card) => {
      if (!cardBelongsToGroup(card, value)) return false
      if (selectedPlatform.value && card.platform !== selectedPlatform.value) return false
      if (selectedBillingType.value && !cardHasBillingType(card, selectedBillingType.value)) return false
      return true
    }).length,
  })),
)

const platformEntries = computed<FilterEntry[]>(() =>
  platformOptions.value.map((value) => ({
    value,
    count: searchedCards.value.filter((card) => {
      if (card.platform !== value) return false
      if (selectedGroup.value && !cardBelongsToGroup(card, selectedGroup.value)) return false
      if (selectedBillingType.value && !cardHasBillingType(card, selectedBillingType.value)) return false
      return true
    }).length,
  })),
)

const billingTypeEntries = computed<FilterEntry[]>(() =>
  billingTypeOptions.value.map((value) => ({
    value,
    count: searchedCards.value.filter((card) => {
      if (!cardHasBillingType(card, value)) return false
      if (selectedGroup.value && !cardBelongsToGroup(card, selectedGroup.value)) return false
      if (selectedPlatform.value && card.platform !== selectedPlatform.value) return false
      return true
    }).length,
  })),
)

const activeCard = computed(() => {
  if (!activeCardKey.value) return null
  return cards.value.find((card) => cardKey(card) === activeCardKey.value) ?? null
})

function cardKey(card: Pick<ModelMarketplaceCard, 'platform' | 'model_name'>): string {
  return `${card.platform}::${card.model_name}`.toLowerCase()
}

function marketplaceSearchText(card: ModelMarketplaceCard): string {
  return [
    card.model_name,
    card.platform,
    card.group_name,
    ...card.groups.map((group) => group.group_name),
  ]
    .join(' ')
    .toLowerCase()
}

function cardMatchesSelectedFilters(card: ModelMarketplaceCard): boolean {
  if (selectedGroup.value && !cardBelongsToGroup(card, selectedGroup.value)) return false
  if (selectedPlatform.value && card.platform !== selectedPlatform.value) return false
  if (selectedBillingType.value && !cardHasBillingType(card, selectedBillingType.value)) return false
  return true
}

function cardBelongsToGroup(card: ModelMarketplaceCard, groupName: string): boolean {
  return card.groups.some((group) => group.group_name === groupName)
}

function cardHasBillingType(card: ModelMarketplaceCard, billingType: string): boolean {
  return card.groups.some((group) => normalizeBillingType(group.billing_type) === billingType)
}

function compareCardsByDisplayPrice(left: ModelMarketplaceCard, right: ModelMarketplaceCard): number {
  const leftPrice = primaryPriceDescriptor(left.current_pricing, left.billing_type)?.value
  const rightPrice = primaryPriceDescriptor(right.current_pricing, right.billing_type)?.value

  if (leftPrice == null && rightPrice != null) return 1
  if (leftPrice != null && rightPrice == null) return -1
  if (leftPrice != null && rightPrice != null && leftPrice !== rightPrice) {
    return leftPrice - rightPrice
  }
  if (left.platform !== right.platform) return left.platform.localeCompare(right.platform)
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

function togglePlatform(value: string) {
  selectedPlatform.value = selectedPlatform.value === value ? '' : value
}

function toggleBillingType(value: string) {
  selectedBillingType.value = selectedBillingType.value === value ? '' : value
}

function resetCategoryFilters() {
  selectedGroup.value = ''
  selectedPlatform.value = ''
  selectedBillingType.value = ''
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

function formatPrimaryPrice(pricing: ModelMarketplacePricing | null, billingType: string): string {
  const descriptor = primaryPriceDescriptor(pricing, billingType)
  if (!descriptor) return t('modelMarketplace.priceUnavailable', '暂无价格')
  return formatFieldPrice(descriptor.field, descriptor.value)
}

function primaryPriceLabel(pricing: ModelMarketplacePricing | null, billingType: string): string {
  const descriptor = primaryPriceDescriptor(pricing, billingType)
  if (!descriptor) return t('modelMarketplace.priceUnavailable', '暂无价格')
  return fieldLabel(descriptor.field)
}

function primaryPriceDescriptor(pricing: ModelMarketplacePricing | null, billingType: string): PriceDescriptor | null {
  if (!pricing) return null

  const directFields =
    normalizeBillingType(billingType) === BILLING_MODE_TOKEN
      ? ['input_price', 'output_price', 'cache_write_price', 'cache_read_price', 'image_output_price']
      : ['per_request_price']

  for (const field of directFields as PriceDescriptor['field'][]) {
    const value = pricing[field]
    if (value != null) {
      return { field, value }
    }
  }

  const intervalCandidate = pricing.intervals
    .map((interval) => intervalPriceDescriptor(interval, billingType))
    .filter((entry): entry is PriceDescriptor => entry !== null)
    .sort((left, right) => left.value - right.value)[0]

  if (intervalCandidate) {
    return intervalCandidate
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
): PriceDescriptor | null {
  if (normalizeBillingType(billingType) !== BILLING_MODE_TOKEN && interval.per_request_price != null) {
    return { field: 'per_request_price', value: interval.per_request_price }
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

function buildPricingLines(pricing: ModelMarketplacePricing | null): PriceSummaryLine[] {
  if (!pricing) return []

  const lines: PriceSummaryLine[] = []
  const directFields: PriceDescriptor['field'][] = [
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
    const descriptor = intervalPriceDescriptor(interval, pricing.billing_mode)
    if (!descriptor) return
    lines.push({
      label: `${t('modelMarketplace.intervalPricing', '区间价格')} ${intervalLabel(interval)}`,
      value: formatFieldPrice(descriptor.field, descriptor.value),
    })
  })

  return lines
}

function intervalLabel(interval: ModelMarketplacePricingInterval): string {
  if (interval.tier_label) return interval.tier_label
  const upper = interval.max_tokens == null ? 'max' : String(interval.max_tokens)
  return `${interval.min_tokens}-${upper}`
}

function formatRate(value: number): string {
  return `${value.toFixed(2)}x`
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
