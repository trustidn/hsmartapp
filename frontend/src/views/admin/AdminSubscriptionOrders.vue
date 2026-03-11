<template>
  <div class="min-h-full">
    <h1 class="text-xl font-bold text-gray-800">Order Langganan</h1>
    <p class="text-gray-500 text-sm mt-1">Verifikasi dan setujui order upgrade/perpanjang dari tenant</p>

    <div class="mt-4 flex gap-2">
      <button
        @click="statusFilter = ''"
        :class="statusFilter === '' ? 'bg-primary-600 text-white' : 'bg-white border border-gray-200 text-gray-600'"
        class="px-3 py-2 rounded-xl font-medium text-sm"
      >
        Semua
      </button>
      <button
        @click="statusFilter = 'pending'"
        :class="statusFilter === 'pending' ? 'bg-primary-600 text-white' : 'bg-white border border-gray-200 text-gray-600'"
        class="px-3 py-2 rounded-xl font-medium text-sm"
      >
        Belum bayar
      </button>
      <button
        @click="statusFilter = 'paid'"
        :class="statusFilter === 'paid' ? 'bg-primary-600 text-white' : 'bg-white border border-gray-200 text-gray-600'"
        class="px-3 py-2 rounded-xl font-medium text-sm"
      >
        Sudah bayar
      </button>
      <button
        @click="statusFilter = 'approved'"
        :class="statusFilter === 'approved' ? 'bg-primary-600 text-white' : 'bg-white border border-gray-200 text-gray-600'"
        class="px-3 py-2 rounded-xl font-medium text-sm"
      >
        Disetujui
      </button>
      <button
        @click="statusFilter = 'rejected'"
        :class="statusFilter === 'rejected' ? 'bg-primary-600 text-white' : 'bg-white border border-gray-200 text-gray-600'"
        class="px-3 py-2 rounded-xl font-medium text-sm"
      >
        Ditolak
      </button>
    </div>

    <div v-if="loading" class="mt-6 space-y-3">
      <div v-for="i in 5" :key="i" class="h-24 bg-gray-200 rounded-2xl animate-pulse" />
    </div>
    <div v-else-if="error" class="mt-6 p-4 bg-red-50 border border-red-200 rounded-xl text-red-700">{{ error }}</div>
    <div v-else class="mt-6 space-y-3">
      <div
        v-for="o in orders"
        :key="o.id"
        class="p-4 bg-white rounded-2xl shadow-sm border border-gray-100"
      >
        <div class="flex flex-wrap justify-between gap-3">
          <div>
            <p class="font-semibold text-gray-800">{{ o.tenant_name || '-' }}</p>
            <p class="text-sm text-gray-500">{{ planLabel(o.plan_slug) }} · Rp {{ (o.amount_rupiah || 0).toLocaleString('id-ID') }}</p>
            <p v-if="o.payment_note" class="text-xs text-gray-500 mt-1">Catatan: {{ o.payment_note }}</p>
            <p class="text-xs text-gray-400 mt-0.5">{{ formatDate(o.created_at) }}</p>
          </div>
          <div class="flex items-center gap-2">
            <span
            :class="{
              'bg-amber-100 text-amber-700': o.status === 'pending',
              'bg-blue-100 text-blue-700': o.status === 'paid',
              'bg-green-100 text-green-700': o.status === 'approved',
              'bg-red-100 text-red-700': o.status === 'rejected',
            }"
            class="px-2 py-1 rounded-lg text-xs font-medium"
          >
            {{ statusLabel(o.status) }}
            </span>
            <p v-if="o.payment_proof_url" class="text-xs mt-1">
              <a :href="o.payment_proof_url" target="_blank" rel="noopener" class="text-primary-600 hover:underline">Lihat bukti</a>
            </p>
            <template v-if="o.status === 'pending' || o.status === 'paid'">
              <button
                @click="approveOrder(o.id)"
                :disabled="actionLoading === o.id"
                class="px-3 py-1.5 rounded-lg bg-green-600 text-white text-sm font-medium hover:bg-green-700 disabled:opacity-50"
              >
                {{ actionLoading === o.id ? '...' : 'Setujui' }}
              </button>
              <button
                @click="rejectOrder(o.id)"
                :disabled="actionLoading === o.id"
                class="px-3 py-1.5 rounded-lg bg-red-600 text-white text-sm font-medium hover:bg-red-700 disabled:opacity-50"
              >
                Tolak
              </button>
            </template>
          </div>
        </div>
        <p v-if="o.status === 'rejected' && o.rejection_reason" class="text-xs text-red-600 mt-2">Alasan: {{ o.rejection_reason }}</p>
      </div>
      <p v-if="orders.length === 0" class="text-gray-500 text-center py-8">Tidak ada order.</p>
    </div>

    <div v-if="total > limit" class="mt-4 flex items-center gap-4 text-sm text-gray-600">
      <span>Total {{ total }} order</span>
      <div class="flex gap-2">
        <button
          :disabled="offset === 0"
          @click="goPage(-1)"
          class="px-3 py-2 rounded-xl bg-white border border-gray-200 disabled:opacity-50 hover:bg-gray-50"
        >
          Sebelumnya
        </button>
        <button
          :disabled="offset + limit >= total"
          @click="goPage(1)"
          class="px-3 py-2 rounded-xl bg-white border border-gray-200 disabled:opacity-50 hover:bg-gray-50"
        >
          Selanjutnya
        </button>
      </div>
    </div>

    <!-- Reject modal -->
    <Teleport to="body">
      <div
        v-if="rejectModal"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4"
        @click.self="rejectModal = null"
      >
        <div class="bg-white rounded-2xl p-5 w-full max-w-sm">
          <h3 class="font-semibold text-gray-800">Tolak Order</h3>
          <p class="text-sm text-gray-500 mt-1">Alasan penolakan (opsional)</p>
          <input
            v-model="rejectReason"
            type="text"
            placeholder="Contoh: Bukti transfer tidak valid"
            class="w-full mt-2 px-4 py-2 rounded-xl border border-gray-200"
            @keyup.enter="confirmReject"
          />
          <div class="mt-4 flex gap-3">
            <button @click="rejectModal = null" class="flex-1 py-2 rounded-xl border border-gray-200">Batal</button>
            <button
              @click="confirmReject"
              :disabled="actionLoading"
              class="flex-1 py-2 rounded-xl bg-red-600 text-white font-medium disabled:opacity-50"
            >
              Tolak
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { adminApi } from '../../lib/adminApi'
import { useToast } from '../../composables/useToast'

