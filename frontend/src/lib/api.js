const BASE = '/api'

function getAuth() {
  const token = localStorage.getItem('token')
  const tenantId = localStorage.getItem('tenantId')
  const headers = { 'Content-Type': 'application/json' }
  if (token) headers['Authorization'] = `Bearer ${token}`
  if (tenantId) headers['X-Tenant-ID'] = tenantId
  return headers
}

async function request(method, path, body) {
  const opts = { method, headers: getAuth() }
  if (body) opts.body = JSON.stringify(body)
  const res = await fetch(BASE + path, opts)
  if (!res.ok) {
    const t = await res.text()
    let err
    try {
      err = JSON.parse(t)
    } catch {
      err = { error: t || res.statusText }
    }
    throw new Error(err.error || 'Request failed')
  }
  if (res.status === 204) return null
  return res.json()
}

export const api = {
  auth: {
    login: (phone, password) => request('POST', '/auth/login', { phone, password }),
    register: (data) => request('POST', '/register', data),
  },
  products: {
    list: (activeOnly = true) => request('GET', '/products' + (activeOnly ? '?active=true' : '')),
    create: (data) => request('POST', '/products', data),
    get: (id) => request('GET', '/products/get?id=' + id),
    update: (id, data) => request('PUT', '/products?id=' + id, data),
    delete: (id) => request('DELETE', '/products?id=' + id),
  },
  sales: {
    create: (data) => request('POST', '/sales', data),
    list: (from, to, tz = 'Asia/Jakarta', limit, offset) => {
      const params = new URLSearchParams()
      if (from) params.set('from', from)
      if (to) params.set('to', to)
      if (tz) params.set('tz', tz)
      if (limit != null) params.set('limit', String(limit))
      if (offset != null) params.set('offset', String(offset))
      const q = params.toString()
      return request('GET', '/sales' + (q ? '?' + q : ''))
    },
    get: (id) => request('GET', '/sales/get?id=' + id),
  },
  expenses: {
    create: (data) => request('POST', '/expenses', data),
    list: (from, to, tz = 'Asia/Jakarta', limit, offset) => {
      const params = new URLSearchParams()
      if (from) params.set('from', from)
      if (to) params.set('to', to)
      if (tz) params.set('tz', tz)
      if (limit != null) params.set('limit', String(limit))
      if (offset != null) params.set('offset', String(offset))
      const q = params.toString()
      return request('GET', '/expenses' + (q ? '?' + q : ''))
    },
    delete: (id) => request('DELETE', '/expenses?id=' + id),
  },
  report: {
    daily: (date, tz = 'Asia/Jakarta') => {
      const p = new URLSearchParams()
      if (date) p.set('date', date)
      if (tz) p.set('tz', tz)
      return request('GET', '/report/daily' + (p.toString() ? '?' + p : ''))
    },
    ranking: (date, tz = 'Asia/Jakarta') => {
      const p = new URLSearchParams()
      if (date) p.set('date', date)
      if (tz) p.set('tz', tz)
      return request('GET', '/report/ranking' + (p.toString() ? '?' + p : ''))
    },
    dashboard: (date, tz = 'Asia/Jakarta') => {
      const p = new URLSearchParams()
      if (date) p.set('date', date)
      if (tz) p.set('tz', tz)
      return request('GET', '/report/dashboard' + (p.toString() ? '?' + p : ''))
    },
    dashboardRange: (from, to, tz = 'Asia/Jakarta') => {
      const p = new URLSearchParams({ from, to, tz })
      return request('GET', '/report/dashboard?' + p)
    },
  },
  subscription: {
    get: () => request('GET', '/subscription'),
  },
  tenant: {
    getSettings: () => request('GET', '/tenant/settings'),
    updateSettings: (data) => request('PUT', '/tenant/settings', data),
  },
}

export function isOnline() {
  return navigator.onLine
}

/** Local date string YYYY-MM-DD (avoids timezone issues) */
export function toLocalDateStr(d = new Date()) {
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}
