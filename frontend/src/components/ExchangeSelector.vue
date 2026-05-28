<template>
  <select v-model="selected" @change="onChange" class="exch-select">
    <option v-for="ex in market.exchanges" :key="ex.id" :value="ex.id">{{ ex.code }} ({{ ex.country }})</option>
  </select>
</template>
<script setup>
import { ref, watch } from 'vue'
import { useMarketStore } from '../stores/market'
import { useWatchlistStore } from '../stores/watchlist'
const market = useMarketStore(); const watchlist = useWatchlistStore(); const selected = ref(0)
watch(() => market.activeExchange, (ex) => { if (ex) selected.value = ex.id }, { immediate: true })
function onChange() { const ex = market.exchanges.find(e => e.id === selected.value); if (ex) { market.setExchange(ex); watchlist.fetchWatchlist(ex.id) } }
</script>
<style scoped>
.exch-select { padding: 4px 8px; background: var(--bg); border: 1px solid var(--border); color: var(--text); font-family: inherit; font-size: 12px; border-radius: 4px; cursor: pointer; color-scheme: dark; }
.exch-select:focus { outline: 2px solid #58a6ff; outline-offset: -1px; }
.exch-select option { background: #161b22; color: #e6edf3; }
</style>
