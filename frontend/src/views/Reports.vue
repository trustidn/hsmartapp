<template>
  <div class="min-h-full">
    <div v-if="!isOnline()" class="bg-amber-50 border border-amber-200 rounded-xl p-3 mb-4 text-amber-800 text-sm">
      Offline — laporan memerlukan koneksi internet.
    </div>

    <!-- Banner limit / expired (hanya setelah data subscription ter-load) -->
    <div v-if="!planLimits.loading && (planLimits.reportLimitMessage || planLimits.isExpired)" class="mb-4 space-y-2">
      <div
        v-if="planLimits.isExpired"
        class="bg-amber-50 border border-amber-200 rounded-xl p-4 flex items-center justify-between gap-3"
      >
        <p class="text-amber-800 text-sm">Langganan Anda telah kadaluarsa. Laporan dibatasi.</p>
        <router-link to="/app/subscription" class="shrink-0 px-4 py-2 rounded-xl bg-amber-600 text-white text-sm font-medium">
          Perpanjang
        </router-link>
      </div>
      <div
        v-else-if="planLimits.reportLimitMessage"
        class="bg-blue-50 border border-blue-200 rounded-xl p-4 flex items-center justify-between gap-3"
      >
        <p class="text-blue-800 text-sm">{{ planLimits.reportLimitMessage }}</p>
        <router-link to="/app/subscription" class="shrink-0 px-4 py-2 rounded-xl bg-primary-600 text-white text-sm font-medium">
          Upgrade
        </router-link>
      </div>
    </div>

    <!-- Period filter -->
    <div class="flex gap-2 mb-4 overflow-x-auto pb-2">
      <button
        v-for="f in reportFilters"
        :key="f.key"
        type="button"
        class="flex-none px-4 py-2 rounded-xl font-medium text-sm whitespace-nowrap"
        :class="[
          filter === f.key ? 'bg-primary-600 text-white' : 'bg-white text-gray-600 border border-gray-200',
          !isOnline() && f.key !== 'today' ? 'opacity-50 pointer-events-none' : '',
        ]"
        @click="filter = f.key; load()"
      >
        {{ f.label }}
      </button>
    </div>

    <div v-if="loading" class="animate-pulse space-y-4">
      <section class="bg-white rounded-2xl p-5">
        <div class="h-4 w-24 bg-gray-200 rounded mb-3" />
        <div class="grid grid-cols-3 gap-4">
          <div v-for="i in 3" :key="i"><div class="h-3 w-12 bg-gray-200 rounded mb-2" /><div class="h-6 w-20 bg-gray-200 rounded" /></div>
        </div>
      </section>
      <section class="bg-white rounded-2xl p-5">
        <div class="h-4 w-28 bg-gray-200 rounded mb-3" />
        <div class="space-y-2">
          <div v-for="i in 4" :key="i" class="h-14 bg-gray-200 rounded-xl" />
        </div>
      </section>
      <section class="bg-white rounded-2xl p-5">
        <div class="h-4 w-24 bg-gray-200 rounded mb-3" />
        <div class="space-y-2">
          <div v-for="i in 3" :key="i" class="h-12 bg-gray-200 rounded-xl" />
        </div>
      </section>
    </div>

    <template v-else>
      <!-- Summary section -->
      <section class="bg-white rounded-2xl shadow-sm p-5 mb-4">
        <div class="flex flex-wrap items-center justify-between gap-3 mb-3">
          <h2 class="text-sm font-medium text-gray-500">Ringkasan {{ filterLabel }}</h2>
          <div class="flex gap-2">
            <button
              type="button"
              class="px-4 py-2 rounded-xl text-sm font-medium border border-gray-200 text-gray-700 hover:bg-gray-50 disabled:opacity-50"
              :disabled="exporting"
              @click="exportExcel"
            >
              {{ exporting ? '...' : 'Ekspor Excel' }}
            </button>
            <button
              type="button"
              class="px-4 py-2 rounded-xl text-sm font-medium bg-red-600 text-white hover:bg-red-700 disabled:opacity-50"
              :disabled="exporting"
              @click="exportPdf"
            >
              {{ exporting ? '...' : 'Ekspor PDF' }}
            </button>
          </div>
        </div>
        <div class="grid grid-cols-3 gap-4">
          <div>
            <p class="text-xs text-gray-500">Penjualan</p>
            <p class="text-lg font-bold text-gray-900">Rp {{ formatNum(summary?.sales_total ?? 0) }}</p>
          </div>
          <div>
            <p class="text-xs text-gray-500">Pengeluaran</p>
            <p class="text-lg font-bold text-gray-900">Rp {{ formatNum(summary?.expense_total ?? 0) }}</p>
          </div>
          <div>
            <p class="text-xs text-gray-500">Untung</p>
            <p class="text-lg font-bold" :class="profitPct >= 0 ? 'text-primary-600' : 'text-red-600'">
              Rp {{ formatNum(summary?.profit ?? 0) }}
            </p>
            <p class="text-xs" :class="profitPct >= 0 ? 'text-primary-600' : 'text-red-600'">
              ({{ profitPct.toFixed(1) }}%)
            </p>
          </div>
        </div>
      </section>

      <!-- Daftar Transaksi -->
      <section class="bg-white rounded-2xl shadow-sm p-4 mb-4">
        <h2 class="text-sm font-medium text-gray-500 mb-3">Daftar Transaksi</h2>
        <div v-if="displayedTransactions.length" class="space-y-2">
          <div
            v-for="s in displayedTransactions"
            :key="s.id"
            class="flex justify-between items-center p-4 rounded-xl border border-gray-100"
          >
            <div>
              <p class="font-semibold text-gray-900">Rp {{ formatNum(s.total) }}</p>
              <p class="text-xs text-gray-500">{{ formatDateTime(s.created_at) }} · {{ (s.payment_method || 'cash') }}</p>
            </div>
            <button
              type="button"
              class="min-h-touch px-3 py-2 rounded-lg border border-primary-200 text-primary-600 text-sm font-medium"
              @click="viewReceipt(s)"
            >
              Struk
            </button>
          </div>
        </div>
        <p v-else class="text-center text-gray-400 py-6">Belum ada transaksi.</p>

        <!-- Pagination transaksi -->
        <div v-if="(displayedTransactions.length > 0 || txPage > 1) && (txPage > 1 || hasMoreTx)" class="flex justify-center gap-2 mt-4">
          <button
            type="button"
            class="px-4 py-2 rounded-lg border text-sm font-medium disabled:opacity-50"
            :disabled="txPage <= 1"
            @click="txPage--; loadTransactions()"
          >
            ← Sebelumnya
          </button>
          <span class="py-2 text-sm text-gray-500">Halaman {{ txPage }}</span>
          <button
            type="button"
            class="px-4 py-2 rounded-lg border text-sm font-medium disabled:opacity-50"
            :disabled="!hasMoreTx"
            @click="txPage++; loadTransactions()"
          >
            Selanjutnya →
          </button>
        </div>
      </section>

      <!-- Daftar Pengeluaran -->
      <section class="bg-white rounded-2xl shadow-sm p-4">
        <h2 class="text-sm font-medium text-gray-500 mb-3">Daftar Pengeluaran</h2>
        <div v-if="displayedExpenses.length" class="space-y-2">
          <div
            v-for="e in displayedExpenses"
            :key="e.id"
            class="flex justify-between items-center py-3 px-4 rounded-xl border border-gray-100"
          >
            <div>
              <p class="font-medium text-gray-900">{{ e.name }}</p>
              <p class="text-xs text-gray-500">{{ formatDateTime(e.created_at) }}</p>
            </div>
            <span class="font-semibold text-red-600">- Rp {{ formatNum(e.amount) }}</span>
          </div>
        </div>
        <p v-else class="text-center text-gray-400 py-6">Belum ada pengeluaran.</p>

        <!-- Pagination pengeluaran -->
        <div v-if="(displayedExpenses.length > 0 || expPage > 1) && (expPage > 1 || hasMoreExp)" class="flex justify-center gap-2 mt-4">
          <button
            type="button"
            class="px-4 py-2 rounded-lg border text-sm font-medium disabled:opacity-50"
            :disabled="expPage <= 1"
            @click="expPage--; loadExpenses()"
          >
            ← Sebelumnya
          </button>
          <span class="py-2 text-sm text-gray-500">Halaman {{ expPage }}</span>
          <button
            type="button"
            class="px-4 py-2 rounded-lg border text-sm font-medium disabled:opacity-50"
            :disabled="!hasMoreExp"
            @click="expPage++; loadExpenses()"
          >
            Selanjutnya →
          </button>
        </div>
      </section>
    </template>

    <!-- Receipt modal -->
    <div
      v-if="selectedSale"
      class="fixed inset-0 z-50 bg-black/50 flex items-end sm:items-center justify-center p-4"
      @click.self="selectedSale = null"
    >
      <div class="bg-white rounded-2xl max-w-sm w-full max-h-[80vh] overflow-auto p-5">
        <div class="flex justify-between items-center mb-4">
          <h3 class="font-semibold text-lg">Struk</h3>
          <button type="button" class="p-2 text-gray-500" @click="selectedSale = null">✕</button>
        </div>
        <div id="receipt-print" class="print:bg-white print:text-black">
          <ReceiptContent :sale="selectedSale" :settings="settings" />
        </div>
        <div class="mt-4 flex gap-2">
          <button
            type="button"
            class="flex-1 py-2 rounded-xl border border-gray-300 text-sm font-medium"
            @click="printReceipt"
          >
            Cetak
          </button>
          <button
            type="button"
            class="flex-1 py-2 rounded-xl border border-gray-300 text-sm font-medium"
            @click="selectedSale = null"
          >
            Tutup
          </button>
          <button
            type="button"
            class="flex-1 py-2 rounded-xl bg-primary-600 text-white text-sm font-medium"
            @click="downloadPdf"
          >
            PDF
          </button>
          <button
            type="button"
            class="flex-1 py-2 rounded-xl bg-green-600 text-white text-sm font-medium"
            @click="shareWhatsApp"
          >
            WhatsApp
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useSettingsStore } from '../stores/settings'
import { useProductsStore } from '../stores/products'
import { usePlanLimits } from '../composables/usePlanLimits'
import { api, isOnline, toLocalDateStr } from '../lib/api'
import * as db from '../lib/db'
import ReceiptContent from '../components/ReceiptContent.vue'
import { exportReportPdf, exportReportCsv } from '../lib/exportReport'

