<template>
  <div class="min-h-dvh flex flex-col bg-[#f8faf8] safe-top safe-bottom">
    <!-- Header: same style as tenant -->
    <header class="sticky top-0 z-20 flex items-center justify-between h-14 px-4 bg-white/95 backdrop-blur border-b border-gray-100/80 safe-top">
      <div class="flex items-center gap-2 shrink-0 min-w-0">
        <span class="w-8 h-8 rounded-lg bg-primary-100 flex items-center justify-center text-primary-600 font-bold text-sm">
          A
        </span>
        <span class="text-base font-bold text-primary-600 tracking-tight truncate">HSmart Admin</span>
      </div>
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
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAdminAuthStore } from '../stores/adminAuth'

const navIcons = {
  dashboard: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.8" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M3.75 6A2.25 2.25 0 016 3.75h2.25A2.25 2.25 0 0110.5 6v2.25a2.25 2.25 0 01-2.25 2.25H6a2.25 2.25 0 01-2.25-2.25V6zM3.75 15.75A2.25 2.25 0 016 13.5h2.25a2.25 2.25 0 012.25 2.25V18a2.25 2.25 0 01-2.25 2.25H6A2.25 2.25 0 013.75 18v-2.25zM13.5 6a2.25 2.25 0 012.25-2.25H18A2.25 2.25 0 0120.25 6v2.25A2.25 2.25 0 0118 10.5h-2.25a2.25 2.25 0 01-2.25-2.25V6zM13.5 15.75a2.25 2.25 0 012.25-2.25H18a2.25 2.25 0 012.25 2.25V18A2.25 2.25 0 0118 20.25h-2.25A2.25 2.25 0 0113.5 18v-2.25z"/></svg>',
  tenants: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.8" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M15 19.128a9.38 9.38 0 002.625.372 9.337 9.337 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 018.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0111.964-3.07M12 6.375a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zm8.25 2.25a2.625 2.625 0 11-5.25 0 2.625 2.625 0 015.25 0z"/></svg>',
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
]

function isActive(item) {
  if (item.to === '/admin') return route.path === '/admin' || route.path === '/admin/'
  return route.path.startsWith(item.to)
}

function handleLogout() {
  adminAuth.logout()
  router.replace('/admin/login')
}
</script>
