import { defineStore } from 'pinia'
import { GetWatchlist, AddSymbol, RemoveSymbol } from '../../wailsjs/go/main/App'
import { EventsOn } from '../../wailsjs/runtime/runtime'

export const useWatchlistStore = defineStore('watchlist', {
  state: () => ({ stocks: [], loading: false, error: null }),
  getters: { stockCount: (s) => s.stocks.length },
  actions: {
    initEvents() { EventsOn('stock:updated', (data) => { this.stocks = data; this.loading = false }) },
    async fetchWatchlist(exchangeId) {
      this.loading = true; this.error = null
      try { this.stocks = await GetWatchlist(exchangeId) } catch(e) { this.error = e.message || 'Failed' }
      this.loading = false
    },
    async addStock(symbol, country) {
      try { return await AddSymbol(symbol, country) } catch(e) { this.error = e.message || 'Failed'; throw e }
    },
    async removeStock(id) {
      try { await RemoveSymbol(id); this.stocks = this.stocks.filter(s => s.stock.id !== id) } catch(e) { this.error = e.message || 'Failed' }
    },
  },
})
