<template>
  <div class="min-h-dvh flex flex-col bg-gradient-to-b from-primary-50/60 via-white to-gray-50/40">
    <main class="flex-1 flex flex-col items-center justify-center px-6 py-10">
      <!-- Logo & nama -->
      <div class="flex flex-col items-center mb-8">
        <img
          v-if="logoUrl"
          :src="logoUrl"
          :alt="appName"
          class="w-16 h-16 object-contain rounded-xl shadow-sm mb-3"
        />
        <span
          v-else
          class="w-16 h-16 rounded-xl bg-primary-100 flex items-center justify-center text-primary-600 font-bold text-xl mb-3 shadow-sm"
        >
          {{ (appName || 'H').charAt(0) }}
        </span>
        <h1 class="text-xl font-bold text-gray-900">{{ appName || 'HSmart' }}</h1>
        <p class="text-gray-500 text-sm mt-1">Login ke akun Anda</p>
      </div>

      <!-- Form card -->
      <div class="w-full max-w-sm bg-white rounded-2xl shadow-lg border border-gray-100 p-6">
        <form @submit.prevent="submit" class="space-y-4">
          <div>
            <input
              v-model="phone"
              type="tel"
              placeholder="Nomor HP (08xxxxxxxxxx)"
              class="w-full min-h-touch px-4 py-3 rounded-xl border border-gray-200 focus:ring-2 focus:ring-primary-500/30 focus:border-primary-500 text-base"
              required
            />
          </div>
          <div>
            <input
              v-model="password"
              type="password"
              placeholder="Password"
              class="w-full min-h-touch px-4 py-3 rounded-xl border border-gray-200 focus:ring-2 focus:ring-primary-500/30 focus:border-primary-500 text-base"
              required
            />
          </div>
          <p v-if="error" class="text-red-600 text-sm">{{ error }}</p>
          <button
            type="submit"
            class="w-full min-h-touch py-3.5 rounded-xl bg-primary-600 text-white font-semibold text-base hover:bg-primary-700 active:scale-[0.98] transition-all disabled:opacity-50"
            :disabled="loading"
          >
            {{ loading ? 'Memproses...' : 'Masuk' }}
          </button>
        </form>
        <p class="mt-5 text-center text-sm text-gray-500">
          Belum punya akun?
          <router-link to="/register" class="text-primary-600 font-semibold hover:underline">Daftar</router-link>
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
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useBranding } from '../composables/useBranding'

const router = useRouter()
const auth = useAuthStore()
const { appName, logoUrl } = useBranding()
const phone = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

async function submit() {
  error.value = ''
  loading.value = true
  try {
    await auth.login(phone.value, password.value)
    router.replace('/app')
  } catch (e) {
    error.value = e.message || 'Login gagal'
  } finally {
    loading.value = false
  }
}
</script>
