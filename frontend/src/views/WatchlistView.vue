<template>
  <watchlist-table :selected-id="chart.selectedStock?.stock?.id" :loading="watchlist.loading" @select="onSelect" @add="showAdd = true" @remove="onRemove" />
  <stock-chart />
  <add-stock-dialog :visible="showAdd" @close="showAdd = false" @added="onAdded" />
</template>
<script setup>
import { ref, onMounted } from 'vue'; import WatchlistTable from '../components/WatchlistTable.vue'; import StockChart from '../components/StockChart.vue'; import AddStockDialog from '../components/AddStockDialog.vue'
import { useWatchlistStore } from '../stores/watchlist'; import { useChartStore } from '../stores/chart'; import { useMarketStore } from '../stores/market'
const watchlist = useWatchlistStore(); const chart = useChartStore(); const market = useMarketStore(); const showAdd = ref(false)
watchlist.initEvents()
onMounted(async () => { if (market.activeExchangeId) await watchlist.fetchWatchlist(market.activeExchangeId) })
function onSelect(swq) { chart.selectStock(swq) }
async function onRemove(id) { if (!id) return; await watchlist.removeStock(id); if (chart.selectedStock?.stock?.id === id) chart.selectStock(null) }
async function onAdded() { await watchlist.fetchWatchlist(market.activeExchangeId) }
</script>
