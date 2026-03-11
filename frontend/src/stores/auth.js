import { defineStore } from 'pinia'
import { api } from '../lib/api'
import { useSubscriptionStore } from './subscription'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token'),
    tenantId: localStorage.getItem('tenantId'),
    userId: localStorage.getItem('userId'),
    name: localStorage.getItem('name') || '',
    role: localStorage.getItem('role') || '',
  }),
  getters: {
    isLoggedIn: (s) => !!s.token && !!s.tenantId,
  },
  actions: {
    setSession(data) {
      this.token = data.token
      this.tenantId = data.tenant_id
      this.userId = data.user_id
      this.name = data.name || ''
      this.role = data.role || ''
      localStorage.setItem('token', data.token)
      localStorage.setItem('tenantId', data.tenant_id)
      localStorage.setItem('userId', data.user_id)
      localStorage.setItem('name', this.name)
      localStorage.setItem('role', this.role)
    },
    async login(phone, password) {
      const data = await api.auth.login(phone, password)
      this.setSession({
        token: data.token,
        tenant_id: data.tenant_id,
        user_id: data.user_id,
        name: data.name,
        role: data.role,
      })
      return data
    },
    async register(phone, password, name = '') {
      await api.auth.register({ phone, password, name: name || phone })
    },
    logout() {
      this.token = null
      this.tenantId = null
      this.userId = null
      this.name = ''
      this.role = ''
      localStorage.removeItem('token')
      localStorage.removeItem('tenantId')
      localStorage.removeItem('userId')
      localStorage.removeItem('name')
      localStorage.removeItem('role')
      useSubscriptionStore().$patch({ data: null })
    },
  },
})
