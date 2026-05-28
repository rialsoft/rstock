<template>
  <div class="cv"><button class="back" @click="$router.push('/')">← Back to Watchlist</button><stock-chart v-if="chart.selectedStock" /><div v-else class="empty">No stock selected</div></div>
</template>
<script setup>import StockChart from '../components/StockChart.vue'; import { useChartStore } from '../stores/chart'; const chart = useChartStore()</script>
<style scoped>.cv { display:flex;flex-direction:column;height:100%;grid-column:1/-1; }.back { padding:6px 16px;margin:8px;background:var(--surface);border:1px solid var(--border);color:var(--text);font-family:inherit;font-size:12px;cursor:pointer;border-radius:4px;align-self:flex-start; }.empty { flex:1;display:flex;align-items:center;justify-content:center;color:var(--text-secondary); }</style>
