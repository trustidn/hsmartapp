<template>
  <div class="min-h-full flex flex-col justify-center p-6 bg-white">
    <div class="w-full max-w-xs mx-auto">
      <h1 class="text-2xl font-bold text-center text-primary-600 mb-8">Daftar HSmart</h1>
      <p class="text-gray-500 text-center text-sm mb-6">Nomor HP & password saja</p>
      <form @submit.prevent="submit" class="space-y-4">
        <input
          v-model="phone"
          type="tel"
          placeholder="08xxxxxxxxxx"
          class="w-full min-h-touch px-4 rounded-xl border border-gray-300 text-lg"
          required
        />
        <input
          v-model="name"
          type="text"
          placeholder="Nama (opsional)"
          class="w-full min-h-touch px-4 rounded-xl border border-gray-300 text-lg"
        />
        <input
          v-model="password"
          type="password"
          placeholder="Password"
          class="w-full min-h-touch px-4 rounded-xl border border-gray-300 text-lg"
          required
          minlength="6"
        />
        <p v-if="error" class="text-red-600 text-sm">{{ error }}</p>
        <p v-if="success" class="text-green-600 text-sm">{{ success }}</p>
        <button
          type="submit"
          class="w-full min-h-touch rounded-xl bg-primary-600 text-white font-semibold text-lg"
          :disabled="loading"
        >
          {{ loading ? '...' : 'Daftar' }}
        </button>
      </form>
      <p class="mt-6 text-center text-sm text-gray-500">
        Sudah punya akun?
        <router-link to="/login" class="text-primary-600 font-medium">Masuk</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
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
