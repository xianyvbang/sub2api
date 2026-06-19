<template>
  <div v-if="homeContent" class="min-h-screen">
    <iframe
      v-if="isHomeContentUrl"
      :src="homeContent.trim()"
      class="h-screen w-full border-0"
      allowfullscreen
    ></iframe>
    <!-- HTML mode - SECURITY: homeContent is admin-only setting, XSS risk is acceptable -->
    <div v-else v-html="homeContent"></div>
  </div>

  <div v-else class="xybbz-home" :class="themeClass">
    <header class="site-header">
      <div class="container">
        <nav class="site-nav">
          <div class="brand reveal" style="--delay: 60ms">
            <div class="brand-mark">
              <img :src="siteLogo || '/logo.png'" :alt="`${siteName} logo`" />
            </div>
            <div class="brand-copy">
              <span class="eyebrow">{{ pageCopy.eyebrow }}</span>
              <h1 class="brand-title">{{ siteName }}</h1>
              <p class="brand-subtitle">{{ siteSubtitle }}</p>
            </div>
          </div>

          <div class="nav-actions reveal" style="--delay: 140ms">
            <div class="nav-tools">
              <div class="tool-shell">
                <LocaleSwitcher />
              </div>

              <a
                v-if="docUrl"
                :href="docUrl"
                target="_blank"
                rel="noopener noreferrer"
                class="icon-control"
                :title="pageCopy.viewDocs"
              >
                <Icon name="book" size="sm" />
              </a>

              <button
                type="button"
                class="icon-control"
                :title="isDark ? pageCopy.switchToLight : pageCopy.switchToDark"
                :aria-label="isDark ? pageCopy.switchToLight : pageCopy.switchToDark"
                @click="toggleTheme"
              >
                <Icon v-if="isDark" name="sun" size="sm" />
                <Icon v-else name="moon" size="sm" />
              </button>
            </div>

            <router-link :to="primaryPath" class="nav-button">
              {{ primaryNavLabel }}
            </router-link>
          </div>
        </nav>
      </div>
    </header>

    <main>
      <section class="hero">
        <div class="container">
          <div class="hero-grid">
            <article class="hero-card reveal" style="--delay: 180ms">
              <span class="hero-badge">{{ heroBadge }}</span>
              <h2 class="hero-title">{{ siteName }}</h2>
              <p class="hero-subtitle">{{ siteSubtitle }}</p>
              <p class="hero-description">{{ pageCopy.heroDescription }}</p>

              <div class="hero-actions">
                <router-link :to="primaryPath" class="cta-button">
                  {{ primaryHeroLabel }}
                  <Icon name="arrowRight" size="sm" :stroke-width="2" />
                </router-link>

                <a
                  v-if="docUrl"
                  :href="docUrl"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="ghost-button"
                >
                  {{ pageCopy.viewDocs }}
                </a>
                <a v-else href="#abilities" class="ghost-button">
                  {{ pageCopy.viewCapabilities }}
                </a>
              </div>

              <div class="hero-notes">
                <div
                  v-for="(note, index) in heroNotes"
                  :key="note"
                  class="hero-note reveal"
                  :style="{ '--delay': `${260 + index * 80}ms` }"
                >
                  <span class="note-dot"></span>
                  {{ note }}
                </div>
              </div>
            </article>

            <aside class="hero-panel reveal" style="--delay: 260ms">
              <div class="console-head">
                <div>
                  <div class="console-label">{{ pageCopy.consoleLabel }}</div>
                  <div class="status-chip">{{ pageCopy.consoleStatus }}</div>
                </div>

                <div class="console-dots" aria-hidden="true">
                  <span></span>
                  <span></span>
                  <span></span>
                </div>
              </div>

              <div class="console-grid">
                <section class="console-card">
                  <div class="label">{{ pageCopy.siteTitleLabel }}</div>
                  <div class="value">{{ siteName }}</div>
                  <div class="desc">{{ pageCopy.siteTitleDescription }}</div>
                </section>

                <section class="console-card">
                  <div class="label">{{ pageCopy.marketplaceLabel }}</div>
                  <div class="value">{{ marketplaceConsoleValue }}</div>
                  <div class="desc">{{ pageCopy.marketplaceDescription }}</div>
                </section>
              </div>

              <div class="console-terminal">
                <div class="terminal-line">
                  <span class="terminal-prompt">$</span>
                  <span class="terminal-code">{{ pageCopy.terminalBrowse }}</span>
                </div>
                <div class="terminal-line">
                  <span class="terminal-muted">{{ pageCopy.terminalHint }}</span>
                </div>
                <div class="terminal-line">
                  <span class="terminal-prompt">$</span>
                  <span class="terminal-code">{{ pageCopy.terminalStart }}</span>
                </div>
                <div class="terminal-line">
                  <span class="terminal-ok">{{ authConsoleText }}</span>
                </div>
              </div>
            </aside>
          </div>
        </div>
      </section>

      <section class="section-block">
        <div class="container">
          <article class="marketplace-strip reveal" style="--delay: 340ms">
            <div class="marketplace-copy">
              <span class="eyebrow">{{ pageCopy.marketplaceEyebrow }}</span>
              <h2>{{ pageCopy.marketplaceTitle }}</h2>
              <p>{{ pageCopy.marketplaceLead }}</p>

              <div class="marketplace-meta">
                <span v-for="item in marketplaceMeta" :key="item" class="meta-pill">
                  {{ item }}
                </span>
              </div>
            </div>

            <div class="marketplace-action">
              <div class="marketplace-status" :class="marketplaceStatusClass">
                {{ marketplaceStatusText }}
              </div>

              <router-link
                v-if="isMarketplaceEnabled"
                :to="marketplacePath"
                class="marketplace-button"
              >
                {{ marketplaceButtonLabel }}
              </router-link>
              <span v-else class="marketplace-button button-disabled">
                {{ marketplaceButtonLabel }}
              </span>
            </div>
          </article>
        </div>
      </section>

      <section id="abilities" class="section-block">
        <div class="container">
          <div class="section-heading reveal" style="--delay: 420ms">
            <div>
              <span class="eyebrow">{{ pageCopy.whyLabel }}</span>
              <h2>{{ pageCopy.featuresTitle }}</h2>
              <p>{{ pageCopy.featuresLead }}</p>
            </div>
          </div>

          <div class="feature-grid">
            <article
              v-for="(feature, index) in features"
              :key="feature.title"
              class="feature-card reveal"
              :style="{ '--delay': `${480 + index * 90}ms` }"
            >
              <div class="feature-icon" :class="feature.variant">{{ feature.index }}</div>
              <h3>{{ feature.title }}</h3>
              <p>{{ feature.description }}</p>
              <span class="feature-highlight">{{ feature.highlight }}</span>
            </article>
          </div>
        </div>
      </section>
    </main>

    <footer class="site-footer">
      <div class="container">
        <div class="footer-panel reveal" style="--delay: 620ms">
          <div class="footer-row">
            <div class="footer-copy">
              <strong>{{ siteName }}</strong>
              <span>{{ pageCopy.footerTagline }}</span>
              <br />
              <span>{{ pageCopy.footerDescription }}</span>
              <br />
              <span>&copy; {{ currentYear }} {{ siteName }}</span>
            </div>

            <div class="footer-links">
              <a
                v-if="docUrl"
                :href="docUrl"
                target="_blank"
                rel="noopener noreferrer"
                class="footer-link"
              >
                {{ pageCopy.viewDocs }}
              </a>

              <a
                v-if="contactInfo && contactHref"
                :href="contactHref"
                :target="isContactExternal ? '_blank' : undefined"
                :rel="isContactExternal ? 'noopener noreferrer' : undefined"
                class="footer-link"
              >
                {{ contactLabel }}
              </a>

              <span v-else-if="contactInfo" class="footer-link">
                {{ contactLabel }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore, useAppStore } from '@/stores'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import Icon from '@/components/icons/Icon.vue'

