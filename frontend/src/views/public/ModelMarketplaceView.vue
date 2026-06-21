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
      <nav class="mx-auto flex max-w-7xl items-center justify-between gap-4">
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
      <div class="mx-auto max-w-7xl">
        <section
          class="rounded-[2rem] border border-white/70 bg-white/80 p-8 shadow-[0_25px_80px_-40px_rgba(15,23,42,0.45)] backdrop-blur dark:border-dark-700/70 dark:bg-dark-900/75 dark:shadow-[0_25px_80px_-40px_rgba(0,0,0,0.55)]"
        >
          <div class="flex flex-col gap-8 xl:flex-row xl:items-end xl:justify-between">
            <div class="max-w-3xl">
              <p class="mb-3 text-xs font-semibold uppercase tracking-[0.3em] text-amber-700/80 dark:text-amber-300/80">
                {{ t('nav.modelMarketplace', 'Model Marketplace') }}
              </p>
              <h2 class="text-4xl font-black leading-tight text-slate-950 dark:text-white md:text-5xl">
                {{ t('modelMarketplace.heroTitle', 'Browse available models by provider and group') }}
              </h2>
              <p class="mt-4 max-w-2xl text-base leading-7 text-slate-600 dark:text-dark-300">
                {{
                  t(
                    'modelMarketplace.heroDescription',
                    'Each card shows one group, provider, and model with billing type and explicit pricing.',
                  )
                }}
              </p>
            </div>

            <div class="grid gap-4 sm:grid-cols-3 xl:min-w-[24rem]">
              <div class="rounded-2xl border border-amber-100 bg-amber-50/90 p-4 dark:border-amber-500/20 dark:bg-amber-500/10">
                <p class="text-xs font-semibold uppercase tracking-[0.2em] text-amber-700 dark:text-amber-300">Cards</p>
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
                :placeholder="t('modelMarketplace.searchPlaceholder', 'Search model names')"
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

        <section class="mt-8 grid gap-6 xl:grid-cols-[280px_minmax(0,1fr)]">
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
                    <span class="truncate pr-3">{{ billingType.value }}</span>
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
                    {{ billingType }}
                  </option>
                </select>
              </label>
            </div>
          </aside>

          <div class="min-w-0">
            <section v-if="loading" class="grid gap-5 md:grid-cols-2 2xl:grid-cols-3">
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

            <section v-else class="grid gap-5 md:grid-cols-2 2xl:grid-cols-3">
              <article
                v-for="card in filteredCards"
                :key="`${card.group_id}-${card.platform}-${card.model_name}`"
                class="rounded-[1.75rem] border border-white/70 bg-white/88 p-6 shadow-[0_20px_70px_-45px_rgba(15,23,42,0.5)] backdrop-blur transition hover:-translate-y-1 hover:shadow-[0_30px_80px_-45px_rgba(15,23,42,0.55)] dark:border-dark-700/70 dark:bg-dark-900/82 dark:shadow-[0_20px_70px_-45px_rgba(0,0,0,0.6)] dark:hover:shadow-[0_30px_80px_-45px_rgba(0,0,0,0.7)]"
              >
                <div class="flex items-start justify-between gap-4">
                  <div class="min-w-0">
                    <p class="truncate text-xs font-semibold uppercase tracking-[0.2em] text-slate-400 dark:text-dark-400">
                      {{ card.platform }}
                    </p>
                    <div class="mt-2 flex min-w-0 items-start gap-2">
                      <h3 class="min-w-0 break-words text-[1.125rem] font-bold leading-snug text-slate-950 dark:text-white">
                        {{ card.model_name }}
                      </h3>
                      <button
                        type="button"
                        :data-testid="`copy-model-name-${card.model_name}`"
                        class="shrink-0 rounded-full border border-slate-200/80 bg-white/90 p-1.5 text-slate-500 shadow-sm transition hover:border-slate-300 hover:text-slate-900 dark:border-dark-700 dark:bg-dark-800/90 dark:text-dark-300 dark:hover:border-dark-600 dark:hover:text-white"
                        :title="copiedModelName === card.model_name ? t('common.copied') : t('common.copy')"
                        :aria-label="copiedModelName === card.model_name ? t('common.copied') : t('common.copy')"
                        @click.stop="copyModelName(card.model_name)"
                      >
                        <Icon v-if="copiedModelName === card.model_name" name="check" size="xs" :stroke-width="2" />
                        <Icon v-else name="copy" size="xs" :stroke-width="2" />
                      </button>
                    </div>
                  </div>
                  <span
                    class="shrink-0 rounded-full px-3 py-1 text-xs font-semibold"
                    :class="
                      card.group_is_exclusive
                        ? 'bg-slate-900 text-white dark:bg-primary-500 dark:text-white'
                        : 'bg-emerald-100 text-emerald-700 dark:bg-emerald-500/15 dark:text-emerald-300'
                    "
                  >
                    {{
                      card.group_is_exclusive
                        ? t('modelMarketplace.exclusiveGroup', 'Exclusive')
                        : t('modelMarketplace.publicGroup', 'Public')
                    }}
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
                    {{ card.subscription_type || '-' }}
                  </span>
                  <span
                    class="rounded-full bg-violet-100 px-3 py-1 text-xs font-medium text-violet-800 dark:bg-violet-500/15 dark:text-violet-300"
                  >
                    {{ card.billing_type || t('modelMarketplace.unknownBilling', 'Unknown') }}
                  </span>
                </div>

                <dl class="mt-6 space-y-3 text-sm">
                  <div class="flex items-center justify-between gap-4">
                    <dt class="text-slate-500 dark:text-dark-400">{{ t('availableChannels.columns.platform', 'Platform') }}</dt>
                    <dd class="font-medium text-slate-900 dark:text-white">{{ card.platform }}</dd>
                  </div>
                  <div class="flex items-center justify-between gap-4">
                    <dt class="text-slate-500 dark:text-dark-400">{{ t('modelMarketplace.groupLabel', '分组') }}</dt>
                    <dd class="font-medium text-slate-900 dark:text-white">{{ card.group_name }}</dd>
                  </div>
                  <div class="flex items-center justify-between gap-4">
                    <dt class="text-slate-500 dark:text-dark-400">{{ t('modelMarketplace.groupRate', 'Rate') }}</dt>
                    <dd class="font-medium text-slate-900 dark:text-white">{{ formatRate(card.group_rate) }}</dd>
                  </div>
                </dl>

                <div class="mt-6 rounded-2xl bg-slate-50 p-4 dark:bg-dark-950/70">
                  <p class="text-xs font-semibold uppercase tracking-[0.2em] text-slate-400 dark:text-dark-400">
                    {{ t('modelMarketplace.pricing', 'Pricing') }}
                  </p>
                  <div v-if="card.pricing" class="mt-3 space-y-2 text-sm text-slate-700 dark:text-dark-200">
                    <p v-if="card.pricing.input_price != null">
                      {{ t('modelMarketplace.inputPrice', 'Input') }}: {{ formatTokenPrice(card.pricing.input_price) }}
                    </p>
                    <p v-if="card.pricing.output_price != null">
                      {{ t('modelMarketplace.outputPrice', 'Output') }}: {{ formatTokenPrice(card.pricing.output_price) }}
                    </p>
                    <p v-if="card.pricing.cache_write_price != null">
                      {{ t('modelMarketplace.cacheWritePrice', 'Cache Write') }}:
                      {{ formatTokenPrice(card.pricing.cache_write_price) }}
                    </p>
                    <p v-if="card.pricing.cache_read_price != null">
                      {{ t('modelMarketplace.cacheReadPrice', 'Cache Read') }}:
                      {{ formatTokenPrice(card.pricing.cache_read_price) }}
                    </p>
                    <p v-if="card.pricing.per_request_price != null">
                      {{ t('modelMarketplace.requestPrice', 'Per Request') }}:
                      {{ formatTokenPrice(card.pricing.per_request_price) }}
                    </p>
                    <p v-if="card.pricing.image_output_price != null">
                      {{ t('modelMarketplace.imagePrice', 'Image Output') }}:
                      {{ formatTokenPrice(card.pricing.image_output_price) }}
                    </p>
                    <details
                      v-if="card.pricing.intervals.length > 0"
                      class="group rounded-xl border border-slate-200 bg-white p-3 dark:border-dark-700 dark:bg-dark-900"
                    >
                      <summary class="cursor-pointer list-none text-sm font-medium text-slate-800 dark:text-white">
                        {{ t('modelMarketplace.intervalPricing', 'Tier Pricing') }}
                      </summary>
                      <div class="mt-3 space-y-2 text-xs text-slate-600 dark:text-dark-300">
                        <div
                          v-for="interval in card.pricing.intervals"
                          :key="`${interval.min_tokens}-${interval.max_tokens}-${interval.tier_label || ''}`"
                          class="rounded-xl bg-slate-50 px-3 py-2 dark:bg-dark-950"
                        >
                          <p>{{ describeInterval(interval) }}</p>
                        </div>
                      </div>
                    </details>
                  </div>
                  <p v-else class="mt-3 text-sm text-slate-500 dark:text-dark-400">
                    {{ t('availableChannels.noPricing', 'No pricing configured') }}
                  </p>
                </div>
              </article>
            </section>
          </div>
        </section>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import { useClipboard } from '@/composables/useClipboard'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import { formatScaled } from '@/utils/pricing'
