<template>
  <div class="min-h-full">
    <h1 class="text-xl font-bold text-gray-800">Pengaturan SaaS</h1>
    <p class="text-gray-500 text-sm mt-1">Konfigurasi global: nama, logo, kontak, rekening pembayaran</p>

    <div v-if="loading" class="mt-6 space-y-4">
      <div v-for="i in 6" :key="i" class="h-14 bg-gray-200 rounded-xl animate-pulse" />
    </div>
    <div v-else-if="error" class="mt-6 p-4 bg-red-50 border border-red-200 rounded-xl text-red-700">{{ error }}</div>
    <form v-else @submit.prevent="save" class="mt-6 space-y-4">
      <section class="p-5 bg-white rounded-2xl shadow-sm border border-gray-100">
        <h2 class="text-sm font-medium text-gray-500 mb-4">Tampilan</h2>
        <div class="space-y-3">
          <div>
            <label class="block text-xs text-gray-600 mb-1">Nama Aplikasi</label>
            <input
              v-model="form.app_name"
              type="text"
              placeholder="HSmart"
              class="w-full px-4 py-2 rounded-xl border border-gray-200 focus:ring-2 focus:ring-primary-400"
            />
          </div>
          <div>
            <label class="block text-xs text-gray-600 mb-1">Logo</label>
            <div class="flex items-center gap-3">
              <div v-if="form.logo_url" class="shrink-0 w-16 h-16 rounded-xl border border-gray-200 overflow-hidden bg-gray-50">
                <img :src="logoDisplayUrl" alt="Logo" class="w-full h-full object-contain" />
              </div>
              <div class="flex-1 min-w-0">
                <input
                  ref="logoInput"
                  type="file"
                  accept=".png,.jpg,.jpeg,.svg,.webp"
                  class="hidden"
                  @change="onLogoFileSelect"
                />
                <button
                  type="button"
                  @click="logoInput?.click()"
                  class="px-4 py-2 rounded-xl border border-gray-200 text-sm font-medium text-gray-700 hover:bg-gray-50"
                >
                  {{ logoUploading ? 'Mengupload...' : 'Unggah Logo' }}
                </button>
                <p class="text-xs text-gray-500 mt-1">PNG, JPG, SVG, WebP (maks 2MB)</p>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section class="p-5 bg-white rounded-2xl shadow-sm border border-gray-100">
        <h2 class="text-sm font-medium text-gray-500 mb-4">Kontak Admin</h2>
        <div>
          <label class="block text-xs text-gray-600 mb-1">Email / Telepon / WhatsApp</label>
          <input
            v-model="form.admin_contact"
            type="text"
            placeholder="admin@example.com / 08123456789"
            class="w-full px-4 py-2 rounded-xl border border-gray-200 focus:ring-2 focus:ring-primary-400"
          />
        </div>
      </section>

      <section class="p-5 bg-white rounded-2xl shadow-sm border border-gray-100">
        <h2 class="text-sm font-medium text-gray-500 mb-4">Rekening Pembayaran</h2>
        <p class="text-xs text-gray-500 mb-3">Informasi ini ditampilkan ke tenant saat membuat order langganan</p>
        <div class="space-y-3">
          <div>
            <label class="block text-xs text-gray-600 mb-1">Nama Bank</label>
            <input
              v-model="form.bank_name"
              type="text"
              placeholder="BCA / Mandiri / ..."
              class="w-full px-4 py-2 rounded-xl border border-gray-200 focus:ring-2 focus:ring-primary-400"
            />
          </div>
          <div>
            <label class="block text-xs text-gray-600 mb-1">Nomor Rekening</label>
            <input
              v-model="form.bank_account_number"
              type="text"
              placeholder="1234567890"
              class="w-full px-4 py-2 rounded-xl border border-gray-200 focus:ring-2 focus:ring-primary-400"
            />
          </div>
          <div>
            <label class="block text-xs text-gray-600 mb-1">Nama Pemilik Rekening</label>
            <input
              v-model="form.bank_account_name"
              type="text"
              placeholder="PT. HSmart / Nama Pemilik"
              class="w-full px-4 py-2 rounded-xl border border-gray-200 focus:ring-2 focus:ring-primary-400"
            />
          </div>
        </div>
      </section>

      <button
        type="submit"
        :disabled="saving"
        class="w-full sm:w-auto px-6 py-2.5 rounded-xl bg-primary-600 text-white font-medium hover:bg-primary-700 disabled:opacity-50"
      >
        {{ saving ? 'Menyimpan...' : 'Simpan' }}
      </button>
    </form>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { adminApi } from '../../lib/adminApi'
import { useToast } from '../../composables/useToast'

const { show: toast } = useToast()
const loading = ref(true)
const error = ref('')
const saving = ref(false)
const logoUploading = ref(false)
const logoInput = ref(null)
const form = ref({
  app_name: '',
  logo_url: '',
  admin_contact: '',
  bank_name: '',
  bank_account_number: '',
  bank_account_name: '',
})

const logoDisplayUrl = computed(() => form.value.logo_url || '')

async function onLogoFileSelect(e) {
  const file = e.target.files?.[0]
  if (!file) return
  if (file.size > 2 * 1024 * 1024) {
    toast('File maksimal 2MB', 'error')
    return
  }
  logoUploading.value = true
  try {
    const url = await adminApi.upload.logo(file)
    form.value.logo_url = url
    toast('Logo berhasil diupload. Klik Simpan untuk menyimpan pengaturan.')
  } catch (err) {
    toast(err.message || 'Gagal upload logo', 'error')
  } finally {
    logoUploading.value = false
    e.target.value = ''
  }
}

onMounted(async () => {
  loading.value = true
  error.value = ''
  try {
    const res = await adminApi.saasSettings.get()
    form.value = {
      app_name: res.app_name ?? '',
      logo_url: res.logo_url ?? '',
      admin_contact: res.admin_contact ?? '',
      bank_name: res.bank_name ?? '',
      bank_account_number: res.bank_account_number ?? '',
      bank_account_name: res.bank_account_name ?? '',
    }
  } catch (e) {
    error.value = e.message || 'Gagal memuat pengaturan'
  } finally {
    loading.value = false
  }
})

async function save() {
  saving.value = true
  try {
    await adminApi.saasSettings.update(form.value)
    toast('Pengaturan berhasil disimpan')
  } catch (e) {
    toast(e.message || 'Gagal menyimpan', 'error')
  } finally {
    saving.value = false
  }
}
</script>
