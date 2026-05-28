import { defineStore } from 'pinia'
import { GetPortfolios, CreatePortfolio, DeletePortfolio, AddTransaction, GetTransactions, GetPortfolioSummary } from '../../wailsjs/go/main/App'

export const usePortfolioStore = defineStore('portfolio', {
  state: () => ({ portfolios: [], activePortfolioId: null, holdings: [], transactions: [], summary: null, loading: false, error: null }),
  actions: {
    async loadPortfolios() {
      this.loading = true
      try {
        this.portfolios = await GetPortfolios()
        if (this.portfolios.length && !this.activePortfolioId) { this.activePortfolioId = this.portfolios[0].id; await this.loadSummary(this.activePortfolioId) }
      } catch(e) { this.error = e.message || 'Failed' }
      this.loading = false
    },
    async createPortfolio(name, currency, cash) {
      try { const p = await CreatePortfolio(name, currency, cash); await this.loadPortfolios(); this.activePortfolioId = p.id; await this.loadSummary(p.id) } catch(e) { this.error = e.message || 'Failed' }
    },
    async deletePortfolio(id) {
      try { await DeletePortfolio(id); await this.loadPortfolios(); if (this.activePortfolioId === id) { this.activePortfolioId = this.portfolios[0]?.id || null; this.loadSummary(this.activePortfolioId) } } catch(e) { this.error = e.message || 'Failed' }
    },
    async addTransaction(tx) {
      try { await AddTransaction(tx); await this.loadSummary(this.activePortfolioId); await this.loadTransactions(this.activePortfolioId) } catch(e) { this.error = e.message || 'Failed' }
    },
    async loadSummary(pid) {
      this.loading = true
      try { this.summary = await GetPortfolioSummary(pid); this.holdings = this.summary?.holdings || [] } catch(e) { this.error = e.message || 'Failed' }
      this.loading = false
    },
    async loadTransactions(pid) { try { this.transactions = await GetTransactions(pid) } catch(e) { this.error = e.message || 'Failed' } },
    setActivePortfolio(id) { this.activePortfolioId = id; if (id) { this.loadSummary(id); this.loadTransactions(id) } },
  },
})
