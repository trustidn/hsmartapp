<template>
  <div class="min-h-full">
      <!-- Banner limit / expired (hanya setelah data subscription ter-load) -->
      <div v-if="!planLimits.loading && (planLimits.productLimitMessage || planLimits.isExpired)" class="mb-4 space-y-2">
        <div
          v-if="planLimits.isExpired"
          class="bg-amber-50 border border-amber-200 rounded-xl p-4 flex items-center justify-between gap-3"
        >
          <p class="text-amber-800 text-sm">Langganan Anda telah kadaluarsa. Beberapa fitur dibatasi.</p>
          <router-link to="/app/subscription" class="shrink-0 px-4 py-2 rounded-xl bg-amber-600 text-white text-sm font-medium">
            Perpanjang
          </router-link>
        </div>
        <div
          v-else-if="planLimits.productLimitMessage"
          :class="[
            planLimits.productLimitReached ? 'bg-amber-50 border-amber-200' : 'bg-blue-50 border-blue-200',
            'border rounded-xl p-4 flex items-center justify-between gap-3',
          ]"
        >
          <p class="text-sm" :class="planLimits.productLimitReached ? 'text-amber-800' : 'text-blue-800'">
            {{ planLimits.productLimitMessage }}
          </p>
          <router-link
            v-if="planLimits.productLimitReached"
            to="/app/subscription"
            class="shrink-0 px-4 py-2 rounded-xl bg-primary-600 text-white text-sm font-medium"
          >
            Upgrade
          </router-link>
        </div>
      </div>

      <div class="flex justify-end mb-4">
        <button
          type="button"
          class="min-h-touch px-4 rounded-xl bg-primary-600 text-white font-semibold disabled:opacity-50 disabled:cursor-not-allowed"
          :disabled="!planLimits.canAddProduct"
          @click="openNewForm"
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
          class="bg-white rounded-xl border border-gray-100 px-4 py-3 flex items-center gap-3"
        >
          <div class="shrink-0 w-12 h-12 rounded-lg overflow-hidden bg-gray-100 flex items-center justify-center">
            <img
              v-if="p.image_url"
              :src="imageUrl(p.image_url)"
              :alt="p.name"
              class="w-full h-full object-cover"
              loading="lazy"
            />
            <svg v-else class="w-6 h-6 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg>
          </div>
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

      <!-- Empty: sample produk untuk tenant baru -->
      <div v-else class="text-center py-8">
        <p class="text-gray-600 mb-4">Belum ada produk. Mulai dengan menambah produk manual atau gunakan sample untuk mencoba.</p>
        <div class="flex flex-wrap gap-3 justify-center">
          <button
            type="button"
            class="px-5 py-2.5 rounded-xl bg-primary-600 text-white font-medium hover:bg-primary-700"
            @click="showForm = true; editProduct = null"
          >
            + Tambah Produk
          </button>
          <button
            type="button"
            class="px-5 py-2.5 rounded-xl border border-primary-600 text-primary-600 font-medium hover:bg-primary-50 disabled:opacity-50"
            :disabled="addingSamples"
            @click="addSampleProducts"
          >
            {{ addingSamples ? '...' : 'Tambah Sample Produk' }}
          </button>
        </div>
      </div>

      <!-- Modal form -->
      <div
        v-if="showForm"
        class="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-10"
        @click.self="showForm = false"
      >
        <div class="bg-white rounded-2xl p-5 w-full max-w-sm max-h-[90vh] overflow-auto">
          <h3 class="font-semibold text-lg mb-4">{{ editProduct ? 'Edit Produk' : 'Tambah Produk' }}</h3>
          <form @submit.prevent="saveProduct" class="space-y-3">
            <!-- Foto produk -->
            <div>
              <label class="block text-xs text-gray-600 mb-1">Foto (opsional)</label>
              <div class="flex items-center gap-3">
                <div class="shrink-0 w-16 h-16 rounded-xl overflow-hidden bg-gray-100 flex items-center justify-center border border-gray-200">
                  <img
                    v-if="form.imagePreview"
                    :src="form.imagePreview"
                    alt="Preview"
                    class="w-full h-full object-cover"
                  />
                  <svg v-else class="w-8 h-8 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg>
                </div>
                <div class="flex-1 min-w-0">
                  <input
                    ref="photoInput"
                    type="file"
                    accept="image/*"
                    class="hidden"
                    @change="onPhotoSelect"
                  />
                  <button
                    type="button"
                    @click="photoInput?.click()"
                    class="min-h-touch px-4 py-2 rounded-xl border border-gray-300 text-sm font-medium"
                  >
                    {{ form.photoFile ? 'Ganti foto' : 'Pilih foto' }}
                  </button>
                  <button
                    v-if="form.photoFile || form.image_url"
                    type="button"
                    @click="clearPhoto"
                    class="ml-2 min-h-touch px-3 py-2 rounded-xl text-xs text-red-600 hover:bg-red-50"
                  >
                    Hapus
                  </button>
                  <p class="text-xs text-gray-500 mt-1">Otomatis dikompres untuk hemat storage</p>
                </div>
              </div>
            </div>
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
import { usePlanLimits } from '../composables/usePlanLimits'
import { api } from '../lib/api'
import { compressImageToFile } from '../lib/imageCompress'