const { locale } = useI18n()

const authStore = useAuthStore()
const appStore = useAppStore()

const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'Sub2API')
const siteLogo = computed(() => appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '')
const siteSubtitle = computed(() => {
  const configured = appStore.cachedPublicSettings?.site_subtitle?.trim()
  return configured || pageCopy.value.defaultSubtitle
})
const docUrl = computed(() => appStore.cachedPublicSettings?.doc_url || appStore.docUrl || '')
const contactInfo = computed(() => appStore.cachedPublicSettings?.contact_info?.trim() || '')
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')

const isHomeContentUrl = computed(() => {
  const content = homeContent.value.trim()
  return content.startsWith('http://') || content.startsWith('https://')
})

const isDark = ref(document.documentElement.classList.contains('dark'))
const isAuthenticated = computed(() => authStore.isAuthenticated)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => (isAdmin.value ? '/admin/dashboard' : '/dashboard'))
const primaryPath = computed(() => (isAuthenticated.value ? dashboardPath.value : '/login'))
const isZh = computed(() => locale.value.toLowerCase().startsWith('zh'))
const currentYear = computed(() => new Date().getFullYear())
const themeClass = computed(() => (isDark.value ? 'theme-dark' : 'theme-light'))

const isMarketplaceEnabled = computed(
  () => appStore.cachedPublicSettings?.model_marketplace_enabled === true
)
const marketplaceRequiresLogin = computed(
  () => appStore.cachedPublicSettings?.model_marketplace_requires_login !== false
)

