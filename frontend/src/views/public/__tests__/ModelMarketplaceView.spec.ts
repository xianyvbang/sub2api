import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { createMemoryHistory, createRouter } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import ModelMarketplaceView from '../ModelMarketplaceView.vue'
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
          'common.close': '关闭',
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
        group_id: 2,
        group_name: 'Pro',
        group_platform: 'openai',
        group_rate: 0.812345,
        group_is_exclusive: false,
        subscription_type: 'standard',
        model_name: 'gpt-4o',
        platform: 'openai',
        billing_type: 'token',
        pricing_source: 'group',
        original_pricing: {
          billing_mode: 'token',
          input_price: 0.000003,
          output_price: 0.000015,
          cache_write_price: null,
          cache_read_price: null,
          image_output_price: null,
          per_request_price: null,
          intervals: [],
        },
        current_pricing: {
          billing_mode: 'token',
          input_price: 0.0000024,
          output_price: 0.000012,
          cache_write_price: null,
          cache_read_price: null,
          image_output_price: null,
          per_request_price: null,
          intervals: [],
        },
        groups: [
          {
            group_id: 1,
            group_name: 'Public',
            group_platform: 'openai',
            group_rate: 1.234567,
            group_is_exclusive: false,
            subscription_type: 'standard',
            model_name: 'gpt-4o',
            platform: 'openai',
            billing_type: 'token',
            pricing_source: 'group',
            original_pricing: {
              billing_mode: 'token',
              input_price: 0.000003,
              output_price: 0.000015,
              cache_write_price: null,
              cache_read_price: null,
              image_output_price: null,
              per_request_price: null,
              intervals: [],
            },
            current_pricing: {
              billing_mode: 'token',
              input_price: 0.0000036,
              output_price: 0.000018,
              cache_write_price: null,
              cache_read_price: null,
              image_output_price: null,
              per_request_price: null,
              intervals: [],
            },
          },
          {
            group_id: 2,
            group_name: 'Pro',
            group_platform: 'openai',
            group_rate: 0.812345,
            group_is_exclusive: false,
            subscription_type: 'standard',
            model_name: 'gpt-4o',
            platform: 'openai',
            billing_type: 'token',
            pricing_source: 'group',
            original_pricing: {
              billing_mode: 'token',
              input_price: 0.000003,
              output_price: 0.000015,
              cache_write_price: null,
              cache_read_price: null,
              image_output_price: null,
              per_request_price: null,
              intervals: [],
            },
            current_pricing: {
              billing_mode: 'token',
              input_price: 0.0000024,
              output_price: 0.000012,
              cache_write_price: null,
              cache_read_price: null,
              image_output_price: null,
              per_request_price: null,
              intervals: [],
            },
          },
        ],
      },
      {
        group_id: 3,
        group_name: 'Claude',
        group_platform: 'anthropic',
        group_rate: 1.5,
        group_is_exclusive: false,
        subscription_type: 'subscription',
        model_name: 'claude-sonnet-4',
        platform: 'anthropic',
        billing_type: 'per_request',
        pricing_source: 'channel',
        original_pricing: {
          billing_mode: 'per_request',
          input_price: null,
          output_price: null,
          cache_write_price: null,
          cache_read_price: null,
          image_output_price: null,
          per_request_price: 0.02,
          intervals: [],
        },
        current_pricing: {
          billing_mode: 'per_request',
          input_price: null,
          output_price: null,
          cache_write_price: null,
          cache_read_price: null,
          image_output_price: null,
          per_request_price: 0.02,
          intervals: [],
        },
        groups: [
          {
            group_id: 3,
            group_name: 'Claude',
            group_platform: 'anthropic',
            group_rate: 1.5,
            group_is_exclusive: false,
            subscription_type: 'subscription',
            model_name: 'claude-sonnet-4',
            platform: 'anthropic',
            billing_type: 'per_request',
            pricing_source: 'channel',
            original_pricing: {
              billing_mode: 'per_request',
              input_price: null,
              output_price: null,
              cache_write_price: null,
              cache_read_price: null,
              image_output_price: null,
              per_request_price: 0.02,
              intervals: [],
            },
            current_pricing: {
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
        ],
      },
    ]),
  },
}))

