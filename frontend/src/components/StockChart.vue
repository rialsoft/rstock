<template>
  <div class="chart-panel">
    <div v-if="!store.selectedStock" class="empty">Select a stock</div>
    <div v-else-if="store.loading && !store.chartData.length" class="empty">{{ store.selectedStock.stock.symbol }} — loading data...</div>
    <div v-else-if="store.error" class="empty" style="padding:20px;text-align:center;">{{ store.error }}</div>
    <div v-else-if="!store.chartData.length" class="empty">No price history for {{ store.selectedStock?.stock?.symbol }}</div>
    <div class="chart-wrap" ref="chartRef" v-show="store.selectedStock && store.chartData.length"></div>
    <div v-if="store.selectedStock" class="timeframe">
      <button v-for="p in tfs" :key="p.v" class="tf-btn" :class="{a:store.chartPeriod===p.v}" @click="store.setPeriod(p.v)">{{ p.l }}</button>
    </div>
  </div>
</template>
<script setup>
import { ref, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { createChart } from 'lightweight-charts'
import { useChartStore } from '../stores/chart'
const store = useChartStore()
const chartRef = ref(null)
let tvChart = null, candle = null, volume = null
const tfs = [{l:'1D',v:'1d'},{l:'1W',v:'1w'},{l:'1M',v:'1mo'},{l:'3M',v:'3mo'},{l:'6M',v:'6mo'},{l:'1Y',v:'1y'},{l:'5Y',v:'5y'},{l:'MAX',v:'max'}]

function create() {
  if (!chartRef.value || tvChart) return
  const w = chartRef.value.clientWidth || chartRef.value.parentElement.clientWidth
  const h = chartRef.value.parentElement.clientHeight * 0.75
  if (w <= 0 || h <= 0) return
  tvChart = createChart(chartRef.value, { width: w, height: h, layout: { background: { color: '#0d1117' }, textColor: '#8b949e' }, grid: { vertLines: { color: '#161b22' }, horzLines: { color: '#161b22' } }, crosshair: { mode: 0 }, rightPriceScale: { borderColor: '#30363d' }, timeScale: { borderColor: '#30363d', timeVisible: true } })
  candle = tvChart.addCandlestickSeries({ upColor: '#3fb950', downColor: '#f85149', borderDownColor: '#f85149', borderUpColor: '#3fb950', wickDownColor: '#f85149', wickUpColor: '#3fb950' })
  volume = tvChart.addHistogramSeries({ priceFormat: { type: 'volume' }, priceScaleId: '' })
  volume.priceScale().applyOptions({ scaleMargins: { top: 0.8, bottom: 0 } })
}

function destroy() {
  try { tvChart?.remove() } catch(e) {}
  tvChart = null; candle = null; volume = null
}

watch(() => store.chartData, async (data) => {
  if (!data || !data.length) return
  await nextTick()
  if (!tvChart) create()
  if (candle && data.length) {
    try {
      candle.setData(data.map(d => ({ time: d.date, open: d.open, high: d.high, low: d.low, close: d.close })))
      volume.setData(data.map(d => ({ time: d.date, value: d.volume, color: d.close >= d.open ? '#3fb95033' : '#f8514933' })))
    } catch(e) {
      store.error = 'Chart error: ' + (e?.message || 'invalid data')
    }
  }
})

watch(() => store.selectedStock, async (stock) => {
  destroy()
  if (!stock) return
  await nextTick()
  create()
})

let ro = null
onMounted(() => {
  const p = chartRef.value?.parentElement
  if (p) { ro = new ResizeObserver(() => { if (tvChart && chartRef.value) tvChart.applyOptions({ width: chartRef.value.parentElement.clientWidth, height: chartRef.value.parentElement.clientHeight * 0.75 }) }); ro.observe(p) }
  if (store.selectedStock && store.chartData.length) { create() }
})

onUnmounted(() => { ro?.disconnect(); destroy() })
</script>
<style scoped>
.chart-panel { display: flex; flex-direction: column; height: 100%; background: var(--bg); }
.chart-wrap { flex: 1; min-height: 0; }
.empty { flex: 1; display: flex; align-items: center; justify-content: center; color: var(--text-secondary); font-size: 14px; }
.timeframe { display: flex; gap: 4px; padding: 8px 12px; border-top: 1px solid var(--border); }
.tf-btn { padding: 2px 10px; background: var(--surface); border: 1px solid var(--border); color: var(--text-secondary); font-family: inherit; font-size: 11px; border-radius: 4px; cursor: pointer; }
.tf-btn:hover { color: var(--text); } .tf-btn.a { background: #1f6feb33; border-color: #58a6ff; color: #58a6ff; }
</style>
