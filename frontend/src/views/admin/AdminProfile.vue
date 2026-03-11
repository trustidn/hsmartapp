<template>
  <div class="min-h-full">
    <h1 class="text-xl font-bold text-gray-800">Profil Admin</h1>
    <p class="text-gray-500 text-sm mt-1">Update nama, email, dan password</p>

    <div v-if="loading" class="mt-6 space-y-4">
      <div v-for="i in 4" :key="i" class="h-14 bg-gray-200 rounded-xl animate-pulse" />
    </div>
    <div v-else-if="loadError" class="mt-6 p-4 bg-red-50 border border-red-200 rounded-xl text-red-700">{{ loadError }}</div>
    <form v-else @submit.prevent="save" class="mt-6 space-y-6">
      <section class="p-5 bg-white rounded-2xl shadow-sm border border-gray-100">
        <h2 class="text-sm font-medium text-gray-500 mb-4">Informasi Profil</h2>
        <div class="space-y-3">
          <div>
            <label class="block text-xs text-gray-600 mb-1">Nama</label>
            <input
              v-model="form.name"
              type="text"
              placeholder="Nama Admin"
              class="w-full px-4 py-2 rounded-xl border border-gray-200 focus:ring-2 focus:ring-primary-400"
            />
          </div>
          <div>
            <label class="block text-xs text-gray-600 mb-1">Email</label>
            <input
              v-model="form.email"
              type="email"
              placeholder="admin@hsmart.app"
              class="w-full px-4 py-2 rounded-xl border border-gray-200 focus:ring-2 focus:ring-primary-400"
            />
          </div>
        </div>
      </section>

      <section class="p-5 bg-white rounded-2xl shadow-sm border border-gray-100">
        <h2 class="text-sm font-medium text-gray-500 mb-4">Ubah Password</h2>
        <p class="text-xs text-gray-500 mb-3">Kosongkan jika tidak ingin mengubah password</p>
        <div class="space-y-3">
          <div>
            <label class="block text-xs text-gray-600 mb-1">Password Saat Ini</label>
            <input
              v-model="form.current_password"
              type="password"
              placeholder="••••••••"
              class="w-full px-4 py-2 rounded-xl border border-gray-200 focus:ring-2 focus:ring-primary-400"
            />
          </div>
          <div>
            <label class="block text-xs text-gray-600 mb-1">Password Baru</label>
            <input
              v-model="form.new_password"
              type="password"
              placeholder="Minimal 6 karakter"
              class="w-full px-4 py-2 rounded-xl border border-gray-200 focus:ring-2 focus:ring-primary-400"
            />
          </div>
          <div>
            <label class="block text-xs text-gray-600 mb-1">Konfirmasi Password Baru</label>
            <input
              v-model="form.new_password_confirm"
              type="password"
              placeholder="Ulangi password baru"
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
import { ref, onMounted } from 'vue'
import { adminApi } from '../../lib/adminApi'
import { useAdminAuthStore } from '../../stores/adminAuth'
import { useToast } from '../../composables/useToast'

const { show: toast } = useToast()
const adminAuth = useAdminAuthStore()
const loading = ref(true)
const loadError = ref('')
const saving = ref(false)
const form = ref({
  name: '',
  email: '',
  current_password: '',
  new_password: '',
  new_password_confirm: '',
})

onMounted(async () => {
  loading.value = true
  loadError.value = ''
  try {
    const p = await adminApi.profile.get()
    form.value.name = p.name ?? ''
    form.value.email = p.email ?? ''
  } catch (e) {
    loadError.value = e.message || 'Gagal memuat profil'
  } finally {
    loading.value = false
  }
})

async function save() {
  const f = form.value
  if (!f.name?.trim() && !f.email?.trim() && !f.new_password?.trim()) {
    toast('Isi nama, email, atau password baru', 'error')
    return
  }
  if (f.new_password?.trim()) {
    if (!f.current_password?.trim()) {
      toast('Masukkan password saat ini untuk mengubah password', 'error')
      return
    }
    if (f.new_password.length < 6) {
      toast('Password baru minimal 6 karakter', 'error')
      return
    }
    if (f.new_password !== f.new_password_confirm) {
      toast('Konfirmasi password tidak cocok', 'error')
      return
    }
  }
  saving.value = true
  try {
    const payload = {}
    if (f.name?.trim()) payload.name = f.name.trim()
    if (f.email?.trim()) payload.email = f.email.trim()
    if (f.new_password?.trim()) {
      payload.current_password = f.current_password
      payload.new_password = f.new_password
    }
    if (Object.keys(payload).length === 0) {
      toast('Tidak ada perubahan', 'error')
      saving.value = false
      return
    }
    const res = await adminApi.profile.update(payload)
    adminAuth.updateProfile({ name: res.name })
    toast('Profil berhasil diperbarui')
    form.value.current_password = ''
    form.value.new_password = ''
    form.value.new_password_confirm = ''
  } catch (e) {
    toast(e.message || 'Gagal menyimpan', 'error')
  } finally {
    saving.value = false
  }
}
</script>
