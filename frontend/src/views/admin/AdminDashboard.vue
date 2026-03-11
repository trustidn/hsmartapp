<template>
  <div class="min-h-full">
    <!-- Selamat datang -->
    <section class="bg-gradient-to-r from-primary-500 to-primary-600 rounded-2xl p-4 mb-4 text-white shadow-sm">
      <p class="text-sm opacity-90">Selamat datang,</p>
      <p class="text-lg font-bold mt-0.5">{{ adminAuth.name }}</p>
    </section>

    <div class="grid gap-4 sm:grid-cols-2">
      <router-link
        to="/admin/tenants"
        class="block p-6 bg-white rounded-2xl shadow-sm border border-gray-100 hover:shadow-md transition-shadow"
      >
        <h2 class="font-semibold text-gray-800">Daftar Tenant</h2>
        <p class="text-sm text-gray-500 mt-1">Kelola merchant / UMKM terdaftar, status, dan subscription</p>
      </router-link>
      <router-link
        to="/admin/subscription-orders"
        class="block p-6 bg-white rounded-2xl shadow-sm border border-gray-100 hover:shadow-md transition-shadow"
      >
        <h2 class="font-semibold text-gray-800">Order Langganan</h2>
        <p class="text-sm text-gray-500 mt-1">Verifikasi dan setujui order upgrade/perpanjang dari tenant</p>
      </router-link>
      <router-link
        to="/admin/plans"
        class="block p-6 bg-white rounded-2xl shadow-sm border border-gray-100 hover:shadow-md transition-shadow"
      >
        <h2 class="font-semibold text-gray-800">Pengaturan Plan</h2>
        <p class="text-sm text-gray-500 mt-1">Atur max produk dan hari laporan per plan (Free, Premium 1/3/6/12 bln)</p>
      </router-link>
      <router-link
        to="/admin/saas-settings"
        class="block p-6 bg-white rounded-2xl shadow-sm border border-gray-100 hover:shadow-md transition-shadow"
      >
        <h2 class="font-semibold text-gray-800">Pengaturan SaaS</h2>
        <p class="text-sm text-gray-500 mt-1">Nama, logo, kontak admin, rekening pembayaran</p>
      </router-link>
    </div>

    <button
      @click="logout"
      class="mt-6 px-4 py-2 rounded-xl bg-white border border-gray-200 text-gray-700 font-medium hover:bg-gray-50"
    >
      Logout
    </button>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useAdminAuthStore } from '../../stores/adminAuth'

const router = useRouter()
const adminAuth = useAdminAuthStore()

function logout() {
  adminAuth.logout()
  router.replace('/admin/login')
}
</script>
