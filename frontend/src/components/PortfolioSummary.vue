<template>
  <div v-if="store.summary" class="cards">
    <div class="card"><span class="lbl">Market Value</span><span class="val">{{ fmt(store.summary.marketValue) }}</span></div>
    <div class="card"><span class="lbl">Total Cost</span><span class="val">{{ fmt(store.summary.totalCost) }}</span></div>
    <div class="card"><span class="lbl">Gain/Loss</span><span class="val" :class="glc(store.summary.totalGL)">{{ fmtGL(store.summary.totalGL) }}</span></div>
    <div class="card"><span class="lbl">G/L %</span><span class="val" :class="glc(store.summary.totalGL)">{{ store.summary.totalGLPct.toFixed(2) }}%</span></div>
    <div class="card"><span class="lbl">Realized G/L</span><span class="val" :class="glc(store.summary.realizedGL)">{{ fmtGL(store.summary.realizedGL) }}</span></div>
  </div>
</template>
<script setup>
import { usePortfolioStore } from '../stores/portfolio'; const store = usePortfolioStore()
function fmt(n) { return Intl.NumberFormat('en-US',{style:'currency',currency:'USD',minimumFractionDigits:2}).format(n) }
function fmtGL(n) { return (n>=0?'+':'')+n.toFixed(2) }
function glc(n) { return n>=0?'pos':'neg' }
</script>
<style scoped>
.cards { display:flex;gap:12px;padding:12px;flex-wrap:wrap; }
.card { flex:1;min-width:130px;padding:12px 16px;background:var(--surface);border:1px solid var(--border);border-radius:6px; }
.lbl { display:block;font-size:11px;color:var(--text-secondary);margin-bottom:4px; }
.val { font-size:16px;font-weight:600;font-variant-numeric:tabular-nums; }
.pos { color:var(--green); } .neg { color:var(--red); }
</style>
