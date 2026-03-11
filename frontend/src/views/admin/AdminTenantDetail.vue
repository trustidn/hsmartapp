<template>
  <div class="min-h-full">
    <router-link
      to="/admin/tenants"
      class="inline-flex items-center gap-1 text-primary-600 hover:text-primary-700 text-sm font-medium mb-4"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
      </svg>
      Kembali ke Daftar Tenant
    </router-link>

    <div v-if="loading" class="mt-6 space-y-4">
      <div class="h-8 w-48 bg-gray-200 rounded-xl animate-pulse" />
      <div class="grid gap-4 sm:grid-cols-2">
        <div class="h-32 bg-gray-200 rounded-2xl animate-pulse" />
        <div class="h-32 bg-gray-200 rounded-2xl animate-pulse" />
      </div>
    </div>
    <div v-else-if="error" class="mt-6 p-4 bg-red-50 border border-red-200 rounded-xl text-red-700">{{ error }}</div>
    <template v-else-if="tenant">
      <!-- Dashboard-style header -->
      <section class="bg-gradient-to-r from-primary-500 to-primary-600 rounded-2xl p-5 mb-6 text-white shadow-sm">
        <p class="text-sm opacity-90">Detail Tenant</p>
        <h1 class="text-xl font-bold mt-0.5">{{ tenant.name }}</h1>
        <p class="text-sm opacity-90 mt-1">{{ tenant.phone }}</p>
      </section>

      <div class="grid gap-4 md:grid-cols-2">
        <!-- Info card -->
        <section class="p-5 bg-white rounded-2xl shadow-sm border border-gray-100">
          <h2 class="text-sm font-medium text-gray-500 mb-3">Informasi</h2>
          <dl class="space-y-2 text-sm">
            <div class="flex justify-between">
              <dt class="text-gray-500">Status</dt>
              <dd>
                <span
                  :class="{
                    'bg-green-100 text-green-700': tenant.status === 'active',
                    'bg-amber-100 text-amber-700': tenant.status === 'suspended',
                    'bg-gray-100 text-gray-600': tenant.status === 'inactive',
                  }"
                  class="px-2 py-1 rounded-lg text-xs font-medium"
                >
                  {{ tenant.status }}
                </span>
              </dd>
            </div>
            <div class="flex justify-between">
              <dt class="text-gray-500">Terdaftar</dt>
              <dd class="text-gray-800">{{ formatDate(tenant.created_at) }}</dd>
            </div>
          </dl>
        </section>

        <!-- Plan & masa aktif -->
        <section class="p-5 bg-white rounded-2xl shadow-sm border border-gray-100">
          <h2 class="text-sm font-medium text-gray-500 mb-3">Plan & Masa Aktif</h2>
          <dl class="space-y-2 text-sm mb-4">
            <div class="flex justify-between">
              <dt class="text-gray-500">Plan saat ini</dt>
              <dd :class="tenant.plan && tenant.plan !== 'free' ? 'text-primary-600 font-medium' : 'text-gray-800'">
                {{ planLabel(tenant.plan) }}
              </dd>
            </div>
            <div v-if="tenant.subscription" class="flex justify-between">
              <dt class="text-gray-500">Masa aktif</dt>
              <dd class="text-gray-800">
                {{ tenant.subscription.expired_at ? formatDate(tenant.subscription.expired_at) : 'Tak terbatas' }}
              </dd>
            </div>
          </dl>

          <!-- Upgrade / Perpanjang: tambah langganan (masa terakumulasi) -->
          <div>
            <label class="block text-sm text-gray-600 mb-2">Upgrade / Perpanjang Langganan</label>
            <p class="text-xs text-gray-500 mb-2">Pilih plan lalu tambah. Masa aktif akan terakumulasi dengan langganan sebelumnya.</p>
            <select
              v-model="subForm.plan"
              class="w-full px-4 py-2.5 rounded-xl border border-gray-200 focus:ring-2 focus:ring-primary-400 bg-white"
            >
              <option v-for="opt in planOptions" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
            </select>
            <button
              @click="addSubscription"
              :disabled="subLoading"
              class="mt-3 w-full px-4 py-2 rounded-xl bg-primary-600 text-white font-medium hover:bg-primary-700 disabled:opacity-50"
            >
              {{ subLoading ? '...' : 'Tambah' }}
            </button>
          </div>
        </section>
      </div>

      <!-- Subscription history -->
      <section v-if="tenant.subscription_history?.length" class="mt-6 p-5 bg-white rounded-2xl shadow-sm border border-gray-100">
        <h2 class="text-sm font-medium text-gray-500 mb-3">Riwayat Langganan</h2>
        <div class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead>
              <tr class="border-b border-gray-200">
                <th class="text-left py-2 font-medium text-gray-600">Plan</th>
                <th class="text-left py-2 font-medium text-gray-600">Mulai</th>
                <th class="text-left py-2 font-medium text-gray-600">Kadaluarsa</th>
                <th class="text-left py-2 font-medium text-gray-600">Status</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(s, i) in tenant.subscription_history" :key="i" class="border-b border-gray-100 last:border-0">
                <td class="py-2">{{ planLabel(s.plan) }}</td>
                <td class="py-2 text-gray-600">{{ formatDate(s.started_at) }}</td>
                <td class="py-2 text-gray-600">{{ s.expired_at ? formatDate(s.expired_at) : '-' }}</td>
                <td class="py-2"><span class="px-2 py-0.5 rounded text-xs" :class="s.status === 'active' ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-600'">{{ s.status }}</span></td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <!-- Status actions -->
      <div class="mt-6 flex gap-3 flex-wrap">
        <button
          v-if="tenant.status !== 'active'"
          @click="updateStatus('active')"
          :disabled="statusLoading"
          class="px-4 py-2 rounded-xl bg-green-600 text-white font-medium hover:bg-green-700 disabled:opacity-50"
        >
          {{ statusLoading ? '...' : 'Aktifkan' }}
        </button>
        <button
          v-if="tenant.status === 'active'"
          @click="updateStatus('suspended')"
          :disabled="statusLoading"
          class="px-4 py-2 rounded-xl bg-amber-600 text-white font-medium hover:bg-amber-700 disabled:opacity-50"
        >
          {{ statusLoading ? '...' : 'Suspensikan' }}
        </button>
        <button
          v-if="tenant.status !== 'inactive'"
          @click="updateStatus('inactive')"
          :disabled="statusLoading"
          class="px-4 py-2 rounded-xl bg-gray-500 text-white font-medium hover:bg-gray-600 disabled:opacity-50"
        >
          {{ statusLoading ? '...' : 'Nonaktifkan' }}
        </button>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { adminApi } from '../../lib/adminApi'
