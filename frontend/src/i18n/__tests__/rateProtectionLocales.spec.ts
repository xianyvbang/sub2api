import { describe, expect, it } from 'vitest'

import en from '../locales/en'
import zh from '../locales/zh'

const rateProtectionKeys = [
  'rateProtection',
  'maxRateMultiplier',
  'maxRateMultiplierPlaceholder',
  'maxRateMultiplierRequired',
  'rateProtectionHint',
] as const

describe.each([
  ['en', en],
  ['zh', zh],
])('API key rate protection locale keys (%s)', (_locale, messages) => {
  it.each(rateProtectionKeys)('contains keys.%s', (key) => {
    expect(messages.keys[key]).toBeTruthy()
  })
})
