import { defineStore } from 'pinia'
import { Ping, GetOHLCV } from '../../wailsjs/go/main/App'

export const useChartStore = defineStore('chart', {
  state: () => ({ selectedStock: null, chartData: [], chartPeriod: '1mo', loading: false, error: null }),
  actions: {
    selectStock(stock) {
      this.selectedStock = stock
      this.error = null
      if (stock?.stock?.id) this.loadChartData(stock.stock.id, this.chartPeriod)
    },
    async loadChartData(stockId, period) {
      this.loading = true
      this.chartPeriod = period
      this.error = null
      try {
        const raw = await GetOHLCV(stockId, period)
        let obj = typeof raw === 'string' ? JSON.parse(raw) : raw
        if (obj?.error) {
          this.chartData = []
          this.error = obj.error + ' stockID=' + obj.stockID + ' symbol=' + obj.symbol + ' suffix=' + obj.suffix + ' db=' + obj.dbEntries
        } else if (Array.isArray(obj)) {
          this.chartData = obj
        } else if (obj?.data) {
          this.chartData = obj.data
        } else {
          this.chartData = []
          this.error = 'Unknown response: ' + (typeof raw === 'string' ? raw.substring(0,200) : JSON.stringify(obj).substring(0,200))
        }
      } catch(e) {
        this.chartData = []
        this.error = e?.message || 'Failed'
      }
      this.loading = false
    },
    setPeriod(p) { this.chartPeriod = p; if (this.selectedStock?.stock?.id) this.loadChartData(this.selectedStock.stock.id, p) },
  },
})