import modelMarketplaceAPI, {
  type ModelMarketplaceCard,
  type ModelMarketplacePricingInterval,
} from '@/api/modelMarketplace'
import { extractApiErrorMessage } from '@/utils/apiError'

type FilterEntry = {
  value: string
  count: number
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

const siteName = computed(() => appStore.siteName || 'Sub2API')
const siteLogo = computed(() => appStore.siteLogo || '')
const isAuthenticated = computed(() => authStore.isAuthenticated)
const dashboardPath = computed(() => (authStore.isAdmin ? '/admin/dashboard' : '/dashboard'))

const groupOptions = computed(() =>
  Array.from(new Set(cards.value.map((card) => card.group_name))).sort((a, b) => a.localeCompare(b)),
)
const platformOptions = computed(() =>
  Array.from(new Set(cards.value.map((card) => card.platform))).sort((a, b) => a.localeCompare(b)),
)
const billingTypeOptions = computed(() =>
  Array.from(new Set(cards.value.map((card) => card.billing_type).filter(Boolean))).sort((a, b) => a.localeCompare(b)),
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
  return cards.value.filter((card) => card.model_name.toLowerCase().includes(query))
})

const filteredCards = computed(() =>
  searchedCards.value.filter((card) => {
    if (selectedGroup.value && card.group_name !== selectedGroup.value) return false
    if (selectedPlatform.value && card.platform !== selectedPlatform.value) return false
    if (selectedBillingType.value && card.billing_type !== selectedBillingType.value) return false
    return true
  }),
)

