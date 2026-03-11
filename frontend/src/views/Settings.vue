<template>
  <div class="min-h-full">
    <div v-if="!isOnline()" class="bg-amber-50 border border-amber-200 rounded-xl p-3 mb-4 text-amber-800 text-sm">
      Offline — pengaturan akan disinkronkan saat online.
    </div>

    <div v-if="loading" class="animate-pulse space-y-4">
      <div class="bg-white rounded-2xl p-5 h-80">
        <div class="h-4 w-24 bg-gray-200 rounded mb-4" />
        <div class="space-y-3">
          <div class="h-4 w-20 bg-gray-200 rounded mb-2" />
          <div class="h-12 bg-gray-200 rounded-xl" />
          <div class="h-4 w-16 bg-gray-200 rounded mb-2" />
          <div class="h-12 bg-gray-200 rounded-xl" />
          <div class="h-4 w-24 bg-gray-200 rounded mb-2" />
          <div class="h-12 bg-gray-200 rounded-xl" />
          <div class="h-12 bg-gray-200 rounded-xl mt-4" />
        </div>
      </div>
    </div>
    <section v-else class="bg-white rounded-2xl shadow-sm p-5 mb-4">
      <h2 class="text-sm font-medium text-gray-500 mb-4">Data Toko</h2>
      <form @submit.prevent="save" class="space-y-3">
        <div>
          <label class="block text-xs text-gray-500 mb-1">Nama Toko</label>
          <input
            v-model="form.name"
            type="text"
            placeholder="Contoh: Warung Gorengan"
            class="w-full min-h-touch px-4 rounded-xl border border-gray-300"
          />
        </div>
        <div>
          <label class="block text-xs text-gray-500 mb-1">No. HP</label>
          <input
            v-model="form.phone"
            type="tel"
            placeholder="08xxxxxxxxxx"
            class="w-full min-h-touch px-4 rounded-xl border border-gray-300 bg-gray-50"
            readonly
          />
          <p class="text-xs text-gray-400 mt-1">No. HP tidak dapat diubah</p>
        </div>
        <div>
          <label class="block text-xs text-gray-500 mb-1">Teks Kaki Struk</label>
          <input
            v-model="form.receipt_footer"
            type="text"
            placeholder="Terima kasih sudah berbelanja"
            class="w-full min-h-touch px-4 rounded-xl border border-gray-300"
          />
        </div>
        <div>
          <label class="block text-xs text-gray-500 mb-1">Metode Bayar Default</label>
          <select
            v-model="form.default_payment"
            class="w-full min-h-touch px-4 rounded-xl border border-gray-300"
          >
            <option value="cash">Tunai</option>
            <option value="qris">QRIS</option>
            <option value="transfer">Transfer</option>
            <option value="ewallet">E-Wallet</option>
          </select>
        </div>
        <div>
          <label class="block text-xs text-gray-500 mb-1">No. WhatsApp (untuk kirim struk)</label>
          <input
            v-model="form.whatsapp_number"
            type="tel"
            placeholder="08xxxxxxxxxx"
            class="w-full min-h-touch px-4 rounded-xl border border-gray-300"
          />
        </div>
        <button
          type="submit"
          class="w-full min-h-touch rounded-xl bg-primary-600 text-white font-semibold"
          :disabled="saving"
        >
          {{ saving ? 'Menyimpan...' : 'Simpan' }}
        </button>
      </form>
    </section>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useSettingsStore } from '../stores/settings'
import { useToastStore } from '../stores/toast'
import { isOnline } from '../lib/api'

const auth = useAuthStore()
const toast = useToastStore()
const settings = useSettingsStore()
const loading = ref(true)
const saving = ref(false)
const form = ref({
  name: '',
  phone: '',
  receipt_footer: '',
  default_payment: 'cash',
  whatsapp_number: '',
})

onMounted(async () => {
  loading.value = true
  try {
    if (auth.tenantId) {
      await settings.load(auth.tenantId)
      if (settings.settings) {
        form.value = {
          name: settings.settings.name ?? '',
          phone: settings.settings.phone ?? '',
          receipt_footer: settings.settings.receipt_footer ?? '',
          default_payment: settings.settings.default_payment ?? 'cash',
          whatsapp_number: settings.settings.whatsapp_number ?? '',
        }
      }
    }
  } finally {
    loading.value = false
  }
})

async function save() {
  saving.value = true
  try {
    await settings.update({
      name: form.value.name,
      receipt_footer: form.value.receipt_footer,
      default_payment: form.value.default_payment,
      whatsapp_number: form.value.whatsapp_number,
    })
    toast.success('Pengaturan berhasil disimpan')
  } catch (e) {
    toast.error(e.message || 'Gagal menyimpan pengaturan')
  } finally {
    saving.value = false
  }
}
</script>
