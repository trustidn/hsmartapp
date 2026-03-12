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
    <aside class="cart-sticky" :class="{ 'cart-expanded': cartExpanded }">
      <!-- Collapse/Expand toggle -->
      <button
        type="button"
        class="w-full flex items-center justify-center gap-2 py-3 text-gray-600 hover:text-gray-800 hover:bg-gray-50 transition-colors touch-manipulation"
        @click="cartExpanded = !cartExpanded"
      >
        <svg
          class="w-5 h-5 transition-transform"
          :class="{ 'rotate-180': cartExpanded }"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
        </svg>
        <span class="text-sm font-semibold">{{ cartExpanded ? 'Sembunyikan' : 'Tampilkan daftar belanjaan' }}</span>
      </button>

      <!-- Cart list (expandable) - max tinggi tidak melebihi layar (penting untuk iPhone) -->
      <Transition name="cart-toggle">
        <div v-show="cartExpanded" class="cart-content">
          <div class="cart-list-scroll px-3 py-2 overflow-y-auto">
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
                class="flex items-center justify-between py-2 px-2 rounded-lg hover:bg-gray-50/80"
              >
                <span class="text-sm text-gray-800 truncate flex-1">{{ i.name }} <span class="text-gray-400 font-normal">×{{ i.qty }}</span></span>
                <div class="flex items-center gap-2 shrink-0">
                  <span class="text-sm font-medium text-gray-900 tabular-nums">Rp {{ formatNum(i.subtotal) }}</span>
                  <button
                    type="button"
                    class="p-2 rounded-lg text-red-600 bg-red-50 hover:bg-red-100 hover:text-red-700 border border-red-200/80 transition-colors touch-manipulation min-w-[2.5rem] min-h-[2.5rem] flex items-center justify-center"
                    @click.stop="pos.removeItem(i.productId)"
                    aria-label="Hapus"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>
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

    <!-- Payment form modal - mobile-friendly, touch targets 48px+ -->
    <div
      v-if="showPaymentForm"
      class="fixed inset-0 z-50 bg-black/50 flex items-end sm:items-center justify-center p-4"
      @click.self="cancelPaymentForm"
    >
      <div class="bg-white rounded-2xl w-full max-w-md p-6 sm:p-6 shadow-xl touch-manipulation" @click.stop>
        <h3 class="font-semibold text-xl mb-5">Pembayaran — {{ paymentLabel(selectedPaymentMethod) }}</h3>
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-600 mb-2">Nama (Opsional)</label>
            <input
              v-model="paymentForm.customerName"
              type="text"
              placeholder="Nama pelanggan"
              class="w-full px-4 py-3.5 min-h-[48px] rounded-xl border border-gray-300 text-base focus:ring-2 focus:ring-primary-400 focus:border-primary-500"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-600 mb-2">No HP (Opsional<span class="text-amber-600 font-normal">, dibutuhkan saat pengiriman via WA</span>)</label>
            <input
              v-model="paymentForm.customerPhone"
              type="tel"
              placeholder="08123456789"
              class="w-full px-4 py-3.5 min-h-[48px] rounded-xl border border-gray-300 text-base focus:ring-2 focus:ring-primary-400 focus:border-primary-500"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-600 mb-2">Jumlah bayar <span class="text-red-500">*</span></label>
            <input
              v-model.number="paymentForm.amountPaid"
              type="number"
              min="0"
              placeholder="0"
              class="w-full px-4 py-3.5 min-h-[48px] rounded-xl border border-gray-300 text-base tabular-nums focus:ring-2 focus:ring-primary-400 focus:border-primary-500"
            />
            <p class="text-sm text-gray-500 mt-2">Total: Rp {{ formatNum(pos.total) }}</p>
          </div>
        </div>
        <div class="flex gap-3 mt-6">
          <button
            type="button"
            class="flex-1 min-h-[52px] py-3.5 rounded-xl border border-gray-300 text-base font-semibold hover:bg-gray-50 active:bg-gray-100 touch-manipulation"
            @click="cancelPaymentForm"
          >
            Batal
          </button>
          <button
            type="button"
            class="flex-1 min-h-[52px] py-3.5 rounded-xl bg-primary-600 text-white text-base font-semibold hover:bg-primary-700 active:bg-primary-800 disabled:opacity-50 disabled:pointer-events-none touch-manipulation"
            :disabled="!isAmountValid || paying"
            @click="submitPayment"
          >
            {{ paying ? 'Memproses...' : 'Bayar' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Receipt modal - mobile-friendly -->
    <div
      v-if="paidReceipt"
      class="fixed inset-0 z-50 bg-black/50 flex items-end sm:items-center justify-center p-4"
      @click.self="closeReceipt"
    >
      <div class="bg-white rounded-2xl w-full max-w-md max-h-[90vh] overflow-auto p-5 shadow-xl touch-manipulation">
        <div class="flex justify-between items-center mb-4">
          <h3 class="font-semibold text-xl">Struk</h3>
          <button type="button" class="p-2.5 min-w-[44px] min-h-[44px] text-gray-500 rounded-xl hover:bg-gray-100 active:bg-gray-200" @click="closeReceipt" aria-label="Tutup">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
          </button>
        </div>
        <ReceiptContent
          :sale="paidReceipt"
          :settings="settings"
          :amount-paid="paidReceipt?.amount_paid"
          :change="paidReceipt?.change"
        />
        <div class="mt-6 flex flex-wrap gap-3">
          <button type="button" class="flex-1 min-w-[100px] min-h-[48px] py-3 rounded-xl border border-gray-300 text-base font-semibold hover:bg-gray-50 active:bg-gray-100" @click="doPrint">Cetak</button>
          <button type="button" class="flex-1 min-w-[100px] min-h-[48px] py-3 rounded-xl bg-primary-600 text-white text-base font-semibold hover:bg-primary-700 active:bg-primary-800" @click="doPdf">PDF</button>
          <button type="button" class="flex-1 min-w-[100px] min-h-[48px] py-3 rounded-xl bg-green-600 text-white text-base font-semibold hover:bg-green-700 active:bg-green-800" @click="doWhatsApp">WhatsApp</button>
          <button type="button" class="flex-1 min-w-[100px] min-h-[48px] py-3 rounded-xl bg-gray-100 text-gray-700 text-base font-semibold hover:bg-gray-200 active:bg-gray-300" @click="closeReceipt">Tutup</button>
        </div>
      </div>
    </div>

    <!-- Area cetak in-page (tanpa popup) - tersembunyi di layar, terlihat saat print -->
    <Teleport to="body">
      <div id="receipt-print-area" class="receipt-print-area" aria-hidden="true" />
    </Teleport>
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

function buildReceiptBodyHtml(sale, s) {
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
    '<div class="store">' + escapeHtml(name) + '</div>' +
    '<div class="date">' + escapeHtml(date) + '</div>' +
    '<div class="sep"></div>' +
    itemsHtml +
    '<div class="sep"></div>' +
    totHtml +
    footer
  )
}

/** Struk dokumen HTML penuh (untuk window baru - fallback mobile) */
function buildReceiptFullDocument(sale, s) {
  const inner = buildReceiptBodyHtml(sale, s)
  const style =
    'body{margin:0;padding:0;font-family:"Courier New",Courier,monospace;font-size:14px;line-height:1.5;color:#000;background:#fff}' +
    '.receipt{width:100%;padding:2mm 0;max-width:80mm;margin:0 auto;box-sizing:border-box}' +
    '.store{font-weight:bold;text-align:center;font-size:18px;margin-bottom:4px;line-height:1.3}' +
    '.date{text-align:center;color:#555;margin-bottom:8px;font-size:13px}' +
    '.sep{border-top:1px dashed #999;margin:10px 0;line-height:0}' +
    '.row{display:flex;justify-content:space-between;align-items:flex-start;gap:8px;margin:4px 0;min-height:1.4em;font-size:14px}' +
    '.row .item,.row span:first-child{flex:1;min-width:0}' +
    '.row .amt,.row span:last-child{flex-shrink:0;text-align:right;white-space:nowrap}' +
    '.row.total{font-weight:bold;margin-top:8px;font-size:16px}' +
    '.row.change{font-weight:bold;color:#15803d;font-size:14px}' +
    '.footer{text-align:center;font-size:12px;color:#666;margin-top:12px}' +
    '@page{size:auto;margin:5mm}' +
    '@media print{body{background:#fff;color:#000;-webkit-print-color-adjust:exact;print-color-adjust:exact}}'
  return '<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width"><title>Struk</title><style>' + style + '</style></head><body><div class="receipt">' + inner + '</div></body></html>'
}

/** Struk untuk cetak in-page (tanpa popup/iframe) */
function buildReceiptForInPagePrint(sale, s) {
  const inner = buildReceiptBodyHtml(sale, s)
  const style =
    '.receipt-print-area .receipt{width:100%;padding:2mm 0;font-family:"Courier New",Courier,monospace;font-size:14px;line-height:1.5;color:#000}' +
    '.receipt-print-area .store{font-weight:bold;text-align:center;font-size:18px;margin-bottom:4px;line-height:1.3}' +
    '.receipt-print-area .date{text-align:center;color:#555;margin-bottom:8px;font-size:13px}' +
    '.receipt-print-area .sep{border-top:1px dashed #999;margin:10px 0;line-height:0}' +
    '.receipt-print-area .row{display:flex;justify-content:space-between;align-items:flex-start;gap:8px;margin:4px 0;min-height:1.4em;font-size:14px}' +
    '.receipt-print-area .row .item,.receipt-print-area .row span:first-child{flex:1;min-width:0;overflow:hidden;text-overflow:ellipsis}' +
    '.receipt-print-area .row .amt,.receipt-print-area .row span:last-child{flex-shrink:0;text-align:right;white-space:nowrap}' +
    '.receipt-print-area .row.total{font-weight:bold;margin-top:8px;font-size:16px}' +
    '.receipt-print-area .row.change{font-weight:bold;color:#15803d;font-size:14px}' +
    '.receipt-print-area .footer{text-align:center;font-size:12px;color:#666;margin-top:12px;line-height:1.4;white-space:pre-wrap}'
  return '<style>' + style + '</style><div class="receipt">' + inner + '</div>'
}

function escapeHtml(str) {
  return String(str).replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;')
}

/** Deteksi perangkat mobile (bukan hanya viewport) */
function isMobileDevice() {
  return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent)
}

