import { defineStore } from 'pinia'
import { api } from '../lib/api'
import { adminApi } from '../lib/adminApi'

/** Store untuk pengaturan SaaS global (logo, dll). Load berbeda untuk admin vs tenant. */
export const useSaasSettingsStore = defineStore('saasSettings', {
  state: () => ({
    logoUrl: '',
    appName: '',
  }),
  actions: {
    async loadForTenant() {
      try {
        const res = await api.saasSettings.get()
        this.logoUrl = res.logo_url || ''
        this.appName = res.app_name || ''
      } catch {
        this.logoUrl = ''
        this.appName = ''
      }
    },
    async loadForAdmin() {
      try {
        const res = await adminApi.saasSettings.get()
        this.logoUrl = res.logo_url || ''
        this.appName = res.app_name || ''
      } catch {
        this.logoUrl = ''
        this.appName = ''
      }
    },
  },
})
