import { describe, expect, it, beforeEach, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { defineComponent } from 'vue'

import KeysView from '../KeysView.vue'

const {
  listKeys,
  createKey,
  updateKey,
  getAvailableGroups,
  getUserGroupRates,
  getDashboardApiKeysUsage,
  getPublicSettings,
  showSuccess,
  showError,
} = vi.hoisted(() => ({
  listKeys: vi.fn(),
  createKey: vi.fn(),
  updateKey: vi.fn(),
  getAvailableGroups: vi.fn(),
  getUserGroupRates: vi.fn(),
  getDashboardApiKeysUsage: vi.fn(),
  getPublicSettings: vi.fn(),
  showSuccess: vi.fn(),
  showError: vi.fn(),
}))

vi.mock('@/api', () => ({
  keysAPI: {
    list: listKeys,
    create: createKey,
    update: updateKey,
    delete: vi.fn(),
    toggleStatus: vi.fn(),
  },
  authAPI: {
    getPublicSettings,
  },
  usageAPI: {
    getDashboardApiKeysUsage,
  },
  userGroupsAPI: {
    getAvailable: getAvailableGroups,
    getUserGroupRates,
  },
}))

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({ showSuccess, showError }),
}))

vi.mock('@/stores/onboarding', () => ({
  useOnboardingStore: () => ({
    isCurrentStep: vi.fn(() => false),
    nextStep: vi.fn(),
  }),
}))

vi.mock('@/composables/usePersistedPageSize', () => ({
  getPersistedPageSize: () => 10,
}))

vi.mock('@/composables/useClipboard', () => ({
  useClipboard: () => ({ copyToClipboard: vi.fn().mockResolvedValue(true) }),
}))

vi.mock('@/utils/ccswitchImport', () => ({
  buildCcSwitchImportDeeplink: vi.fn(() => 'ccswitch://import'),
}))

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string, params?: Record<string, unknown>) =>
        params ? `${key}:${JSON.stringify(params)}` : key,
    }),
  }
})

const SelectStub = defineComponent({
  name: 'Select',
  props: {
    modelValue: [String, Number, Boolean, null],
    options: { type: Array, default: () => [] },
  },
  emits: ['update:modelValue'],
  template: `
    <select
      class="select-stub"
      :value="modelValue ?? ''"
      @change="$emit('update:modelValue', $event.target.value === '' ? null : Number($event.target.value))"
    >
      <option value=""></option>
      <option v-for="option in options" :key="option.value" :value="option.value">
        {{ option.label }}
      </option>
    </select>
  `,
})

const BaseDialogStub = defineComponent({
  name: 'BaseDialog',
  props: { show: Boolean },
  template: '<div v-if="show" class="dialog-stub"><slot /><slot name="footer" /></div>',
})

const DataTableStub = defineComponent({
  name: 'DataTable',
  props: { data: { type: Array, default: () => [] } },
  template: `
    <div>
      <div v-for="row in data" :key="row.id" class="row-stub">
        <slot name="cell-actions" :row="row" />
      </div>
      <slot v-if="data.length === 0" name="empty" />
    </div>
  `,
})

const EmptyStateStub = defineComponent({
  name: 'EmptyState',
  emits: ['action'],
  template: '<button class="empty-action" @click="$emit(\'action\')">empty action</button>',
})

const mountKeysView = async () => {
  const wrapper = mount(KeysView, {
    global: {
      stubs: {
        AppLayout: { template: '<div><slot /></div>' },
        TablePageLayout: {
          template: '<div><slot name="actions" /><slot name="filters" /><slot name="table" /><slot name="pagination" /></div>',
        },
        DataTable: DataTableStub,
        BaseDialog: BaseDialogStub,
        ConfirmDialog: true,
        EmptyState: EmptyStateStub,
        Select: SelectStub,
        SearchInput: true,
        Pagination: true,
        Icon: true,
        UseKeyModal: true,
        EndpointPopover: true,
        GroupBadge: true,
        GroupOptionItem: true,
        Teleport: true,
      },
    },
  })
  await flushPromises()
  return wrapper
}

