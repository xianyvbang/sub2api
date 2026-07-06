import { describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import PaymentMethodSelector from '@/components/payment/PaymentMethodSelector.vue'

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string, fallback?: string) => ({
      'payment.paymentMethod': 'Payment method',
      'payment.methods.alipay': 'Alipay',
      'payment.methods.ustd_usdc': 'USTD/USDC',
      'payment.fee': 'Fee',
    }[key] ?? fallback ?? key),
  }),
}))

describe('PaymentMethodSelector', () => {
  it('renders USTD/USDC as its own selectable method while keeping Alipay separate', async () => {
    const wrapper = mount(PaymentMethodSelector, {
      props: {
        selected: 'ustd_usdc',
        methods: [
          { type: 'alipay', fee_rate: 0, available: true },
          { type: 'ustd_usdc', fee_rate: 1.5, available: true },
        ],
      },
    })

    expect(wrapper.text()).toContain('Alipay')
    expect(wrapper.text()).toContain('USTD/USDC')

    const ustdImage = wrapper.find('img[alt="USTD/USDC"]')
    expect(ustdImage.exists()).toBe(true)

    await wrapper.findAll('button').find(button => button.text().includes('USTD/USDC'))?.trigger('click')
    expect(wrapper.emitted('select')?.[0]).toEqual(['ustd_usdc'])
  })

  it('shows the configured display name for custom EasyPay methods', () => {
    const wrapper = mount(PaymentMethodSelector, {
      props: {
        selected: 'ldc',
        methods: [{ type: 'ldc', display_name: 'LDC Pay', fee_rate: 0, available: true }],
      },
    })

    expect(wrapper.text()).toContain('LDC Pay')
    expect(wrapper.text()).not.toContain('ldc')
    expect(wrapper.text()).not.toContain('payment.methods.ldc')
  })

  it('uses the generic selected style for custom methods that contain built-in names', () => {
    const wrapper = mount(PaymentMethodSelector, {
      props: {
        selected: 'card_alipay',
        methods: [{ type: 'card_alipay', display_name: 'Card Pay', fee_rate: 0, available: true }],
      },
    })

    const button = wrapper.get('button')
    expect(button.classes()).toContain('border-primary-500')
    expect(button.classes()).not.toContain('border-[#02A9F1]')
  })
})
