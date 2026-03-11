<template>
  <div class="min-h-dvh flex flex-col bg-[#f8faf8] safe-top safe-bottom">
    <!-- Header: same style as tenant -->
    <header class="sticky top-0 z-20 flex items-center justify-between h-14 px-4 bg-white/95 backdrop-blur border-b border-gray-100/80 safe-top">
      <router-link to="/admin" class="flex items-center gap-2 shrink-0 min-w-0">
        <img
          v-if="saasSettings.logoUrl"
          :src="saasSettings.logoUrl"
          alt="Logo"
          class="w-8 h-8 object-contain rounded-lg"
        />
        <span v-else class="w-8 h-8 rounded-lg bg-primary-100 flex items-center justify-center text-primary-600 font-bold text-sm">A</span>
        <span class="text-base font-bold text-primary-600 tracking-tight truncate">HSmart Admin</span>
      </router-link>
      <h1 class="flex-1 text-center text-base font-semibold text-gray-800 truncate px-2">{{ pageTitle }}</h1>
      <div class="shrink-0 flex items-center gap-3">
        <span class="text-sm text-gray-500 hidden sm:inline">{{ adminAuth.name }}</span>
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

    <div class="flex flex-1 min-h-0">
      <!-- Sidebar -->
      <aside class="hidden md:flex flex-col w-56 shrink-0 bg-white/95 backdrop-blur border-r border-gray-100/80">
        <nav class="flex-1 py-4 px-3">
          <router-link
            v-for="item in navItems"
            :key="item.to"
            :to="item.to"
            class="flex items-center gap-3 px-3 py-2.5 rounded-xl mb-1 transition-colors"
            :class="isActive(item) ? 'bg-primary-50 text-primary-600 font-medium' : 'text-gray-600 hover:bg-gray-50 hover:text-gray-800'"
          >
            <span class="w-5 h-5 flex items-center justify-center [&_svg]:w-5 [&_svg]:h-5" v-html="item.icon" />
            <span>{{ item.label }}</span>
          </router-link>
        </nav>
      </aside>

      <!-- Mobile nav: bottom dock style (like tenant) when sidebar hidden -->
      <nav class="md:hidden fixed bottom-4 left-1/2 -translate-x-1/2 z-20 safe-bottom" style="bottom: max(1rem, env(safe-area-inset-bottom))">
        <div class="flex items-center gap-1 px-3 py-2 rounded-2xl bg-white/95 backdrop-blur-xl shadow-lg shadow-black/10 border border-gray-100/80" style="max-width: 360px">
          <router-link
            v-for="item in navItems"
            :key="item.to"
            :to="item.to"
            class="flex flex-col items-center justify-center min-w-[48px] py-2 px-1 rounded-xl transition-colors"
            :class="isActive(item) ? 'text-primary-600 bg-primary-50' : 'text-gray-500 hover:bg-gray-50 hover:text-gray-700'"
          >
            <span class="text-lg [&_svg]:w-5 [&_svg]:h-5" v-html="item.icon" />
            <span class="text-[10px] font-medium leading-tight">{{ item.label }}</span>
          </router-link>
        </div>
      </nav>

      <!-- Main content -->
      <main class="flex-1 overflow-auto pb-24 md:pb-6">
        <div class="min-h-full p-4 md:p-6 max-w-4xl mx-auto">
          <router-view />
        </div>
      </main>
    </div>
    <AppToast />
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import AppToast from '../components/AppToast.vue'
import { useAdminAuthStore } from '../stores/adminAuth'
import { useSaasSettingsStore } from '../stores/saasSettings'

const saasSettings = useSaasSettingsStore()
onMounted(() => saasSettings.loadForAdmin())