const { show: toast } = useToast()
const orders = ref([])
const total = ref(0)
const limit = ref(20)
const offset = ref(0)
const statusFilter = ref('pending')
const loading = ref(false)
const error = ref('')
const actionLoading = ref('')
const rejectModal = ref(null)
const rejectReason = ref('')

const planLabels = { free: 'Free', premium_1m: 'Premium 1 Bulan', premium_3m: 'Premium 3 Bulan', premium_6m: 'Premium 6 Bulan', premium_1y: 'Premium 1 Tahun', platinum: 'Platinum' }
function planLabel(slug) {
  return planLabels[slug] || slug || '-'
}
function statusLabel(s) {
  return { pending: 'Belum bayar', paid: 'Sudah bayar', approved: 'Disetujui', rejected: 'Ditolak' }[s] || s
}
function formatDate(s) {
  if (!s) return '-'
  try {
    return new Date(s).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' })
  } catch {
    return s
  }
}

async function fetchOrders() {
  loading.value = true
  error.value = ''
  try {
    const res = await adminApi.subscriptionOrders.list({
      limit: limit.value,
      offset: offset.value,
      status: statusFilter.value || undefined,
    })
    orders.value = res.orders || []
    total.value = res.total ?? 0
  } catch (e) {
    error.value = e.message || 'Gagal memuat order'
  } finally {
    loading.value = false
  }
}

watch(statusFilter, () => {
  offset.value = 0
  fetchOrders()
})

function goPage(delta) {
  offset.value = Math.max(0, offset.value + delta * limit.value)
  fetchOrders()
}

async function approveOrder(orderId) {
  actionLoading.value = orderId
  try {
    await adminApi.subscriptionOrders.approve(orderId)
    toast('Order disetujui')
    fetchOrders()
  } catch (e) {
    toast(e.message || 'Gagal menyetujui', 'error')
  } finally {
    actionLoading.value = ''
  }
}

function rejectOrder(orderId) {
  rejectModal.value = orderId
  rejectReason.value = ''
}

async function confirmReject() {
  if (!rejectModal.value) return
  actionLoading.value = rejectModal.value
  try {
    await adminApi.subscriptionOrders.reject(rejectModal.value, rejectReason.value)
    toast('Order ditolak')
    rejectModal.value = null
    fetchOrders()
  } catch (e) {
    toast(e.message || 'Gagal menolak', 'error')
  } finally {
    actionLoading.value = ''
  }
}

fetchOrders()
</script>
