import { defineStore } from 'pinia'

export const useToastStore = defineStore('toast', {
  state: () => ({
    items: [],
  }),
  actions: {
    show(message, type = 'info') {
      const id = Date.now()
      this.items.push({ id, message, type })
      setTimeout(() => {
        this.dismiss(id)
      }, 3500)
      return id
    },
    success(message) {
      return this.show(message, 'success')
    },
    error(message) {
      return this.show(message, 'error')
    },
    dismiss(id) {
      this.items = this.items.filter((t) => t.id !== id)
    },
  },
})
