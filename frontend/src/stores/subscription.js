import { defineStore } from 'pinia'
import { api } from '../lib/api'

export const useSubscriptionStore = defineStore('subscription', {
  state: () => ({
    data: null,
  }),
  getters: {
    isExpired: (s) => {
      const exp = s.data?.expired_at
      if (!exp) return false
      const end = new Date(exp)
      return !Number.isNaN(end.getTime()) && end < new Date()
    },
    plan: (s) => s.data?.plan || 'free',
  },
  actions: {
    async load(tenantId) {
      if (!tenantId) return
      try {
        this.data = await api.subscription.get()
      } catch {
        this.data = null
      }
    },
  },
})
