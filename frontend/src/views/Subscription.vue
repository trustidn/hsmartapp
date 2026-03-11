<template>
  <div class="min-h-full">
    <h1 class="text-xl font-bold text-gray-800">Langganan</h1>
    <p class="text-gray-500 text-sm mt-1">Upgrade plan atau perpanjang layanan. Masa aktif terakumulasi.</p>

    <!-- Current subscription -->
    <section class="mt-6 p-4 bg-white rounded-2xl shadow-sm border border-gray-100">
      <h2 class="text-sm font-medium text-gray-500 mb-2">Plan saat ini</h2>
      <p class="font-semibold text-gray-800">{{ planLabel(sub?.plan) }}</p>
      <p class="text-sm text-gray-500 mt-1">
        Masa aktif: {{ sub?.expired_at ? formatDate(sub.expired_at) : 'Tak terbatas' }}
      </p>
    </section>

    <!-- Available plans -->
    <section class="mt-6">
      <h2 class="text-sm font-medium text-gray-500 mb-3">Upgrade / Perpanjang</h2>
      <div class="space-y-3">
        <div
          v-for="p in availablePlans"
          :key="p.plan_slug"
          class="p-4 bg-white rounded-2xl shadow-sm border border-gray-100"
        >
          <div class="flex flex-wrap items-center justify-between gap-3">
            <div>
              <p class="font-semibold text-gray-800">{{ p.name }}</p>
              <p class="text-sm text-gray-500">
                {{ p.duration_days > 0 ? p.duration_days + ' hari' : 'Unlimited' }} ·
                Rp {{ (p.price_rupiah || 0).toLocaleString('id-ID') }}
              </p>
            </div>
            <button
              @click="openOrder(p)"
              class="px-4 py-2 rounded-xl bg-primary-600 text-white font-medium hover:bg-primary-700 text-sm"
            >
              Pilih
            </button>
          </div>
        </div>
      </div>
    </section>

    <!-- My orders -->
    <section class="mt-6">
      <h2 class="text-sm font-medium text-gray-500 mb-3">Order saya</h2>
      <div v-if="ordersLoading" class="text-sm text-gray-500">Memuat...</div>
      <div v-else-if="ordersError" class="p-4 bg-red-50 border border-red-200 rounded-xl text-red-700 text-sm">{{ ordersError }}</div>
      <div v-else-if="orders.length === 0" class="text-sm text-gray-500">Belum ada order.</div>
      <div v-else class="space-y-2">
        <div
          v-for="o in orders"
          :key="o.id"
          class="p-4 bg-white rounded-2xl shadow-sm border border-gray-100"
        >
          <div class="flex flex-wrap justify-between items-center gap-3">
            <div>
              <p class="font-medium text-gray-800">{{ planLabel(o.plan_slug) }}</p>
              <p class="text-sm text-gray-500">Rp {{ (o.amount_rupiah || 0).toLocaleString('id-ID') }} · {{ formatDate(o.created_at) }}</p>
            </div>
            <span
              :class="{
                'bg-amber-100 text-amber-700': o.status === 'pending',
                'bg-blue-100 text-blue-700': o.status === 'paid',
                'bg-green-100 text-green-700': o.status === 'approved',
                'bg-red-100 text-red-700': o.status === 'rejected',
              }"
              class="px-2 py-1 rounded-lg text-xs font-medium shrink-0"
            >
              {{ statusLabel(o.status) }}
            </span>
          </div>
          <!-- Upload bukti untuk order pending -->
          <div v-if="o.status === 'pending'" class="mt-3 pt-3 border-t border-gray-100">
            <p class="text-xs text-gray-500 mb-2">Transfer ke rekening di atas, lalu upload bukti pembayaran:</p>
            <div class="flex flex-wrap items-center gap-2">
              <input
                :ref="el => { if (el) proofFileRefs[o.id] = el }"
                type="file"
                accept=".png,.jpg,.jpeg,.webp,.pdf"
                class="hidden"
                @change="(e) => onProofFileSelect(e, o.id)"
              />
              <button
                type="button"
                @click="proofFileRefs[o.id]?.click?.()"
                :disabled="proofSubmitting === o.id"
                class="px-4 py-2 rounded-xl border border-gray-200 text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50"
              >
                {{ proofSubmitting === o.id ? 'Mengupload...' : 'Pilih File' }}
              </button>
              <span v-if="proofFileNames[o.id]" class="text-xs text-gray-600 truncate max-w-[140px]">{{ proofFileNames[o.id] }}</span>
            </div>
          </div>
          <img v-if="o.payment_proof_url" :src="o.payment_proof_url" alt="Bukti" class="mt-2 w-24 h-24 object-cover rounded-lg" />
        </div>
      </div>
    </section>

    <!-- Modal: Create order -->
    <Teleport to="body">
      <div
        v-if="orderModal"
        class="fixed inset-0 z-50 flex items-end sm:items-center justify-center bg-black/40 p-4 overflow-y-auto"
        @click.self="orderModal = null"
      >
        <div class="bg-white rounded-2xl p-5 w-full max-w-md shadow-xl my-4">
          <h3 class="font-semibold text-gray-800">Buat Order</h3>
          <p v-if="orderModal" class="text-sm text-gray-500 mt-1">
            {{ orderModal.name }} · Rp {{ (orderModal.price_rupiah || 0).toLocaleString('id-ID') }}
          </p>

          <!-- Rekening pembayaran -->
          <div v-if="paymentInfo.bank_name || paymentInfo.bank_account_number" class="mt-4 p-3 bg-gray-50 rounded-xl">
            <p class="text-xs font-medium text-gray-600 mb-2">Transfer ke:</p>
            <p class="text-sm font-semibold">{{ paymentInfo.bank_name }}</p>
            <p class="text-sm">No. Rek: {{ paymentInfo.bank_account_number }}</p>
            <p class="text-sm">a.n. {{ paymentInfo.bank_account_name }}</p>
          </div>

          <div class="mt-4">
            <label class="block text-xs text-gray-500 mb-1">Catatan pembayaran (opsional)</label>
            <textarea
              v-model="orderForm.payment_note"
              rows="2"
              placeholder="Contoh: Transfer BCA 123..."
              class="w-full px-4 py-2 rounded-xl border border-gray-200 text-sm"
            />
          </div>
          <div class="mt-4 flex gap-3">
            <button
              @click="orderModal = null"
              class="flex-1 py-2 rounded-xl border border-gray-200 text-gray-700 font-medium"
            >
              Batal
            </button>
            <button
              @click="submitOrder"
              :disabled="orderSubmitting"
              class="flex-1 py-2 rounded-xl bg-primary-600 text-white font-medium disabled:opacity-50"
            >
              {{ orderSubmitting ? '...' : 'Buat Order' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { api } from '../lib/api'
import { useToastStore } from '../stores/toast'

const toast = useToastStore()
const sub = ref(null)
const plans = ref([])
const orders = ref([])
const ordersLoading = ref(false)
const ordersError = ref('')
const orderModal = ref(null)
const orderForm = ref({ payment_note: '' })
const orderSubmitting = ref(false)
const proofFileRefs = ref({})
const proofFileNames = ref({})
const proofSubmitting = ref('')
const paymentInfo = ref({})

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
    return new Date(s).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' })
  } catch {
    return s
  }
}

