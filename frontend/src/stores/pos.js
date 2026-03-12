import { defineStore } from 'pinia'
import { api } from '../lib/api'
import { isOnline } from '../lib/api'
import * as db from '../lib/db'

export const usePosStore = defineStore('pos', {
  state: () => ({
    cart: [], // { productId, name, price, qty, subtotal }
    paymentMethod: 'cash',
  }),
  getters: {
    total: (s) => s.cart.reduce((sum, i) => sum + (i.subtotal || i.price * i.qty), 0),
    count: (s) => s.cart.reduce((n, i) => n + i.qty, 0),
  },
  actions: {
    addItem(product, qty = 1) {
      const existing = this.cart.find((i) => i.productId === product.id)
      if (existing) {
        existing.qty += qty
        existing.subtotal = existing.price * existing.qty
      } else {
        this.cart.push({
          productId: product.id,
          name: product.name,
          price: product.price,
          qty,
          subtotal: product.price * qty,
          isCustom: false,
        })
      }
    },
    addCustomItem(name, price, qty = 1) {
      const id = 'custom-' + (typeof crypto !== 'undefined' && crypto.randomUUID ? crypto.randomUUID() : Date.now() + '-' + Math.random().toString(36).slice(2))
      this.cart.push({
        productId: id,
        name: name.trim() || 'Item',
        price: Number(price) || 0,
        qty,
        subtotal: (Number(price) || 0) * qty,
        isCustom: true,
      })
    },
    setQty(productId, qty) {
      if (qty <= 0) {
        this.cart = this.cart.filter((i) => i.productId !== productId)
        return
      }
      const i = this.cart.find((x) => x.productId === productId)
      if (i) {
        i.qty = qty
        i.subtotal = i.price * qty
      }
    },
    removeItem(productId) {
      this.cart = this.cart.filter((i) => i.productId !== productId)
    },
    clearCart() {
      this.cart = []
    },
    async pay(tenantId) {
      const items = this.cart.map((i) => ({
        product_id: i.isCustom ? '' : i.productId,
        product_name: i.isCustom ? i.name : '',
        qty: i.qty,
        price: i.price,
        subtotal: i.subtotal,
      }))
      const total = this.total
      const paymentMethod = this.paymentMethod
      const cartCopy = this.cart.map((i) => ({
        product_name: i.name,
        qty: i.qty,
        subtotal: i.subtotal,
        price: i.price,
      }))
      const payload = { total, payment_method: paymentMethod, items }
      if (isOnline()) {
        const sale = await api.sales.create(payload)
        this.clearCart()
        return { ok: true, sale }
      }
      await db.addPendingSale(tenantId, { payload, total, paymentMethod })
      const localReceipt = {
        id: 'local-' + Date.now(),
        total,
        payment_method: paymentMethod,
        created_at: new Date().toISOString(),
        items: cartCopy,
      }
      this.clearCart()
      return { ok: true, offline: true, sale: localReceipt }
    },
  },
})
