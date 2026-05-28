<template>
  <div class="pv">
    <portfolio-list-view />
    <template v-if="!store.portfolios.length">
      <div class="empty-full">
        <p>No portfolios</p>
        <button class="btn-empty" @click="showNew = true">+ Create Portfolio</button>
      </div>
    </template>
    <template v-else>
      <portfolio-summary />
      <div class="body"><div class="header"><h2>{{ store.activePortfolio?.name || 'Portfolio' }}</h2><button class="btn" @click="showTx=true">+ Add Transaction</button></div><holdings-table /><transaction-list /></div>
    </template>
    <transaction-form :visible="showTx" @close="showTx=false" @added="showTx=false" />
    <!-- Inline New Portfolio form for empty state -->
    <div v-if="showNew" class="overlay" @click.self="showNew=false">
      <div class="form">
        <h4>New Portfolio</h4>
        <label>Name</label>
        <input v-model="newName" placeholder="e.g. IDX Trading" class="inp" />
        <label>Currency</label>
        <input v-model="newCurrency" placeholder="IDR, USD" class="inp" />
        <label>Initial Cash</label>
        <input v-model.number="newCash" type="number" placeholder="0" class="inp" />
        <div class="acts"><button class="btn-c" @click="showNew=false">Cancel</button><button class="btn-s" @click="create" :disabled="!newName">Create</button></div>
      </div>
    </div>
  </div>
</template>
<script setup>
import { ref, onMounted } from 'vue'
import PortfolioListView from '../components/PortfolioListView.vue'
import PortfolioSummary from '../components/PortfolioSummary.vue'
import HoldingsTable from '../components/HoldingsTable.vue'
import TransactionList from '../components/TransactionList.vue'
import TransactionForm from '../components/TransactionForm.vue'
import { usePortfolioStore } from '../stores/portfolio'
import { useWatchlistStore } from '../stores/watchlist'
import { useMarketStore } from '../stores/market'
const store = usePortfolioStore(); const watchlist = useWatchlistStore(); const market = useMarketStore()
const showTx = ref(false); const showNew = ref(false)
const newName = ref(''); const newCurrency = ref('USD'); const newCash = ref(0)
onMounted(async () => { await market.loadExchanges(); if (market.activeExchangeId) await watchlist.fetchWatchlist(market.activeExchangeId); await store.loadPortfolios() })
async function create() { if (!newName.value) return; await store.createPortfolio(newName.value, newCurrency.value, newCash.value||0); showNew.value=false; newName.value='' }
</script>
<style scoped>
.pv { display:flex;flex-direction:column;height:100%;overflow-y:auto;grid-column:1/-1; }
.empty-full { flex:1;display:flex;flex-direction:column;align-items:center;justify-content:center;gap:16px; }
.empty-full p { color:var(--text-secondary);font-size:14px;margin:0; }
.btn-empty { padding:6px 20px;background:#1f6feb22;border:1px solid #58a6ff;color:#58a6ff;font-family:inherit;font-size:13px;border-radius:4px;cursor:pointer; }
.btn-empty:hover { background:#1f6feb33; }
.body { flex:1;padding-bottom:24px; }
.header { display:flex;align-items:center;justify-content:space-between;padding:8px 12px; }
.header h2 { font-size:16px;margin:0; }
.btn { padding:6px 16px;background:#1f6feb22;border:1px solid #58a6ff;color:#58a6ff;font-family:inherit;font-size:12px;border-radius:4px;cursor:pointer; }
.overlay { position:fixed;inset:0;background:#0d1117cc;display:flex;align-items:center;justify-content:center;z-index:100; }
.form { background:var(--surface);border:1px solid var(--border);border-radius:8px;padding:20px;width:320px; }
h4 { margin:0 0 12px;font-size:14px;color:var(--text); }
label { display:block;font-size:11px;color:var(--text-secondary);margin:8px 0 2px;text-transform:uppercase; }
.inp { width:100%;padding:6px 10px;margin-bottom:4px;background:var(--bg);border:1px solid var(--border);color:var(--text);font-family:inherit;font-size:13px;border-radius:4px;box-sizing:border-box; }
.acts { display:flex;gap:8px;justify-content:flex-end;margin-top:16px; }
.btn-c { padding:6px 16px;background:var(--bg);border:1px solid var(--border);color:var(--text);font-family:inherit;font-size:12px;border-radius:4px;cursor:pointer; }
.btn-s { padding:6px 16px;background:#1f6feb22;border:1px solid #58a6ff;color:#58a6ff;font-family:inherit;font-size:12px;border-radius:4px;cursor:pointer; }
.btn-s:disabled { opacity:0.4;cursor:default; }
</style>
