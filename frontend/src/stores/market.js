import { defineStore } from 'pinia'
import { GetExchanges } from '../../wailsjs/go/main/App'

export const useMarketStore = defineStore('market', {
  state: () => ({ exchanges: [], activeExchange: null, polling: false, pollingInterval: 300, connectionStatus: 'online' }),
  getters: { activeExchangeId: (s) => s.activeExchange?.id || 1 },
  actions: {
    async loadExchanges() {
      try {
        this.exchanges = await GetExchanges()
        if (this.exchanges.length && !this.activeExchange) this.setExchange(this.exchanges[0])
      } catch(e) { console.error(e) }
    },
    setExchange(ex) { this.activeExchange = ex; localStorage.setItem('rstock:exchange', String(ex.id)) },
    togglePolling() { this.polling = !this.polling },
    setPollingInterval(s) { this.pollingInterval = s },
    setConnectionStatus(s) { this.connectionStatus = s },
  },
})
