import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useAdminAuthStore } from '../stores/adminAuth'
import TenantLayout from '../layouts/TenantLayout.vue'
import AdminLayout from '../layouts/AdminLayout.vue'
import POS from '../views/POS.vue'
import Dashboard from '../views/Dashboard.vue'
import Products from '../views/Products.vue'
import Expenses from '../views/Expenses.vue'
import Reports from '../views/Reports.vue'
import Settings from '../views/Settings.vue'
import More from '../views/More.vue'

const routes = [
  { path: '/', name: 'Welcome', component: () => import('../views/Welcome.vue'), meta: { public: true } },
  { path: '/login', name: 'Login', component: () => import('../views/Login.vue'), meta: { guest: true } },
  { path: '/register', name: 'Register', component: () => import('../views/Register.vue'), meta: { guest: true } },
  { path: '/admin/login', name: 'AdminLogin', component: () => import('../views/admin/AdminLogin.vue'), meta: { adminGuest: true } },
  {
    path: '/admin',
    component: AdminLayout,
    meta: { requiresAdmin: true },
    children: [
      { path: '', name: 'AdminDashboard', component: () => import('../views/admin/AdminDashboard.vue'), meta: { title: 'Admin' } },
      { path: 'tenants', name: 'AdminTenants', component: () => import('../views/admin/AdminTenants.vue'), meta: { title: 'Tenants' } },
      { path: 'tenants/:id', name: 'AdminTenantDetail', component: () => import('../views/admin/AdminTenantDetail.vue'), meta: { title: 'Detail Tenant' } },
      { path: 'plans', name: 'AdminPlans', component: () => import('../views/admin/AdminPlans.vue'), meta: { title: 'Pengaturan Plan' } },
      { path: 'subscription-orders', name: 'AdminSubscriptionOrders', component: () => import('../views/admin/AdminSubscriptionOrders.vue'), meta: { title: 'Order Langganan' } },
      { path: 'saas-settings', name: 'AdminSaasSettings', component: () => import('../views/admin/AdminSaasSettings.vue'), meta: { title: 'Pengaturan SaaS' } },
      { path: 'profile', name: 'AdminProfile', component: () => import('../views/admin/AdminProfile.vue'), meta: { title: 'Profil Admin' } },
    ],
  },
  {
    path: '/app',
    component: TenantLayout,
    meta: { requiresAuth: true },
    children: [
      { path: '', name: 'Dashboard', component: Dashboard, meta: { title: 'Dashboard' } },
      { path: 'pos', name: 'POS', component: POS, meta: { title: 'POS' } },
      { path: 'products', name: 'Products', component: Products, meta: { title: 'Produk' } },
      { path: 'expenses', name: 'Expenses', component: Expenses, meta: { title: 'Pengeluaran' } },
      { path: 'reports', name: 'Reports', component: Reports, meta: { title: 'Laporan' } },
      { path: 'settings', name: 'Settings', component: Settings, meta: { title: 'Pengaturan' } },
      { path: 'subscription', name: 'Subscription', component: () => import('../views/Subscription.vue'), meta: { title: 'Langganan' } },
      { path: 'more', name: 'More', component: More, meta: { title: 'Lainnya' } },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach(async (to) => {
  const auth = useAuthStore()
  const adminAuth = useAdminAuthStore()
  if (to.meta.requiresAdmin && !adminAuth.isAdminLoggedIn) return { name: 'AdminLogin' }
  if (to.meta.adminGuest && adminAuth.isAdminLoggedIn) return { name: 'AdminDashboard' }
  if (to.meta.requiresAuth && !auth.isLoggedIn) return { name: 'Welcome' }
  if ((to.meta.guest || to.meta.public) && auth.isLoggedIn && to.name !== 'Welcome') return { name: 'Dashboard' }
  if (to.name === 'Welcome' && auth.isLoggedIn) return { name: 'Dashboard' }
  return true
})

export default router
