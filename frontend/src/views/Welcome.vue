<template>
  <div class="welcome-page min-h-dvh flex flex-col relative overflow-hidden">
    <!-- Dynamic background: gradient + soft orbs (CSS-only, GPU-friendly) -->
    <div class="welcome-bg" aria-hidden="true">
      <div class="welcome-gradient" />
      <div class="welcome-orb welcome-orb-1" />
      <div class="welcome-orb welcome-orb-2" />
      <div class="welcome-orb welcome-orb-3" />
    </div>

    <main class="flex-1 flex flex-col items-center justify-center px-6 py-12 relative z-10">
      <!-- Logo & nama -->
      <div class="flex flex-col items-center mb-10">
        <img
          v-if="logoUrl"
          :src="logoUrl"
          :alt="appName"
          class="w-20 h-20 object-contain rounded-2xl shadow-lg mb-4 bg-white/80 backdrop-blur-sm"
        />
        <span
          v-else
          class="w-20 h-20 rounded-2xl bg-primary-100/90 flex items-center justify-center text-primary-600 font-bold text-2xl mb-4 shadow-lg backdrop-blur-sm"
        >
          {{ (appName || 'H').charAt(0) }}
        </span>
        <h1 class="text-2xl font-bold text-gray-900 tracking-tight drop-shadow-sm">{{ appName || 'HSmart' }}</h1>
      </div>

      <!-- Copywriting -->
      <div class="max-w-sm text-center space-y-4 mb-12">
        <p class="text-gray-700 text-base leading-relaxed font-medium">
          Kelola penjualan, produk, dan laporan keuangan toko Anda dengan mudah.
        </p>
        <p class="text-gray-600 text-sm leading-relaxed">
          POS sederhana, cepat, dan bisa dipakai offline. Cocok untuk UMKM, warung, dan pedagang jalanan.
        </p>
      </div>

      <!-- CTA -->
      <div class="w-full max-w-xs space-y-3">
        <router-link
          to="/login"
          class="block w-full py-3.5 rounded-xl bg-primary-600 text-white font-semibold text-center hover:bg-primary-700 active:scale-[0.98] transition-all shadow-lg shadow-primary-600/25"
        >
          Masuk
        </router-link>
        <router-link
          to="/register"
          class="block w-full py-3.5 rounded-xl border-2 border-primary-600 text-primary-600 font-semibold text-center hover:bg-primary-50 active:scale-[0.98] transition-all bg-white/80 backdrop-blur-sm"
        >
          Daftar
        </router-link>
      </div>
    </main>
  </div>
</template>

<script setup>
import { useBranding } from '../composables/useBranding'

const { appName, logoUrl } = useBranding()
</script>

<style scoped>
.welcome-page {
  --primary-50: #f0fdf4;
  --primary-100: #dcfce7;
  --primary-600: #16a34a;
}

.welcome-bg {
  position: absolute;
  inset: 0;
  pointer-events: none;
}

@media (prefers-reduced-motion: reduce) {
  .welcome-gradient,
  .welcome-orb-1,
  .welcome-orb-2,
  .welcome-orb-3 {
    animation: none !important;
  }
  .welcome-orb {
    opacity: 0.35;
  }
}

.welcome-gradient {
  position: absolute;
  inset: 0;
  background:
    linear-gradient(
      135deg,
      rgba(34, 197, 94, 0.08) 0%,
      rgba(240, 253, 244, 0.6) 25%,
      rgba(255, 255, 255, 0.95) 50%,
      rgba(241, 245, 241, 0.7) 75%,
      rgba(34, 197, 94, 0.06) 100%
    );
  background-size: 400% 400%;
  animation: welcome-shift 18s ease-in-out infinite;
}

.welcome-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(60px);
  opacity: 0.5;
  will-change: transform;
}

.welcome-orb-1 {
  width: 280px;
  height: 280px;
  background: rgba(34, 197, 94, 0.2);
  top: -10%;
  left: -15%;
  animation: welcome-float-1 22s ease-in-out infinite;
}

.welcome-orb-2 {
  width: 200px;
  height: 200px;
  background: rgba(22, 163, 74, 0.15);
  bottom: 10%;
  right: -10%;
  animation: welcome-float-2 20s ease-in-out infinite;
}

.welcome-orb-3 {
  width: 150px;
  height: 150px;
  background: rgba(34, 197, 94, 0.12);
  top: 45%;
  left: 50%;
  animation: welcome-float-3 25s ease-in-out infinite;
}

@keyframes welcome-shift {
  0%, 100% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
}

@keyframes welcome-float-1 {
  0%, 100% { transform: translate(0, 0) scale(1); }
  33% { transform: translate(30px, -20px) scale(1.05); }
  66% { transform: translate(-15px, 15px) scale(0.95); }
}

@keyframes welcome-float-2 {
  0%, 100% { transform: translate(0, 0) scale(1); }
  50% { transform: translate(-25px, -25px) scale(1.1); }
}

@keyframes welcome-float-3 {
  0%, 100% { transform: translate(-50%, -50%) scale(1); }
  50% { transform: translate(calc(-50% + 20px), calc(-50% - 15px)) scale(1.15); }
}
</style>
