import { defineStore } from 'pinia'
import { GetExchanges } from '../../wailsjs/go/main/App'

export const useMarketStore = defineStore('market', {
  state: () => ({ exchanges: [], activeExchange: null, polling: false, pollingInterval: 300, connectionStatus: 'online' }),
  getters: { activeExchangeId: (s) => s.activeExchange?.id || 1 },
  actions: {
    async loadExchanges() {
      try {
        this.exchanges = await GetExchanges()
        if (!this.exchanges.length) return
        const saved = localStorage.getItem('rstock:exchange')
        const fallback = saved ? this.exchanges.find(e => e.id === Number(saved)) : null
        this.setExchange(fallback || this.exchanges[0])
      } catch(e) { console.error(e) }
    },
    setExchange(ex) { this.activeExchange = ex; localStorage.setItem('rstock:exchange', String(ex.id)) },
    togglePolling() { this.polling = !this.polling },
    setPollingInterval(s) { this.pollingInterval = s },
    setConnectionStatus(s) { this.connectionStatus = s },
  },
})
