import Dexie from 'dexie'

export const db = new Dexie('HSmartDB')
db.version(1).stores({
  pending_sales: '++id, tenantId, createdAt, synced',
  pending_expenses: '++id, tenantId, createdAt, synced',
  products_cache: 'id, tenantId, updatedAt',
})
// v2: compound indexes agar query [tenantId+synced] cepat; products_cache pakai put() (upsert)
db.version(2).stores({
  pending_sales: '++id, tenantId, createdAt, synced, [tenantId+synced]',
  pending_expenses: '++id, tenantId, createdAt, synced, [tenantId+synced]',
  products_cache: 'id, tenantId, updatedAt',
})

export async function addPendingSale(tenantId, sale) {
  return db.pending_sales.add({
    tenantId,
    ...sale,
    createdAt: Date.now(),
    synced: false,
  })
}

export async function getPendingSales(tenantId) {
  if (!tenantId) return []
  const rows = await db.pending_sales.toArray()
  return rows.filter((r) => r.tenantId === tenantId && r.synced === false)
}

export async function markSaleSynced(localId) {
  return db.pending_sales.update(localId, { synced: true })
}

export async function addPendingExpense(tenantId, expense) {
  return db.pending_expenses.add({
    tenantId,
    ...expense,
    createdAt: Date.now(),
    synced: false,
  })
}

export async function getPendingExpenses(tenantId) {
  if (!tenantId) return []
  const rows = await db.pending_expenses.toArray()
  return rows.filter((r) => r.tenantId === tenantId && r.synced === false)
}

export async function markExpenseSynced(localId) {
  return db.pending_expenses.update(localId, { synced: true })
}

export async function deletePendingExpense(localId) {
  return db.pending_expenses.delete(localId)
}

/** Offline today expenses list (semua pengeluaran hari ini di lokal) */
export async function getTodayLocalExpenses(tenantId) {
  if (!tenantId) return []
  const today = new Date().toDateString()
  const rows = await db.pending_expenses.toArray()
  return rows
    .filter((r) => r.tenantId === tenantId && new Date(r.createdAt).toDateString() === today)
    .sort((a, b) => b.createdAt - a.createdAt)
    .map((r) => ({
      id: 'local-' + r.id,
      name: r.payload?.name || 'Pengeluaran',
      amount: r.payload?.amount || 0,
      note: r.payload?.note,
      created_at: new Date(r.createdAt).toISOString(),
    }))
}

export async function cacheProducts(tenantId, products) {
  const arr = Array.isArray(products) ? products : []
  const updatedAt = Date.now()
  await db.products_cache.where('tenantId').equals(tenantId).delete()
  for (const p of arr) {
    await db.products_cache.put({ ...p, tenantId, updatedAt })
  }
}

export async function getCachedProducts(tenantId) {
  return db.products_cache.where('tenantId').equals(tenantId).toArray()
}

/** Offline today summary from pending sales/expenses */
export async function getTodayLocalSummary(tenantId) {
  if (!tenantId) return null
  const today = new Date().toDateString()
  const sales = await db.pending_sales.toArray()
  const expenses = await db.pending_expenses.toArray()
  const todaySales = sales.filter((r) => r.tenantId === tenantId && new Date(r.createdAt).toDateString() === today)
  const todayExpenses = expenses.filter((r) => r.tenantId === tenantId && new Date(r.createdAt).toDateString() === today)
  const salesTotal = todaySales.reduce((s, r) => s + (r.total || 0), 0)
  const expenseTotal = todayExpenses.reduce((s, r) => s + (r.payload?.amount || 0), 0)
  return {
    sales_total: salesTotal,
    expense_total: expenseTotal,
    profit: salesTotal - expenseTotal,
    transactions: todaySales.length,
    date: new Date().toISOString().slice(0, 10),
  }
}

/** Offline today sales list (for Dashboard & Transactions) */
export async function getTodayLocalSales(tenantId, limit = 10) {
  if (!tenantId) return []
  const today = new Date().toDateString()
  const [rows, products] = await Promise.all([
    db.pending_sales.toArray(),
    db.products_cache.where('tenantId').equals(tenantId).toArray(),
  ])
  const nameMap = Object.fromEntries((products || []).map((p) => [p.id, p.name || 'Produk']))
  const todaySales = rows
    .filter((r) => r.tenantId === tenantId && new Date(r.createdAt).toDateString() === today)
    .sort((a, b) => b.createdAt - a.createdAt)
    .slice(0, limit)
  return todaySales.map((r) => ({
    id: 'local-' + r.id,
    total: r.total || 0,
    payment_method: r.paymentMethod || 'cash',
    created_at: new Date(r.createdAt).toISOString(),
    items: (r.payload?.items || []).map((i) => ({
      ...i,
      product_name: nameMap[i.product_id] || i.product_name || 'Produk',
    })),
  }))
}

/** Offline today product ranking from pending sales items */
export async function getTodayLocalProductRank(tenantId, limit = 10) {
  if (!tenantId) return []
  const products = await db.products_cache.where('tenantId').equals(tenantId).toArray()
  const nameMap = Object.fromEntries(products.map((p) => [p.id, p.name || 'Produk']))
  const today = new Date().toDateString()
  const sales = await db.pending_sales.toArray()
  const todaySales = sales.filter((r) => r.tenantId === tenantId && new Date(r.createdAt).toDateString() === today)
  const byProduct = {}
  for (const s of todaySales) {
    const items = s.payload?.items || []
    for (const i of items) {
      const key = i.product_id || ('custom-' + (i.product_name || 'Item'))
      if (!byProduct[key]) byProduct[key] = { product_id: key, product_name: nameMap[key] || i.product_name || 'Produk', qty: 0, total: 0 }
      byProduct[key].qty += i.qty || 1
      byProduct[key].total += i.subtotal || (i.price || 0) * (i.qty || 1)
    }
  }
  return Object.values(byProduct)
    .sort((a, b) => b.qty - a.qty || b.total - a.total)
    .slice(0, limit)
}