const auth = useAuthStore()
const toast = useToastStore()
const products = useProductsStore()
const planLimits = usePlanLimits(auth, products)
const productList = computed(() => (products.items ?? []))
const showForm = ref(false)
const editProduct = ref(null)
const saving = ref(false)
const addingSamples = ref(false)
const form = ref({ name: '', price: 0, image_url: '', imagePreview: '', photoFile: null })
const photoInput = ref(null)
const loadError = ref('')

function imageUrl(url) {
  if (!url) return ''
  return url.startsWith('http') ? url : (url.startsWith('/') ? url : '/' + url)
}

async function onPhotoSelect(e) {
  const file = e.target?.files?.[0]
  if (!file || !file.type.startsWith('image/')) return
  try {
    const compressed = await compressImageToFile(file)
    form.value.photoFile = compressed
    form.value.imagePreview = URL.createObjectURL(compressed)
  } catch (err) {
    toast.error(err.message || 'Gagal kompresi gambar')
  }
  e.target.value = ''
}

function clearPhoto() {
  form.value.photoFile = null
  form.value.image_url = ''
  if (form.value.imagePreview) URL.revokeObjectURL(form.value.imagePreview)
  form.value.imagePreview = ''
}

const SAMPLE_PRODUCTS = [
  { name: 'Sample Produk 1', price: 5000 },
  { name: 'Sample Produk 2', price: 10000 },
  { name: 'Sample Produk 3', price: 15000 },
]

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

onMounted(async () => {
  await planLimits.load()
  await loadProducts()
})

function openNewForm() {
  editProduct.value = null
  form.value = { name: '', price: 0, image_url: '', imagePreview: '', photoFile: null }
  showForm.value = true
}

function startEdit(p) {
  editProduct.value = p
  form.value = {
    name: p.name,
    price: p.price,
    image_url: p.image_url || '',
    imagePreview: p.image_url ? imageUrl(p.image_url) : '',
    photoFile: null,
  }
  showForm.value = true
}

async function saveProduct() {
  if (!form.value.name || form.value.price < 0) return
  saving.value = true
  try {
    let imageUrlToSave = form.value.image_url
    if (form.value.photoFile) {
      imageUrlToSave = await api.products.uploadImage(form.value.photoFile)
    } else if (form.value.photoFile === null && editProduct.value && form.value.image_url === '') {
      imageUrlToSave = ''
    }
    if (editProduct.value) {
      await api.products.update(editProduct.value.id, {
        name: form.value.name,
        price: form.value.price,
        image_url: imageUrlToSave || undefined,
      })
      toast.success('Produk berhasil diubah')
    } else {
      await products.create(auth.tenantId, {
        name: form.value.name,
        price: form.value.price,
        image_url: imageUrlToSave || undefined,
      })
      toast.success('Produk berhasil ditambah')
    }
    await loadProducts()
    showForm.value = false
    editProduct.value = null
    form.value = { name: '', price: 0, image_url: '', imagePreview: '', photoFile: null }
  } catch (e) {
    const msg = e.message || 'Gagal menyimpan produk'
    toast.error(msg)
    if (msg.toLowerCase().includes('batas') || msg.toLowerCase().includes('limit')) {
      planLimits.load()
    }
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

async function addSampleProducts() {
  if (!auth.tenantId) return
  addingSamples.value = true
  try {
    for (const sp of SAMPLE_PRODUCTS) {
      await products.create(auth.tenantId, { name: sp.name, price: sp.price })
    }
    await loadProducts()
    toast.success('Sample produk berhasil ditambahkan')
    planLimits.load()
  } catch (e) {
    toast.error(e.message || 'Gagal menambah sample produk')
  } finally {
    addingSamples.value = false
  }
}

function formatNum(n) {
  return Number(n).toLocaleString('id-ID')
}
</script>
