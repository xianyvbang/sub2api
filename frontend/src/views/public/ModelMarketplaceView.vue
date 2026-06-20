<template>
  <div class="relative min-h-screen overflow-hidden bg-[radial-gradient(circle_at_top_left,_rgba(245,158,11,0.14),_transparent_28%),radial-gradient(circle_at_bottom_right,_rgba(14,165,233,0.16),_transparent_32%),linear-gradient(180deg,_#fffaf0_0%,_#f8fafc_40%,_#eef2ff_100%)] text-gray-900">
    <div class="pointer-events-none absolute inset-0 bg-[linear-gradient(rgba(15,23,42,0.04)_1px,transparent_1px),linear-gradient(90deg,rgba(15,23,42,0.04)_1px,transparent_1px)] bg-[size:30px_30px]" />

    <header class="relative z-10 px-6 py-5">
      <nav class="mx-auto flex max-w-7xl items-center justify-between gap-4">
        <router-link to="/home" class="flex items-center gap-3">
          <div class="h-11 w-11 overflow-hidden rounded-2xl bg-white shadow-lg shadow-amber-100">
            <img :src="siteLogo || '/logo.png'" alt="Logo" class="h-full w-full object-contain" />
          </div>
          <div>
            <p class="text-xs font-semibold uppercase tracking-[0.24em] text-amber-700/80">
              {{ t('nav.modelMarketplace', '模型广场') }}
            </p>
            <h1 class="text-lg font-semibold text-slate-900">{{ siteName }}</h1>
          </div>
        </router-link>

        <div class="flex items-center gap-3">
          <router-link
            v-if="isAuthenticated"
            :to="dashboardPath"
            class="rounded-full border border-slate-200/80 bg-white/90 px-4 py-2 text-sm font-medium text-slate-700 shadow-sm backdrop-blur transition hover:border-slate-300 hover:text-slate-900"
          >
            {{ t('home.dashboard') }}
          </router-link>
          <router-link
            v-else
            to="/login"
            class="rounded-full bg-slate-900 px-4 py-2 text-sm font-medium text-white shadow-lg shadow-slate-300 transition hover:bg-slate-800"
          >
            {{ t('home.login') }}
          </router-link>
        </div>
      </nav>
    </header>

    <main class="relative z-10 px-6 pb-16 pt-6">
      <div class="mx-auto max-w-7xl">
        <section class="mb-10 grid gap-8 lg:grid-cols-[1.15fr_0.85fr]">
          <div class="rounded-[2rem] border border-white/70 bg-white/85 p-8 shadow-[0_25px_80px_-40px_rgba(15,23,42,0.45)] backdrop-blur">
            <p class="mb-3 text-xs font-semibold uppercase tracking-[0.3em] text-amber-700/80">{{ t('nav.modelMarketplace', '模型广场') }}</p>
            <h2 class="max-w-3xl text-4xl font-black leading-tight text-slate-950 md:text-5xl">
              {{ t('modelMarketplace.heroTitle', '按模型、供应商和分组快速浏览可用能力') }}
            </h2>
            <p class="mt-4 max-w-2xl text-base leading-7 text-slate-600">
              {{ t('modelMarketplace.heroDescription', '每张卡片对应一个分组 + 供应商 + 模型，展示可用令牌分组、计费类型和明确价格。') }}
            </p>

            <div class="mt-8 grid gap-4 sm:grid-cols-3">
              <div class="rounded-2xl border border-amber-100 bg-amber-50/80 p-4">
                <p class="text-xs font-semibold uppercase tracking-[0.2em] text-amber-700">Cards</p>
                <p class="mt-2 text-3xl font-bold text-slate-900">{{ filteredCards.length }}</p>
              </div>
              <div class="rounded-2xl border border-sky-100 bg-sky-50/80 p-4">
                <p class="text-xs font-semibold uppercase tracking-[0.2em] text-sky-700">{{ t('common.groups', '分组') }}</p>
                <p class="mt-2 text-3xl font-bold text-slate-900">{{ groupOptions.length }}</p>
              </div>
              <div class="rounded-2xl border border-violet-100 bg-violet-50/80 p-4">
                <p class="text-xs font-semibold uppercase tracking-[0.2em] text-violet-700">{{ t('availableChannels.columns.platform', '供应商') }}</p>
                <p class="mt-2 text-3xl font-bold text-slate-900">{{ platformOptions.length }}</p>
              </div>
            </div>
          </div>

          <div class="rounded-[2rem] border border-slate-200/70 bg-slate-950 p-8 text-slate-50 shadow-[0_25px_80px_-40px_rgba(15,23,42,0.65)]">
            <p class="text-xs font-semibold uppercase tracking-[0.24em] text-slate-400">{{ t('common.filter', '筛选') }}</p>
            <div class="mt-5 space-y-4">
              <div>
                <label class="mb-2 block text-sm font-medium text-slate-200">{{ t('common.search', '搜索') }}</label>
                <input
                  v-model="searchQuery"
                  type="text"
                  :placeholder="t('modelMarketplace.searchPlaceholder', '搜索模型名称')"
                  class="w-full rounded-2xl border border-slate-700 bg-slate-900 px-4 py-3 text-sm text-white placeholder:text-slate-500 outline-none transition focus:border-amber-400"
                />
              </div>

              <div class="grid gap-4 sm:grid-cols-3">
                <div>
                  <label class="mb-2 block text-sm font-medium text-slate-200">{{ t('common.groups', '分组') }}</label>
                  <select v-model="selectedGroup" class="w-full rounded-2xl border border-slate-700 bg-slate-900 px-4 py-3 text-sm text-white outline-none transition focus:border-amber-400">
                    <option value="">{{ t('common.all', '全部') }}</option>
                    <option v-for="group in groupOptions" :key="group" :value="group">{{ group }}</option>
                  </select>
                </div>

                <div>
                  <label class="mb-2 block text-sm font-medium text-slate-200">{{ t('availableChannels.columns.platform', '供应商') }}</label>
                  <select v-model="selectedPlatform" class="w-full rounded-2xl border border-slate-700 bg-slate-900 px-4 py-3 text-sm text-white outline-none transition focus:border-amber-400">
                    <option value="">{{ t('common.all', '全部') }}</option>
                    <option v-for="platform in platformOptions" :key="platform" :value="platform">{{ platform }}</option>
                  </select>
                </div>

                <div>
                  <label class="mb-2 block text-sm font-medium text-slate-200">{{ t('modelMarketplace.billingType', '计费类型') }}</label>
                  <select v-model="selectedBillingType" class="w-full rounded-2xl border border-slate-700 bg-slate-900 px-4 py-3 text-sm text-white outline-none transition focus:border-amber-400">
                    <option value="">{{ t('common.all', '全部') }}</option>
                    <option v-for="billingType in billingTypeOptions" :key="billingType" :value="billingType">{{ billingType }}</option>
                  </select>
                </div>
              </div>

              <button
                @click="loadMarketplace"
                :disabled="loading"
                class="inline-flex items-center rounded-full bg-amber-400 px-5 py-2.5 text-sm font-semibold text-slate-950 transition hover:bg-amber-300 disabled:cursor-not-allowed disabled:opacity-60"
              >
                {{ loading ? t('common.loading', '加载中...') : t('common.refresh', '刷新') }}
              </button>
            </div>
          </div>
        </section>

        <section v-if="loading" class="grid gap-5 md:grid-cols-2 xl:grid-cols-3">
          <div v-for="index in 6" :key="index" class="h-72 animate-pulse rounded-[1.75rem] border border-white/70 bg-white/70" />
        </section>

        <section v-else-if="filteredCards.length === 0" class="rounded-[2rem] border border-dashed border-slate-300 bg-white/70 px-8 py-20 text-center shadow-sm">
          <h3 class="text-2xl font-bold text-slate-900">{{ t('modelMarketplace.emptyTitle', '暂无匹配模型') }}</h3>
          <p class="mt-3 text-sm text-slate-500">{{ t('modelMarketplace.emptyDescription', '调整搜索词或筛选条件后再试。') }}</p>
        </section>

        <section v-else class="grid gap-5 md:grid-cols-2 xl:grid-cols-3">
          <article
            v-for="card in filteredCards"
            :key="`${card.group_id}-${card.platform}-${card.model_name}`"
            class="rounded-[1.75rem] border border-white/70 bg-white/88 p-6 shadow-[0_20px_70px_-45px_rgba(15,23,42,0.5)] backdrop-blur transition hover:-translate-y-1 hover:shadow-[0_30px_80px_-45px_rgba(15,23,42,0.55)]"
          >
            <div class="flex items-start justify-between gap-4">
              <div>
                <p class="text-xs font-semibold uppercase tracking-[0.2em] text-slate-400">{{ card.platform }}</p>
                <h3 class="mt-2 text-2xl font-bold text-slate-950">{{ card.model_name }}</h3>
              </div>
              <span class="rounded-full px-3 py-1 text-xs font-semibold" :class="card.group_is_exclusive ? 'bg-slate-900 text-white' : 'bg-emerald-100 text-emerald-700'">
                {{ card.group_is_exclusive ? t('modelMarketplace.exclusiveGroup', '专属') : t('modelMarketplace.publicGroup', '公开') }}
              </span>
            </div>

            <div class="mt-5 flex flex-wrap gap-2">
              <span class="rounded-full bg-amber-100 px-3 py-1 text-xs font-medium text-amber-800">{{ card.group_name }}</span>
              <span class="rounded-full bg-sky-100 px-3 py-1 text-xs font-medium text-sky-800">{{ card.subscription_type || '-' }}</span>
              <span class="rounded-full bg-violet-100 px-3 py-1 text-xs font-medium text-violet-800">{{ card.billing_type || t('modelMarketplace.unknownBilling', '未标注') }}</span>
            </div>

            <dl class="mt-6 space-y-3 text-sm">
              <div class="flex items-center justify-between gap-4">
                <dt class="text-slate-500">{{ t('availableChannels.columns.platform', '供应商') }}</dt>
                <dd class="font-medium text-slate-900">{{ card.platform }}</dd>
              </div>
              <div class="flex items-center justify-between gap-4">
                <dt class="text-slate-500">{{ t('common.groups', '分组') }}</dt>
                <dd class="font-medium text-slate-900">{{ card.group_name }}</dd>
              </div>
              <div class="flex items-center justify-between gap-4">
                <dt class="text-slate-500">{{ t('modelMarketplace.groupRate', '倍率') }}</dt>
                <dd class="font-medium text-slate-900">{{ formatRate(card.group_rate) }}</dd>
              </div>
            </dl>

            <div class="mt-6 rounded-2xl bg-slate-50 p-4">
              <p class="text-xs font-semibold uppercase tracking-[0.2em] text-slate-400">{{ t('modelMarketplace.pricing', '价格') }}</p>
              <div v-if="card.pricing" class="mt-3 space-y-2 text-sm text-slate-700">
                <p v-if="card.pricing.input_price != null">{{ t('modelMarketplace.inputPrice', '输入') }}: {{ formatPrice(card.pricing.input_price) }}</p>
                <p v-if="card.pricing.output_price != null">{{ t('modelMarketplace.outputPrice', '输出') }}: {{ formatPrice(card.pricing.output_price) }}</p>
                <p v-if="card.pricing.cache_write_price != null">{{ t('modelMarketplace.cacheWritePrice', '缓存写入') }}: {{ formatPrice(card.pricing.cache_write_price) }}</p>
                <p v-if="card.pricing.cache_read_price != null">{{ t('modelMarketplace.cacheReadPrice', '缓存读取') }}: {{ formatPrice(card.pricing.cache_read_price) }}</p>
                <p v-if="card.pricing.per_request_price != null">{{ t('modelMarketplace.requestPrice', '按次') }}: {{ formatPrice(card.pricing.per_request_price) }}</p>
                <p v-if="card.pricing.image_output_price != null">{{ t('modelMarketplace.imagePrice', '图片输出') }}: {{ formatPrice(card.pricing.image_output_price) }}</p>
                <details v-if="card.pricing.intervals.length > 0" class="group rounded-xl border border-slate-200 bg-white p-3">
                  <summary class="cursor-pointer list-none text-sm font-medium text-slate-800">
                    {{ t('modelMarketplace.intervalPricing', '区间价格') }}
                  </summary>
                  <div class="mt-3 space-y-2 text-xs text-slate-600">
                    <div v-for="interval in card.pricing.intervals" :key="`${interval.min_tokens}-${interval.max_tokens}-${interval.tier_label || ''}`" class="rounded-xl bg-slate-50 px-3 py-2">
                      <p>{{ describeInterval(interval) }}</p>
                    </div>
                  </div>
                </details>
              </div>
              <p v-else class="mt-3 text-sm text-slate-500">{{ t('availableChannels.noPricing', '未配置价格') }}</p>
            </div>
          </article>
        </section>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import modelMarketplaceAPI, { type ModelMarketplaceCard, type ModelMarketplacePricingInterval } from '@/api/modelMarketplace'