const pageCopy = computed(() => {
  if (isZh.value) {
    return {
      eyebrow: 'MODEL RELAY',
      defaultSubtitle: '一条咸鱼开的中转站',
      heroDescription:
        '为模型调用准备的一站式中转首页。你可以在这里快速进入控制台、查看模型广场入口，并更清楚地了解这个站点的统一入口、稳定转发与透明计费能力。',
      viewDocs: '查看文档',
      viewCapabilities: '查看能力',
      switchToLight: '切换浅色模式',
      switchToDark: '切换深色模式',
      login: '进入控制台',
      dashboard: '前往控制台',
      adminDashboard: '前往管理台',
      getStarted: '立即开始',
      goToDashboard: '进入控制台',
      goToAdminDashboard: '进入管理台',
      consoleLabel: 'Control Surface',
      consoleStatus: 'Ready To Explore',
      siteTitleLabel: '站点名称',
      siteTitleDescription: '这里展示当前站点品牌与统一入口，帮助你更快确认访问位置。',
      marketplaceLabel: '模型广场',
      marketplaceDescription: '这里会提示模型广场当前是否可进入，帮助你快速找到模型能力入口。',
      terminalBrowse: '浏览模型广场与站点入口',
      terminalHint: '快速查看文档、能力说明和常用入口',
      terminalStart: '登录后即可开始使用',
      authGuest: '未登录，建议先前往登录页开始使用。',
      authAdmin: '已识别管理员，可直接进入管理台。',
      authUser: '已识别用户，可直接前往控制台。',
      marketplaceEyebrow: 'Model Marketplace',
      marketplaceTitle: '模型广场入口',
      marketplaceLead:
        '快速查看可用模型与能力信息，按需进入模型广场，找到更适合你的使用入口与价格方案。',
      marketplaceAvailabilityMeta: '聚合可用模型入口',
      marketplaceAccessMetaClosed: '当前暂未开放',
      marketplaceAccessMetaPublic: '无需登录即可访问',
      marketplaceAccessMetaPrivate: '登录后进入模型广场',
      marketplaceDocMeta: '按需跳转到文档与入口',
      marketplaceStatusClosed: '模型广场暂未开放',
      marketplaceStatusLocked: '已开放，访问前需先登录',
      marketplaceStatusOpen: '模型广场已开放',
      marketplaceButtonClosed: '模型广场暂未开放',
      marketplaceButtonLocked: '登录后进入模型广场',
      marketplaceButtonOpen: '进入模型广场',
      marketplaceConsoleClosed: '关闭',
      marketplaceConsoleLocked: '登录后访问',
      marketplaceConsoleOpen: '开放',
      whyLabel: 'Why Xybbz',
      featuresTitle: '一个入口，轻松连接常用模型',
      featuresLead:
        '快速进入模型广场、控制台与文档，按需查看模型能力、入口与价格信息，更省心地开始使用。',
      footerTagline: ' · 一条咸鱼开的中转站',
      footerDescription: '在这里快速了解站点能力，并进入你需要的模型与控制台入口。',
      contactPrefix: '联系: '
    }
  }

  return {
    eyebrow: 'MODEL RELAY',
    defaultSubtitle: 'A compact relay station for everyday model access',
    heroDescription:
      'A focused landing page for model routing. Jump into the dashboard, inspect the marketplace entry, and understand the site with a clearer view of unified access, steady relay, and transparent billing.',
    viewDocs: 'View Docs',
    viewCapabilities: 'View Capabilities',
    switchToLight: 'Switch to light mode',
    switchToDark: 'Switch to dark mode',
    login: 'Open Console',
    dashboard: 'Go to Dashboard',
    adminDashboard: 'Go to Admin',
    getStarted: 'Get Started',
    goToDashboard: 'Open Dashboard',
    goToAdminDashboard: 'Open Admin',
    consoleLabel: 'Control Surface',
    consoleStatus: 'Ready To Explore',
    siteTitleLabel: 'Site Title',
    siteTitleDescription: 'This card reflects the current brand and main entry point for the site.',
    marketplaceLabel: 'Marketplace',
    marketplaceDescription: 'This area shows whether the model marketplace is currently reachable.',
    terminalBrowse: 'browse marketplace and entry surfaces',
    terminalHint: 'quick access to docs, capability notes, and common entry points',
    terminalStart: 'sign in to start using the platform',
    authGuest: 'You are not signed in yet. Start from the login page.',
    authAdmin: 'Admin account detected. You can go straight to the admin console.',
    authUser: 'Signed-in user detected. You can go straight to the dashboard.',
    marketplaceEyebrow: 'Model Marketplace',
    marketplaceTitle: 'Marketplace Entry',
    marketplaceLead:
      'Check available models and capability information quickly, then enter the marketplace when it fits your workflow.',
    marketplaceAvailabilityMeta: 'Aggregated model entry surface',
    marketplaceAccessMetaClosed: 'Currently unavailable',
    marketplaceAccessMetaPublic: 'Open without sign-in',
    marketplaceAccessMetaPrivate: 'Marketplace after sign-in',
    marketplaceDocMeta: 'Docs and entry points in one place',
    marketplaceStatusClosed: 'Marketplace is not open yet',
    marketplaceStatusLocked: 'Marketplace is open but requires sign-in',
    marketplaceStatusOpen: 'Marketplace is open',
    marketplaceButtonClosed: 'Marketplace unavailable',
    marketplaceButtonLocked: 'Sign in for marketplace',
    marketplaceButtonOpen: 'Enter Marketplace',
    marketplaceConsoleClosed: 'Closed',
    marketplaceConsoleLocked: 'Login First',
    marketplaceConsoleOpen: 'Open',
    whyLabel: 'Why Xybbz',
    featuresTitle: 'One entry point, easier access to everyday models',
    featuresLead:
      'Reach the marketplace, dashboard, and docs faster, then choose the capability and pricing path that fits best.',
    footerTagline: ' · A compact relay station for everyday model access',
    footerDescription: 'Understand the site quickly and move into the model or console entry you need.',
    contactPrefix: 'Contact: '
  }
})