describe('ModelMarketplaceView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    i18n.global.locale.value = 'zh'
    copyToClipboardMock.mockReset()
    copyToClipboardMock.mockResolvedValue(true)
  })

  async function mountView() {
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
    return wrapper
  }

  it('filters aggregated cards by model name and group', async () => {
    const wrapper = await mountView()

    expect(wrapper.text()).toContain('gpt-4o')
    expect(wrapper.text()).toContain('claude-sonnet-4')

    const searchInput = wrapper.find('input[type="text"]')
    await searchInput.setValue('claude')
    await flushPromises()

    expect(wrapper.text()).toContain('claude-sonnet-4')
    expect(wrapper.text()).not.toContain('gpt-4o')

    await searchInput.setValue('')
    await flushPromises()

    const selects = wrapper.findAll('select')
    await selects[0].setValue('Public')
    await flushPromises()

    expect(wrapper.text()).toContain('gpt-4o')
    expect(wrapper.text()).not.toContain('claude-sonnet-4')

    const groupedCard = wrapper.get('[data-testid="marketplace-card-openai-gpt-4o"]')
    expect(groupedCard.text()).toContain('Public')
    expect(groupedCard.text()).toContain('当前展示分组')
    expect(groupedCard.text()).toContain('1.234567x')
    expect(groupedCard.text()).toContain('$3.6 / 1M Tokens')
    expect(groupedCard.text()).toContain('$18 / 1M Tokens')
    expect(groupedCard.text()).not.toContain('$2.4 / 1M Tokens')
    expect(groupedCard.text()).not.toContain('$12 / 1M Tokens')
  })

  it('renders one aggregated card, shows original/current prices, and opens drawer for copyable details', async () => {
    const wrapper = await mountView()

    const aggregatedCard = wrapper.get('[data-testid="marketplace-card-openai-gpt-4o"]')
    expect(wrapper.findAll('[data-testid="marketplace-card-openai-gpt-4o"]').length).toBe(1)
    expect(aggregatedCard.text()).toContain('输入价格')
    expect(aggregatedCard.text()).toContain('输出价格')
    expect(aggregatedCard.text()).toContain('原价')
    expect(aggregatedCard.text()).toContain('$2.4 / 1M Tokens')
    expect(aggregatedCard.text()).toContain('$12 / 1M Tokens')
    expect(aggregatedCard.text()).toContain('$3 / 1M Tokens')
    expect(aggregatedCard.text()).toContain('$15 / 1M Tokens')
    expect(aggregatedCard.text()).toContain('0.812345x')
    expect(aggregatedCard.findAll('.line-through').length).toBeGreaterThan(0)

    await aggregatedCard.trigger('click')
    await flushPromises()

    const drawer = wrapper.get('[data-testid="marketplace-detail-drawer"]')
    expect(drawer.text()).toContain('Public')
    expect(drawer.text()).toContain('Pro')
    expect(drawer.text()).toContain('输入价格')
    expect(drawer.text()).toContain('输出价格')
    expect(drawer.text()).toContain('原价')
    expect(drawer.text()).toContain('1.234567x')
    expect(drawer.text()).toContain('0.812345x')
    expect(drawer.findAll('.line-through').length).toBeGreaterThan(0)

    const copyButton = wrapper.get('[data-testid="copy-model-name-gpt-4o"]')
    await copyButton.trigger('click')
    await flushPromises()

    expect(copyToClipboardMock).toHaveBeenCalledWith('gpt-4o', '已复制到剪贴板')
    expect(copyButton.attributes('aria-label')).toBe('已复制')
  })
})
