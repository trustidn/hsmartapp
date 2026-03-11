<template>
  <div class="min-h-full">
    <div v-if="isOffline" class="bg-amber-50 border border-amber-200 rounded-xl p-4 mb-4 text-center">
      <p class="text-amber-800 text-sm">Anda offline. Hanya data hari ini yang tersedia.</p>
    </div>
    <!-- Selamat datang + Plan -->
    <section class="bg-gradient-to-r from-primary-500 to-primary-600 rounded-2xl p-4 mb-4 text-white shadow-sm">
      <p class="text-sm opacity-90">Selamat datang,</p>
      <p class="text-lg font-bold mt-0.5">{{ tenantName }}</p>
      <div class="flex flex-wrap gap-3 mt-3 pt-3 border-t border-white/20">
        <div>
          <p class="text-xs opacity-80">Plan</p>
          <p class="font-semibold capitalize">{{ planLabel }}</p>
        </div>
        <div>
          <p class="text-xs opacity-80">Sisa waktu</p>
          <p class="font-semibold">{{ planRemaining }}</p>
        </div>
      </div>
    </section>
    <!-- Order belum selesai -->
    <section v-if="pendingOrders.length" class="bg-white rounded-2xl shadow-sm p-4 mb-4 border border-amber-100">
      <h2 class="text-sm font-medium text-gray-700 mb-3">Order belum selesai</h2>
      <div class="space-y-2">
        <router-link
          v-for="o in pendingOrders"
          :key="o.id"
          to="/subscription"
          class="flex items-center justify-between gap-3 p-3 rounded-xl bg-gray-50 hover:bg-gray-100 transition-colors"
        >
          <div class="min-w-0">
            <p class="font-medium text-gray-800">{{ planLabelFromSlug(o.plan_slug) }}</p>
            <p class="text-xs text-gray-500">Rp {{ formatNum(o.amount_rupiah ?? 0) }} · {{ formatDate(o.created_at) }}</p>
          </div>
          <span
            :class="{
              'bg-amber-100 text-amber-700': o.status === 'pending',
              'bg-blue-100 text-blue-700': o.status === 'paid',
            }"
            class="px-2 py-1 rounded-lg text-xs font-medium shrink-0"
          >
            {{ o.status === 'pending' ? 'Belum bayar' : 'Menunggu verifikasi' }}
          </span>
        </router-link>
      </div>
      <router-link to="/subscription" class="block mt-3 text-center text-sm font-medium text-primary-600 hover:text-primary-700">
        Lihat semua →
      </router-link>
    </section>
    <!-- Filters: Today, 7 Days, 30 Days, 12 Months. Offline: only Today -->
    <div class="flex gap-2 mb-4 overflow-x-auto pb-2">
      <button
        v-for="f in filters"
        :key="f.key"
        type="button"
        class="flex-none px-4 py-2 rounded-xl font-medium text-sm whitespace-nowrap"
        :class="[
          filter === f.key ? 'bg-primary-600 text-white' : 'bg-white text-gray-600 border border-gray-200',
          isOffline && f.key !== 'today' ? 'opacity-50 pointer-events-none' : '',
        ]"
        @click="filter = f.key; load()"
      >
        {{ f.label }}
      </button>
    </div>
    <div v-if="loading" class="animate-pulse space-y-4">
      <div class="h-24 bg-primary-500/30 rounded-2xl" />
      <div class="flex gap-2">
        <div v-for="i in 4" :key="i" class="h-10 w-20 bg-gray-200 rounded-xl" />
      </div>
      <section class="bg-white rounded-2xl p-5">
        <div class="h-4 w-28 bg-gray-200 rounded mb-3" />
        <div class="grid grid-cols-2 gap-4">
          <div><div class="h-3 w-16 bg-gray-200 rounded mb-2" /><div class="h-6 w-24 bg-gray-200 rounded" /></div>
          <div><div class="h-3 w-20 bg-gray-200 rounded mb-2" /><div class="h-6 w-24 bg-gray-200 rounded" /></div>
          <div class="col-span-2 pt-2 border-t"><div class="h-3 w-12 bg-gray-200 rounded mb-2" /><div class="h-8 w-28 bg-gray-200 rounded" /></div>
          <div><div class="h-3 w-16 bg-gray-200 rounded mb-2" /><div class="h-5 w-8 bg-gray-200 rounded" /></div>
        </div>
      </section>
      <section class="bg-white rounded-2xl p-5 h-44">
        <div class="h-4 w-32 bg-gray-200 rounded mb-3" />
        <div class="h-32 bg-gray-200 rounded-xl" />
      </section>
      <section class="bg-white rounded-2xl p-5">
        <div class="h-4 w-28 bg-gray-200 rounded mb-3" />
        <div class="space-y-2">
          <div v-for="i in 4" :key="i" class="h-10 bg-gray-200 rounded-xl" />
        </div>
      </section>
    </div>
    <template v-else>
      <section class="bg-white rounded-2xl shadow-sm p-5 mb-4">
        <h2 class="text-sm font-medium text-gray-500 mb-3">Ringkasan {{ filterLabel }}</h2>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <p class="text-xs text-gray-500">Penjualan</p>
            <p class="text-xl font-bold text-gray-900">Rp {{ formatNum(data?.sales_total ?? 0) }}</p>
          </div>
          <div>
            <p class="text-xs text-gray-500">Pengeluaran</p>
            <p class="text-xl font-bold text-gray-900">Rp {{ formatNum(data?.expense_total ?? 0) }}</p>
          </div>
          <div class="col-span-2 pt-2 border-t">
            <p class="text-xs text-gray-500">Untung</p>
            <p class="text-2xl font-bold text-primary-600">Rp {{ formatNum(data?.profit ?? 0) }}</p>
          </div>
          <div>
            <p class="text-xs text-gray-500">Transaksi</p>
            <p class="text-lg font-semibold">{{ data?.transactions ?? 0 }}</p>
          </div>
        </div>
      </section>
      <section v-if="chartData?.length" class="bg-white rounded-2xl shadow-sm p-5 mb-4">
        <h2 class="text-sm font-medium text-gray-500 mb-3">Grafik Penjualan Minggu Ini (Minggu–Sabtu)</h2>
        <div class="h-40 flex items-end gap-2">
          <div
            v-for="(d, i) in chartData"
            :key="d.date"
            class="flex-1 min-w-0 flex flex-col items-center gap-1"
            :title="d.date + ': Rp ' + formatNum(d.total)"
          >
            <span class="text-[10px] font-medium text-gray-600 h-4 leading-tight text-center truncate w-full">
              {{ d.total > 0 ? 'Rp ' + formatCompact(d.total) : '' }}
            </span>
            <div
              class="w-full rounded-t bg-primary-500/80 min-h-[4px]"
              :style="{ height: Math.max(4, barHeight(d)) + 'px' }"
            />
          </div>
        </div>
        <div class="flex justify-between mt-2 text-xs text-gray-500">
          <span>{{ chartData[0]?.date ?? '' }}</span>
          <span>{{ chartData[chartData.length - 1]?.date ?? '' }}</span>
        </div>
      </section>
      <section class="bg-white rounded-2xl shadow-sm p-5">
        <h2 class="text-sm font-medium text-gray-500 mb-3">Produk Terlaris</h2>
        <ul v-if="ranking?.length" class="space-y-2">
          <li
            v-for="(r, i) in ranking"
            :key="r.product_id"
            class="flex items-center justify-between gap-3 py-2 border-b border-gray-100 last:border-0"
          >
            <span class="font-medium">{{ i + 1 }}. {{ r.product_name }}</span>
            <div class="flex items-center gap-3 shrink-0">
              <span class="text-xs text-gray-500">{{ r.qty ?? 0 }} terjual</span>
              <span class="text-primary-600 font-medium">Rp {{ formatNum(r.total) }}</span>
            </div>
          </li>
        </ul>
        <p v-else class="text-gray-400 text-sm">Belum ada penjualan.</p>
      </section>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useSettingsStore } from '../stores/settings'
