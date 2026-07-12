import { describe, expect, it } from 'vitest'
import en from '@/i18n/locales/en'
import zh from '@/i18n/locales/zh'

describe.each([
  ['en', en],
  ['zh', zh],
])('payment locale keys (%s)', (_locale, messages) => {
  it('contains provider empty-state and USDT/USDC labels', () => {
    expect(messages.admin.settings.payment.noSupportedTypesSelected).toBeTruthy()
    expect(messages.payment.methods.ustd_usdc).toBe('USDT/USDC')
  })
})
