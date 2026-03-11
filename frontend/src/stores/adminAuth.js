import { defineStore } from 'pinia'

const ADMIN_TOKEN = 'admin_token'
const ADMIN_ID = 'admin_id'
const ADMIN_NAME = 'admin_name'
const ADMIN_ROLE = 'admin_role'

export const useAdminAuthStore = defineStore('adminAuth', {
  state: () => ({
    token: localStorage.getItem(ADMIN_TOKEN),
    adminId: localStorage.getItem(ADMIN_ID),
    name: localStorage.getItem(ADMIN_NAME) || '',
    role: localStorage.getItem(ADMIN_ROLE) || '',
  }),
  getters: {
    isAdminLoggedIn: (s) => !!s.token && s.role === 'superadmin',
  },
  actions: {
    setSession(data) {
      this.token = data.token
      this.adminId = data.admin_id
      this.name = data.name || ''
      this.role = data.role || ''
      localStorage.setItem(ADMIN_TOKEN, data.token)
      localStorage.setItem(ADMIN_ID, data.admin_id)
      localStorage.setItem(ADMIN_NAME, this.name)
      localStorage.setItem(ADMIN_ROLE, this.role)
    },
    async login(email, password) {
      const data = await fetch('/api/admin/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
      }).then(async (res) => {
        if (!res.ok) {
          const t = await res.text()
          let err
          try { err = JSON.parse(t) } catch { err = { error: t || res.statusText } }
          throw new Error(err.error || 'Login gagal')
        }
        return res.json()
      })
      this.setSession({
        token: data.token,
        admin_id: data.admin_id,
        name: data.name,
        role: data.role,
      })
      return data
    },
    logout() {
      this.token = null
      this.adminId = null
      this.name = ''
      this.role = ''
      localStorage.removeItem(ADMIN_TOKEN)
      localStorage.removeItem(ADMIN_ID)
      localStorage.removeItem(ADMIN_NAME)
      localStorage.removeItem(ADMIN_ROLE)
    },
  },
})
