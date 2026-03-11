import { ref, onMounted, watch } from 'vue'
import { api } from '../lib/api'

const CACHE_KEY = 'hsmart_branding'
const CACHE_TTL_MS = 24 * 60 * 60 * 1000 // 24 jam
const DEFAULT_FAVICON = '/favicon.svg'

function getCached() {
  try {
    const raw = localStorage.getItem(CACHE_KEY)
    if (!raw) return null
    const data = JSON.parse(raw)
    if (data.cachedAt && Date.now() - data.cachedAt < CACHE_TTL_MS) {
      return { appName: data.app_name || 'HSmart', logoUrl: data.logo_url || '' }
    }
  } catch {}
  return null
}

function setCache(data) {
  try {
    localStorage.setItem(CACHE_KEY, JSON.stringify({
      app_name: data.app_name || 'HSmart',
      logo_url: data.logo_url || '',
      cachedAt: Date.now(),
    }))
  } catch {}
}

function applyDocumentHead(appName, logoUrl) {
  document.title = appName || 'HSmart'

  let linkIcon = document.querySelector('link[rel="icon"]')
  let linkApple = document.querySelector('link[rel="apple-touch-icon"]')
  const href = logoUrl || DEFAULT_FAVICON

  if (!linkIcon) {
    linkIcon = document.createElement('link')
    linkIcon.rel = 'icon'
    linkIcon.type = href.endsWith('.svg') ? 'image/svg+xml' : 'image/png'
    document.head.appendChild(linkIcon)
  }
  linkIcon.href = href
  if (href.endsWith('.svg')) linkIcon.type = 'image/svg+xml'
  else linkIcon.type = 'image/png'

  if (!linkApple) {
    linkApple = document.createElement('link')
    linkApple.rel = 'apple-touch-icon'
    document.head.appendChild(linkApple)
  }
  linkApple.href = href
}

/**
 * Composable untuk branding (app name, logo) - dipakai di Welcome, Login, Register.
 * Update document title, favicon, apple-touch-icon. Cache localStorage untuk offline.
 */
export function useBranding() {
  const appName = ref('HSmart')
  const logoUrl = ref('')
  const loading = ref(false)

  watch([appName, logoUrl], ([name, logo]) => {
    applyDocumentHead(name, logo)
  }, { immediate: true })

  onMounted(async () => {
    const cached = getCached()
    if (cached) {
      appName.value = cached.appName
      logoUrl.value = cached.logoUrl
    }
    loading.value = true
    try {
      const data = await api.publicBranding.get()
      appName.value = data.app_name || 'HSmart'
      logoUrl.value = data.logo_url || ''
      setCache(data)
    } catch {
      if (!cached) {
        appName.value = 'HSmart'
        logoUrl.value = ''
      }
    } finally {
      loading.value = false
    }
  })

  return { appName, logoUrl, loading }
}