const auth = useAuthStore()
const settingsStore = useSettingsStore()
const productsStore = useProductsStore()
const planLimits = usePlanLimits(auth, productsStore)
const settings = computed(() => settingsStore.settings)
const loading = ref(true)
const filter = ref('today')
const summary = ref(null)
const transactions = ref([])
const expenses = ref([])
const allTransactions = ref([])
const allExpenses = ref([])
const selectedSale = ref(null)
const txPage = ref(1)
const exporting = ref(false)
const expPage = ref(1)
const pageSize = 10

const reportFilters = computed(() => {
  const allowed = planLimits.allowedReportFilters.value
  if (allowed.length > 0) return allowed
  return [
    { key: 'today', label: 'Hari Ini', days: 1 },
    { key: '7d', label: '7 Hari', days: 7 },
    { key: '30d', label: '30 Hari', days: 30 },
    { key: '12m', label: '12 Bulan', days: 365 },
  ]
})

const filterLabel = computed(() => reportFilters.value.find((f) => f.key === filter.value)?.label ?? 'Hari Ini')

const displayedTransactions = computed(() => {
  if (allTransactions.value.length > 0) {
    const start = (txPage.value - 1) * pageSize
    return allTransactions.value.slice(start, start + pageSize)
  }
  return transactions.value
})
const displayedExpenses = computed(() => {
  if (allExpenses.value.length > 0) {
    const start = (expPage.value - 1) * pageSize
    return allExpenses.value.slice(start, start + pageSize)
  }
  return expenses.value
})
const hasMoreTx = computed(() => {
  if (allTransactions.value.length > 0) {
    return txPage.value * pageSize < allTransactions.value.length
  }
  return transactions.value.length >= pageSize
})
const hasMoreExp = computed(() => {
  if (allExpenses.value.length > 0) {
    return expPage.value * pageSize < allExpenses.value.length
  }
  return expenses.value.length >= pageSize
})

