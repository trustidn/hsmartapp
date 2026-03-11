<template>
  <div class="min-h-dvh flex flex-col bg-gradient-to-b from-primary-50/60 via-white to-gray-50/40">
    <main class="flex-1 flex flex-col items-center justify-center px-6 py-10">
      <!-- Logo & nama -->
      <div class="flex flex-col items-center mb-6">
        <img
          v-if="logoUrl"
          :src="logoUrl"
          :alt="appName"
          class="w-14 h-14 object-contain rounded-xl shadow-sm mb-2"
        />
        <span
          v-else
          class="w-14 h-14 rounded-xl bg-primary-100 flex items-center justify-center text-primary-600 font-bold text-lg mb-2 shadow-sm"
        >
          {{ (appName || 'H').charAt(0) }}
        </span>
        <h1 class="text-lg font-bold text-gray-900">{{ appName || 'HSmart' }}</h1>
        <p class="text-gray-500 text-sm">Daftar akun baru</p>
      </div>

      <!-- Form card -->
      <div class="w-full max-w-sm bg-white rounded-2xl shadow-lg border border-gray-100 p-6">
        <form @submit.prevent="submit" class="space-y-4">
          <input
            v-model="phone"
            type="tel"
            placeholder="Nomor HP (08xxxxxxxxxx)"
            class="w-full min-h-touch px-4 py-3 rounded-xl border border-gray-200 focus:ring-2 focus:ring-primary-500/30 focus:border-primary-500 text-base"
            required
          />
          <input
            v-model="name"
            type="text"
            placeholder="Nama toko (opsional)"
            class="w-full min-h-touch px-4 py-3 rounded-xl border border-gray-200 focus:ring-2 focus:ring-primary-500/30 focus:border-primary-500 text-base"
          />
          <input
            v-model="password"
            type="password"
            placeholder="Password (min 6 karakter)"
            class="w-full min-h-touch px-4 py-3 rounded-xl border border-gray-200 focus:ring-2 focus:ring-primary-500/30 focus:border-primary-500 text-base"
            required
            minlength="6"
          />
          <p v-if="error" class="text-red-600 text-sm">{{ error }}</p>
          <p v-if="success" class="text-green-600 text-sm">{{ success }}</p>
          <button
            type="submit"
            class="w-full min-h-touch py-3.5 rounded-xl bg-primary-600 text-white font-semibold text-base hover:bg-primary-700 active:scale-[0.98] transition-all disabled:opacity-50"
            :disabled="loading"
          >
            {{ loading ? 'Memproses...' : 'Daftar' }}
          </button>
        </form>
        <p class="mt-5 text-center text-sm text-gray-500">
          Sudah punya akun?
          <router-link to="/login" class="text-primary-600 font-semibold hover:underline">Masuk</router-link>
        </p>
      </div>

      <router-link to="/" class="mt-6 text-sm text-gray-500 hover:text-gray-700">
        ← Kembali
      </router-link>
    </main>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useBranding } from '../composables/useBranding'

const auth = useAuthStore()
const { appName, logoUrl } = useBranding()
const phone = ref('')
const name = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')
const success = ref('')

async function submit() {
  error.value = ''
  success.value = ''
  loading.value = true
  try {
    await auth.register(phone.value, password.value, name.value)
    success.value = 'Berhasil daftar. Silakan login.'
  } catch (e) {
    error.value = e.message || 'Daftar gagal'
  } finally {
    loading.value = false
  }
}
</script>
