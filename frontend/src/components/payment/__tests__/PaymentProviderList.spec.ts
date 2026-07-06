import { describe, expect, it, vi } from 'vitest'
import { defineComponent } from 'vue'
import { mount } from '@vue/test-utils'
import PaymentProviderList from '@/components/payment/PaymentProviderList.vue'
import type { ProviderInstance } from '@/types/payment'

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string, fallback?: string) => ({
      'admin.settings.payment.providerManagement': 'Provider Management',
      'admin.settings.payment.providerManagementDesc': 'Manage providers',
      'admin.settings.payment.createProvider': 'Create Provider',
      'admin.settings.payment.noProviders': 'No providers',
      'admin.settings.payment.enableTypesFirst': 'Enable types first',
      'common.refresh': 'Refresh',
      'payment.methods.alipay': 'Alipay',
      'payment.methods.wxpay': 'WeChat Pay',
      'payment.methods.ustd_usdc': 'USTD/USDC',
    }[key] ?? fallback ?? key),
  }),
}))

function providerFactory(overrides: Partial<ProviderInstance> = {}): ProviderInstance {
  return {
    id: 1,
    provider_key: 'easypay',
    name: 'EasyPay',
    config: {},
    supported_types: [],
    enabled: true,
    payment_mode: 'qrcode',
    refund_enabled: false,
    allow_user_refund: false,
    limits: '',
    sort_order: 0,
    ...overrides,
  }
}

function mountList(provider: ProviderInstance) {
  return mount(PaymentProviderList, {
    props: {
      providers: [provider],
      loading: false,
      canCreate: true,
      enabledPaymentTypes: ['easypay'],
      allPaymentTypes: [
        { value: 'alipay', label: 'Alipay' },
        { value: 'wxpay', label: 'WeChat Pay' },
        { value: 'ustd_usdc', label: 'USTD/USDC' },
      ],
      redirectLabel: 'Redirect',
    },
    global: {
      stubs: {
        Icon: true,
        VueDraggable: defineComponent({
          name: 'VueDraggable',
          template: '<div><slot /></div>',
        }),
        ProviderCard: defineComponent({
          name: 'ProviderCard',
          props: {
            provider: { type: Object, required: true },
            enabled: { type: Boolean, required: true },
            availableTypes: { type: Array, required: true },
          },
          template: `
            <div class="provider-card">
              <span
                v-for="type in availableTypes"
                :key="type.value"
                class="available-type"
              >{{ type.label }}|{{ type.value }}</span>
            </div>
          `,
        }),
      },
    },
  })
}

describe('PaymentProviderList', () => {
  it('shows EasyPay custom method display names on provider cards', () => {
    const wrapper = mountList(providerFactory({
      config: {
        customMethods: '[{"type":"ldc","upstreamType":"epay","displayName":"LDC"}]',
      },
      supported_types: ['alipay', 'ldc'],
    }))

    expect(wrapper.text()).toContain('Alipay|alipay')
    expect(wrapper.text()).toContain('LDC|ldc')
  })

  it('keeps saved EasyPay custom supported types visible without display metadata', () => {
    const wrapper = mountList(providerFactory({
      supported_types: ['alipay', 'usdt_trc20'],
    }))

    expect(wrapper.text()).toContain('Alipay|alipay')
    expect(wrapper.text()).toContain('usdt_trc20|usdt_trc20')
  })
})