const profitPct = computed(() => {
  const s = summary.value
  if (!s || (s.sales_total ?? 0) <= 0) return 0
  const profit = (s.profit ?? 0)
  return (profit / s.sales_total) * 100
})

function normalizeSalesList(res) {
  if (Array.isArray(res)) return res
  if (res && typeof res === 'object') {
    const arr = res.data ?? res.sales ?? res.list
    return Array.isArray(arr) ? arr : []
  }
  return []
}

function getDateRange() {
  const today = new Date()
  let from = new Date(today)
  if (filter.value === 'today') {
    return { fromStr: toLocalDateStr(today), toStr: toLocalDateStr(today) }
  }
  if (filter.value === '7d') from.setDate(from.getDate() - 7)
  else if (filter.value === '30d') from.setDate(from.getDate() - 30)
  else if (filter.value === '12m') from.setMonth(from.getMonth() - 12)
  return { fromStr: toLocalDateStr(from), toStr: toLocalDateStr(today) }
}

async function load() {
  if (!auth.tenantId) return
  if (!isOnline() && filter.value !== 'today') {
    summary.value = null
    transactions.value = []
    expenses.value = []
    allTransactions.value = []
    allExpenses.value = []
    loading.value = false
    return
  }
  txPage.value = 1
  expPage.value = 1
  loading.value = true
  try {
    const { fromStr, toStr } = getDateRange()
    if (!isOnline()) {
      const [s, tx, exp] = await Promise.all([
        db.getTodayLocalSummary(auth.tenantId),
        db.getTodayLocalSales(auth.tenantId, 50),
        db.getTodayLocalExpenses(auth.tenantId),
      ])
      summary.value = s
      const txList = (tx || []).map((r) => ({ ...r, payment_method: r.payment_method || 'cash' }))
      const expList = (exp || []).map((e) => ({ ...e, created_at: e.created_at }))
      allTransactions.value = txList
      allExpenses.value = expList
      transactions.value = []
      expenses.value = []
    } else {
      const [dash, txRes, expRes] = await Promise.all([
        filter.value === 'today'
          ? api.report.dashboard(fromStr)
          : api.report.dashboardRange(fromStr, toStr),
        api.sales.list(fromStr, toStr, undefined, pageSize, 0).catch(() => []),
        api.expenses.list(fromStr, toStr, undefined, pageSize, 0).catch(() => []),
      ])
      summary.value = filter.value === 'today' ? dash?.today : dash?.range
      allTransactions.value = []
      allExpenses.value = []
      transactions.value = normalizeSalesList(txRes)
      expenses.value = Array.isArray(expRes) ? expRes : []
    }
  } catch {
    summary.value = null
    transactions.value = []
    expenses.value = []
    allTransactions.value = []
    allExpenses.value = []
  } finally {
    loading.value = false
  }
}