const navIcons = {
  dashboard: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.8" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M3.75 6A2.25 2.25 0 016 3.75h2.25A2.25 2.25 0 0110.5 6v2.25a2.25 2.25 0 01-2.25 2.25H6a2.25 2.25 0 01-2.25-2.25V6zM3.75 15.75A2.25 2.25 0 016 13.5h2.25a2.25 2.25 0 012.25 2.25V18a2.25 2.25 0 01-2.25 2.25H6A2.25 2.25 0 013.75 18v-2.25zM13.5 6a2.25 2.25 0 012.25-2.25H18A2.25 2.25 0 0120.25 6v2.25A2.25 2.25 0 0118 10.5h-2.25a2.25 2.25 0 01-2.25-2.25V6zM13.5 15.75a2.25 2.25 0 012.25-2.25H18a2.25 2.25 0 012.25 2.25V18A2.25 2.25 0 0118 20.25h-2.25A2.25 2.25 0 0113.5 18v-2.25z"/></svg>',
  tenants: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.8" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M15 19.128a9.38 9.38 0 002.625.372 9.337 9.337 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 018.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0111.964-3.07M12 6.375a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zm8.25 2.25a2.625 2.625 0 11-5.25 0 2.625 2.625 0 015.25 0z"/></svg>',
  plans: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.8" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M9.594 3.94c.09-.542.56-.94 1.11-.94h2.593c.55 0 1.02.398 1.11.94l.213 1.281c.063.374.313.686.645.87.074.04.147.083.22.127.325.196.72.257 1.075.124l1.217-.456a1.125 1.125 0 011.37.49l1.296 2.247a1.125 1.125 0 01-.26 1.431l-1.003.827c-.293.241-.438.613-.431.992a7.723 7.723 0 010 .255c-.007.378.138.75.43.991l1.004.827c.424.35.534.955.26 1.43l-1.298 2.247a1.125 1.125 0 01-1.369.491l-1.217-.456c-.355-.133-.75-.072-1.076.124a6.47 6.47 0 01-.22.128c-.331.183-.581.495-.644.869l-.213 1.281c-.09.543-.56.94-1.11.94h-2.594c-.55 0-1.019-.398-1.11-.94l-.213-1.281c-.062-.374-.312-.686-.644-.87a6.52 6.52 0 01-.22-.127c-.325-.196-.72-.257-1.076-.124l-1.217.456a1.125 1.125 0 01-1.369-.49l-1.297-2.247a1.125 1.125 0 01.26-1.431l1.004-.827c.292-.24.437-.613.43-.991a6.932 6.932 0 010-.255c.007-.379-.138-.75-.43-.992l-1.004-.827a1.125 1.125 0 01-.26-1.43l1.297-2.247a1.125 1.125 0 011.37-.491l1.216.456c.356.133.751.072 1.076-.124.072-.044.146-.087.22-.128.332-.183.582-.495.644-.869l.214-1.28z"/><path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/></svg>',
  orders: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.8" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M15.75 10.5V4.5a3.75 3.75 0 10-7.5 0v6m11.356-1.993l1.263 12c.07.665-.45 1.243-1.119 1.243H4.25a1.125 1.125 0 01-1.12-1.243l1.264-12A1.125 1.125 0 015.513 7.5h12.974c.576 0 1.059.435 1.119 1.007zM8.625 10.5a.375.375 0 11-.75 0 .375.375 0 01.75 0zm7.5 0a.375.375 0 11-.75 0 .375.375 0 01.75 0z"/></svg>',
  settings: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.8" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M9.594 3.94c.09-.542.56-.94 1.11-.94h2.593c.55 0 1.02.398 1.11.94l.213 1.281c.063.374.313.686.645.87.074.04.147.083.22.127.325.196.72.257 1.075.124l1.217-.456a1.125 1.125 0 011.37.49l1.296 2.247a1.125 1.125 0 01-.26 1.431l-1.003.827c-.293.241-.438.613-.431.992a7.723 7.723 0 010 .255c-.007.378.138.75.43.991l1.004.827c.424.35.534.955.26 1.43l-1.298 2.247a1.125 1.125 0 01-1.369.491l-1.217-.456c-.355-.133-.75-.072-1.076.124a6.47 6.47 0 01-.22.128c-.331.183-.581.495-.644.869l-.213 1.281c-.09.543-.56.94-1.11.94h-2.594c-.55 0-1.019-.398-1.11-.94l-.213-1.281c-.062-.374-.312-.686-.644-.87a6.52 6.52 0 01-.22-.127c-.325-.196-.72-.257-1.076-.124l-1.217.456a1.125 1.125 0 01-1.369-.49l-1.297-2.247a1.125 1.125 0 01.26-1.431l1.004-.827c.292-.24.437-.613.43-.991a6.932 6.932 0 010-.255c.007-.379-.138-.75-.43-.992l-1.004-.827a1.125 1.125 0 01-.26-1.43l1.297-2.247a1.125 1.125 0 011.37-.491l1.216.456c.356.133.751.072 1.076-.124.072-.044.146-.087.22-.128.332-.183.582-.495.644-.869l.214-1.28z"/><path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/></svg>',
}

const route = useRoute()
const router = useRouter()
const adminAuth = useAdminAuthStore()

const pageTitle = computed(() => {
  const m = route.matched[route.matched.length - 1]?.meta
  return m?.title ?? 'Admin'
})

const navItems = [
  { to: '/admin', label: 'Dashboard', icon: navIcons.dashboard },
  { to: '/admin/tenants', label: 'Tenants', icon: navIcons.tenants },
  { to: '/admin/subscription-orders', label: 'Order Langganan', icon: navIcons.orders },
  { to: '/admin/plans', label: 'Pengaturan Plan', icon: navIcons.plans },
  { to: '/admin/saas-settings', label: 'Pengaturan SaaS', icon: navIcons.settings },
]

function isActive(item) {
  if (item.to === '/admin') return route.path === '/admin' || route.path === '/admin/'
  if (item.to === '/admin/tenants') return route.path.startsWith('/admin/tenants')
  return route.path.startsWith(item.to)
}

function handleLogout() {
  adminAuth.logout()
  router.replace('/admin/login')
}
</script>
