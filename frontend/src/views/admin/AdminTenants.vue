<template>
  <div class="min-h-full">
    <h1 class="text-xl font-bold text-gray-800">Daftar Tenant</h1>
    <p class="text-gray-500 text-sm mt-1">Kelola merchant / UMKM terdaftar</p>

    <div class="mt-6 flex flex-wrap gap-3 items-center">
      <input
        v-model="search"
        type="search"
        placeholder="Cari nama atau nomor HP..."
        class="flex-1 min-w-[200px] px-4 py-2 rounded-xl border border-gray-200 focus:ring-2 focus:ring-primary-400 focus:border-transparent bg-white"
        @input="debouncedFetch"
      />
      <button
        @click="fetchTenants"
        class="px-4 py-2 rounded-xl bg-primary-600 text-white font-medium hover:bg-primary-700"
      >
        Cari
      </button>
      <button
        @click="exportCsv"
        class="px-4 py-2 rounded-xl bg-green-600 text-white font-medium hover:bg-green-700"
      >
        Ekspor Excel
      </button>
    </div>

    <div v-if="loading" class="mt-6 space-y-4">
      <div class="h-24 bg-primary-500/20 rounded-2xl animate-pulse" />
      <div class="grid gap-3 md:hidden">
        <div v-for="i in 3" :key="i" class="h-32 bg-gray-200 rounded-2xl animate-pulse" />
      </div>
      <div class="hidden md:block bg-white rounded-2xl p-5 shadow-sm">
        <div v-for="i in 5" :key="i" class="h-12 bg-gray-100 rounded-xl mb-2 animate-pulse" />
      </div>
    </div>
    <div v-else-if="error" class="mt-6 p-4 bg-red-50 border border-red-200 rounded-xl text-red-700">{{ error }}</div>
    <template v-else>
      <!-- Desktop: table -->
      <section class="hidden md:block mt-6 overflow-x-auto">
        <div class="bg-white rounded-2xl shadow-sm overflow-hidden">
          <table class="w-full text-sm">
            <thead class="bg-gray-50">
              <tr>
                <th class="text-left p-3 font-medium text-gray-700">Nama</th>
                <th class="text-left p-3 font-medium text-gray-700">HP</th>
                <th class="text-left p-3 font-medium text-gray-700">Plan</th>
                <th class="text-left p-3 font-medium text-gray-700">Status</th>
                <th class="text-left p-3 font-medium text-gray-700">Transaksi</th>
                <th class="text-left p-3 font-medium text-gray-700">Transaksi terakhir</th>
                <th class="text-left p-3 font-medium text-gray-700">Daftar</th>
                <th class="text-left p-3 font-medium text-gray-700">Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="t in tenants" :key="t.id" class="border-b border-gray-100 last:border-0 hover:bg-gray-50/50">
                <td class="p-3 font-medium text-gray-800">{{ t.name }}</td>
                <td class="p-3 text-gray-600">{{ t.phone }}</td>
                <td class="p-3">
                  <span :class="t.plan && t.plan !== 'free' ? 'text-primary-600 font-medium' : 'text-gray-600'">
                    {{ planLabel(t.plan) }}
                  </span>
                </td>
                <td class="p-3">
                  <span
                    :class="{
                      'bg-green-100 text-green-700': t.status === 'active',
                      'bg-amber-100 text-amber-700': t.status === 'suspended',
                      'bg-gray-100 text-gray-600': t.status === 'inactive',
                    }"
                    class="px-2 py-1 rounded-lg text-xs font-medium"
                  >
                    {{ t.status }}
                  </span>
                </td>
                <td class="p-3 text-gray-600">{{ t.transactions ?? 0 }}</td>
                <td class="p-3 text-gray-500 text-xs">{{ formatDate(t.last_transaction_at) }}</td>
                <td class="p-3 text-gray-500 text-xs">{{ formatDate(t.created_at) }}</td>
                <td class="p-3">
                  <router-link
                    :to="{ name: 'AdminTenantDetail', params: { id: t.id } }"
                    class="text-primary-600 font-medium hover:text-primary-700"
                  >
                    Detail
                  </router-link>
                </td>
              </tr>
            </tbody>
          </table>
          <p v-if="tenants.length === 0" class="p-6 text-gray-500 text-center">Belum ada tenant.</p>
        </div>
      </section>

      <!-- Mobile: cards -->
      <div class="mt-6 md:hidden space-y-3">
        <div
          v-for="t in tenants"
          :key="t.id"
          class="bg-white rounded-2xl shadow-sm border border-gray-100 p-4"
        >
          <div class="flex items-start justify-between gap-2">
            <div class="min-w-0">
              <h3 class="font-semibold text-gray-800 truncate">{{ t.name }}</h3>
              <p class="text-sm text-gray-500 mt-0.5">{{ t.phone }}</p>
            </div>
            <span
              :class="{
                'bg-green-100 text-green-700': t.status === 'active',
                'bg-amber-100 text-amber-700': t.status === 'suspended',
                'bg-gray-100 text-gray-600': t.status === 'inactive',
              }"
              class="px-2 py-1 rounded-lg text-xs font-medium shrink-0"
            >
              {{ t.status }}
            </span>
          </div>
          <div class="mt-3 flex items-center justify-between text-sm text-gray-600">
            <span>
              Plan: <span :class="t.plan && t.plan !== 'free' ? 'text-primary-600 font-medium' : ''">{{ planLabel(t.plan) }}</span>
            </span>
          </div>
          <div class="mt-2 flex items-center gap-4 text-xs text-gray-500">
            <span>{{ t.transactions ?? 0 }} transaksi</span>
            <span>Terakhir: {{ formatDate(t.last_transaction_at) }}</span>
          </div>
          <router-link
            :to="{ name: 'AdminTenantDetail', params: { id: t.id } }"
            class="mt-3 block text-center py-2 rounded-xl bg-primary-50 text-primary-600 font-medium text-sm"
          >
            Detail
          </router-link>
        </div>
        <p v-if="tenants.length === 0" class="p-6 text-gray-500 text-center">Belum ada tenant.</p>
      </div>
    </template>

    <div v-if="total > limit" class="mt-4 flex items-center gap-4 text-sm text-gray-600">
      <span>Total {{ total }} tenant</span>
      <div class="flex gap-2">
        <button
          :disabled="offset === 0"
          @click="goPage(-1)"
          class="px-3 py-2 rounded-xl bg-white border border-gray-200 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
        >
          Sebelumnya
        </button>
        <button
          :disabled="offset + limit >= total"
          @click="goPage(1)"
          class="px-3 py-2 rounded-xl bg-white border border-gray-200 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
        >
          Selanjutnya
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { adminApi } from '../../lib/adminApi'