const heroBadge = computed(() =>
  isZh.value ? `${siteName.value} 统一模型入口` : `${siteName.value} Model Relay Surface`
)

const primaryNavLabel = computed(() => {
  if (!isAuthenticated.value) return pageCopy.value.login
  return isAdmin.value ? pageCopy.value.adminDashboard : pageCopy.value.dashboard
})

const primaryHeroLabel = computed(() => {
  if (!isAuthenticated.value) return pageCopy.value.getStarted
  return isAdmin.value ? pageCopy.value.goToAdminDashboard : pageCopy.value.goToDashboard
})

const heroNotes = computed(() =>
  isZh.value
    ? [
        '聚合多平台模型入口',
        '登录后可快速进入控制台',
        '模型广场统一展示能力入口'
      ]
    : [
        'Unified access across multiple model platforms',
        'Fast dashboard entry after sign-in',
        'Marketplace shown as one clear capability surface'
      ]
)

const features = computed(() =>
  isZh.value
    ? [
        {
          index: '01',
          variant: 'gateway',
          title: '统一入口',
          description:
            '常用模型能力集中展示，少找路径，少切页面，更快进入你需要的使用入口。',
          highlight: '快速开始'
        },
        {
          index: '02',
          variant: 'routing',
          title: '稳定转发',
          description:
            '使用体验更顺滑，入口更清楚，让你把精力放在模型本身，而不是繁琐切换。',
          highlight: '省心稳定'
        },
        {
          index: '03',
          variant: 'billing',
          title: '透明计费',
          description:
            '价格与能力入口更清晰，先看可用模型，再按需选择更适合自己的使用方式。',
          highlight: '价格清晰'
        }
      ]
    : [
        {
          index: '01',
          variant: 'gateway',
          title: 'Unified Entry',
          description:
            'Common model capabilities are grouped in one place so you spend less time hunting for paths.',
          highlight: 'Quick Start'
        },
        {
          index: '02',
          variant: 'routing',
          title: 'Stable Relay',
          description:
            'A clearer surface and smoother flow help you stay focused on the model instead of the plumbing.',
          highlight: 'Reliable Flow'
        },
        {
          index: '03',
          variant: 'billing',
          title: 'Clear Billing',
          description:
            'Capability and pricing entry points stay visible first, so you can choose the right path with less guesswork.',
          highlight: 'Clear Pricing'
        }
      ]
)

const marketplacePath = computed(() => {
  if (!isMarketplaceEnabled.value) return '/model-marketplace'
  if (marketplaceRequiresLogin.value && !isAuthenticated.value) {
    return `/login?redirect=${encodeURIComponent('/model-marketplace')}`
  }
  return '/model-marketplace'
})

const marketplaceStatusText = computed(() => {
  if (!isMarketplaceEnabled.value) return pageCopy.value.marketplaceStatusClosed
  if (marketplaceRequiresLogin.value && !isAuthenticated.value) {
    return pageCopy.value.marketplaceStatusLocked
  }
  return pageCopy.value.marketplaceStatusOpen
})

const marketplaceStatusClass = computed(() => ({
  'is-open': isMarketplaceEnabled.value && (!marketplaceRequiresLogin.value || isAuthenticated.value),
  'is-locked': isMarketplaceEnabled.value && marketplaceRequiresLogin.value && !isAuthenticated.value
}))

const marketplaceButtonLabel = computed(() => {
  if (!isMarketplaceEnabled.value) return pageCopy.value.marketplaceButtonClosed
  if (marketplaceRequiresLogin.value && !isAuthenticated.value) {
    return pageCopy.value.marketplaceButtonLocked
  }
  return pageCopy.value.marketplaceButtonOpen
})

const marketplaceConsoleValue = computed(() => {
  if (!isMarketplaceEnabled.value) return pageCopy.value.marketplaceConsoleClosed
  if (marketplaceRequiresLogin.value && !isAuthenticated.value) {
    return pageCopy.value.marketplaceConsoleLocked
  }
  return pageCopy.value.marketplaceConsoleOpen
})

const marketplaceMeta = computed(() => [
  pageCopy.value.marketplaceAvailabilityMeta,
  !isMarketplaceEnabled.value
    ? pageCopy.value.marketplaceAccessMetaClosed
    : marketplaceRequiresLogin.value
      ? pageCopy.value.marketplaceAccessMetaPrivate
      : pageCopy.value.marketplaceAccessMetaPublic,
  pageCopy.value.marketplaceDocMeta
])

