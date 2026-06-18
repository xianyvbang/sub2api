import { describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import PaymentMethodSelector from '../PaymentMethodSelector.vue'

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string) => ({
      'payment.paymentMethod': 'Payment method',
      'payment.methods.alipay': 'Alipay',
      'payment.methods.ustd_usdc': 'USTD/USDC',
      'payment.fee': 'Fee',
    }[key] ?? key),
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
})
