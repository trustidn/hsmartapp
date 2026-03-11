import { defineStore } from 'pinia'
import { api } from '../lib/api'
import { isOnline } from '../lib/api'
import * as db from '../lib/db'

export const useProductsStore = defineStore('products', {
  state: () => ({
    items: [],
    loading: false,
  }),
  actions: {
    async load(tenantId) {
      if (!tenantId) return
      this.loading = true
      try {
        if (isOnline()) {
          const list = await api.products.list(true)
          const arr = Array.isArray(list) ? list : []
          this.items = arr
          await db.cacheProducts(tenantId, arr)
        } else {
          const cached = await db.getCachedProducts(tenantId)
          this.items = Array.isArray(cached) ? cached : []
        }
      } finally {
        this.loading = false
      }
    },
    async create(tenantId, data) {
      const p = await api.products.create(data)
      this.items.push(p)
      return p
    },
    async refresh(tenantId) {
      await this.load(tenantId)
    },
  },
})