const authConsoleText = computed(() => {
  if (!isAuthenticated.value) return pageCopy.value.authGuest
  return isAdmin.value ? pageCopy.value.authAdmin : pageCopy.value.authUser
})

const contactHref = computed(() => {
  if (!contactInfo.value) return ''
  if (/^https?:\/\//i.test(contactInfo.value)) return contactInfo.value
  if (/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(contactInfo.value)) return `mailto:${contactInfo.value}`
  return ''
})

const isContactExternal = computed(() => /^https?:\/\//i.test(contactInfo.value))
const contactLabel = computed(() => `${pageCopy.value.contactPrefix}${contactInfo.value}`)

let themeObserver: MutationObserver | null = null

function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

function syncThemeFromDocument() {
  isDark.value = document.documentElement.classList.contains('dark')
}

function initTheme() {
  const savedTheme = localStorage.getItem('theme')
  if (
    savedTheme === 'dark' ||
    (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)
  ) {
    isDark.value = true
    document.documentElement.classList.add('dark')
    return
  }

  isDark.value = false
  document.documentElement.classList.remove('dark')
}

onMounted(() => {
  initTheme()
  syncThemeFromDocument()
  authStore.checkAuth()

  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }

  themeObserver = new MutationObserver(() => {
    syncThemeFromDocument()
  })

  themeObserver.observe(document.documentElement, {
    attributes: true,
    attributeFilter: ['class']
  })
})

onBeforeUnmount(() => {
  themeObserver?.disconnect()
  themeObserver = null
})
</script>

<style scoped>
.xybbz-home {
  --line: rgba(55, 65, 81, 0.12);
  --text: #172033;
  --text-soft: #51607d;
  --text-faint: #75819b;
  --accent: #0f8ba8;
  --accent-strong: #0a6f88;
  --warm: #ea8d2b;
  --success: #157347;
  --panel: rgba(255, 255, 255, 0.82);
  --panel-strong: rgba(255, 255, 255, 0.95);
  --shadow-lg: 0 28px 80px rgba(19, 34, 56, 0.14);
  --shadow-md: 0 18px 50px rgba(19, 34, 56, 0.1);
  --radius-xl: 32px;
  --container: 1180px;
  position: relative;
  isolation: isolate;
  min-height: 100vh;
  color: var(--text);
  font-family: 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', sans-serif;
  background:
    radial-gradient(circle at 12% 12%, rgba(15, 139, 168, 0.12), transparent 26%),
    radial-gradient(circle at 84% 16%, rgba(234, 141, 43, 0.1), transparent 24%),
    radial-gradient(circle at 78% 82%, rgba(30, 64, 175, 0.08), transparent 30%),
    linear-gradient(180deg, #fff8f1 0%, #f7f1e8 48%, #eef3f7 100%);
  overflow-x: hidden;
}

.xybbz-home::before {
  content: '';
  position: fixed;
  inset: 0;
  pointer-events: none;
  background-image:
    linear-gradient(rgba(23, 32, 51, 0.03) 1px, transparent 1px),
    linear-gradient(90deg, rgba(23, 32, 51, 0.03) 1px, transparent 1px);
  background-size: 44px 44px;
  mask-image: linear-gradient(180deg, rgba(0, 0, 0, 0.48), transparent 90%);
}

.xybbz-home.theme-dark {
  --line: rgba(148, 163, 184, 0.16);
  --text: #edf4fb;
  --text-soft: #b4c1d5;
  --text-faint: #8ea0b8;
  --accent: #5dcce6;
  --accent-strong: #85e7ff;
  --warm: #f4b24f;
  --success: #46ba7d;
  --panel: rgba(11, 19, 33, 0.82);
  --panel-strong: rgba(11, 19, 33, 0.95);
  --shadow-lg: 0 28px 80px rgba(2, 8, 22, 0.42);
  --shadow-md: 0 18px 50px rgba(2, 8, 22, 0.3);
  background:
    radial-gradient(circle at 12% 12%, rgba(93, 204, 230, 0.12), transparent 26%),
    radial-gradient(circle at 84% 16%, rgba(244, 178, 79, 0.1), transparent 24%),
    radial-gradient(circle at 78% 82%, rgba(59, 130, 246, 0.08), transparent 30%),
    linear-gradient(180deg, #09111d 0%, #0d1726 48%, #111d2f 100%);
}

.xybbz-home.theme-dark::before {
  background-image:
    linear-gradient(rgba(148, 163, 184, 0.05) 1px, transparent 1px),
    linear-gradient(90deg, rgba(148, 163, 184, 0.05) 1px, transparent 1px);
}

.container {
  width: min(var(--container), calc(100vw - 32px));
  margin: 0 auto;
}

.site-header {
  position: sticky;
  top: 0;
  z-index: 40;
  backdrop-filter: blur(16px);
  background: rgba(247, 241, 232, 0.72);
  border-bottom: 1px solid rgba(55, 65, 81, 0.08);
}

.xybbz-home.theme-dark .site-header {
  background: rgba(7, 14, 25, 0.78);
  border-bottom-color: rgba(148, 163, 184, 0.12);
}

.site-nav {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 18px;
  padding: 18px 0;
}

.brand {
  display: flex;
  align-items: center;
  gap: 14px;
  min-width: 0;
}

.brand-mark {
  width: 52px;
  height: 52px;
  border-radius: 18px;
  overflow: hidden;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.96), rgba(244, 248, 252, 0.92));
  border: 1px solid rgba(55, 65, 81, 0.08);
  box-shadow: 0 12px 28px rgba(15, 139, 168, 0.14);
  flex: 0 0 auto;
}

.brand-mark img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  display: block;
}

