import { defineStore } from 'pinia'
import { GetAlerts, CreateAlert, ToggleAlert, DeleteAlert } from '../../wailsjs/go/main/App'
import { EventsOn } from '../../wailsjs/runtime/runtime'

export const useAlertStore = defineStore('alerts', {
  state: () => ({ alerts: [], notifications: [], loading: false, error: null }),
  actions: {
    initEvents() {
      EventsOn('alert:triggered', (data) => {
        this.notifications.unshift({ ...data, time: new Date().toLocaleTimeString() })
        try { const ctx = new (window.AudioContext || window.webkitAudioContext)(); const osc = ctx.createOscillator(); const g = ctx.createGain(); osc.connect(g); g.connect(ctx.destination); osc.frequency.value = 880; g.gain.value = 0.3; osc.start(); osc.stop(ctx.currentTime + 0.3) } catch(e) {}
      })
    },
    async loadAlerts() {
      this.loading = true
      try { this.alerts = await GetAlerts() } catch(e) { this.error = e.message || 'Failed' }
      this.loading = false
    },
    async createAlert(stockId, condition, price) { try { await CreateAlert(stockId, condition, price); await this.loadAlerts() } catch(e) { this.error = e.message || 'Failed' } },
    async toggleAlert(id, enabled) { try { await ToggleAlert(id, enabled); await this.loadAlerts() } catch(e) { this.error = e.message || 'Failed' } },
    async deleteAlert(id) { try { await DeleteAlert(id); this.alerts = this.alerts.filter(a => a.id !== id) } catch(e) { this.error = e.message || 'Failed' } },
    clearNotifications() { this.notifications = [] },
  },
})
