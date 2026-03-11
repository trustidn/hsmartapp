import { defineStore } from 'pinia'
import { api } from '../lib/api'

export const useSettingsStore = defineStore('settings', {
  state: () => ({
    settings: null,
    loading: false,
  }),
  getters: {
    merchantName: (s) => s.settings?.name ?? '',
    defaultPayment: (s) => s.settings?.default_payment ?? 'cash',
    receiptFooter: (s) => s.settings?.receipt_footer ?? '',
    whatsappNumber: (s) => s.settings?.whatsapp_number ?? '',
  },
  actions: {
    async load(tenantId) {
      if (!tenantId) return
      this.loading = true
      try {
        this.settings = await api.tenant.getSettings()
      } catch {
        this.settings = null
      } finally {
        this.loading = false
      }
    },
    async update(data) {
      await api.tenant.updateSettings(data)
      if (this.settings) this.settings = { ...this.settings, ...data }
    },
  },
})