const availablePlans = ref([])

onMounted(async () => {
  try {
    sub.value = await api.subscription.get()
  } catch {
    sub.value = null
  }
  try {
    const res = await api.saasSettings.get()
    paymentInfo.value = {
      bank_name: res.bank_name || '',
      bank_account_number: res.bank_account_number || '',
      bank_account_name: res.bank_account_name || '',
    }
  } catch {
    paymentInfo.value = {}
  }
  try {
    const res = await api.plans.list()
    const active = (res.plans || []).filter((p) => p.plan_slug !== 'free')
    plans.value = res.plans || []
    availablePlans.value = active
  } catch {
    availablePlans.value = [
      { plan_slug: 'premium_1m', name: 'Premium 1 Bulan', duration_days: 30, price_rupiah: 10000 },
      { plan_slug: 'premium_3m', name: 'Premium 3 Bulan', duration_days: 90, price_rupiah: 25000 },
      { plan_slug: 'premium_6m', name: 'Premium 6 Bulan', duration_days: 180, price_rupiah: 45000 },
      { plan_slug: 'premium_1y', name: 'Premium 1 Tahun', duration_days: 365, price_rupiah: 80000 },
      { plan_slug: 'platinum', name: 'Platinum', duration_days: 365, price_rupiah: 150000 },
    ]
  }
  fetchOrders()
})

function openOrder(p) {
  orderModal.value = p
  orderForm.value.payment_note = ''
}

async function submitOrder() {
  if (!orderModal.value) return
  orderSubmitting.value = true
  try {
    await api.subscriptionOrders.create({
      plan_slug: orderModal.value.plan_slug,
      payment_note: orderForm.value.payment_note || undefined,
    })
    toast.success('Order berhasil dibuat. Silakan transfer lalu upload bukti pembayaran.')
    orderModal.value = null
    fetchOrders()
  } catch (e) {
    toast.error(e.message || 'Gagal membuat order')
  } finally {
    orderSubmitting.value = false
  }
}

async function onProofFileSelect(e, orderId) {
  const file = e.target.files?.[0]
  if (!file) return
  if (file.size > 5 * 1024 * 1024) {
    toast.error('File maksimal 5MB')
    return
  }
  proofFileNames.value[orderId] = file.name
  proofSubmitting.value = orderId
  try {
    const url = await api.subscriptionOrders.uploadPaymentProof(file)
    await api.subscriptionOrders.setPaymentProof(orderId, url)
    toast.success('Bukti pembayaran berhasil diupload')
    proofFileNames.value[orderId] = ''
    fetchOrders()
  } catch (err) {
    toast.error(err.message || 'Gagal upload bukti')
  } finally {
    proofSubmitting.value = ''
  }
  e.target.value = ''
}

async function fetchOrders() {
  ordersLoading.value = true
  ordersError.value = ''
  try {
    const res = await api.subscriptionOrders.list()
    orders.value = res.orders || []
  } catch (e) {
    orders.value = []
    ordersError.value = e.message || 'Gagal memuat order'
  } finally {
    ordersLoading.value = false
  }
}
</script>
