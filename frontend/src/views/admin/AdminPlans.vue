<template>
  <div class="min-h-full">
    <h1 class="text-xl font-bold text-gray-800">Pengaturan Plan</h1>
    <p class="text-gray-500 text-sm mt-1">Atur durasi, max produk, hari laporan, dan harga. -1 = unlimited.</p>

    <div v-if="loading" class="mt-6 space-y-4">
      <div v-for="i in 5" :key="i" class="h-24 bg-gray-200 rounded-2xl animate-pulse" />
    </div>
    <div v-else-if="error" class="mt-6 p-4 bg-red-50 border border-red-200 rounded-xl text-red-700">{{ error }}</div>
    <div v-else class="mt-6 space-y-4">
      <section
        v-for="p in plans"
        :key="p.plan_slug"
        class="p-4 sm:p-5 bg-white rounded-2xl shadow-sm border border-gray-100"
        :class="{ 'opacity-60': !p.is_active }"
      >
        <div class="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-3 mb-4">
          <h2 class="font-semibold text-gray-800">
            {{ p.name }}
            <span class="text-gray-500 font-normal text-sm">({{ p.plan_slug }})</span>
            <span v-if="!p.is_active" class="ml-2 text-amber-600 text-xs font-medium">Nonaktif</span>
          </h2>
          <div class="flex gap-2 shrink-0">
            <button
              v-if="p.is_active && p.plan_slug !== 'free'"
              @click="deletePlan(p.plan_slug)"
              :disabled="deleting === p.plan_slug"
              class="px-3 py-1.5 rounded-lg text-red-600 text-sm font-medium hover:bg-red-50 disabled:opacity-50"
            >
              {{ deleting === p.plan_slug ? '...' : 'Nonaktifkan' }}
            </button>
            <button
              v-if="!p.is_active"
              @click="restorePlan(p.plan_slug)"
              :disabled="restoring === p.plan_slug"
              class="px-3 py-1.5 rounded-lg text-green-600 text-sm font-medium hover:bg-green-50 disabled:opacity-50"
            >
              {{ restoring === p.plan_slug ? '...' : 'Aktifkan' }}
            </button>
          </div>
        </div>
        <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-5">
          <div>
            <label class="block text-xs text-gray-600 mb-1">Urutan</label>
            <input
              v-model.number="edits[p.plan_slug].sort_order"
              type="number"
              min="0"
              class="w-full px-3 py-2 rounded-xl border border-gray-200 focus:ring-2 focus:ring-primary-400 text-sm"
            />
          </div>
          <div>
            <label class="block text-xs text-gray-600 mb-1">Durasi (hari, 0 = unlimited)</label>
            <input
              v-model.number="edits[p.plan_slug].duration_days"
              type="number"
              min="-1"
              class="w-full px-3 py-2 rounded-xl border border-gray-200 focus:ring-2 focus:ring-primary-400 text-sm"
            />
          </div>
          <div>
            <label class="block text-xs text-gray-600 mb-1">Max Produk (-1 = unlimited)</label>
            <input
              v-model.number="edits[p.plan_slug].max_products"
              type="number"
              min="-1"
              class="w-full px-3 py-2 rounded-xl border border-gray-200 focus:ring-2 focus:ring-primary-400 text-sm"
            />
          </div>
          <div>
            <label class="block text-xs text-gray-600 mb-1">Hari Laporan (-1 = unlimited)</label>
            <input
              v-model.number="edits[p.plan_slug].report_days"
              type="number"
              min="-1"
              class="w-full px-3 py-2 rounded-xl border border-gray-200 focus:ring-2 focus:ring-primary-400 text-sm"
            />
          </div>
          <div>
            <label class="block text-xs text-gray-600 mb-1">Harga (Rp)</label>
            <input
              v-model.number="edits[p.plan_slug].price_rupiah"
              type="number"
              min="0"
              class="w-full px-3 py-2 rounded-xl border border-gray-200 focus:ring-2 focus:ring-primary-400 text-sm"
            />
          </div>
        </div>
        <button
          @click="savePlan(p.plan_slug)"
          :disabled="saving === p.plan_slug"
          class="mt-4 w-full sm:w-auto px-4 py-2 rounded-xl bg-primary-600 text-white font-medium hover:bg-primary-700 disabled:opacity-50 text-sm"
        >
          {{ saving === p.plan_slug ? '...' : 'Simpan' }}
        </button>
      </section>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { adminApi } from '../../lib/adminApi'
import { useToast } from '../../composables/useToast'

const { show: toast } = useToast()
const plans = ref([])
const edits = ref({})
const loading = ref(true)
const error = ref('')
const saving = ref('')
const deleting = ref('')
const restoring = ref('')

async function fetchPlans() {
  loading.value = true
  error.value = ''
  try {
    const res = await adminApi.plans.list()
    plans.value = res.plans || []
    const e = {}
    for (const p of plans.value) {
      e[p.plan_slug] = {
        sort_order: p.sort_order ?? 999,
        duration_days: p.duration_days ?? 0,
        max_products: p.max_products,
        report_days: p.report_days,
        price_rupiah: p.price_rupiah ?? 0,
      }
    }
    edits.value = e
  } catch (e) {
    error.value = e.message || 'Gagal memuat plan'
  } finally {
    loading.value = false
  }
}

async function savePlan(planSlug) {
  saving.value = planSlug
  try {
    await adminApi.plans.update(planSlug, {
      sort_order: edits.value[planSlug].sort_order,
      duration_days: edits.value[planSlug].duration_days,
      max_products: edits.value[planSlug].max_products,
      report_days: edits.value[planSlug].report_days,
      price_rupiah: edits.value[planSlug].price_rupiah,
    })
    toast('Plan berhasil disimpan')
  } catch (e) {
    toast(e.message || 'Gagal menyimpan', 'error')
  } finally {
    saving.value = ''
  }
}

async function deletePlan(planSlug) {
  if (!confirm('Nonaktifkan plan ini? Tenant yang sudah memakai plan ini tetap berjalan.')) return
  deleting.value = planSlug
  try {
    await adminApi.plans.delete(planSlug)
    await fetchPlans()
    toast('Plan berhasil dinonaktifkan')
  } catch (e) {
    toast(e.message || 'Gagal menonaktifkan', 'error')
  } finally {
    deleting.value = ''
  }
}

async function restorePlan(planSlug) {
  restoring.value = planSlug
  try {
    await adminApi.plans.restore(planSlug)
    await fetchPlans()
    toast('Plan berhasil diaktifkan')
  } catch (e) {
    toast(e.message || 'Gagal mengaktifkan', 'error')
  } finally {
    restoring.value = ''
  }
}

onMounted(fetchPlans)
</script>
