<template>
  <div class="p"><h3>Transaction History</h3>
    <div class="h"><span>Date</span><span>Type</span><span>Symbol</span><span class="num">Qty</span><span class="num">Price</span><span class="num">Fees</span><span class="num">Total</span></div>
    <div v-if="!store.transactions.length" class="e">No transactions</div>
    <div v-for="tx in store.transactions" :key="tx.id" class="r"><span>{{ tx.date }}</span><span :class="'t'+tx.type">{{ tx.type.toUpperCase() }}</span><span>{{ sym(tx.stockId) }}</span><span class="num">{{ tx.quantity }}</span><span class="num">{{ tx.price.toFixed(2) }}</span><span class="num">{{ tx.fees.toFixed(2) }}</span><span class="num">{{ (tx.quantity*tx.price+tx.fees).toFixed(2) }}</span></div>
  </div>
</template>
<script setup>
import { usePortfolioStore } from '../stores/portfolio'; import { useWatchlistStore } from '../stores/watchlist'
const store = usePortfolioStore(); const ws = useWatchlistStore(); const sym = id => ws.stocks.find(s=>s.stock.id===id)?.stock?.symbol||'#'+id
</script>
<style scoped>
.p { padding:0 12px;margin-top:16px; } h3 { font-size:14px;margin:0 0 8px; }
.h,.r { display:grid;grid-template-columns:90px 60px 70px 60px 75px 60px 85px;padding:4px 0;font-size:12px;border-bottom:1px solid var(--border); }
.h { color:var(--text-secondary);font-size:11px;text-transform:uppercase; }
.num { text-align:right;font-variant-numeric:tabular-nums; } .tbuy { color:var(--green);font-weight:600; } .tsell { color:var(--red);font-weight:600; } .e { padding:24px;text-align:center;color:var(--text-secondary); }
</style>