import { api, isOnline, toLocalDateStr } from '../lib/api'
import * as db from '../lib/db'

const auth = useAuthStore()
const settings = useSettingsStore()
const loading = ref(true)
const filter = ref('today')
const dashboardData = ref(null)
const subscription = ref(null)
const pendingOrders = ref([])

const filters = [
  { key: 'today', label: 'Hari Ini' },
  { key: '7d', label: '7 Hari' },
  { key: '30d', label: '30 Hari' },
  { key: '12m', label: '12 Bulan' },
]

const isOffline = computed(() => !isOnline())
const filterLabel = computed(() => filters.find(f => f.key === filter.value)?.label ?? 'Hari Ini')
const tenantName = computed(() => settings.settings?.name || auth.name || 'Toko')
const planLabel = computed(() => (subscription.value?.plan || 'free').replace(/^./, (c) => c.toUpperCase()))
const planRemaining = computed(() => {
  const exp = subscription.value?.expired_at
  if (!exp) return 'Tanpa batas'
  const end = new Date(exp)
  const now = new Date()
  if (end < now) return 'Kadaluarsa'
  const ms = end - now
  const days = Math.floor(ms / (24 * 60 * 60 * 1000))
  const hours = Math.floor((ms % (24 * 60 * 60 * 1000)) / (60 * 60 * 1000))
  if (days > 0) return `${days} hari`
  return `${hours} jam`
})

