<template>
  <div class="min-h-full">
    <template v-if="loading">
      <div class="animate-pulse space-y-4">
        <div class="flex items-center gap-3 p-4 mb-4 bg-white rounded-2xl shadow-sm border border-gray-100">
          <div class="w-12 h-12 rounded-full bg-gray-200" />
          <div class="flex-1 space-y-2">
            <div class="h-4 w-32 bg-gray-200 rounded" />
            <div class="h-3 w-20 bg-gray-200 rounded" />
          </div>
        </div>
        <div class="space-y-3">
          <div v-for="i in 4" :key="i" class="flex items-center gap-4 p-4 bg-white rounded-2xl">
            <div class="w-8 h-8 bg-gray-200 rounded" />
            <div class="flex-1 space-y-1">
              <div class="h-4 w-24 bg-gray-200 rounded" />
              <div class="h-3 w-32 bg-gray-200 rounded" />
            </div>
          </div>
        </div>
      </div>
    </template>
    <template v-else>
    <!-- User info -->
  

    <div class="space-y-3">
      <router-link
        v-for="item in menuItems"
        :key="item.to"
        :to="item.to"
        class="flex items-center gap-4 p-4 bg-white rounded-2xl shadow-sm border border-gray-100 active:bg-gray-50 transition-colors"
      >
        <span class="w-8 h-8 flex items-center justify-center shrink-0 text-gray-500" v-html="item.icon" />
        <div class="flex-1">
          <p class="font-semibold text-gray-900">{{ item.label }}</p>
          <p class="text-xs text-gray-500">{{ item.desc }}</p>
        </div>
        <span class="text-gray-400 text-lg">›</span>
      </router-link>

      <!-- Sync -->
      <button
        type="button"
        class="flex items-center gap-4 w-full p-4 bg-white rounded-2xl shadow-sm border border-gray-100 active:bg-gray-50 transition-colors text-left disabled:opacity-60"
        :disabled="!isOnline || syncing"
        @click="manualSync"
      >
        <span class="w-8 h-8 flex items-center justify-center shrink-0 text-gray-500" v-html="syncing ? syncIcons.loading : syncIcons.sync" />
        <div class="flex-1">
          <p class="font-semibold text-gray-900">{{ syncing ? 'Menyalin...' : 'Sinkronkan data offline' }}</p>
          <p class="text-xs text-gray-500">{{ isOnline ? 'Kirim data lokal ke server' : 'Perlu koneksi internet' }}</p>
        </div>
      </button>

      <!-- Logout -->
      <button
        type="button"
        class="flex items-center gap-4 w-full p-4 bg-white rounded-2xl shadow-sm border border-gray-100 active:bg-red-50 transition-colors text-left"
        @click="handleLogout"
      >
        <span class="w-8 h-8 flex items-center justify-center shrink-0 text-red-500" v-html="syncIcons.logout" />
        <div class="flex-1">
          <p class="font-semibold text-red-600">Keluar</p>
          <p class="text-xs text-gray-500">Logout dari akun</p>
        </div>
      </button>
    </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useSettingsStore } from '../stores/settings'
import { syncPending } from '../stores/sync'

const router = useRouter()
const auth = useAuthStore()
const settingsStore = useSettingsStore()
const settings = computed(() => settingsStore.settings)
const loading = ref(true)
const syncing = ref(false)
const isOnline = ref(typeof navigator !== 'undefined' ? navigator.onLine : true)

function updateOnline() {
  isOnline.value = typeof navigator !== 'undefined' ? navigator.onLine : true
}
onMounted(async () => {
  window.addEventListener('online', updateOnline)
  window.addEventListener('offline', updateOnline)
  try {
    if (auth.tenantId) await settingsStore.load(auth.tenantId)
  } finally {
    loading.value = false
  }
})
onUnmounted(() => {
  window.removeEventListener('online', updateOnline)
  window.removeEventListener('offline', updateOnline)
})

const menuIcons = {
  chart: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.8" stroke="currentColor" class="w-5 h-5"><path stroke-linecap="round" stroke-linejoin="round" d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 013 19.875v-6.75zM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V8.625zM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V4.125z"/></svg>',
  settings: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.8" stroke="currentColor" class="w-5 h-5"><path stroke-linecap="round" stroke-linejoin="round" d="M9.594 3.94c.09-.542.56-.94 1.11-.94h2.593c.55 0 1.02.398 1.11.94l.213 1.281c.063.374.313.686.645.87.074.04.147.083.22.127.324.196.72.257 1.075.124l1.217-.456a1.125 1.125 0 011.37.49l1.296 2.247a1.125 1.125 0 01-.26 1.431l-1.003.827c-.293.24-.438.613-.431.992a6.759 6.759 0 010 .255c-.007.378.138.75.43.99l1.005.828c.424.35.534.954.26 1.43l-1.298 2.247a1.125 1.125 0 01-1.369.491l-1.217-.456c-.355-.133-.75-.072-1.076.124a6.57 6.57 0 01-.22.128c-.331.183-.581.495-.644.869l-.213 1.28c-.09.543-.56.941-1.11.941h-2.594c-.55 0-1.02-.398-1.11-.94l-.213-1.281c-.062-.374-.312-.686-.644-.87a6.52 6.52 0 01-.22-.127c-.325-.196-.72-.257-1.076-.124l-1.217.456a1.125 1.125 0 01-1.369-.49l-1.297-2.247a1.125 1.125 0 01.26-1.431l1.004-.827c.292-.24.437-.613.43-.992a6.932 6.932 0 010-.255c.007-.378-.138-.75-.43-.99l-1.004-.828a1.125 1.125 0 01-.26-1.43l1.297-2.247a1.125 1.125 0 011.37-.491l1.216.456c.356.133.751.072 1.076-.124.072-.044.146-.087.22-.128.332-.183.582-.495.644-.869l.213-1.281z"/><path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/></svg>',
  subscription: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.8" stroke="currentColor" class="w-5 h-5"><path stroke-linecap="round" stroke-linejoin="round" d="M3.75 13.5l10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75z"/></svg>',
}
const syncIcons = {
  sync: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.8" stroke="currentColor" class="w-5 h-5"><path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M2.985 19.644l3.181-3.183m0 0l-3.181 3.183m3.183-3.183v4.992"/></svg>',
  loading: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.8" stroke="currentColor" class="w-5 h-5 animate-spin"><path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M2.985 19.644l3.181-3.183m0 0l-3.181 3.183m3.183-3.183v4.992"/></svg>',
  logout: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.8" stroke="currentColor" class="w-5 h-5"><path stroke-linecap="round" stroke-linejoin="round" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"/></svg>',
}

const menuItems = [
  { to: '/app/subscription', label: 'Langganan', icon: menuIcons.subscription, desc: 'Upgrade plan atau perpanjang layanan' },
  { to: '/app/reports', label: 'Laporan', icon: menuIcons.chart, desc: 'Transaksi, pengeluaran & analisis' },
  { to: '/app/settings', label: 'Pengaturan', icon: menuIcons.settings, desc: 'Konfigurasi toko' },
]

const avatarLetter = computed(() => {
  const n = auth.name || settingsStore.settings?.name || 'U'
  return (n + '').charAt(0).toUpperCase()
})

async function manualSync() {
  if (!auth.tenantId || !isOnline.value) return
  syncing.value = true
  try {
    await syncPending(auth.tenantId)
  } finally {
    syncing.value = false
  }
}

function handleLogout() {
  auth.logout()
  router.replace('/')
}
</script>
