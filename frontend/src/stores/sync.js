import { api } from '../lib/api'
import { isOnline } from '../lib/api'
import * as db from '../lib/db'

export async function syncPending(tenantId) {
  if (!isOnline() || !tenantId) return
  const sales = await db.getPendingSales(tenantId)
  for (const row of sales) {
    try {
      await api.sales.create(row.payload)
      await db.markSaleSynced(row.id)
    } catch (err) {
      if (import.meta.env.DEV) console.warn('Sync penjualan gagal:', err)
    }
  }
  const expenses = await db.getPendingExpenses(tenantId)
  for (const row of expenses) {
    try {
      await api.expenses.create(row.payload)
      await db.markExpenseSynced(row.id)
    } catch (err) {
      if (import.meta.env.DEV) console.warn('Sync pengeluaran gagal:', err)
    }
  }
}