const groupAllCount = computed(
  () =>
    searchedCards.value.filter((card) => {
      if (selectedPlatform.value && card.platform !== selectedPlatform.value) return false
      if (selectedBillingType.value && card.billing_type !== selectedBillingType.value) return false
      return true
    }).length,
)

const platformAllCount = computed(
  () =>
    searchedCards.value.filter((card) => {
      if (selectedGroup.value && card.group_name !== selectedGroup.value) return false
      if (selectedBillingType.value && card.billing_type !== selectedBillingType.value) return false
      return true
    }).length,
)

const billingTypeAllCount = computed(
  () =>
    searchedCards.value.filter((card) => {
      if (selectedGroup.value && card.group_name !== selectedGroup.value) return false
      if (selectedPlatform.value && card.platform !== selectedPlatform.value) return false
      return true
    }).length,
)

const groupEntries = computed<FilterEntry[]>(() =>
  buildEntries(groupOptions.value, (card) => {
    if (selectedPlatform.value && card.platform !== selectedPlatform.value) return false
    if (selectedBillingType.value && card.billing_type !== selectedBillingType.value) return false
    return true
  }, 'group_name'),
)

const platformEntries = computed<FilterEntry[]>(() =>
  buildEntries(platformOptions.value, (card) => {
    if (selectedGroup.value && card.group_name !== selectedGroup.value) return false
    if (selectedBillingType.value && card.billing_type !== selectedBillingType.value) return false
    return true
  }, 'platform'),
)

const billingTypeEntries = computed<FilterEntry[]>(() =>
  buildEntries(billingTypeOptions.value, (card) => {
    if (selectedGroup.value && card.group_name !== selectedGroup.value) return false
    if (selectedPlatform.value && card.platform !== selectedPlatform.value) return false
    return true
  }, 'billing_type'),
)

function buildEntries(
  options: string[],
  predicate: (card: ModelMarketplaceCard) => boolean,
  field: 'group_name' | 'platform' | 'billing_type',
): FilterEntry[] {
  const countMap = new Map<string, number>()

  searchedCards.value.forEach((card) => {
    if (!predicate(card)) return
    const value = card[field]
    if (!value) return
    countMap.set(value, (countMap.get(value) || 0) + 1)
  })

  return options.map((value) => ({
    value,
    count: countMap.get(value) || 0,
  }))
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

function formatTokenPrice(value: number | null): string {
  if (value == null) return '-'
  return `${formatScaled(value, 1_000_000)} / 1M Tokens`
}

function formatRate(value: number): string {
  return `${value.toFixed(2)}x`
}

function describeInterval(interval: ModelMarketplacePricingInterval): string {
  const upper = interval.max_tokens == null ? 'max' : interval.max_tokens
  const price =
    interval.per_request_price ??
    interval.input_price ??
    interval.output_price ??
    interval.cache_write_price ??
    interval.cache_read_price
  return `${interval.tier_label || `${interval.min_tokens}-${upper}`} - ${formatTokenPrice(price ?? null)}`
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
