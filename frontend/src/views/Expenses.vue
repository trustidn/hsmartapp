<template>
  <div class="min-h-full">
      <div v-if="!isOnline()" class="bg-amber-50 border border-amber-200 rounded-xl p-3 mb-4 text-amber-800 text-sm">
        Offline — pengeluaran disimpan lokal dan akan disinkronkan saat online.
      </div>
      <template v-if="loading">
        <div class="space-y-4 animate-pulse">
          <div class="bg-white rounded-2xl p-4 h-36">
            <div class="h-4 w-16 bg-gray-200 rounded mb-3" />
            <div class="grid grid-cols-4 gap-2">
              <div v-for="i in 4" :key="i" class="h-14 bg-gray-200 rounded-xl" />
            </div>
          </div>
          <div class="bg-white rounded-2xl p-4 h-64">
            <div class="h-4 w-32 bg-gray-200 rounded mb-3" />
            <div class="space-y-3">
              <div class="h-12 bg-gray-200 rounded-xl" />
              <div class="h-12 bg-gray-200 rounded-xl" />
              <div class="h-12 bg-gray-200 rounded-xl" />
              <div class="h-12 bg-gray-200 rounded-xl" />
            </div>
          </div>
          <div class="bg-white rounded-2xl p-4 h-48">
            <div class="h-4 w-48 bg-gray-200 rounded mb-3" />
            <div class="space-y-2">
              <div v-for="i in 3" :key="i" class="h-14 bg-gray-200 rounded-xl" />
            </div>
          </div>
        </div>
      </template>
      <template v-else>
      <!-- Quick expense buttons -->
      <section class="bg-white rounded-2xl shadow-sm p-4 mb-4">
        <h2 class="text-sm font-medium text-gray-500 mb-3">Quick</h2>
        <div class="grid grid-cols-4 gap-2 mb-4">
          <button
            v-for="qb in quickButtons"
            :key="qb.key"
            type="button"
            class="min-h-touch rounded-xl border-2 border-gray-200 text-gray-700 font-medium text-sm active:scale-95 flex flex-col items-center justify-center py-3"
            @click="quickAdd(qb)"
          >
            <span class="text-lg mb-0.5">{{ qb.icon }}</span>
            <span>{{ qb.label }}</span>
          </button>
        </div>
        <div v-if="quickName" class="flex gap-2 items-center">
          <span class="text-sm text-gray-600 shrink-0">{{ quickName }}:</span>
          <input
            v-model.number="quickAmount"
            type="number"
            min="0"
            class="flex-1 min-h-touch px-4 rounded-xl border border-gray-300"
            placeholder="Rp"
          />
          <button
            type="button"
            class="min-h-touch px-4 rounded-xl bg-primary-600 text-white font-semibold shrink-0"
            :disabled="!quickAmount || quickAmount <= 0"
            @click="submitQuick"
          >
            Simpan
          </button>
        </div>
      </section>
      <!-- Manual add -->
      <section class="bg-white rounded-2xl shadow-sm p-4 mb-4">
        <h2 class="text-sm font-medium text-gray-500 mb-3">Catat Pengeluaran</h2>
        <form @submit.prevent="addExpense" class="space-y-3">
          <input
            v-model="form.name"
            type="text"
            placeholder="Contoh: beli minyak"
            class="w-full min-h-touch px-4 rounded-xl border border-gray-300"
            required
          />
          <input
            v-model.number="form.amount"
            type="number"
            min="0"
            placeholder="Jumlah (Rp)"
            class="w-full min-h-touch px-4 rounded-xl border border-gray-300"
            required
          />
          <input
            v-model="form.note"
            type="text"
            placeholder="Catatan (opsional)"
            class="w-full min-h-touch px-4 rounded-xl border border-gray-300"
          />
          <button
            type="submit"
            class="w-full min-h-touch rounded-xl bg-primary-600 text-white font-semibold"
            :disabled="saving"
          >
            {{ saving ? '...' : 'Simpan' }}
          </button>
        </form>
      </section>
      <!-- List -->
      <section class="space-y-2 px-2">
        <h2 class="text-xs font-medium text-gray-500 uppercase tracking-wider mb-3">Daftar Pengeluaran Hari Ini</h2>
        <div v-if="expenseList.length" class="space-y-2">
          <div
            v-for="e in expenseList"
            :key="e.id"
            class="bg-white rounded-xl border border-gray-100 px-4 py-3 flex justify-between items-center group"
          >
            <div class="min-w-0 flex-1">
              <p class="font-medium text-gray-900 truncate text-sm">{{ e.name }}</p>
              <p class="text-xs text-gray-400 mt-0.5">{{ formatDate(e.created_at) }}</p>
            </div>
            <div class="flex items-center gap-2 shrink-0">
              <span class="text-sm font-medium text-red-600 tabular-nums">−{{ formatNum(e.amount) }}</span>
              <button
                v-if="canDeleteExpense(e)"
                type="button"
                class="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50"
                :disabled="deletingId === e.id"
                aria-label="Hapus"
                @click="deleteExpense(e)"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg>
              </button>
            </div>
          </div>
        </div>
        <p v-else class="text-gray-400 text-sm py-6">Belum ada pengeluaran hari ini.</p>
      </section>
      </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useToastStore } from '../stores/toast'