function doPrint() {
  if (!paidReceipt.value) return
  const s = settings.value
  const sale = paidReceipt.value

  // Di ponsel: fallback ke window baru dengan dokumen HTML penuh (lebih andal di Chrome/Safari mobile)
  if (isMobileDevice()) {
    const html = buildReceiptFullDocument(sale, s)
    const w = window.open('', '_blank', 'noopener,noreferrer,width=320,height=500')
    if (w) {
      w.document.write(html)
      w.document.close()
      w.focus()
      setTimeout(() => {
        try {
          w.print()
          w.onafterprint = () => w.close()
        } catch (e) {
          alert('Gagal membuka dialog cetak.')
          w.close()
        }
      }, 400)
    } else {
      alert('Izinkan pop-up untuk mencetak struk.')
    }
    return
  }

  // Desktop: cetak in-page
  const el = document.getElementById('receipt-print-area')
  if (!el) return
  el.innerHTML = buildReceiptForInPagePrint(sale, s)
  const cleanup = () => {
    el.innerHTML = ''
    window.onafterprint = null
  }
  window.onafterprint = cleanup
  const runPrint = () => {
    try {
      window.print()
    } catch (e) {
      alert('Gagal membuka dialog cetak.')
      cleanup()
    }
  }
  requestAnimationFrame(() => setTimeout(runPrint, 200))
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

/* Product card: accent lembut untuk membedakan dari background */
.product-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 5.25rem;
  padding: 0.875rem 0.625rem;
  background: linear-gradient(145deg, #f0fdf4 0%, #ecfdf5 50%, #d1fae5 100%);
  border: 1px solid rgb(187 247 208);
  border-radius: 1rem;
  box-shadow: 0 1px 2px rgba(34, 197, 94, 0.08);
  transition: all 0.2s ease;
  -webkit-tap-highlight-color: transparent;
  touch-action: manipulation;
}
.product-card:hover {
  border-color: rgb(134 239 172);
  background: linear-gradient(145deg, #ecfdf5 0%, #d1fae5 100%);
  box-shadow: 0 2px 8px rgba(34, 197, 94, 0.12);
  transform: translateY(-1px);
}
.product-card:active {
  transform: scale(0.97) translateY(0);
  background: linear-gradient(145deg, #dcfce7 0%, #bbf7d0 100%);
  border-color: rgba(22, 163, 74, 0.4);
  box-shadow: 0 0 0 2px rgba(22, 163, 74, 0.2);
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
  background: linear-gradient(90deg, rgb(240 253 244) 25%, rgb(209 250 229) 50%, rgb(240 253 244) 75%);
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

/* Cart: lebar sama dengan bottom nav, tinggi tidak melebihi layar */
.cart-sticky {
  position: fixed;
  left: 50%;
  transform: translateX(-50%);
  width: 100%;
  max-width: 360px;
  bottom: 8rem;
  max-height: min(45vh, calc(100vh - 14rem));
  border-radius: 1rem;
  z-index: 15;
  overflow: hidden;
  background: white;
  box-shadow: 0 -4px 20px rgba(0,0,0,0.08);
  border: 1px solid rgb(229 231 235);
  display: flex;
  flex-direction: column;
  transition: max-height 0.25s ease;
}
.cart-sticky.cart-expanded {
  max-height: min(70vh, calc(100vh - 12rem));
}

@media (min-width: 768px) {
  .cart-sticky {
    bottom: 9rem;
  }
  .cart-sticky.cart-expanded {
    max-height: min(75vh, calc(100vh - 12rem));
  }
}

.cart-content {
  border-top: 1px solid rgb(243 244 246);
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.cart-list-scroll {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
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

<!-- Area cetak in-page: proporsional di kertas ukuran apapun, perbaikan Android -->
<style>
/* Tersembunyi di layar tapi tetap dalam viewport (Android perlu element "on-page" agar tercetak) */
.receipt-print-area {
  position: fixed;
  left: 0;
  top: 0;
  width: 80mm;
  max-width: 80mm;
  visibility: hidden;
  pointer-events: none;
  background: white;
  z-index: -1;
}
/* Saat cetak: sembunyikan semua selain receipt (display lebih andal di Android daripada visibility) */
@media print {
  /* Sembunyikan app utama, hanya tampilkan area struk */
  body > *:not(#receipt-print-area) {
    display: none !important;
    visibility: hidden !important;
  }
  #receipt-print-area {
    display: block !important;
    visibility: visible !important;
    position: fixed !important;
    left: 0 !important;
    top: 0 !important;
    right: 0 !important;
    width: 100% !important;
    max-width: 100%;
    min-width: 80mm;
    box-sizing: border-box;
    padding: 5mm;
    margin: 0 !important;
    background: white !important;
    color: #000 !important;
    -webkit-print-color-adjust: exact;
    print-color-adjust: exact;
  }
  #receipt-print-area * {
    visibility: visible !important;
    color: inherit;
  }
  @page {
    size: auto;
    margin: 5mm;
  }
}
</style>
