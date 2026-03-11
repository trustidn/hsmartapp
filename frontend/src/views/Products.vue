<template>
  <div class="min-h-full">
      <div class="flex justify-end mb-4">
        <button
          type="button"
          class="min-h-touch px-4 rounded-xl bg-primary-600 text-white font-semibold"
          @click="showForm = true; editProduct = null"
        >
          + Tambah Produk
        </button>
      </div>

      <!-- Loading skeleton -->
      <div v-if="products.loading" class="space-y-2 animate-pulse">
        <div v-for="i in 6" :key="i" class="bg-white rounded-xl p-4 flex items-center gap-4">
          <div class="flex-1">
            <div class="h-4 w-32 bg-gray-200 rounded mb-2" />
            <div class="h-3 w-20 bg-gray-200 rounded" />
          </div>
          <div class="flex gap-2">
            <div class="h-10 w-14 bg-gray-200 rounded-lg" />
            <div class="h-10 w-14 bg-gray-200 rounded-lg" />
          </div>
        </div>
      </div>

      <!-- Error -->
      <div v-else-if="loadError" class="bg-red-50 border border-red-200 rounded-xl p-4 text-center">
        <p class="text-red-700 text-sm">{{ loadError }}</p>
        <button
          type="button"
          class="mt-3 min-h-touch px-4 rounded-xl bg-red-600 text-white font-medium text-sm"
          @click="retryLoad"
        >
          Coba lagi
        </button>
      </div>

      <!-- List -->
      <div v-else-if="productList.length" class="space-y-2 px-2">
        <div
          v-for="p in productList"
          :key="p.id"
          class="bg-white rounded-xl border border-gray-100 px-4 py-3 flex items-center justify-between"
        >
          <div class="min-w-0 flex-1">
            <p class="font-medium text-gray-900 text-sm">{{ p.name }}</p>
            <p class="text-sm text-primary-600 font-medium tabular-nums mt-0.5">Rp {{ formatNum(p.price) }}</p>
          </div>
          <div class="flex gap-1.5 shrink-0">
            <button
              type="button"
              class="px-3 py-1.5 rounded-lg text-xs font-medium text-gray-500 hover:bg-gray-100 active:bg-gray-100"
              @click="startEdit(p)"
            >
              Edit
            </button>
            <button
              type="button"
              class="px-3 py-1.5 rounded-lg text-xs font-medium text-red-500 hover:bg-red-50 active:bg-red-50"
              @click="removeProduct(p)"
            >
              Hapus
            </button>
          </div>
        </div>
      </div>

      <!-- Empty -->
      <p v-else class="text-gray-400 text-center py-8">Belum ada produk. Klik &quot;+ Tambah Produk&quot; untuk menambah.</p>

      <!-- Modal form -->
      <div
        v-if="showForm"
        class="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-10"
        @click.self="showForm = false"
      >
        <div class="bg-white rounded-2xl p-5 w-full max-w-sm">
          <h3 class="font-semibold text-lg mb-4">{{ editProduct ? 'Edit Produk' : 'Tambah Produk' }}</h3>
          <form @submit.prevent="saveProduct" class="space-y-3">
            <input
              v-model="form.name"
              type="text"
              placeholder="Nama produk"
              class="w-full min-h-touch px-4 rounded-xl border border-gray-300"
              required
            />
            <input
              v-model.number="form.price"
              type="number"
              min="0"
              placeholder="Harga (Rp)"
              class="w-full min-h-touch px-4 rounded-xl border border-gray-300"
              required
            />
            <div class="flex gap-2 pt-2">
              <button
                type="button"
                class="flex-1 min-h-touch rounded-xl border border-gray-300"
                @click="showForm = false"
              >
                Batal
              </button>
              <button
                type="submit"
                class="flex-1 min-h-touch rounded-xl bg-primary-600 text-white font-semibold"
                :disabled="saving"
              >
                {{ saving ? '...' : 'Simpan' }}
              </button>
            </div>
          </form>
        </div>
      </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useProductsStore } from '../stores/products'
import { useToastStore } from '../stores/toast'
import { api } from '../lib/api'

const auth = useAuthStore()
const toast = useToastStore()
const products = useProductsStore()
const productList = computed(() => (products.items ?? []))
const showForm = ref(false)
const editProduct = ref(null)
const saving = ref(false)
const form = ref({ name: '', price: 0 })
const loadError = ref('')

async function loadProducts() {
  loadError.value = ''
  if (!auth.tenantId) {
    loadError.value = 'Sesi tidak valid. Silakan login lagi.'
    return
  }
  try {
    await products.load(auth.tenantId)
  } catch (e) {
    loadError.value = e.message || 'Gagal memuat produk.'
  }
}

function retryLoad() {
  loadProducts()
}

onMounted(loadProducts)

function startEdit(p) {
  editProduct.value = p
  form.value = { name: p.name, price: p.price }
  showForm.value = true
}

async function saveProduct() {
  if (!form.value.name || form.value.price < 0) return
  saving.value = true
  try {
    if (editProduct.value) {
      await api.products.update(editProduct.value.id, { name: form.value.name, price: form.value.price })
      toast.success('Produk berhasil diubah')
    } else {
      await products.create(auth.tenantId, { name: form.value.name, price: form.value.price })
      toast.success('Produk berhasil ditambah')
    }
    await loadProducts()
    showForm.value = false
    editProduct.value = null
    form.value = { name: '', price: 0 }
  } catch (e) {
    toast.error(e.message || 'Gagal menyimpan produk')
  } finally {
    saving.value = false
  }
}

async function removeProduct(p) {
  if (!confirm('Hapus produk ini?')) return
  try {
    await api.products.delete(p.id)
    await loadProducts()
    toast.success('Produk berhasil dihapus')
  } catch (e) {
    toast.error(e.message || 'Gagal menghapus produk.')
  }
}

function formatNum(n) {
  return Number(n).toLocaleString('id-ID')
}
</script>
