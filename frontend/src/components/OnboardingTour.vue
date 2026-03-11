<template>
  <Teleport to="body">
    <div
      v-if="visible"
      class="fixed inset-0 z-[100] flex items-center justify-center bg-black/50 p-4"
      @click.self="skip"
    >
      <div class="bg-white rounded-2xl shadow-xl max-w-sm w-full p-6" @click.stop>
        <h3 class="font-semibold text-gray-900 text-lg">{{ steps[current].title }}</h3>
        <p class="text-gray-600 text-sm mt-2">{{ steps[current].desc }}</p>
        <div class="flex gap-3 mt-6">
          <button
            v-if="current > 0"
            type="button"
            class="flex-1 py-2.5 rounded-xl border border-gray-200 text-gray-700 font-medium text-sm"
            @click="prev"
          >
            Kembali
          </button>
          <button
            type="button"
            class="flex-1 py-2.5 rounded-xl bg-primary-600 text-white font-medium text-sm hover:bg-primary-700"
            @click="current < steps.length - 1 ? next() : finish()"
          >
            {{ current < steps.length - 1 ? 'Lanjut' : 'Selesai' }}
          </button>
          <button
            v-if="current === 0"
            type="button"
            class="py-2.5 px-3 rounded-xl text-gray-500 text-sm hover:bg-gray-100"
            @click="skip"
          >
            Lewati
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps({
  tenantId: { type: String, default: '' },
})

const router = useRouter()
const current = ref(0)
const visible = ref(false)

const steps = [
  { title: 'Selamat datang!', desc: 'HSmart membantu Anda mengelola penjualan, produk, dan laporan keuangan toko dengan mudah.' },
  { title: 'Dashboard', desc: 'Lihat ringkasan penjualan hari ini, grafik, dan produk terlaris di halaman utama.' },
  { title: 'POS', desc: 'Lakukan transaksi penjualan dengan cepat. Pilih produk, jumlah, lalu selesaikan pembayaran.' },
  { title: 'Produk', desc: 'Kelola daftar produk Anda. Tambah, edit, atau nonaktifkan produk kapan saja.' },
  { title: 'Laporan', desc: 'Akses laporan penjualan, pengeluaran, dan profit. Ekspor ke Excel atau PDF.' },
  { title: 'Langganan', desc: 'Upgrade plan untuk fitur lebih: jumlah produk tak terbatas dan laporan hingga 12 bulan.' },
]

const storageKey = computed(() => `onboarding_done_${props.tenantId}`)

onMounted(() => {
  if (!props.tenantId) return
  try {
    if (localStorage.getItem(storageKey.value)) return
    visible.value = true
  } catch {
    visible.value = false
  }
})

function next() {
  if (current.value < steps.length - 1) current.value++
}

function prev() {
  if (current.value > 0) current.value--
}

function finish() {
  try {
    localStorage.setItem(storageKey.value, '1')
  } catch {}
  visible.value = false
}

function skip() {
  finish()
}

defineExpose({ show: () => (visible.value = true), hide: () => (visible.value = false) })
</script>
