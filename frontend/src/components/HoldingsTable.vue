<template>
  <div class="p"><h3>Current Holdings</h3>
    <div class="h"><span class="cs">Symbol</span><span class="cq num">Qty</span><span class="ca num">Avg</span><span class="cm num">Market</span><span class="cv num">Value</span><span class="cg num">G/L</span><span class="cgl num">%</span></div>
    <div v-if="!store.holdings.length" class="e">No holdings</div>
    <div v-for="h in store.holdings" :key="h.stockId" class="r"><span class="cs">{{ h.symbol }}</span><span class="cq num">{{ h.totalQty }}</span><span class="ca num">{{ h.avgCost.toFixed(2) }}</span><span class="cm num">{{ h.marketPrice?h.marketPrice.toFixed(2):'—' }}</span><span class="cv num">{{ h.marketValue?.toFixed(2)||'—' }}</span><span class="cg num" :class="glc(h.gainLoss)">{{ (h.gainLoss>=0?'+':'')+h.gainLoss.toFixed(2) }}</span><span class="cgl num" :class="glc(h.gainLoss)">{{ h.gainLossPct?h.gainLossPct.toFixed(2)+'%':'—' }}</span></div>
  </div>
</template>
<script setup>import { usePortfolioStore } from '../stores/portfolio'; const store = usePortfolioStore(); function glc(n) { return n>=0?'pos':'neg' }</script>
<style scoped>
.p { padding:0 12px; } h3 { font-size:14px;margin:0 0 8px; }
.h,.r { display:grid;grid-template-columns:80px 60px 80px 80px 90px 90px 70px;padding:4px 0;font-size:12px;border-bottom:1px solid var(--border); }
.h { color:var(--text-secondary);font-size:11px;text-transform:uppercase; } .num { text-align:right;font-variant-numeric:tabular-nums; } .cs { font-weight:600; } .pos { color:var(--green); } .neg { color:var(--red); } .e { padding:24px;text-align:center;color:var(--text-secondary); }
</style>
