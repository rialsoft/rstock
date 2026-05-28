<template>
  <footer class="status-bar">
    <span>{{ watchlist.stockCount }} stocks</span>
    <span v-if="lastUpdate">Last: {{ lastUpdate }}</span>
    <span><span class="dot" :class="market.connectionStatus"></span>{{ market.connectionStatus }}</span>
  </footer>
</template>
<script setup>
import { computed } from 'vue'
import { useWatchlistStore } from '../stores/watchlist'; import { useMarketStore } from '../stores/market'
const watchlist = useWatchlistStore(); const market = useMarketStore()
const lastUpdate = computed(() => { const s = watchlist.stocks; return s.length ? new Date(s[0]?.quote?.fetchedAt).toLocaleTimeString() : null })
</script>
<style scoped>
.status-bar { display: flex; align-items: center; justify-content: space-between; padding: 0 12px; background: var(--surface); border-top: 1px solid var(--border); font-size: 11px; color: var(--text-secondary); }
.dot { display: inline-block; width: 8px; height: 8px; border-radius: 50%; margin-right: 4px; vertical-align: middle; }
.dot.online { background: #3fb950; } .dot.delayed { background: #d29922; } .dot.offline { background: #f85149; }
</style>