import { useToast } from '../../composables/useToast'

const route = useRoute()
const { show: toast } = useToast()
const tenant = ref(null)
const loading = ref(true)
const error = ref('')
const statusLoading = ref(false)
const subLoading = ref(false)
const subForm = ref({ plan: 'free' })

const id = computed(() => route.params.id)
const planOptions = ref([])

onMounted(async () => {
  await fetchTenant()
  try {
    const res = await adminApi.plans.list()
    const active = (res.plans || []).filter((p) => p.is_active !== false && p.plan_slug !== 'free')
    planOptions.value = active.map((p) => ({ value: p.plan_slug, label: p.name }))
    if (planOptions.value.length && !planOptions.value.find((o) => o.value === subForm.value.plan)) {
      subForm.value.plan = planOptions.value[0].value
    }
  } catch {
    planOptions.value = [
      { value: 'premium_1m', label: 'Premium 1 Bulan' },
      { value: 'premium_3m', label: 'Premium 3 Bulan' },
      { value: 'premium_6m', label: 'Premium 6 Bulan' },
      { value: 'premium_1y', label: 'Premium 1 Tahun' },
      { value: 'platinum', label: 'Platinum' },
    ]
    subForm.value.plan = 'premium_1m'
  }
})

function planLabel(slug) {
  const o = planOptions.value.find((p) => p.value === slug)
  if (o) return o.label
  const fallback = { free: 'Free', premium_1m: 'Premium 1 Bulan', premium_3m: 'Premium 3 Bulan', premium_6m: 'Premium 6 Bulan', premium_1y: 'Premium 1 Tahun', platinum: 'Platinum' }
  return fallback[slug] || slug || '-'
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

async function fetchTenant() {
  if (!id.value) return
  loading.value = true
  error.value = ''
  try {
    tenant.value = await adminApi.tenants.get(id.value)
    if (tenant.value) {
      subForm.value.plan = subForm.value.plan || tenant.value.plan || tenant.value.subscription?.plan || 'premium_1m'
    }
  } catch (e) {
    error.value = e.message || 'Gagal memuat tenant'
  } finally {
    loading.value = false
  }
}

async function updateStatus(status) {
  statusLoading.value = true
  try {
    await adminApi.tenants.updateStatus(id.value, status)
    tenant.value = { ...tenant.value, status }
    toast('Status berhasil diubah')
  } catch (e) {
    toast(e.message || 'Gagal mengubah status', 'error')
  } finally {
    statusLoading.value = false
  }
}

async function addSubscription() {
  subLoading.value = true
  try {
    await adminApi.tenants.updateSubscription(id.value, { plan: subForm.value.plan })
    await fetchTenant()
    toast('Langganan berhasil ditambah')
  } catch (e) {
    toast(e.message || 'Gagal menambah langganan', 'error')
  } finally {
    subLoading.value = false
  }
}
</script>
