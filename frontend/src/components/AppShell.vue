<template>
  <div class="app-shell">
    <header class="toolbar">
      <div class="toolbar-left">
        <exchange-selector />
        <nav class="nav-links">
          <router-link to="/" class="nav-link" active-class="nav-active">Watchlist</router-link>
          <router-link to="/portfolio" class="nav-link" active-class="nav-active">Portfolio</router-link>
          <router-link to="/alerts" class="nav-link" active-class="nav-active">Alerts</router-link>
        </nav>
      </div>
      <div class="toolbar-actions">
        <button class="btn" @click="polling.toggle" :class="{ active: polling.isPolling.value }">{{ polling.isPolling.value ? 'Pause' : 'Auto' }}</button>
        <button class="btn" @click="manualRefresh">&#x21bb;</button>
      </div>
    </header>
    <main class="content"><slot /></main>
    <status-bar />
  </div>
</template>
<script setup>
import ExchangeSelector from './ExchangeSelector.vue'
import StatusBar from './StatusBar.vue'
import { usePolling } from '../composables/usePolling'
import { useWatchlistStore } from '../stores/watchlist'
import { useMarketStore } from '../stores/market'
const polling = usePolling()
const watchlist = useWatchlistStore()
const market = useMarketStore()
async function manualRefresh() { await watchlist.fetchWatchlist(market.activeExchangeId) }
</script>
<style scoped>
.app-shell { display: grid; grid-template-rows: 44px 1fr 26px; height: 100vh; background: var(--bg); color: var(--text); font-family: 'JetBrains Mono', monospace; font-size: 13px; }
.toolbar { display: flex; align-items: center; justify-content: space-between; padding: 0 12px; background: var(--surface); border-bottom: 1px solid var(--border); gap: 12px; }
.toolbar-left { display: flex; align-items: center; gap: 16px; }
.nav-links { display: flex; gap: 4px; }
.nav-link { padding: 4px 12px; color: var(--text-secondary); text-decoration: none; font-size: 12px; border-radius: 4px; }
.nav-link:hover { color: var(--text); }
.nav-active { color: #58a6ff; background: #1f6feb22; }
.toolbar-actions { display: flex; gap: 8px; }
.btn { padding: 4px 12px; background: var(--surface); border: 1px solid var(--border); color: var(--text); font-family: inherit; font-size: 12px; border-radius: 4px; cursor: pointer; }
.btn:hover { background: var(--border); }
.btn.active { background: #1f6feb33; border-color: #1f6feb; color: #58a6ff; }
.content { display: grid; grid-template-columns: 380px 1fr; overflow: hidden; }
@media (max-width: 800px) { .content { grid-template-columns: 1fr; } }
</style>