const tenants = ref([])
const total = ref(0)
const limit = ref(20)
const offset = ref(0)
const search = ref('')
const loading = ref(false)
const error = ref('')
let debounceTimer = null

const planLabels = { free: 'Free', premium: 'Premium', premium_1m: '1 Bulan', premium_3m: '3 Bulan', premium_6m: '6 Bulan', premium_1y: '1 Tahun', platinum: 'Platinum' }
function planLabel(slug) {
  return planLabels[slug] || slug || '-'
}
function formatDate(s) {
  if (!s) return '-'
  try {
    const d = new Date(s)
    return d.toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' })
  } catch {
    return s
  }
}

async function fetchTenants() {
  loading.value = true
  error.value = ''
  try {
    const res = await adminApi.tenants.list({ limit: limit.value, offset: offset.value, search: search.value })
    tenants.value = res.tenants || []
    total.value = res.total ?? 0
  } catch (e) {
    error.value = e.message || 'Gagal memuat tenant'
  } finally {
    loading.value = false
  }
}

function debouncedFetch() {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    offset.value = 0
    fetchTenants()
  }, 300)
}

function goPage(delta) {
  offset.value = Math.max(0, offset.value + delta * limit.value)
  fetchTenants()
}

function exportCsv() {
  const rows = [['Nama', 'HP', 'Plan', 'Status', 'Transaksi', 'Transaksi Terakhir', 'Daftar']]
  for (const t of tenants.value) {
    rows.push([
      t.name || '',
      t.phone || '',
      planLabel(t.plan),
      t.status || '',
      t.transactions ?? 0,
      formatDate(t.last_transaction_at),
      formatDate(t.created_at),
    ])
  }
  const csv = rows.map((r) => r.map((c) => `"${String(c).replace(/"/g, '""')}"`).join(',')).join('\n')
  const blob = new Blob(['\uFEFF' + csv], { type: 'text/csv;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `tenant-${new Date().toISOString().slice(0, 10)}.csv`
  a.click()
  URL.revokeObjectURL(url)
}

onMounted(fetchTenants)
</script>