async function loadTransactions() {
  if (!auth.tenantId || !isOnline()) return
  const { fromStr, toStr } = getDateRange()
  const offset = (txPage.value - 1) * pageSize
  try {
    const res = await api.sales.list(fromStr, toStr, undefined, pageSize, offset)
    transactions.value = normalizeSalesList(res)
  } catch {
    transactions.value = []
  }
}

async function loadExpenses() {
  if (!auth.tenantId || !isOnline()) return
  const { fromStr, toStr } = getDateRange()
  const offset = (expPage.value - 1) * pageSize
  try {
    const res = await api.expenses.list(fromStr, toStr, undefined, pageSize, offset)
    expenses.value = Array.isArray(res) ? res : []
  } catch {
    expenses.value = []
  }
}

async function viewReceipt(sale) {
  if (String(sale.id).startsWith('local')) {
    selectedSale.value = sale
    return
  }
  try {
    const full = await api.sales.get(sale.id)
    selectedSale.value = full
  } catch {
    selectedSale.value = sale
  }
}

function formatNum(n) {
  return Number(n).toLocaleString('id-ID')
}
function formatDateTime(s) {
  if (!s) return ''
  return new Date(s).toLocaleString('id-ID', {
    day: 'numeric',
    month: 'short',
    hour: '2-digit',
    minute: '2-digit',
  })
}

