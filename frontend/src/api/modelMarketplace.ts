import { apiClient } from './client'
import type { BillingMode } from '@/constants/channel'

export interface ModelMarketplacePricingInterval {
  min_tokens: number
  max_tokens: number | null
  tier_label?: string
  input_price: number | null
  output_price: number | null
  cache_write_price: number | null
  cache_read_price: number | null
  per_request_price: number | null
}

export interface ModelMarketplacePricing {
  billing_mode: BillingMode | string
  input_price: number | null
  output_price: number | null
  cache_write_price: number | null
  cache_read_price: number | null
  image_output_price: number | null
  per_request_price: number | null
  intervals: ModelMarketplacePricingInterval[]
}

export interface ModelMarketplaceCard {
  group_id: number
  group_name: string
  group_platform: string
  group_rate: number
  group_is_exclusive: boolean
  subscription_type: string
  model_name: string
  platform: string
  billing_type: string
  pricing: ModelMarketplacePricing | null
}

export async function getModelMarketplace(options?: { signal?: AbortSignal }): Promise<ModelMarketplaceCard[]> {
  const { data } = await apiClient.get<ModelMarketplaceCard[]>('/model-marketplace', {
    signal: options?.signal,
  })
  return data
}

export const modelMarketplaceAPI = {
  getModelMarketplace,
}

export default modelMarketplaceAPI
