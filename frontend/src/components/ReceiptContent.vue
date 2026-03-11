<template>
  <div class="font-mono text-sm space-y-1" v-if="sale">
    <p class="font-sans font-semibold text-center text-lg">{{ (settings?.name || 'HSmart POS') }}</p>
    <p class="text-gray-500 text-center">{{ formatDate(sale.created_at) }}</p>
    <p v-if="sale.customer_name" class="text-center">Pelanggan: {{ sale.customer_name }}</p>
    <div class="border-t border-b border-gray-200 my-3 py-2">
      <div v-for="i in (sale.items || [])" :key="i.id || i.product_id" class="flex justify-between">
        <span>{{ i.product_name || i.name || 'Item' }} × {{ i.qty }}</span>
        <span>Rp {{ formatNum(i.subtotal || i.price * i.qty) }}</span>
      </div>
    </div>
    <div class="flex justify-between font-semibold">
      <span>Total</span>
      <span>Rp {{ formatNum(sale.total) }}</span>
    </div>
    <template v-if="amountPaid != null && amountPaid > 0">
      <div class="flex justify-between">
        <span>Jumlah bayar</span>
        <span>Rp {{ formatNum(amountPaid) }}</span>
      </div>
      <div v-if="change != null && change > 0" class="flex justify-between font-semibold text-primary-600">
        <span>Kembalian</span>
        <span>Rp {{ formatNum(change) }}</span>
      </div>
    </template>
    <p class="text-gray-500">Bayar: {{ paymentLabel(sale.payment_method) }}</p>
    <p v-if="settings?.receipt_footer" class="text-center text-xs text-gray-500 mt-2">{{ settings.receipt_footer }}</p>
  </div>
</template>

<script setup>
defineProps({
  sale: { type: Object, default: null },
  settings: { type: Object, default: null },
  amountPaid: { type: Number, default: null },
  change: { type: Number, default: null },
})

function formatNum(n) {
  return Number(n || 0).toLocaleString('id-ID')
}
function formatDate(s) {
  if (!s) return ''
  return new Date(s).toLocaleString('id-ID')
}
function paymentLabel(m) {
  const map = { cash: 'Tunai', qris: 'QRIS', transfer: 'Transfer', ewallet: 'E-Wallet' }
  return map[m] || m || 'Tunai'
}
</script>