import { extractApiErrorMessage } from '@/utils/apiError'

const { t } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()

const loading = ref(false)
const cards = ref<ModelMarketplaceCard[]>([])
const searchQuery = ref('')
const selectedGroup = ref('')
const selectedPlatform = ref('')
const selectedBillingType = ref('')

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

const filteredCards = computed(() => {
  const query = searchQuery.value.trim().toLowerCase()
  return cards.value.filter((card) => {
    if (query && !card.model_name.toLowerCase().includes(query)) return false
    if (selectedGroup.value && card.group_name !== selectedGroup.value) return false
    if (selectedPlatform.value && card.platform !== selectedPlatform.value) return false
    if (selectedBillingType.value && card.billing_type !== selectedBillingType.value) return false
    return true
  })
})

function formatPrice(value: number | null): string {
  if (value == null) return '-'
  return `$${value.toFixed(value < 0.01 ? 6 : 4)}`
}

function formatRate(value: number): string {
  return `${value.toFixed(2)}x`
}

function describeInterval(interval: ModelMarketplacePricingInterval): string {
  const upper = interval.max_tokens == null ? '∞' : interval.max_tokens
  const price =
    interval.per_request_price ??
    interval.input_price ??
    interval.output_price ??
    interval.cache_write_price ??
    interval.cache_read_price
  return `${interval.tier_label || `${interval.min_tokens}-${upper}`} · ${formatPrice(price ?? null)}`
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