const makeApiKey = (overrides: Partial<Record<string, unknown>> = {}) => ({
  id: 1,
  user_id: 7,
  key: 'sk-test',
  name: 'Existing',
  group_id: 1,
  status: 'active',
  ip_whitelist: [],
  ip_blacklist: [],
  last_used_at: null,
  quota: 0,
  quota_used: 0,
  expires_at: null,
  created_at: '2026-01-01T00:00:00Z',
  updated_at: '2026-01-01T00:00:00Z',
  rate_limit_5h: 0,
  rate_limit_1d: 0,
  rate_limit_7d: 0,
  usage_5h: 0,
  usage_1d: 0,
  usage_7d: 0,
  window_5h_start: null,
  window_1d_start: null,
  window_7d_start: null,
  reset_5h_at: null,
  reset_1d_at: null,
  reset_7d_at: null,
  rate_protection_enabled: true,
  max_rate_multiplier: 1.1,
  ...overrides,
})

describe('KeysView rate protection', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    listKeys.mockResolvedValue({ items: [], total: 0, pages: 0 })
    createKey.mockResolvedValue(makeApiKey())
    updateKey.mockResolvedValue(makeApiKey())
    getAvailableGroups.mockResolvedValue([
      {
        id: 1,
        name: 'Default',
        description: '',
        platform: 'openai',
        subscription_type: 'standard',
        rate_multiplier: 1.25,
      },
      {
        id: 2,
        name: 'VIP',
        description: '',
        platform: 'openai',
        subscription_type: 'standard',
        rate_multiplier: 1.5,
      },
    ])
    getUserGroupRates.mockResolvedValue({ 2: 2.75 })
    getDashboardApiKeysUsage.mockResolvedValue({ stats: {} })
    getPublicSettings.mockResolvedValue({ hide_ccs_import_button: true })
  })

  it('defaults max multiplier to effective user group rate after selecting a group and includes protection fields on create', async () => {
    const wrapper = await mountKeysView()

    await wrapper.find('[data-tour="keys-create-btn"]').trigger('click')
    const selects = wrapper.findAll('select.select-stub')
    const groupSelect = selects[2]
    await groupSelect.setValue('2')

    const maxInput = wrapper.find('input[type="number"][step="0.0001"]')
    expect((maxInput.element as HTMLInputElement).value).toBe('2.75')

    await wrapper.find('input[data-tour="key-form-name"]').setValue('Created')
    await wrapper.find('form#key-form').trigger('submit')
    await flushPromises()

    expect(createKey).toHaveBeenCalledWith(
      'Created',
      2,
      undefined,
      [],
      [],
      0,
      undefined,
      { rate_limit_5h: 0, rate_limit_1d: 0, rate_limit_7d: 0 },
      { rate_protection_enabled: true, max_rate_multiplier: 2.75 },
    )
  })

  it('echoes existing protection fields in edit modal and submits updates', async () => {
    listKeys.mockResolvedValue({
      items: [
        makeApiKey({
          id: 9,
          group_id: 1,
          rate_protection_enabled: true,
          max_rate_multiplier: 1.1,
        }),
      ],
      total: 1,
      pages: 1,
    })
    const wrapper = await mountKeysView()

    const editButton = wrapper
      .findAll('.row-stub button')
      .find((button) => button.text().includes('common.edit'))
    expect(editButton).toBeTruthy()
    await editButton!.trigger('click')
    const maxInput = wrapper.find('input[type="number"][step="0.0001"]')
    expect((maxInput.element as HTMLInputElement).value).toBe('1.1')

    await maxInput.setValue('1.8')
    await wrapper.find('form#key-form').trigger('submit')
    await flushPromises()

    expect(updateKey).toHaveBeenCalledWith(9, expect.objectContaining({
      rate_protection_enabled: true,
      max_rate_multiplier: 1.8,
    }))
  })
})
