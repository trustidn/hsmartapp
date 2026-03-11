const BASE = '/api'

function getAdminHeaders() {
  const token = localStorage.getItem('admin_token')
  const headers = { 'Content-Type': 'application/json' }
  if (token) headers['Authorization'] = `Bearer ${token}`
  return headers
}

function getAdminHeadersNoContentType() {
  const token = localStorage.getItem('admin_token')
  const headers = {}
  if (token) headers['Authorization'] = `Bearer ${token}`
  return headers
}

async function adminRequest(method, path, body) {
  const opts = { method, headers: getAdminHeaders() }
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
    const msg = err.detail ? `${err.error || 'Request failed'}: ${err.detail}` : (err.error || 'Request failed')
    throw new Error(msg)
  }
  if (res.status === 204) return null
  return res.json()
}

export const adminApi = {
  dashboard: {
    stats: () => adminRequest('GET', '/admin/dashboard/stats'),
  },
  tenants: {
    list: (params = {}) => {
      const p = new URLSearchParams()
      if (params.limit != null) p.set('limit', String(params.limit))
      if (params.offset != null) p.set('offset', String(params.offset))
      if (params.search) p.set('search', params.search)
      const q = p.toString()
      return adminRequest('GET', '/admin/tenants' + (q ? '?' + q : ''))
    },
    get: (id) => adminRequest('GET', '/admin/tenants/get?id=' + encodeURIComponent(id)),
    updateStatus: (id, status) => adminRequest('PATCH', '/admin/tenants/status', { id, status }),
    updateSubscription: (id, data) =>
      adminRequest('PATCH', '/admin/tenants/subscription', { id, ...(data || {}) }),
    revokeSubscription: (id) =>
      adminRequest('POST', '/admin/tenants/subscription/revoke', { id }),
  },
  saasSettings: {
    get: () => adminRequest('GET', '/admin/saas-settings'),
    update: (data) => adminRequest('PATCH', '/admin/saas-settings', data),
  },
  upload: {
    logo: async (file) => {
      const form = new FormData()
      form.append('file', file)
      const res = await fetch(BASE + '/admin/upload/logo', {
        method: 'POST',
        headers: getAdminHeadersNoContentType(),
        body: form,
      })
      if (!res.ok) {
        const t = await res.text()
        let err
        try {
          err = JSON.parse(t)
        } catch {
          err = { error: t || res.statusText }
        }
        throw new Error(err.error || 'Upload gagal')
      }
      const data = await res.json()
      return data.url
    },
  },
  subscriptionOrders: {
    list: (params = {}) => {
      const p = new URLSearchParams()
      if (params.limit != null) p.set('limit', String(params.limit))
      if (params.offset != null) p.set('offset', String(params.offset))
      if (params.status) p.set('status', params.status)
      const q = p.toString()
      return adminRequest('GET', '/admin/subscription-orders' + (q ? '?' + q : ''))
    },
    approve: (orderId) => adminRequest('POST', '/admin/subscription-orders/approve', { order_id: orderId }),
    reject: (orderId, reason) => adminRequest('POST', '/admin/subscription-orders/reject', { order_id: orderId, reason }),
  },
  plans: {
    list: () => adminRequest('GET', '/admin/plans'),
    update: (planSlug, data) =>
      adminRequest('PATCH', '/admin/plans', { plan_slug: planSlug, ...data }),
    delete: (planSlug) =>
      adminRequest('DELETE', '/admin/plans?plan_slug=' + encodeURIComponent(planSlug)),
    restore: (planSlug) =>
      adminRequest('POST', '/admin/plans/restore', { plan_slug: planSlug }),
  },
}
