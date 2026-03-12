<template>
  <div class="min-h-dvh flex flex-col bg-[#f8faf8] safe-top safe-bottom">
    <!-- Header: Logo + nama app kiri, judul tengah, exit kanan -->
    <header class="sticky top-0 z-20 flex items-center justify-between h-14 px-4 bg-white/95 backdrop-blur border-b border-gray-100/80 safe-top">
      <router-link to="/app" class="flex items-center gap-2 shrink-0 min-w-0">
        <img
          v-if="saasSettings.logoUrl"
          :src="saasSettings.logoUrl"
          alt="Logo"
          class="w-8 h-8 object-contain rounded-lg"
        />
        <span v-else class="w-8 h-8 rounded-lg bg-primary-100 flex items-center justify-center text-primary-600 font-bold text-sm">
          {{ appInitial }}
        </span>
        <span class="text-base font-bold text-primary-600 tracking-tight truncate max-w-[120px]">{{ appName }}</span>
      </router-link>
      <h1 class="flex-1 text-center text-base font-semibold text-gray-800 truncate px-2">{{ pageTitle }}</h1>
      <div class="shrink-0 flex items-center gap-1">
        <!-- Menu Pengaturan (gear) dengan dropdown -->
        <div class="relative">
          <button
            type="button"
            class="p-2 -mr-1 rounded-lg text-gray-500 hover:bg-gray-100 hover:text-gray-700"
            aria-label="Pengaturan"
            aria-haspopup="true"
            :aria-expanded="showSettingsMenu"
            @click="showSettingsMenu = !showSettingsMenu"
            @blur="onSettingsBlur"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
          </button>
          <Transition name="dropdown">
            <div
              v-show="showSettingsMenu"
              class="absolute right-0 top-full mt-1 w-56 py-1 bg-white rounded-xl shadow-lg border border-gray-100 z-30"
            >
              <router-link
                to="/app/subscription"
                class="flex items-center gap-3 px-4 py-2.5 text-sm text-gray-700 hover:bg-gray-50"
                @click="showSettingsMenu = false"
              >
                <svg class="w-4 h-4 text-gray-500 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3.75 13.5l10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75z"/></svg>
                Langganan
              </router-link>
              <router-link
                to="/app/settings"
                class="flex items-center gap-3 px-4 py-2.5 text-sm text-gray-700 hover:bg-gray-50"
                @click="showSettingsMenu = false"
              >
                <svg class="w-4 h-4 text-gray-500 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/></svg>
                Pengaturan
              </router-link>
              <button
                type="button"
                class="flex items-center gap-3 w-full px-4 py-2.5 text-sm text-gray-700 hover:bg-gray-50 text-left disabled:opacity-60"
                :disabled="!isOnlineVal || syncing"
                @click="manualSync"
              >
                <svg v-if="syncing" class="w-4 h-4 text-gray-500 shrink-0 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/></svg>
                <svg v-else class="w-4 h-4 text-gray-500 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/></svg>
                {{ syncing ? 'Menyalin...' : 'Sinkronkan data offline' }}
              </button>
              <hr class="my-1 border-gray-100" />
              <button
                type="button"
                class="flex items-center gap-3 w-full px-4 py-2.5 text-sm text-red-600 hover:bg-red-50 text-left"
                @click="handleLogout"
              >
                <svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"/></svg>
                Keluar
              </button>
            </div>
          </Transition>
        </div>
      </div>
    </header>

    <!-- Main content -->
    <main class="flex-1 overflow-auto pb-28 md:pb-32">
      <div class="min-h-full p-4 md:p-6 max-w-4xl mx-auto">
        <router-view />
      </div>
    </main>

    <!-- Bottom nav: diperbesar untuk navigasi jari (touch-friendly) -->
    <nav
      class="fixed bottom-2 left-1/2 -translate-x-1/2 z-20 safe-bottom nav-dock"
      style="bottom: max(0.5rem, env(safe-area-inset-bottom))"
    >
      <div
        class="flex items-center justify-between w-full px-2 py-3 rounded-2xl bg-white/95 backdrop-blur-xl shadow-lg shadow-black/10 border border-gray-100/80 nav-dock-inner"
        style="max-width: 400px"
      >
        <router-link
          v-for="item in bottomNavItems"
          :key="item.to"
          :to="item.to"
          class="flex flex-col items-center justify-center min-w-[4rem] w-16 py-3 px-1 rounded-xl transition-all duration-200 gap-1.5 touch-manipulation"
          :class="[
            item.featured
              ? isActive(item)
                ? 'text-primary-600 bg-primary-50 ring-1 ring-primary-200/60'
                : 'text-primary-500 hover:bg-primary-50/60 hover:text-primary-600 ring-1 ring-primary-100/40'
              : isActive(item)
                ? 'text-primary-600 bg-primary-50/80'
                : 'text-gray-500 hover:bg-gray-50 hover:text-gray-700',
          ]"
        >
          <span class="leading-none [&_svg]:inline-block flex-shrink-0" :class="item.featured ? '[&_svg]:w-7 [&_svg]:h-7' : '[&_svg]:w-6 [&_svg]:h-6'" v-html="item.icon" />
          <span class="text-[11px] font-semibold leading-tight text-center truncate w-full px-0.5">{{ item.label }}</span>
        </router-link>
      </div>
    </nav>

    <!-- Onboarding tour untuk tenant baru -->
    <OnboardingTour :tenant-id="auth.tenantId || ''" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import OnboardingTour from '../components/OnboardingTour.vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useSettingsStore } from '../stores/settings'