.brand-copy {
  min-width: 0;
}

.eyebrow {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
  color: var(--text-faint);
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.22em;
  text-transform: uppercase;
}

.eyebrow::before {
  content: '';
  width: 20px;
  height: 1px;
  background: rgba(117, 129, 155, 0.58);
}

.brand-title {
  margin: 0;
  font-family: 'Trebuchet MS', 'Avenir Next', 'PingFang SC', sans-serif;
  font-size: 20px;
  line-height: 1.1;
  font-weight: 800;
  letter-spacing: 0.02em;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.brand-subtitle {
  margin: 5px 0 0;
  color: var(--text-soft);
  font-size: 13px;
}

.nav-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
  flex-wrap: wrap;
}

.nav-tools {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.tool-shell {
  display: inline-flex;
  align-items: center;
  min-height: 44px;
  padding: 4px 8px;
  border-radius: 16px;
  border: 1px solid rgba(55, 65, 81, 0.08);
  background: rgba(255, 255, 255, 0.6);
  backdrop-filter: blur(12px);
}

.icon-control,
.nav-button,
.cta-button,
.ghost-button,
.marketplace-button,
.footer-link {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  transition:
    transform 180ms ease,
    box-shadow 180ms ease,
    border-color 180ms ease,
    background-color 180ms ease,
    color 180ms ease,
    opacity 180ms ease;
}

.icon-control:hover,
.nav-button:hover,
.cta-button:hover,
.ghost-button:hover,
.marketplace-button:hover,
.footer-link:hover {
  transform: translateY(-1px);
}

.icon-control {
  min-height: 44px;
  padding: 10px 14px;
  border-radius: 16px;
  border: 1px solid rgba(55, 65, 81, 0.08);
  background: rgba(255, 255, 255, 0.6);
  color: var(--text-soft);
  cursor: pointer;
  backdrop-filter: blur(12px);
}

.nav-button,
.cta-button,
.marketplace-button {
  border-radius: 999px;
  padding: 12px 20px;
  color: #f6fbff;
  background: linear-gradient(135deg, #173c58, #0f8ba8 75%, #4eb3c9);
  box-shadow: 0 16px 38px rgba(15, 139, 168, 0.24);
}

.ghost-button,
.footer-link {
  border-radius: 999px;
  padding: 12px 18px;
  color: var(--text);
  background: rgba(255, 255, 255, 0.78);
  border: 1px solid rgba(55, 65, 81, 0.12);
}

.xybbz-home.theme-dark .tool-shell,
.xybbz-home.theme-dark .icon-control,
.xybbz-home.theme-dark .ghost-button,
.xybbz-home.theme-dark .footer-link {
  background: rgba(12, 21, 34, 0.76);
  border-color: rgba(148, 163, 184, 0.16);
  color: var(--text-soft);
}

.button-disabled {
  opacity: 0.55;
  cursor: not-allowed;
  pointer-events: none;
  box-shadow: none;
}

.hero {
  padding: 54px 0 26px;
}

.hero-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.14fr) minmax(320px, 0.86fr);
  gap: 28px;
  align-items: stretch;
}

.hero-card,
.hero-panel,
.feature-card,
.marketplace-strip,
.footer-panel {
  position: relative;
  overflow: hidden;
  border-radius: var(--radius-xl);
  border: 1px solid var(--line);
  background: var(--panel);
  box-shadow: var(--shadow-lg);
  backdrop-filter: blur(18px);
}

.hero-card {
  padding: 40px;
}

.hero-card::before,
.hero-panel::before,
.feature-card::before,
.marketplace-strip::before,
.footer-panel::before {
  content: '';
  position: absolute;
  inset: 0;
  pointer-events: none;
  background: linear-gradient(145deg, rgba(255, 255, 255, 0.22), transparent 42%);
}

.hero-badge {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding: 10px 16px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.78);
  border: 1px solid rgba(15, 139, 168, 0.18);
  color: var(--accent-strong);
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.04em;
}

.hero-badge::before {
  content: '';
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--warm), #f3b454);
  box-shadow: 0 0 0 8px rgba(234, 141, 43, 0.14);
}

