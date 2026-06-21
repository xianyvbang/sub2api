import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import ModelMarketplaceView from '../ModelMarketplaceView.vue'
import { createRouter, createMemoryHistory } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import { i18n } from '@/i18n'

const { copyToClipboardMock } = vi.hoisted(() => ({
  copyToClipboardMock: vi.fn().mockResolvedValue(true),
}))

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string, fallback?: string) => {
        const translations: Record<string, string> = {
          'common.copy': '复制',
          'common.copied': '已复制',
          'common.copiedToClipboard': '已复制到剪贴板',
        }
        return translations[key] ?? fallback ?? key
      },
    }),
  }
})

vi.mock('@/composables/useClipboard', () => ({
  useClipboard: () => ({
    copyToClipboard: copyToClipboardMock,
  }),
}))

vi.mock('@/api/modelMarketplace', () => ({
  default: {
    getModelMarketplace: vi.fn().mockResolvedValue([
      {
        group_id: 1,
        group_name: 'Public',
        group_platform: 'openai',
        group_rate: 1.2,
        group_is_exclusive: false,
        subscription_type: 'standard',
        model_name: 'gpt-4o',
        platform: 'openai',
        billing_type: 'per_request',
        pricing: {
          billing_mode: 'per_request',
          input_price: null,
          output_price: null,
          cache_write_price: null,
          cache_read_price: null,
          image_output_price: null,
          per_request_price: 0.02,
          intervals: [],
        },
      },
      {
        group_id: 2,
        group_name: 'Pro',
        group_platform: 'anthropic',
        group_rate: 1.8,
        group_is_exclusive: true,
        subscription_type: 'subscription',
        model_name: 'claude-sonnet-4',
        platform: 'anthropic',
        billing_type: 'token',
        pricing: {
          billing_mode: 'token',
          input_price: 0.000003,
          output_price: 0.000015,
          cache_write_price: null,
          cache_read_price: null,
          image_output_price: null,
          per_request_price: null,
          intervals: [],
        },
      },
    ]),
  },
}))

describe('ModelMarketplaceView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    i18n.global.locale.value = 'zh'
  })

  it('filters cards by model name and group', async () => {
    const router = createRouter({
      history: createMemoryHistory(),
      routes: [
        { path: '/home', component: { template: '<div />' } },
        { path: '/login', component: { template: '<div />' } },
        { path: '/model-marketplace', component: ModelMarketplaceView },
      ],
    })
    await router.push('/model-marketplace')
    await router.isReady()

    const appStore = useAppStore()
    appStore.publicSettingsLoaded = true
    appStore.siteName = 'Sub2API'
    appStore.siteLogo = ''
    appStore.cachedPublicSettings = {
      registration_enabled: false,
      email_verify_enabled: false,
      force_email_on_third_party_signup: false,
      registration_email_suffix_whitelist: [],
      promo_code_enabled: true,
      password_reset_enabled: false,
      invitation_code_enabled: false,
      turnstile_enabled: false,
      turnstile_site_key: '',
      site_name: 'Sub2API',
      site_logo: '',
      site_subtitle: '',
      api_base_url: '',
      contact_info: '',
      doc_url: '',
      home_content: '',
      hide_ccs_import_button: false,
      payment_enabled: false,
      risk_control_enabled: false,
      table_default_page_size: 20,
      table_page_size_options: [10, 20, 50],
      custom_menu_items: [],
      custom_endpoints: [],
      linuxdo_oauth_enabled: false,
      wechat_oauth_enabled: false,
      oidc_oauth_enabled: false,
      oidc_oauth_provider_name: 'OIDC',
      github_oauth_enabled: false,
      google_oauth_enabled: false,
      backend_mode_enabled: false,
      version: '',
      balance_low_notify_enabled: false,
      account_quota_notify_enabled: false,
      balance_low_notify_threshold: 0,
      channel_monitor_enabled: true,
      channel_monitor_default_interval_seconds: 60,
      available_channels_enabled: false,
      model_marketplace_enabled: true,
      model_marketplace_requires_login: false,
      service_quota_enabled: false,
      affiliate_enabled: false,
    } as any

    const authStore = useAuthStore()
    authStore.$patch({
      token: null,
      user: null,
    } as never)

    const wrapper = mount(ModelMarketplaceView, {
      global: {
        plugins: [router, i18n],
        stubs: {
          RouterLink: {
            props: ['to'],
            template: '<a><slot /></a>',
          },
        },
      },
    })

    await flushPromises()

    expect(wrapper.text()).toContain('gpt-4o')
    expect(wrapper.text()).toContain('claude-sonnet-4')

    const searchInput = wrapper.find('input[type="text"]')
    await searchInput.setValue('claude')
    await flushPromises()
    expect(wrapper.text()).toContain('claude-sonnet-4')
    expect(wrapper.text()).not.toContain('gpt-4o')

    const selects = wrapper.findAll('select')
    await selects[0].setValue('Pro')
    await flushPromises()
    expect(wrapper.text()).toContain('claude-sonnet-4')
  })

  it('shows token prices per 1M tokens and supports copying model names', async () => {
    const router = createRouter({
      history: createMemoryHistory(),
      routes: [
        { path: '/home', component: { template: '<div />' } },
        { path: '/login', component: { template: '<div />' } },
        { path: '/model-marketplace', component: ModelMarketplaceView },
      ],
    })
    await router.push('/model-marketplace')
    await router.isReady()

    const appStore = useAppStore()
    appStore.publicSettingsLoaded = true
    appStore.siteName = 'Sub2API'
    appStore.siteLogo = ''

    const authStore = useAuthStore()
    authStore.$patch({
      token: null,
      user: null,
    } as never)

    const wrapper = mount(ModelMarketplaceView, {
      global: {
        plugins: [router, i18n],
        stubs: {
          RouterLink: {
            props: ['to'],
            template: '<a><slot /></a>',
          },
        },
      },
    })

    await flushPromises()

    expect(wrapper.text()).toContain('$3 / 1M Tokens')
    expect(wrapper.text()).toContain('分组')

    const copyButton = wrapper.find('[data-testid="copy-model-name-gpt-4o"]')
    expect(copyButton.exists()).toBe(true)
    await copyButton.trigger('click')
    await flushPromises()
    expect(copyToClipboardMock).toHaveBeenCalledWith('gpt-4o', '已复制到剪贴板')
    expect(copyButton.attributes('aria-label')).toBe('已复制')
  })
})
