<template>
  <div class="overlay" v-if="visible" @click.self="$emit('close')"><div class="form"><h3>Add Transaction</h3>
    <label>Stock</label>
    <select v-if="stocks.length" v-model="f.stockId" class="inp"><option :value="0" disabled>— Select —</option><option v-for="s in stocks" :key="s.stock.id" :value="s.stock.id">{{ s.stock.symbol }} — {{ s.stock.name }}</option></select>
    <div v-else class="nostock">No stocks in watchlist. <button class="btn-add" @click="showAdd=true">+ Add Stock</button></div>
    <label>Type</label><select v-model="f.type" class="inp"><option value="buy">Buy</option><option value="sell">Sell</option></select>
    <label>Quantity</label><input v-model.number="f.quantity" type="number" min="1" class="inp" />
    <label>Price</label><input v-model.number="f.price" type="number" min="0.01" step="0.01" class="inp" />
    <label>Fees</label><input v-model.number="f.fees" type="number" min="0" step="0.01" class="inp" />
    <label>Date</label><input v-model="f.date" type="date" class="inp" />
    <label>Notes</label><input v-model="f.notes" type="text" class="inp" placeholder="Optional" />
    <div v-if="error" class="err">{{ error }}</div>
    <div class="acts"><button class="btn-c" @click="$emit('close')">Cancel</button><button class="btn-s" @click="submit" :disabled="!valid">Add</button></div>
  </div>
  <add-stock-dialog :visible="showAdd" @close="showAdd=false" @added="onAdded" /></div>
</template>
<script setup>
import { ref, computed, onMounted } from 'vue'; import { usePortfolioStore } from '../stores/portfolio'; import { useWatchlistStore } from '../stores/watchlist'; import { useMarketStore } from '../stores/market'; import AddStockDialog from './AddStockDialog.vue'
defineProps({ visible:Boolean }); const emit = defineEmits(['close','added'])
const portfolio = usePortfolioStore(); const watchlist = useWatchlistStore(); const market = useMarketStore()
const error = ref(null); const showAdd = ref(false)
const f = ref({stockId:0,type:'buy',quantity:1,price:0,fees:0,date:new Date().toISOString().split('T')[0],notes:'',portfolioId:0})
const stocks = computed(() => watchlist.stocks)
const valid = computed(() => f.value.stockId>0 && f.value.quantity>0 && f.value.price>0)
onMounted(() => { f.value.portfolioId = portfolio.activePortfolioId })
async function submit() { if (!valid.value) return; error.value=null; try { await portfolio.addTransaction({portfolioId:portfolio.activePortfolioId,stockId:f.value.stockId,type:f.value.type,quantity:f.value.quantity,price:f.value.price,fees:f.value.fees,date:f.value.date,notes:f.value.notes}); emit('added'); emit('close') } catch(e) { error.value=e.message||'Failed' } }
async function onAdded() { await watchlist.fetchWatchlist(market.activeExchangeId) }
</script>
<style scoped>
.overlay { position:fixed;inset:0;background:#0d1117cc;display:flex;align-items:center;justify-content:center;z-index:100; }
.form { background:var(--surface);border:1px solid var(--border);border-radius:8px;padding:20px;width:400px;max-height:80vh;overflow-y:auto; }
h3 { margin:0 0 12px;font-size:16px; }
label { display:block;font-size:11px;color:var(--text-secondary);margin:8px 0 2px;text-transform:uppercase; }
.inp { width:100%;padding:6px 10px;background:var(--bg);border:1px solid var(--border);color:var(--text);font-family:inherit;font-size:13px;border-radius:4px;box-sizing:border-box;color-scheme:dark; }
.inp:focus { outline:2px solid #58a6ff;outline-offset:-1px; } .inp option { background:#161b22;color:#e6edf3; }
.acts { display:flex;gap:8px;margin-top:16px;justify-content:flex-end; }
.btn-c { padding:6px 16px;background:var(--bg);border:1px solid var(--border);color:var(--text);font-family:inherit;font-size:12px;border-radius:4px;cursor:pointer; }
.btn-s { padding:6px 16px;background:#1f6feb22;border:1px solid #58a6ff;color:#58a6ff;font-family:inherit;font-size:12px;border-radius:4px;cursor:pointer; }
.btn-s:disabled { opacity:0.4;cursor:default; }
.err { color:var(--red);font-size:12px;margin-top:8px; }
.nostock { padding:8px;text-align:center;font-size:12px;color:var(--text-secondary); }
.btn-add { background:#1f6feb22;border:1px solid #58a6ff;color:#58a6ff;font-family:inherit;font-size:12px;padding:4px 8px;border-radius:4px;cursor:pointer; }
</style>