.hero-title {
  margin: 24px 0 14px;
  font-family: 'Trebuchet MS', 'Avenir Next', 'PingFang SC', sans-serif;
  font-size: clamp(42px, 6vw, 72px);
  line-height: 0.98;
  letter-spacing: -0.04em;
  font-weight: 900;
  max-width: 11ch;
}

.hero-subtitle {
  margin: 0 0 14px;
  color: var(--warm);
  font-size: clamp(18px, 2vw, 24px);
  font-weight: 800;
}

.hero-description {
  max-width: 620px;
  margin: 0;
  color: var(--text-soft);
  font-size: 17px;
  line-height: 1.8;
}

.hero-actions {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 14px;
  margin-top: 28px;
}

.hero-notes {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 24px;
}

.hero-note {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  min-height: 44px;
  padding: 11px 14px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.72);
  border: 1px solid rgba(55, 65, 81, 0.08);
  color: var(--text-soft);
  font-size: 13px;
  font-weight: 600;
}

.note-dot {
  width: 9px;
  height: 9px;
  border-radius: 50%;
  background: var(--accent);
  box-shadow: 0 0 0 8px rgba(15, 139, 168, 0.12);
  flex: 0 0 auto;
}

.hero-panel {
  padding: 28px;
  background:
    linear-gradient(180deg, rgba(19, 34, 56, 0.98), rgba(27, 48, 76, 0.96)),
    linear-gradient(180deg, rgba(255, 255, 255, 0.06), transparent);
  color: #edf4fb;
  box-shadow: 0 30px 70px rgba(19, 34, 56, 0.28);
}

.console-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 18px;
}

.console-label {
  font-size: 11px;
  letter-spacing: 0.2em;
  text-transform: uppercase;
  color: rgba(207, 224, 243, 0.72);
}

.console-dots {
  display: flex;
  gap: 8px;
}

.console-dots span {
  width: 10px;
  height: 10px;
  border-radius: 999px;
}

.console-dots span:nth-child(1) {
  background: #ff7b6b;
}

.console-dots span:nth-child(2) {
  background: #ffcc6e;
}

.console-dots span:nth-child(3) {
  background: #4fd489;
}

.status-chip {
  display: inline-flex;
  align-items: center;
  gap: 9px;
  margin-top: 10px;
  padding: 8px 12px;
  border-radius: 999px;
  color: #bee3ff;
  background: rgba(78, 179, 201, 0.14);
  border: 1px solid rgba(78, 179, 201, 0.2);
  font-size: 12px;
  font-weight: 700;
}

.status-chip::before {
  content: '';
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #60d8f2;
  box-shadow: 0 0 0 8px rgba(96, 216, 242, 0.13);
}

.console-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.console-card {
  padding: 16px;
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.console-card .label {
  font-size: 11px;
  text-transform: uppercase;
  letter-spacing: 0.16em;
  color: rgba(191, 207, 226, 0.62);
}

.console-card .value {
  margin-top: 12px;
  font-size: 28px;
  font-weight: 800;
  letter-spacing: -0.03em;
}

.console-card .desc {
  margin-top: 10px;
  color: rgba(191, 207, 226, 0.78);
  font-size: 13px;
  line-height: 1.6;
}

.console-terminal {
  margin-top: 14px;
  padding: 18px;
  border-radius: 22px;
  background: rgba(8, 18, 33, 0.46);
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.terminal-line {
  display: flex;
  gap: 12px;
  align-items: baseline;
  color: #e5eef7;
  font-family: 'Cascadia Mono', 'Consolas', monospace;
  font-size: 13px;
  line-height: 1.9;
}

.terminal-line + .terminal-line {
  margin-top: 6px;
}

.terminal-prompt {
  color: #63d0ff;
}

.terminal-code {
  color: #f3f7fb;
}

.terminal-muted {
  color: #8faac7;
}

.terminal-ok {
  color: #7cf0b2;
}

.section-block {
  padding: 18px 0;
}

.marketplace-strip {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 24px;
  align-items: center;
  padding: 30px 34px;
  background:
    radial-gradient(circle at 0% 50%, rgba(15, 139, 168, 0.12), transparent 30%),
    linear-gradient(135deg, rgba(255, 255, 255, 0.92), rgba(245, 249, 252, 0.86));
}

.xybbz-home.theme-dark .marketplace-strip {
  background:
    radial-gradient(circle at 0% 50%, rgba(93, 204, 230, 0.12), transparent 30%),
    linear-gradient(135deg, rgba(14, 24, 38, 0.92), rgba(18, 30, 46, 0.88));
}

.marketplace-copy h2 {
  margin: 8px 0 0;
  font-size: clamp(26px, 4vw, 40px);
  line-height: 1.08;
  letter-spacing: -0.04em;
}

.marketplace-copy p {
  margin: 14px 0 0;
  max-width: 760px;
  color: var(--text-soft);
  font-size: 16px;
  line-height: 1.8;
}

.marketplace-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 18px;
}

.meta-pill {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding: 10px 14px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.72);
  border: 1px solid rgba(55, 65, 81, 0.08);
  color: var(--text-soft);
  font-size: 13px;
  font-weight: 700;
}

