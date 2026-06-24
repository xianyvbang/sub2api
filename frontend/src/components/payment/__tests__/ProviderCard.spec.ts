import { describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import ProviderCard from '@/components/payment/ProviderCard.vue'
import type { ProviderInstance } from '@/types/payment'

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string, fallback?: string) => ({
      'admin.settings.payment.providerAlipay': 'Alipay Direct',
      'admin.settings.payment.modeQRCode': 'QR Code',
      'admin.settings.payment.noSupportedTypesSelected': 'No payment types selected',
      'admin.settings.payment.refundEnabled': 'Allow Refund',
      'admin.settings.payment.allowUserRefund': 'Allow User Refund',
      'common.enabled': 'Enabled',
      'common.edit': 'Edit',
      'common.delete': 'Delete',
      'payment.methods.alipay': 'Alipay',
      'payment.methods.ustd_usdc': 'USTD/USDC',
    }[key] ?? fallback ?? key),
  }),
}))

function providerFactory(overrides: Partial<ProviderInstance> = {}): ProviderInstance {
  return {
    id: 1,
    provider_key: 'alipay',
    name: 'Official Alipay',
    config: {},
    supported_types: [],
    enabled: true,
    payment_mode: '',
    refund_enabled: false,
    allow_user_refund: false,
    limits: '',
    sort_order: 0,
    ...overrides,
  }
}

function mountCard(options: {
  provider?: Partial<ProviderInstance>
  enabled?: boolean
} = {}) {
  return mount(ProviderCard, {
    props: {
      provider: providerFactory(options.provider),
      enabled: options.enabled ?? true,
      availableTypes: [
        { value: 'alipay', label: 'Alipay' },
        { value: 'ustd_usdc', label: 'USTD/USDC' },
      ],
    },
  })
}

describe('ProviderCard', () => {
  it('keeps an empty supported type provider visible and recoverable', async () => {
    const wrapper = mountCard()

    expect(wrapper.text()).toContain('Official Alipay')
    expect(wrapper.text()).toContain('No payment types selected')
    expect(wrapper.text()).toContain('Alipay')
    expect(wrapper.text()).toContain('USTD/USDC')

    const alipayButton = wrapper
      .findAll('button')
      .find(button => button.text().includes('Alipay'))
    await alipayButton?.trigger('click')

    expect(wrapper.emitted('toggleType')?.[0]).toEqual(['alipay'])
  })

  it('leaves edit and delete actions usable when the provider type is disabled', async () => {
    const wrapper = mountCard({ enabled: false })

    expect(wrapper.classes()).not.toContain('opacity-50')
    expect(wrapper.find('.px-4.py-2\\.5').classes()).not.toContain('pointer-events-none')

    const buttons = wrapper.findAll('button')
    await buttons.find(button => button.text().includes('Edit'))?.trigger('click')
    await buttons.find(button => button.text().includes('Delete'))?.trigger('click')

    expect(wrapper.emitted('edit')).toHaveLength(1)
    expect(wrapper.emitted('delete')).toHaveLength(1)
  })
})
