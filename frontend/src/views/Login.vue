<template>
  <div class="min-h-full flex flex-col justify-center p-6 bg-white">
    <div class="w-full max-w-xs mx-auto">
      <h1 class="text-2xl font-bold text-center text-primary-600 mb-8">HSmart</h1>
      <p class="text-gray-500 text-center text-sm mb-6">Login dengan nomor HP & password</p>
      <form @submit.prevent="submit" class="space-y-4">
        <input
          v-model="phone"
          type="tel"
          placeholder="08xxxxxxxxxx"
          class="w-full min-h-touch px-4 rounded-xl border border-gray-300 text-lg"
          required
        />
        <input
          v-model="password"
          type="password"
          placeholder="Password"
          class="w-full min-h-touch px-4 rounded-xl border border-gray-300 text-lg"
          required
        />
        <p v-if="error" class="text-red-600 text-sm">{{ error }}</p>
        <button
          type="submit"
          class="w-full min-h-touch rounded-xl bg-primary-600 text-white font-semibold text-lg"
          :disabled="loading"
        >
          {{ loading ? '...' : 'Masuk' }}
        </button>
      </form>
      <p class="mt-6 text-center text-sm text-gray-500">
        Belum punya akun?
        <router-link to="/register" class="text-primary-600 font-medium">Daftar</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const auth = useAuthStore()
const phone = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

async function submit() {
  error.value = ''
  loading.value = true
  try {
    await auth.login(phone.value, password.value)
    router.replace('/')
  } catch (e) {
    error.value = e.message || 'Login gagal'
  } finally {
    loading.value = false
  }
}
</script>