.xybbz-home.theme-dark .meta-pill {
  background: rgba(12, 21, 34, 0.76);
  border-color: rgba(148, 163, 184, 0.16);
}

.marketplace-action {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 10px;
  min-width: 220px;
}

.marketplace-status {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  min-height: 20px;
  color: var(--text-soft);
  font-size: 13px;
  font-weight: 600;
}

.marketplace-status::before {
  content: '';
  width: 10px;
  height: 10px;
  border-radius: 999px;
  background: var(--warm);
  box-shadow: 0 0 0 8px rgba(234, 141, 43, 0.12);
}

.marketplace-status.is-open::before {
  background: var(--success);
  box-shadow: 0 0 0 8px rgba(21, 115, 71, 0.12);
}

.marketplace-status.is-locked::before {
  background: var(--accent);
  box-shadow: 0 0 0 8px rgba(15, 139, 168, 0.12);
}

.section-heading {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 20px;
  margin: 20px 0 22px;
}

.section-heading h2 {
  margin: 8px 0 0;
  font-size: clamp(28px, 4vw, 42px);
  line-height: 1.08;
  letter-spacing: -0.04em;
}

.section-heading p {
  margin: 12px 0 0;
  max-width: 760px;
  color: var(--text-soft);
  font-size: 16px;
  line-height: 1.8;
}

.feature-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 18px;
}

.feature-card {
  padding: 28px 26px;
  transition:
    transform 220ms ease,
    border-color 220ms ease,
    box-shadow 220ms ease;
}

.feature-card:hover {
  transform: translateY(-4px);
  border-color: rgba(15, 139, 168, 0.22);
  box-shadow: 0 26px 56px rgba(15, 139, 168, 0.12);
}

.feature-icon {
  width: 54px;
  height: 54px;
  border-radius: 18px;
  display: grid;
  place-items: center;
  font-size: 24px;
  font-weight: 800;
  color: #fff;
  box-shadow: 0 16px 28px rgba(15, 139, 168, 0.18);
}

.feature-icon.gateway {
  background: linear-gradient(135deg, #0f8ba8, #4eb3c9);
}

.feature-icon.routing {
  background: linear-gradient(135deg, #2a4b7c, #597daf);
}

.feature-icon.billing {
  background: linear-gradient(135deg, #ea8d2b, #f3b454);
}

.feature-card h3 {
  margin: 18px 0 10px;
  font-size: 22px;
  letter-spacing: -0.03em;
}

.feature-card p {
  margin: 0;
  color: var(--text-soft);
  font-size: 15px;
  line-height: 1.9;
}

.feature-highlight {
  display: inline-flex;
  margin-top: 18px;
  padding: 8px 12px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(55, 65, 81, 0.08);
  color: var(--text-soft);
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.04em;
}

.xybbz-home.theme-dark .feature-highlight,
.xybbz-home.theme-dark .hero-note {
  background: rgba(12, 21, 34, 0.76);
  border-color: rgba(148, 163, 184, 0.16);
  color: var(--text-soft);
}

.site-footer {
  padding: 18px 0 42px;
}

.footer-panel {
  padding: 24px 28px;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.86), rgba(248, 251, 253, 0.78));
}

.xybbz-home.theme-dark .footer-panel {
  background: linear-gradient(180deg, rgba(11, 19, 33, 0.86), rgba(16, 28, 44, 0.78));
}

.footer-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  flex-wrap: wrap;
}

.footer-copy {
  color: var(--text-soft);
  font-size: 14px;
  line-height: 1.8;
}

.footer-copy strong {
  color: var(--text);
}

.footer-links {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.reveal {
  opacity: 0;
  transform: translateY(22px);
  animation: rise-in 520ms ease forwards;
  animation-delay: var(--delay, 0ms);
}

@keyframes rise-in {
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (max-width: 1080px) {
  .hero-grid,
  .feature-grid,
  .marketplace-strip {
    grid-template-columns: 1fr;
  }

  .marketplace-action {
    align-items: flex-start;
    min-width: 0;
  }
}

@media (max-width: 720px) {
  .site-nav,
  .section-heading,
  .footer-row {
    align-items: flex-start;
    flex-direction: column;
  }

  .nav-actions,
  .hero-actions,
  .hero-notes,
  .marketplace-meta,
  .footer-links,
  .nav-tools {
    width: 100%;
  }

  .nav-actions {
    align-items: stretch;
  }

  .tool-shell,
  .icon-control,
  .nav-button,
  .cta-button,
  .ghost-button,
  .marketplace-button {
    width: 100%;
  }

  .hero-card,
  .hero-panel,
  .marketplace-strip,
  .footer-panel {
    padding-left: 22px;
    padding-right: 22px;
  }

  .console-grid {
    grid-template-columns: 1fr;
  }
}
</style>
