const BASE = '/api'

function getAdminHeaders() {
  const token = localStorage.getItem('admin_token')
  const headers = { 'Content-Type': 'application/json' }
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
    throw new Error(err.error || 'Request failed')
  }
  if (res.status === 204) return null
  return res.json()
}

export const adminApi = {
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
  },
}