const data = computed(() => {
  const d = dashboardData.value
  if (!d) return null
  if (d.today) return d.today
  if (d.range) return d.range
  return null
})

const ranking = computed(() => {
  const d = dashboardData.value
  return d?.product_rank ?? []
})

const chartData = computed(() => {
  const d = dashboardData.value
  return d?.sales_chart ?? []
})

const maxChartTotal = computed(() => {
  const arr = chartData.value
  if (!arr.length) return 1
  return Math.max(...arr.map(x => x.total), 1)
})

function barHeight(item) {
  if (!item || maxChartTotal.value <= 0) return 4
  return Math.max(20, (item.total / maxChartTotal.value) * 120)
}
function formatCompact(n) {
  const num = Number(n)
  if (num >= 1000000) return (num / 1000000).toFixed(1) + 'jt'
  if (num >= 1000) return (num / 1000).toFixed(0) + 'rb'
  return num.toLocaleString('id-ID')
}

onMounted(async () => {
  if (auth.tenantId && !settings.settings) await settings.load(auth.tenantId)
  loadSubscription()
  loadPendingOrders()
  load()
})

async function loadSubscription() {
  if (!auth.tenantId || isOffline.value) return
  try {
    subscription.value = await api.subscription.get()
  } catch {
    subscription.value = null
  }
}

async function loadPendingOrders() {
  if (!auth.tenantId || isOffline.value) return
  try {
    const res = await api.subscriptionOrders.list()
    const orders = res.orders || []
    pendingOrders.value = orders.filter((o) => o.status === 'pending' || o.status === 'paid')
  } catch {
    pendingOrders.value = []
  }
}

const planLabels = { free: 'Free', premium_1m: 'Premium 1 Bulan', premium_3m: 'Premium 3 Bulan', premium_6m: 'Premium 6 Bulan', premium_1y: 'Premium 1 Tahun', platinum: 'Platinum' }
function planLabelFromSlug(slug) {
  return planLabels[slug] || slug || '-'
}
function formatDate(s) {
  if (!s) return '-'
  try {
    return new Date(s).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' })
  } catch {
    return s
  }
}

async function load() {
  if (!auth.tenantId) return
  if (isOffline.value && filter.value !== 'today') {
    dashboardData.value = null
    loading.value = false
    return
  }
  if (isOffline.value) {
    loading.value = true
    try {
      const [summary, productRank] = await Promise.all([
        db.getTodayLocalSummary(auth.tenantId),
        db.getTodayLocalProductRank(auth.tenantId),
      ])
      dashboardData.value = {
        today: summary ? { ...summary, date: summary.date } : null,
        product_rank: productRank,
        sales_chart: [],
      }
    } catch {
      dashboardData.value = null
    } finally {
      loading.value = false
    }
    return
  }
  loading.value = true
  try {
    const today = new Date()
    const dateStr = toLocalDateStr(today)
    if (filter.value === 'today') {
      const dash = await api.report.dashboard(dateStr)
      dashboardData.value = dash
    } else {
      const to = new Date(today)
      const from = new Date(today)
      if (filter.value === '7d') from.setDate(from.getDate() - 7)
      else if (filter.value === '30d') from.setDate(from.getDate() - 30)
      else if (filter.value === '12m') from.setMonth(from.getMonth() - 12)
      const fromStr = toLocalDateStr(from)
      const toStr = toLocalDateStr(to)
      const dash = await api.report.dashboardRange(fromStr, toStr)
      dashboardData.value = dash
    }
  } catch {
    dashboardData.value = null
  } finally {
    loading.value = false
  }
}

function formatNum(n) {
  return Number(n).toLocaleString('id-ID')
}
</script>
