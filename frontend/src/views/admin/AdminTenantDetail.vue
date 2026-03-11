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
        <div class="h-40 bg-gray-200 rounded-2xl animate-pulse" />
        <div class="h-40 bg-gray-200 rounded-2xl animate-pulse" />
      </div>
    </div>
    <div v-else-if="error" class="mt-6 p-4 bg-red-50 border border-red-200 rounded-xl text-red-700">{{ error }}</div>
    <template v-else-if="tenant">
      <h1 class="mt-2 text-xl font-bold text-gray-800">{{ tenant.name }}</h1>
      <p class="text-gray-500 text-sm mt-1">{{ tenant.phone }}</p>

      <div class="mt-6 grid gap-4 sm:grid-cols-2 max-w-2xl">
        <section class="p-5 bg-white rounded-2xl shadow-sm border border-gray-100">
          <h2 class="text-sm font-medium text-gray-500 mb-3">Informasi</h2>
          <dl class="space-y-2 text-sm">
            <div>
              <dt class="text-gray-500">Plan</dt>
              <dd :class="tenant.plan === 'premium' ? 'text-primary-600 font-medium' : 'text-gray-800'">{{ tenant.plan }}</dd>
            </div>
            <div>
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
            <div>
              <dt class="text-gray-500">Terdaftar</dt>
              <dd class="text-gray-800">{{ formatDate(tenant.created_at) }}</dd>
            </div>
          </dl>
        </section>
        <section v-if="tenant.subscription" class="p-5 bg-white rounded-2xl shadow-sm border border-gray-100">
          <h2 class="text-sm font-medium text-gray-500 mb-3">Subscription</h2>
          <dl class="space-y-2 text-sm">
            <div>
              <dt class="text-gray-500">Plan</dt>
              <dd class="text-gray-800">{{ tenant.subscription.plan }}</dd>
            </div>
            <div>
              <dt class="text-gray-500">Status</dt>
              <dd class="text-gray-800">{{ tenant.subscription.status }}</dd>
            </div>
            <div>
              <dt class="text-gray-500">Mulai</dt>
              <dd class="text-gray-800">{{ formatDate(tenant.subscription.started_at) }}</dd>
            </div>
            <div>
              <dt class="text-gray-500">Kadaluarsa</dt>
              <dd class="text-gray-800">{{ tenant.subscription.expired_at ? formatDate(tenant.subscription.expired_at) : 'Tak terbatas' }}</dd>
            </div>
          </dl>
        </section>
      </div>

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

const route = useRoute()
const tenant = ref(null)
const loading = ref(true)
const error = ref('')
const statusLoading = ref(false)

const id = computed(() => route.params.id)

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
  } catch (e) {
    alert(e.message || 'Gagal mengubah status')
  } finally {
    statusLoading.value = false
  }
}

onMounted(fetchTenant)
</script>
