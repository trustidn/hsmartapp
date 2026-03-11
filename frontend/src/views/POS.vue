<template>
  <div class="pos-container w-full">
    <!-- Offline badge -->
    <div v-if="!isOnline" class="py-1.5 mb-2 bg-amber-50 border border-amber-200 rounded-lg text-amber-800 text-xs font-medium text-center">
      Offline — data akan disinkronkan saat online
    </div>

    <!-- Search products -->
    <div class="mb-3">
      <div class="relative">
        <input
          v-model="productSearch"
          type="text"
          placeholder="Cari produk..."
          class="w-full pl-10 pr-4 py-2.5 rounded-xl border border-gray-200 bg-white text-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-primary-500/30 focus:border-primary-500"
        />
        <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
      </div>
    </div>

    <!-- Products: scroll only when needed -->
    <div class="products-scroll flex-1 overflow-y-auto pt-1 pb-64">
      <div v-if="products.loading" class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-3 sm:gap-3">
        <div v-for="i in 8" :key="i" class="product-card-skeleton rounded-2xl" />
      </div>
      <div v-else class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-3 sm:gap-3">
        <button
          v-for="p in filteredProducts"
          :key="p.id"
          type="button"
          class="product-card"
          @click="addProduct(p)"
        >
          <span class="product-name">{{ p.name }}</span>
          <span class="product-price">Rp {{ formatNum(p.price) }}</span>
        </button>
        <p v-if="!products.loading && filteredProducts.length === 0" class="col-span-full text-center text-gray-400 py-8 text-sm">
          {{ productSearch.trim() ? 'Tidak ada produk yang cocok.' : 'Belum ada produk.' }}
        </p>
      </div>
    </div>

    <!-- Cart + Payment: lebar sama dengan bottom nav, jarak dari bottom nav -->
    <aside class="cart-sticky">
      <!-- Collapse/Expand toggle -->
      <button
        type="button"
        class="w-full flex items-center justify-center gap-1.5 py-2 text-gray-500 hover:text-gray-700 transition-colors"
        @click="cartExpanded = !cartExpanded"
      >
        <svg
          class="w-4 h-4 transition-transform"
          :class="{ 'rotate-180': cartExpanded }"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
        </svg>
        <span class="text-xs font-medium">{{ cartExpanded ? 'Sembunyikan' : 'Tampilkan' }}</span>
      </button>

      <!-- Cart list (expandable) - minimalis modern -->
      <Transition name="cart-toggle">
        <div v-show="cartExpanded" class="cart-content">
          <div class="px-3 py-2 max-h-32 overflow-y-auto">
            <div v-if="pos.cart.length === 0" class="flex flex-col items-center justify-center py-6 text-gray-400">
              <svg class="w-10 h-10 mb-1 opacity-60" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 10.5V6a3.75 3.75 0 10-7.5 0v4.5m11.356-1.993l1.263 12c.07.665-.45 1.243-1.119 1.243H4.25a1.125 1.125 0 01-1.12-1.243l1.264-12A1.125 1.125 0 015.513 7.5h12.974c.576 0 1.059.435 1.119 1.007z" />
              </svg>
              <span class="text-xs">Keranjang kosong</span>
            </div>
            <ul v-else class="space-y-1">
              <li
                v-for="i in pos.cart"
                :key="i.productId"
                class="flex items-center justify-between py-1.5 px-2 rounded-lg hover:bg-gray-50/80 group"
              >
                <span class="text-sm text-gray-800 truncate flex-1">{{ i.name }} <span class="text-gray-400 font-normal">×{{ i.qty }}</span></span>
                <div class="flex items-center gap-1">
                  <span class="text-sm font-medium text-gray-900 tabular-nums">Rp {{ formatNum(i.subtotal) }}</span>
                  <button
                    type="button"
                    class="p-1 rounded-md text-gray-400 hover:text-red-500 hover:bg-red-50 opacity-0 group-hover:opacity-100 transition-opacity"
                    @click.stop="pos.removeItem(i.productId)"
                    aria-label="Hapus"
                  >
                    <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
                  </button>
                </div>
              </li>
            </ul>
          </div>
        </div>
      </Transition>

      <!-- Total + Payment -->
      <div class="px-4 pb-4 pt-3 border-t border-gray-100">
        <div class="flex items-center justify-between mb-3">
          <span class="text-sm text-gray-500">Total</span>
          <span class="text-xl font-bold text-gray-900">Rp {{ formatNum(pos.total) }}</span>
        </div>
        <div class="grid grid-cols-4 gap-2">
          <button
            v-for="pm in paymentMethods"
            :key="pm.value"
            type="button"
            class="flex flex-col items-center justify-center min-h-[48px] rounded-xl font-medium text-sm transition-all active:scale-95 disabled:opacity-50 disabled:pointer-events-none"
            :class="pm.bgClass"
            :disabled="pos.cart.length === 0 || paying"
            @click="openPaymentForm(pm.value)"
          >
            <component :is="pm.icon" class="w-5 h-5 mb-0.5" />
            <span>{{ pm.label }}</span>
          </button>
        </div>
      </div>
    </aside>

    <!-- Payment form modal -->
    <div
      v-if="showPaymentForm"
      class="fixed inset-0 z-50 bg-black/50 flex items-end sm:items-center justify-center p-4"
      @click.self="cancelPaymentForm"
    >
      <div class="bg-white rounded-2xl max-w-sm w-full p-5 shadow-xl" @click.stop>
        <h3 class="font-semibold text-lg mb-4">Pembayaran — {{ paymentLabel(selectedPaymentMethod) }}</h3>
        <div class="space-y-3">
          <div>
            <label class="block text-xs text-gray-500 mb-1">Nama (Opsional)</label>
            <input
              v-model="paymentForm.customerName"
              type="text"
              placeholder="Nama pelanggan"
              class="w-full px-3 py-2 rounded-lg border border-gray-300 text-sm"
            />
          </div>
          <div>
            <label class="block text-xs text-gray-500 mb-1">No HP (Opsional<span class="text-amber-600">, dibutuhkan saat pengiriman via WA</span>)</label>
            <input
              v-model="paymentForm.customerPhone"
              type="tel"
              placeholder="08123456789"
              class="w-full px-3 py-2 rounded-lg border border-gray-300 text-sm"
            />
          </div>
          <div>
            <label class="block text-xs text-gray-500 mb-1">Jumlah bayar <span class="text-red-500">*</span></label>
            <input
              v-model.number="paymentForm.amountPaid"
              type="number"
              min="0"
              placeholder="0"
              class="w-full px-3 py-2 rounded-lg border border-gray-300 text-sm tabular-nums"
            />
            <p class="text-xs text-gray-500 mt-1">Total: Rp {{ formatNum(pos.total) }}</p>
          </div>
        </div>
        <div class="flex gap-2 mt-5">
          <button
            type="button"
            class="flex-1 py-2.5 rounded-xl border border-gray-300 text-sm font-medium hover:bg-gray-50"
            @click="cancelPaymentForm"
          >
            Batal
          </button>
          <button
            type="button"
            class="flex-1 py-2.5 rounded-xl bg-primary-600 text-white text-sm font-medium hover:bg-primary-700 disabled:opacity-50 disabled:pointer-events-none"
            :disabled="!isAmountValid || paying"
            @click="submitPayment"
          >
            {{ paying ? 'Memproses...' : 'Bayar' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Receipt modal -->
    <div
      v-if="paidReceipt"
      class="fixed inset-0 z-50 bg-black/50 flex items-end sm:items-center justify-center p-4"
      @click.self="closeReceipt"
    >
      <div class="bg-white rounded-2xl max-w-sm w-full max-h-[85vh] overflow-auto p-5 shadow-xl">
        <div class="flex justify-between items-center mb-4">
          <h3 class="font-semibold text-lg">Struk</h3>
          <button type="button" class="p-2 text-gray-500 rounded-lg hover:bg-gray-100" @click="closeReceipt">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
          </button>
        </div>
        <ReceiptContent
          :sale="paidReceipt"
          :settings="settings"
          :amount-paid="paidReceipt?.amount_paid"
          :change="paidReceipt?.change"
        />
        <div class="mt-4 flex flex-wrap gap-2">
          <button type="button" class="flex-1 min-w-[80px] py-2.5 rounded-xl border border-gray-300 text-sm font-medium hover:bg-gray-50" @click="doPrint">Cetak</button>
          <button type="button" class="flex-1 min-w-[80px] py-2.5 rounded-xl bg-primary-600 text-white text-sm font-medium hover:bg-primary-700" @click="doPdf">PDF</button>
          <button type="button" class="flex-1 min-w-[80px] py-2.5 rounded-xl bg-green-600 text-white text-sm font-medium hover:bg-green-700" @click="doWhatsApp">WhatsApp</button>
          <button type="button" class="flex-1 min-w-[80px] py-2.5 rounded-xl bg-gray-100 text-gray-700 text-sm font-medium hover:bg-gray-200" @click="closeReceipt">Tutup</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useProductsStore } from '../stores/products'
import { usePosStore } from '../stores/pos'
import { useSettingsStore } from '../stores/settings'
import { isOnline as checkOnline } from '../lib/api'
import ReceiptContent from '../components/ReceiptContent.vue'
import { exportReceiptPdf, getReceiptWhatsAppText } from '../lib/receipt'

const auth = useAuthStore()
const products = useProductsStore()
const pos = usePosStore()
const settingsStore = useSettingsStore()
const settings = computed(() => settingsStore.settings)
const isOnline = ref(checkOnline())
const productSearch = ref('')
const paidReceipt = ref(null)
const customerWhatsapp = ref('') // No HP untuk WhatsApp
const cartExpanded = ref(false)

// Product search
const filteredProducts = computed(() => {
  const q = (productSearch.value || '').trim().toLowerCase()
  const items = products.items || []
  if (!q) return items
  return items.filter((p) => (p.name || '').toLowerCase().includes(q))
})

// Payment form
const showPaymentForm = ref(false)
const selectedPaymentMethod = ref('cash')
const paymentForm = ref({
  customerName: '',
  customerPhone: '',
  amountPaid: 0,
})
const isAmountValid = computed(() => {
  const amt = Number(paymentForm.value.amountPaid)
  return !isNaN(amt) && amt >= pos.total && pos.total > 0
})

// Payment icons
const CashIcon = { template: `<svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M2.25 18.75a60.07 60.07 0 0115.797 2.101c.727.198 1.453-.342 1.453-1.096V18.75M3.75 4.5v.75A.75.75 0 013 6h-.75m0 0v-.375c0-.621.504-1.125 1.125-1.125H20.25M2.25 6v9m18-10.5v.75c0 .414.336.75.75.75h.75m-1.5-1.5h.375c.621 0 1.125.504 1.125 1.125v9.75c0 .621-.504 1.125-1.125 1.125h-.375" /></svg>` }
const QrisIcon = { template: `<svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M3.75 4.5v15m9-11.25v11.25m6-13.5v13.5M3.75 9.75h16.5m-16.5 6.75h16.5" /></svg>` }
const TransferIcon = { template: `<svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M7.5 21L3 16.5m0 0L7.5 12M3 16.5h13.5m0-13.5L21 7.5m0 0L16.5 12M21 7.5H7.5" /></svg>` }
const EwalletIcon = { template: `<svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M21 12a2.25 2.25 0 00-2.25-2.25H15a3 3 0 11-6 0H5.25A2.25 2.25 0 003 12m18 0v6a2.25 2.25 0 01-2.25 2.25H5.25A2.25 2.25 0 013 18v-6m18 0V9M3 12V9" /></svg>` }

const paymentMethods = [
  { value: 'cash', label: 'Tunai', icon: CashIcon, bgClass: 'bg-emerald-100 text-emerald-800 hover:bg-emerald-200' },
  { value: 'qris', label: 'QRIS', icon: QrisIcon, bgClass: 'bg-violet-100 text-violet-800 hover:bg-violet-200' },
  { value: 'transfer', label: 'Transfer', icon: TransferIcon, bgClass: 'bg-slate-100 text-slate-700 hover:bg-slate-200' },
  { value: 'ewallet', label: 'E-Wallet', icon: EwalletIcon, bgClass: 'bg-amber-100 text-amber-800 hover:bg-amber-200' },
]

onMounted(async () => {
  if (auth.tenantId) {
    products.load(auth.tenantId)
    await settingsStore.load(auth.tenantId)
    const defaultPay = settingsStore.defaultPayment
    if (defaultPay && paymentMethods.some(pm => pm.value === defaultPay)) {
      pos.paymentMethod = defaultPay
    }
  }
  const update = () => { isOnline.value = checkOnline() }
  window.addEventListener('online', update)
  window.addEventListener('offline', update)
})

function playBeep() {
  try {
    const ctx = new (window.AudioContext || window.webkitAudioContext)()
    const osc = ctx.createOscillator()
    const gain = ctx.createGain()
    osc.connect(gain)
    gain.connect(ctx.destination)
    osc.frequency.value = 880
    osc.type = 'sine'
    gain.gain.setValueAtTime(0.12, ctx.currentTime)
    gain.gain.exponentialRampToValueAtTime(0.01, ctx.currentTime + 0.08)
    osc.start(ctx.currentTime)
    osc.stop(ctx.currentTime + 0.08)
  } catch (_) {}
}

function addProduct(p) {
  playBeep()
  pos.addItem(p, 1)
}

const paying = ref(false)

function openPaymentForm(method) {
  if (pos.cart.length === 0) return
  pos.paymentMethod = method
  selectedPaymentMethod.value = method
  paymentForm.value = {
    customerName: '',
    customerPhone: '',
    amountPaid: pos.total,
  }
  showPaymentForm.value = true
}

function cancelPaymentForm() {
  showPaymentForm.value = false
}

async function submitPayment() {
  if (!isAmountValid.value) return
  const totalBeforePay = pos.total
  const amt = Number(paymentForm.value.amountPaid)
  paying.value = true
  try {
    const result = await pos.pay(auth.tenantId)
    if (result?.ok && result?.sale) {
      const change = Math.max(0, amt - totalBeforePay)
      paidReceipt.value = {
        ...result.sale,
        amount_paid: amt,
        change,
        customer_name: paymentForm.value.customerName,
      }
      customerWhatsapp.value = paymentForm.value.customerPhone || ''
      showPaymentForm.value = false
    }
  } catch {
    // error
  } finally {
    paying.value = false
  }
}

function closeReceipt() {
  paidReceipt.value = null
}

function paymentLabel(m) {
  const map = { cash: 'Tunai', qris: 'QRIS', transfer: 'Transfer', ewallet: 'E-Wallet' }
  return map[m] || m || 'Tunai'
}

function buildReceiptHtml(sale, s) {
  const name = (s?.name || 'HSmart POS')
  const date = sale.created_at ? new Date(sale.created_at).toLocaleString('id-ID') : ''
  const trunc = (str, len) => String(str).slice(0, len)
  const itemsHtml = (sale.items || []).map((i) => {
    const nm = trunc(i.product_name || i.name || 'Item', 22)
    const qty = i.qty || 1
    const sub = i.subtotal ?? (i.price || 0) * qty
    const left = nm + ' x' + qty
    const right = 'Rp ' + formatNum(sub)
    return '<div class="row"><span class="item">' + escapeHtml(left) + '</span><span class="amt">' + escapeHtml(right) + '</span></div>'
  }).join('')
  let totHtml = '<div class="row total"><span>Total</span><span>Rp ' + escapeHtml(formatNum(sale.total)) + '</span></div>'
  if (sale.amount_paid != null && sale.amount_paid > 0) {
    totHtml += '<div class="row"><span>Bayar</span><span>Rp ' + escapeHtml(formatNum(sale.amount_paid)) + '</span></div>'
    if (sale.change != null && sale.change > 0) {
      totHtml += '<div class="row change"><span>Kembalian</span><span>Rp ' + escapeHtml(formatNum(sale.change)) + '</span></div>'
    }
  }
  totHtml += '<div class="row"><span>Metode</span><span>' + escapeHtml(paymentLabel(sale.payment_method)) + '</span></div>'
  const footer = s?.receipt_footer ? '<div class="footer">' + escapeHtml(s.receipt_footer) + '</div>' : ''
  return (
    '<!DOCTYPE html><html><head><meta charset="UTF-8"><title>Struk</title>' +
    '<style>' +
    '@page{size:80mm auto;margin:5mm}' +
    'body{margin:0;padding:10mm;font-family:"Courier New",Courier,monospace;font-size:12px;line-height:1.5;color:#000;width:72mm;box-sizing:border-box}' +
    '*{box-sizing:border-box}' +
    '.receipt{width:100%}' +
    '.store{font-weight:bold;text-align:center;font-size:14px;margin-bottom:2px}' +
    '.date{text-align:center;color:#555;margin-bottom:6px;font-size:11px}' +
    '.sep{border-top:1px dashed #999;margin:8px 0;line-height:0}' +
    '.row{display:flex;justify-content:space-between;align-items:flex-start;gap:6px;margin:3px 0;min-height:1.2em}' +
    '.row .item,.row span:first-child{flex:1;min-width:0;overflow:hidden;text-overflow:ellipsis}' +
    '.row .amt,.row span:last-child{flex-shrink:0;text-align:right;white-space:nowrap}' +
    '.row.total{font-weight:bold;margin-top:6px;font-size:13px}' +
    '.row.change{font-weight:bold;color:#15803d;font-size:12px}' +
    '.footer{text-align:center;font-size:10px;color:#666;margin-top:10px;line-height:1.4;white-space:pre-wrap}' +
    '@media print{body{margin:0;padding:8mm;width:72mm;-webkit-print-color-adjust:exact;print-color-adjust:exact}@page{margin:6mm}}' +
    '</style></head><body>' +
    '<div class="receipt">' +
    '<div class="store">' + escapeHtml(name) + '</div>' +
    '<div class="date">' + escapeHtml(date) + '</div>' +
    '<div class="sep"></div>' +
    itemsHtml +
    '<div class="sep"></div>' +
    totHtml +
    footer +
    '</div>' +
    '<script>window.onload=function(){window.print();setTimeout(window.close,100);}<\/script></body></html>'
  )
}
function escapeHtml(str) {
  return String(str).replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;')
}

function doPrint() {
  if (!paidReceipt.value) return
  const s = settings.value
  const html = buildReceiptHtml(paidReceipt.value, s)

  const blob = new Blob([html], { type: 'text/html;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const printWin = window.open(url, '_blank', 'width=320,height=480')
  if (printWin) {
    setTimeout(() => URL.revokeObjectURL(url), 5000)
  } else {
    URL.revokeObjectURL(url)
    // Fallback: iframe cetak
    const iframe = document.createElement('iframe')
    iframe.style.cssText = 'position:fixed;width:80mm;min-height:200px;left:-9999px;top:0;border:none;'
    document.body.appendChild(iframe)
    const doc = iframe.contentDocument
    if (doc) {
      doc.open()
      doc.write(html)
      doc.close()
      setTimeout(() => {
        try { iframe.contentWindow?.print() } catch { alert('Gagal membuka dialog cetak') }
        setTimeout(() => iframe.parentNode?.removeChild(iframe), 1000)
      }, 300)
    } else {
      alert('Izinkan pop-up untuk mencetak struk')
    }
  }
}

function doPdf() {
  if (!paidReceipt.value) return
  try {
    exportReceiptPdf(paidReceipt.value, settings.value)
  } catch (e) {
    console.error('PDF error:', e)
    alert('Gagal mengunduh PDF')
  }
}

function doWhatsApp() {
  if (!paidReceipt.value) return
  const text = getReceiptWhatsAppText(paidReceipt.value, settings.value)
  const num = (customerWhatsapp.value || settings.value?.whatsapp_number || '').replace(/\D/g, '')
  const url = num
    ? `https://wa.me/${num}?text=${encodeURIComponent(text)}`
    : `https://api.whatsapp.com/send?text=${encodeURIComponent(text)}`
  window.open(url, '_blank', 'noopener')
}

function formatNum(n) {
  return Number(n).toLocaleString('id-ID')
}
</script>

<style scoped>
.pos-container {
  display: flex;
  flex-direction: column;
  min-height: calc(100vh - 3.5rem);
}

.products-scroll {
  min-height: 0;
}

/* Product card: modern, touch-friendly */
.product-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 5.25rem;
  padding: 0.875rem 0.625rem;
  background: linear-gradient(145deg, #ffffff 0%, #f8fafc 100%);
  border: 1px solid rgb(226 232 240);
  border-radius: 1rem;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.04);
  transition: all 0.2s ease;
  -webkit-tap-highlight-color: transparent;
  touch-action: manipulation;
}
.product-card:hover {
  border-color: rgb(203 213 225);
  background: linear-gradient(145deg, #f8fafc 0%, #f1f5f9 100%);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  transform: translateY(-1px);
}
.product-card:active {
  transform: scale(0.97) translateY(0);
  background: linear-gradient(145deg, #ecfdf5 0%, #d1fae5 100%);
  border-color: rgba(22, 163, 74, 0.35);
  box-shadow: 0 0 0 2px rgba(22, 163, 74, 0.15);
}
.product-name {
  font-size: 0.8125rem;
  font-weight: 600;
  color: rgb(30 41 59);
  text-align: center;
  line-height: 1.3;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  width: 100%;
}
.product-price {
  font-size: 0.75rem;
  font-weight: 600;
  color: #16a34a;
  margin-top: 0.375rem;
  letter-spacing: 0.02em;
}
.product-card-skeleton {
  min-height: 5.25rem;
  background: linear-gradient(90deg, rgb(241 245 249) 25%, rgb(226 232 240) 50%, rgb(241 245 249) 75%);
  background-size: 200% 100%;
  animation: product-shimmer 1.2s ease-in-out infinite;
}
@keyframes product-shimmer {
  0% { background-position: 200% 0 }
  100% { background-position: -200% 0 }
}
@media (min-width: 640px) {
  .product-card { min-height: 5.5rem; padding: 1rem 0.75rem; }
  .product-name { font-size: 0.875rem; }
  .product-price { font-size: 0.8125rem; }
}
@media (min-width: 1024px) {
  .product-card { min-height: 6rem; padding: 1.125rem 0.875rem; }
  .product-name { font-size: 0.9375rem; }
  .product-price { font-size: 0.875rem; }
}

/* Cart: lebar sama dengan bottom nav (360px), jarak cukup dari bottom nav */
.cart-sticky {
  position: fixed;
  left: 50%;
  transform: translateX(-50%);
  width: 100%;
  max-width: 360px;
  bottom: 6.5rem;
  max-height: 50vh;
  border-radius: 1rem;
  z-index: 15;
  overflow: hidden;
  background: white;
  box-shadow: 0 -4px 20px rgba(0,0,0,0.08);
  border: 1px solid rgb(229 231 235);
}

@media (min-width: 768px) {
  .cart-sticky {
    bottom: 7rem;
  }
}

.cart-content {
  border-top: 1px solid rgb(243 244 246);
}

.cart-toggle-enter-active,
.cart-toggle-leave-active {
  transition: all 0.2s ease;
}
.cart-toggle-enter-from,
.cart-toggle-leave-to {
  max-height: 0;
  opacity: 0;
  overflow: hidden;
}
</style>
