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
      <div class="shrink-0 w-16 flex justify-end">
      <button
        type="button"
        class="p-2 -mr-2 rounded-lg text-gray-500 hover:bg-gray-100 hover:text-gray-700"
        aria-label="Keluar"
        @click="handleLogout"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
        </svg>
      </button>
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
import { computed, onMounted } from 'vue'
import OnboardingTour from '../components/OnboardingTour.vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useSettingsStore } from '../stores/settings'
import { useSaasSettingsStore } from '../stores/saasSettings'
import { useSubscriptionStore } from '../stores/subscription'

const navIcons = {
  pos: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.8" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M15.75 10.5V6a3.75 3.75 0 10-7.5 0v4.5m11.356-1.993l1.263 12c.07.665-.45 1.243-1.119 1.243H4.25a1.125 1.125 0 01-1.12-1.243l1.264-12A1.125 1.125 0 015.513 7.5h12.974c.576 0 1.059.435 1.119 1.007z"/></svg>',
  chart: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.8" stroke="currentColor" class="w-5 h-5"><path stroke-linecap="round" stroke-linejoin="round" d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 013 19.875v-6.75zM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V8.625zM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V4.125z"/></svg>',
  product: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.8" stroke="currentColor" class="w-5 h-5"><path stroke-linecap="round" stroke-linejoin="round" d="M20.25 7.5l-.625 10.632a2.25 2.25 0 01-2.247 2.118H6.622a2.25 2.25 0 01-2.247-2.118L3.75 7.5M10 11.25h4M3.375 7.5h17.25c.621 0 1.125-.504 1.125-1.125v-1.5c0-.621-.504-1.125-1.125-1.125H3.375c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125z"/></svg>',
  expense: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.8" stroke="currentColor" class="w-5 h-5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 6v12m-3-2.818l.879.659c1.171.879 3.07.879 4.242 0 1.172-.879 1.172-2.303 0-3.182C13.536 12.219 12.768 12 12 12c-.725 0-1.45-.22-2.003-.659-1.106-.879-1.106-2.303 0-3.182s2.9-.879 4.006 0l.415.33M21 12a2.25 2.25 0 00-2.25-2.25H15a3 3 0 11-6 0H5.25A2.25 2.25 0 003 12m18 0v6a2.25 2.25 0 01-2.25 2.25H5.25A2.25 2.25 0 013 18v-6"/></svg>',
  more: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.8" stroke="currentColor" class="w-5 h-5"><path stroke-linecap="round" stroke-linejoin="round" d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"/></svg>',
}

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const settings = useSettingsStore()
const saasSettings = useSaasSettingsStore()
const subscriptionStore = useSubscriptionStore()

const pageTitle = computed(() => {
  const m = route.matched[route.matched.length - 1]?.meta
  return m?.title ?? 'HSmart'
})
const appName = computed(() => settings.settings?.name || 'HSmart')
const appInitial = computed(() => (appName.value?.charAt(0) || 'H').toUpperCase())

onMounted(() => {
  if (auth.tenantId) {
    settings.load(auth.tenantId)
    saasSettings.loadForTenant()
    subscriptionStore.load(auth.tenantId)
  }
})

const bottomNavItems = [
  { to: '/app', label: 'Dashboard', icon: navIcons.chart },
  { to: '/app/expenses', label: 'Pengeluaran', icon: navIcons.expense },
  { to: '/app/pos', label: 'POS', icon: navIcons.pos, featured: true },
  { to: '/app/products', label: 'Produk', icon: navIcons.product },
  { to: '/app/more', label: 'Lainnya', icon: navIcons.more },
]

function isActive(item) {
  if (item.to === '/app') return route.path === '/app' || route.path === '/app/'
  if (item.to === '/app/pos') return route.path === '/app/pos'
  if (item.to === '/app/more') return route.path === '/app/more' || ['/app/reports', '/app/settings', '/app/subscription'].includes(route.path)
  return route.path.startsWith(item.to)
}

function handleLogout() {
  auth.logout()
  router.replace('/login')
}
</script>