import { useSaasSettingsStore } from '../stores/saasSettings'
import { useSubscriptionStore } from '../stores/subscription'
import { syncPending } from '../stores/sync'
import { isOnline } from '../lib/api'
import { useToast } from '../composables/useToast'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const settings = useSettingsStore()
const saasSettings = useSaasSettingsStore()
const subscriptionStore = useSubscriptionStore()
const { show: toast } = useToast()

const showSettingsMenu = ref(false)
const syncing = ref(false)

const navIcons = {
  pos: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.8" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M15.75 10.5V6a3.75 3.75 0 10-7.5 0v4.5m11.356-1.993l1.263 12c.07.665-.45 1.243-1.119 1.243H4.25a1.125 1.125 0 01-1.12-1.243l1.264-12A1.125 1.125 0 015.513 7.5h12.974c.576 0 1.059.435 1.119 1.007z"/></svg>',
  chart: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.8" stroke="currentColor" class="w-5 h-5"><path stroke-linecap="round" stroke-linejoin="round" d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 013 19.875v-6.75zM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V8.625zM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V4.125z"/></svg>',
  product: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.8" stroke="currentColor" class="w-5 h-5"><path stroke-linecap="round" stroke-linejoin="round" d="M20.25 7.5l-.625 10.632a2.25 2.25 0 01-2.247 2.118H6.622a2.25 2.25 0 01-2.247-2.118L3.75 7.5M10 11.25h4M3.375 7.5h17.25c.621 0 1.125-.504 1.125-1.125v-1.5c0-.621-.504-1.125-1.125-1.125H3.375c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125z"/></svg>',
  expense: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.8" stroke="currentColor" class="w-5 h-5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 6v12m-3-2.818l.879.659c1.171.879 3.07.879 4.242 0 1.172-.879 1.172-2.303 0-3.182C13.536 12.219 12.768 12 12 12c-.725 0-1.45-.22-2.003-.659-1.106-.879-1.106-2.303 0-3.182s2.9-.879 4.006 0l.415.33M21 12a2.25 2.25 0 00-2.25-2.25H15a3 3 0 11-6 0H5.25A2.25 2.25 0 003 12m18 0v6a2.25 2.25 0 01-2.25 2.25H5.25A2.25 2.25 0 013 18v-6"/></svg>',
  report: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.8" stroke="currentColor" class="w-5 h-5"><path stroke-linecap="round" stroke-linejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z"/></svg>',
}

const pageTitle = computed(() => {
  const m = route.matched[route.matched.length - 1]?.meta
  return m?.title ?? 'HSmart'
})
// Nama aplikasi dari Pengaturan SaaS (set oleh admin), bukan nama tenant
const appName = computed(() => saasSettings.appName || 'HSmart')
const appInitial = computed(() => (appName.value?.charAt(0) || 'H').toUpperCase())
const isOnlineVal = computed(() => isOnline())

const bottomNavItems = [
  { to: '/app', label: 'Dashboard', icon: navIcons.chart },
  { to: '/app/expenses', label: 'Pengeluaran', icon: navIcons.expense },
  { to: '/app/pos', label: 'POS', icon: navIcons.pos, featured: true },
  { to: '/app/products', label: 'Produk', icon: navIcons.product },
  { to: '/app/reports', label: 'Laporan', icon: navIcons.report },
]

onMounted(() => {
  if (auth.tenantId) {
    settings.load(auth.tenantId)
    saasSettings.loadForTenant()
    subscriptionStore.load(auth.tenantId)
  }
})

function isActive(item) {
  if (item.to === '/app') return route.path === '/app' || route.path === '/app/'
  if (item.to === '/app/pos') return route.path === '/app/pos'
  return route.path.startsWith(item.to)
}

function onSettingsBlur() {
  setTimeout(() => { showSettingsMenu.value = false }, 150)
}

async function manualSync() {
  if (!auth.tenantId || !isOnline()) return
  syncing.value = true
  showSettingsMenu.value = false
  try {
    await syncPending(auth.tenantId)
    toast('Data offline berhasil disinkronkan')
  } catch {
    toast('Gagal sinkronisasi', 'error')
  } finally {
    syncing.value = false
  }
}

function handleLogout() {
  showSettingsMenu.value = false
  auth.logout()
  router.replace('/login')
}
</script>

<style scoped>
.dropdown-enter-active,
.dropdown-leave-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
}
.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>
