<template>
  <div class="min-h-full flex flex-col safe-top safe-bottom">
    <router-view />
    <Toast />
  </div>
</template>

<script setup>
import { onMounted, onUnmounted } from 'vue'
import Toast from './components/Toast.vue'
import { useAuthStore } from './stores/auth'
import { useProductsStore } from './stores/products'
import { syncPending } from './stores/sync'

const auth = useAuthStore()

function doSync() {
  if (auth.isLoggedIn && navigator.onLine && auth.tenantId) {
    syncPending(auth.tenantId).catch(() => {})
  }
}

function onOnline() {
  setTimeout(doSync, 1500)
}

function onVisibilityChange() {
  if (document.visibilityState === 'visible') doSync()
}

onMounted(async () => {
  if (auth.isLoggedIn) {
    try {
      const products = useProductsStore()
      await products.load(auth.tenantId)
      await syncPending(auth.tenantId)
    } catch (_) {
      // Per-page load will show error; avoid crashing app
    }
  }
  window.addEventListener('online', onOnline)
  document.addEventListener('visibilitychange', onVisibilityChange)
})

onUnmounted(() => {
  window.removeEventListener('online', onOnline)
  document.removeEventListener('visibilitychange', onVisibilityChange)
})
</script>
