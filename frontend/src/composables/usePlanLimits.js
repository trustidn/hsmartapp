import { ref, computed } from 'vue'
import { api } from '../lib/api'
import { useSubscriptionStore } from '../stores/subscription'

/**
 * Composable untuk limit plan tenant: max_products, report_days, status expired.
 * Jika subscription expired, otomatis pakai limit plan Free.
 */
export function usePlanLimits(auth, productsStore) {
  const subscriptionStore = useSubscriptionStore()
  const subscription = ref(null)
  const plans = ref([])
  const loading = ref(false)

  /** Cek apakah tanggal expired_at sudah lewat. Invalid/tidak ada = tidak expired. */
  function checkExpired(expiredAt) {
    if (!expiredAt || typeof expiredAt !== 'string') return false
    try {
      const end = new Date(expiredAt)
      if (Number.isNaN(end.getTime())) return false
      return end < new Date()
    } catch {
      return false
    }
  }

  const effectivePlan = computed(() => {
    const sub = subscription.value ?? subscriptionStore.data
    if (!sub?.plan) return 'free'
    if (sub.expired_at && checkExpired(sub.expired_at)) return 'free'
    return sub.plan
  })

  /** Pakai data subscription dari load() - jika plan premium & expired_at masa depan = tidak expired */
  const isExpired = computed(() => {
    const sub = subscription.value ?? subscriptionStore.data
    if (!sub) return false
    if (!sub.expired_at) return false
    if (sub.plan && sub.plan !== 'free' && !checkExpired(sub.expired_at)) return false
    return checkExpired(sub.expired_at)
  })

  const planConfig = computed(() => {
    const slug = effectivePlan.value
    return plans.value.find((p) => p.plan_slug === slug) || null
  })

  const maxProducts = computed(() => {
    const p = planConfig.value
    if (!p || p.max_products == null) return 10
    return p.max_products < 0 ? 999999 : p.max_products
  })

  const reportDays = computed(() => {
    const p = planConfig.value
    const slug = effectivePlan.value
    let days = 7
    if (p && p.report_days != null) {
      days = p.report_days < 0 ? 999999 : p.report_days
    }
    // Plan free selalu minimal 7 hari laporan (sesuai default di plan_config)
    if (slug === 'free' && days < 7) days = 7
    return days
  })

  const productCount = computed(() => (productsStore?.items ?? []).length)

  const productLimitReached = computed(() => maxProducts.value >= 0 && productCount.value >= maxProducts.value)
  const productLimitNear = computed(() => {
    if (maxProducts.value < 0) return false
    const remaining = maxProducts.value - productCount.value
    return remaining <= 2 && remaining >= 0
  })

  const canAddProduct = computed(() => !productLimitReached.value)
  const allowedReportFilters = computed(() => {
    const days = reportDays.value
    const slug = effectivePlan.value
    const all = [
      { key: 'today', label: 'Hari Ini', days: 1 },
      { key: '7d', label: '7 Hari', days: 7 },
      { key: '30d', label: '30 Hari', days: 30 },
      { key: '12m', label: '12 Bulan', days: 365 },
    ]
    let list = all.filter((f) => f.days <= days)
    // Plan free: selalu minimal Hari Ini + 7 Hari (override jika konfig salah/kosong)
    if (slug === 'free' && list.length < 2) {
      list = [all[0], all[1]]
    }
    return list.length > 0 ? list : [{ key: 'today', label: 'Hari Ini', days: 1 }]
  })

  const productLimitMessage = computed(() => {
    if (productLimitReached.value) {
      return `Limit produk tercapai (${productCount.value}/${maxProducts.value}). Upgrade untuk menambah produk.`
    }
    if (productLimitNear.value) {
      const left = maxProducts.value - productCount.value
      return `Hampir mencapai limit (${productCount.value}/${maxProducts.value}). ${left} slot tersisa.`
    }
    return null
  })

  const reportLimitMessage = computed(() => {
    if (reportDays.value >= 365) return null
    return `Plan Anda membatasi laporan hingga ${reportDays.value} hari terakhir. Upgrade untuk mengakses data lebih lama.`
  })

  async function load() {
    if (!auth?.tenantId) return
    loading.value = true
    try {
      const [subRes, plansRes] = await Promise.all([
        api.subscription.get().catch(() => null),
        api.plans.list().catch(() => ({ plans: [] })),
      ])
      subscription.value = subRes
      if (subRes) subscriptionStore.data = subRes
      plans.value = plansRes?.plans ?? []
    } catch {
      subscription.value = null
      plans.value = []
    } finally {
      loading.value = false
    }
  }

  return {
    subscription,
    plans,
    loading,
    effectivePlan,
    isExpired,
    maxProducts,
    reportDays,
    productCount,
    productLimitReached,
    productLimitNear,
    canAddProduct,
    allowedReportFilters,
    productLimitMessage,
    reportLimitMessage,
    load,
  }
}