import { api, isOnline, toLocalDateStr } from '../lib/api'
import * as db from '../lib/db'

const auth = useAuthStore()
const toast = useToastStore()
const list = ref([])
const loading = ref(true)
const expenseList = computed(() => (list.value ?? []))
const saving = ref(false)
const form = ref({ name: '', amount: 0, note: '' })
const quickAmount = ref(0)
const quickName = ref('')
const deletingId = ref(null)

const quickButtons = [
  { key: 'gas', label: 'Gas', icon: '⛽' },
  { key: 'oil', label: 'Minyak', icon: '🛢️' },
  { key: 'ingredients', label: 'Bahan', icon: '🥬' },
  { key: 'other', label: 'Lainnya', icon: '📦' },
]

function quickAdd(qb) {
  quickName.value = qb.label
}

async function saveExpenseNow(name, amount) {
  if (!name || amount <= 0) return
  saving.value = true
  try {
    if (isOnline()) {
      await api.expenses.create({ name, amount })
      await loadList()
      toast.success('Pengeluaran berhasil disimpan')
    } else {
      await db.addPendingExpense(auth.tenantId, { payload: { name, amount } })
      list.value = [{ id: 'local', name, amount, created_at: new Date().toISOString() }, ...(list.value ?? [])]
      toast.success('Pengeluaran disimpan (offline)')
    }
  } catch (e) {
    toast.error(e.message || 'Gagal menyimpan pengeluaran')
  } finally {
    saving.value = false
  }
}

async function submitQuick() {
  if (!quickName.value || quickAmount.value <= 0) return
  await saveExpenseNow(quickName.value, quickAmount.value)
  quickAmount.value = 0
  quickName.value = ''
}

onMounted(loadList)

async function loadList() {
  if (!auth.tenantId) {
    loading.value = false
    return
  }
  loading.value = true
  try {
    const today = new Date()
    const dateStr = toLocalDateStr(today)
    if (!isOnline()) {
      list.value = await db.getTodayLocalExpenses(auth.tenantId)
    } else {
      const res = await api.expenses.list(dateStr, dateStr)
      list.value = Array.isArray(res) ? res : []
    }
  } catch {
    list.value = []
  } finally {
    loading.value = false
  }
}

function canDeleteExpense(e) {
  if (!e.id) return false
  return true // Bisa hapus lokal (offline) atau server (online)
}

async function deleteExpense(e) {
  if (!canDeleteExpense(e)) return
  const confirmed = window.confirm(`Hapus pengeluaran "${e.name}" sebesar Rp ${formatNum(e.amount)}?`)
  if (!confirmed) return
  deletingId.value = e.id
  try {
    const id = String(e.id)
    if (id.startsWith('local-')) {
      const localId = parseInt(id.replace('local-', ''), 10)
      await db.deletePendingExpense(localId)
    } else {
      await api.expenses.delete(e.id)
    }
    list.value = list.value.filter((x) => x.id !== e.id)
    toast.success('Pengeluaran berhasil dihapus')
  } catch (err) {
    console.error('Delete expense failed:', err)
    toast.error(err?.message || 'Gagal menghapus pengeluaran')
  } finally {
    deletingId.value = null
  }
}

async function addExpense() {
  const { name, amount, note } = form.value
  if (!name || amount < 0) return
  saving.value = true
  try {
    if (isOnline()) {
      await api.expenses.create({ name, amount, note })
      await loadList()
      toast.success('Pengeluaran berhasil disimpan')
    } else {
      await db.addPendingExpense(auth.tenantId, { payload: { name, amount, note } })
      list.value = [{ id: 'local', name, amount, note, created_at: new Date().toISOString() }, ...(list.value ?? [])]
      toast.success('Pengeluaran disimpan (offline)')
    }
    form.value = { name: '', amount: 0, note: '' }
  } catch (e) {
    toast.error(e.message || 'Gagal menyimpan pengeluaran')
  } finally {
    saving.value = false
  }
}

function formatNum(n) {
  return Number(n).toLocaleString('id-ID')
}
function formatDate(s) {
  if (!s) return ''
  const d = new Date(s)
  return d.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' })
}
</script>
