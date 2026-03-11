<template>
  <div class="min-h-full">
    <!-- Selamat datang -->
    <section class="bg-gradient-to-r from-primary-500 to-primary-600 rounded-2xl p-4 mb-4 text-white shadow-sm">
      <p class="text-sm opacity-90">Selamat datang,</p>
      <p class="text-lg font-bold mt-0.5">{{ adminAuth.name }}</p>
    </section>

    <!-- Stats -->
    <section v-if="statsLoading" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4 mb-4">
      <div v-for="i in 4" :key="i" class="h-24 bg-gray-200 rounded-2xl animate-pulse" />
    </section>
    <section v-else class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4 mb-6">
      <div class="p-5 bg-white rounded-2xl shadow-sm border border-gray-100">
        <p class="text-xs font-medium text-gray-500">Tenant Aktif</p>
        <p class="text-2xl font-bold text-gray-900 mt-1">{{ formatNum(stats.active_tenants ?? 0) }}</p>
      </div>
      <div class="p-5 bg-white rounded-2xl shadow-sm border border-gray-100">
        <p class="text-xs font-medium text-gray-500">Total Order Langganan</p>
        <p class="text-2xl font-bold text-gray-900 mt-1">{{ formatNum(stats.total_orders ?? 0) }}</p>
      </div>
      <div class="p-5 bg-white rounded-2xl shadow-sm border border-gray-100">
        <p class="text-xs font-medium text-gray-500">Total Revenue</p>
        <p class="text-2xl font-bold text-primary-600 mt-1">Rp {{ formatNum(stats.total_revenue ?? 0) }}</p>
      </div>
      <div class="p-5 bg-white rounded-2xl shadow-sm border border-gray-100">
        <p class="text-xs font-medium text-gray-500">MRR (bulan ini)</p>
        <p class="text-2xl font-bold text-primary-600 mt-1">Rp {{ formatNum(stats.mrr ?? 0) }}</p>
      </div>
    </section>

    <!-- Grafik pertumbuhan -->
    <section v-if="!statsLoading && (stats.tenant_growth?.length || stats.revenue_by_month?.length)" class="grid gap-4 lg:grid-cols-2 mb-6">
      <div class="p-5 bg-white rounded-2xl shadow-sm border border-gray-100">
        <h3 class="text-sm font-medium text-gray-700 mb-3">Pertumbuhan Tenant (6 bulan)</h3>
        <div v-if="stats.tenant_growth?.length" class="h-40 flex items-end gap-2">
          <div
            v-for="d in stats.tenant_growth"
            :key="d.month"
            class="flex-1 min-w-0 flex flex-col items-center gap-1"
          >
            <span class="text-[10px] font-medium text-gray-600 h-4 leading-tight text-center truncate w-full">{{ d.count }}</span>
            <div
              class="w-full rounded-t bg-primary-500/80 min-h-[4px]"
              :style="{ height: Math.max(4, barHeightTenant(d)) + 'px' }"
            />
            <span class="text-[10px] text-gray-500 truncate w-full text-center">{{ d.month }}</span>
          </div>
        </div>
        <p v-else class="text-sm text-gray-400">Belum ada data</p>
      </div>
      <div class="p-5 bg-white rounded-2xl shadow-sm border border-gray-100">
        <h3 class="text-sm font-medium text-gray-700 mb-3">Revenue per Bulan (6 bulan)</h3>
        <div v-if="stats.revenue_by_month?.length" class="h-40 flex items-end gap-2">
          <div
            v-for="d in stats.revenue_by_month"
            :key="d.month"
            class="flex-1 min-w-0 flex flex-col items-center gap-1"
          >
            <span class="text-[10px] font-medium text-gray-600 h-4 leading-tight text-center truncate w-full">{{ formatCompact(d.revenue) }}</span>
            <div
              class="w-full rounded-t bg-green-500/80 min-h-[4px]"
              :style="{ height: Math.max(4, barHeightRevenue(d)) + 'px' }"
            />
            <span class="text-[10px] text-gray-500 truncate w-full text-center">{{ d.month }}</span>
          </div>
        </div>
        <p v-else class="text-sm text-gray-400">Belum ada data</p>
      </div>
    </section>

    <div class="grid gap-4 sm:grid-cols-2">
      <router-link
        to="/admin/tenants"
        class="block p-6 bg-white rounded-2xl shadow-sm border border-gray-100 hover:shadow-md transition-shadow"
      >
        <h2 class="font-semibold text-gray-800">Daftar Tenant</h2>
        <p class="text-sm text-gray-500 mt-1">Kelola merchant / UMKM terdaftar, status, dan subscription</p>
      </router-link>
      <router-link
        to="/admin/subscription-orders"
        class="block p-6 bg-white rounded-2xl shadow-sm border border-gray-100 hover:shadow-md transition-shadow"
      >
        <h2 class="font-semibold text-gray-800">Order Langganan</h2>
        <p class="text-sm text-gray-500 mt-1">Verifikasi dan setujui order upgrade/perpanjang dari tenant</p>
      </router-link>
      <router-link
        to="/admin/plans"
        class="block p-6 bg-white rounded-2xl shadow-sm border border-gray-100 hover:shadow-md transition-shadow"
      >
        <h2 class="font-semibold text-gray-800">Pengaturan Plan</h2>
        <p class="text-sm text-gray-500 mt-1">Atur max produk dan hari laporan per plan (Free, Premium 1/3/6/12 bln)</p>
      </router-link>
      <router-link
        to="/admin/saas-settings"
        class="block p-6 bg-white rounded-2xl shadow-sm border border-gray-100 hover:shadow-md transition-shadow"
      >
        <h2 class="font-semibold text-gray-800">Pengaturan SaaS</h2>
        <p class="text-sm text-gray-500 mt-1">Nama, logo, kontak admin, rekening pembayaran</p>
      </router-link>
    </div>

    <button
      @click="logout"
      class="mt-6 px-4 py-2 rounded-xl bg-white border border-gray-200 text-gray-700 font-medium hover:bg-gray-50"
    >
      Logout
    </button>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAdminAuthStore } from '../../stores/adminAuth'
import { adminApi } from '../../lib/adminApi'

const router = useRouter()
const adminAuth = useAdminAuthStore()
const stats = ref({})
const statsLoading = ref(true)

onMounted(async () => {
  try {
    stats.value = await adminApi.dashboard.stats()
  } catch {
    stats.value = {}
  } finally {
    statsLoading.value = false
  }
})

function formatNum(n) {
  return Number(n).toLocaleString('id-ID')
}
function formatCompact(n) {
  const num = Number(n)
  if (num >= 1000000) return (num / 1000000).toFixed(1) + 'jt'
  if (num >= 1000) return (num / 1000).toFixed(0) + 'rb'
  return num.toLocaleString('id-ID')
}
const maxTenantCount = computed(() => {
  const arr = stats.value?.tenant_growth ?? []
  return Math.max(...arr.map((x) => x.count), 1)
})
const maxRevenue = computed(() => {
  const arr = stats.value?.revenue_by_month ?? []
  return Math.max(...arr.map((x) => x.revenue), 1)
})
function barHeightTenant(d) {
  if (!d || maxTenantCount.value <= 0) return 4
  return Math.max(20, (d.count / maxTenantCount.value) * 120)
}
function barHeightRevenue(d) {
  if (!d || maxRevenue.value <= 0) return 4
  return Math.max(20, (d.revenue / maxRevenue.value) * 120)
}

function logout() {
  adminAuth.logout()
  router.replace('/admin/login')
}
</script>