async function exportPdf() {
  exporting.value = true
  try {
    const { tx, exp } = await fetchAllForExport()
    exportReportPdf(summary.value, tx, exp, filterLabel.value, settings.value)
  } finally {
    exporting.value = false
  }
}

async function exportExcel() {
  exporting.value = true
  try {
    const { tx, exp } = await fetchAllForExport()
    exportReportCsv(summary.value, tx, exp, filterLabel.value)
  } finally {
    exporting.value = false
  }
}

async function fetchAllForExport() {
  const { fromStr, toStr } = getDateRange()
  if (!isOnline()) {
    return {
      tx: allTransactions.value,
      exp: allExpenses.value,
    }
  }
  try {
    const [txRes, expRes] = await Promise.all([
      api.sales.list(fromStr, toStr, undefined, 9999, 0),
      api.expenses.list(fromStr, toStr, undefined, 9999, 0),
    ])
    return {
      tx: normalizeSalesList(txRes),
      exp: Array.isArray(expRes) ? expRes : [],
    }
  } catch {
    return { tx: transactions.value, exp: expenses.value }
  }
}

function printReceipt() {
  if (!selectedSale.value) return
  try {
    window.print()
  } catch (e) {
    alert('Gagal membuka dialog cetak.')
  }
}

function downloadPdf() {
  if (!selectedSale.value) return
  import('../lib/receipt').then(({ exportReceiptPdf }) => {
    exportReceiptPdf(selectedSale.value, settings.value)
  })
}

/** WhatsApp - alur HSGoMart: 62 prefix untuk Indonesia */
function shareWhatsApp() {
  if (!selectedSale.value) return
  import('../lib/receipt').then(({ getReceiptWhatsAppText }) => {
    const text = getReceiptWhatsAppText(selectedSale.value, settings.value)
    const raw = (settingsStore.settings?.whatsapp_number || '').trim().replace(/\D/g, '')
    const waPhone = raw ? (raw.startsWith('62') ? raw : '62' + raw.replace(/^0+/, '')) : ''
    const url = waPhone
      ? `https://wa.me/${waPhone}?text=${encodeURIComponent(text)}`
      : `https://wa.me/?text=${encodeURIComponent(text)}`
    window.open(url, '_blank')
  })
}

onMounted(async () => {
  if (auth.tenantId) await settingsStore.load(auth.tenantId)
  await planLimits.load()
  syncFilterToAllowed()
  load()
})

watch(
  () => planLimits.allowedReportFilters.value,
  () => syncFilterToAllowed(),
  { deep: true }
)
function syncFilterToAllowed() {
  const allowed = planLimits.allowedReportFilters.value || []
  if (allowed.length && !allowed.find((f) => f.key === filter.value)) {
    filter.value = allowed[0].key
    load()
  }
}
</script>
