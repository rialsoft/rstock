<template>
  <div class="av">
    <div class="ah"><h2>Price Alerts</h2><button class="btn" @click="showForm=true">+ New Alert</button></div>
    <div v-if="store.notifications.length" class="notifs"><div v-for="n in store.notifications.slice(0,5)" :key="n.id+'-'+n.time" class="ni"><span class="ni-icon">!</span><div><strong>{{ n.symbol }}</strong> {{ n.condition }} {{ n.targetPrice?.toFixed(2) }} <span class="nic">(now: {{ n.currentPrice?.toFixed(2) }})</span><span class="nit">{{ n.time }}</span></div></div><button class="btn-clr" @click="store.clearNotifications()">Clear</button></div>
    <div v-if="store.loading" class="e">Loading...</div>
    <div v-else-if="!store.alerts.length" class="e">No alerts</div>
    <div v-else class="al"><div v-for="a in store.alerts" :key="a.id" class="ar" :class="{tr:a.triggered,dis:!a.enabled}"><span class="cs">{{ a.symbol }}</span><span class="cn">{{ a.name }}</span><span>{{ a.condition }}</span><span>{{ a.price?.toFixed(2) }}</span><span class="st">{{ a.triggered?'TRIGGERED':a.enabled?'Active':'Paused' }}</span><div class="ca"><button class="bs" @click="store.toggleAlert(a.id,!a.enabled)">{{ a.enabled?'Pause':'Enable' }}</button><button class="bs bd" @click="store.deleteAlert(a.id)">Delete</button></div></div></div>
    <div v-if="showForm" class="overlay" @click.self="showForm=false"><div class="fc"><h3>New Alert</h3>
      <label>Stock</label>
      <select v-if="wl.length" v-model="f.stockId" class="inp"><option :value="0" disabled>— Select —</option><option v-for="s in wl" :key="s.stock.id" :value="s.stock.id">{{ s.stock.symbol }} — {{ s.stock.name }}</option></select>
      <div v-else class="ns">No stocks. <button class="btn-add" @click="showAdd=true">+ Add Stock</button></div>
      <label>Condition</label><select v-model="f.condition" class="inp"><option value="above">Price Above</option><option value="below">Price Below</option></select>
      <label>Price</label><input v-model.number="f.price" type="number" step="0.01" min="0.01" class="inp" />
      <div v-if="fe" class="err">{{ fe }}</div>
      <div class="acts"><button class="btn-c" @click="showForm=false">Cancel</button><button class="btn-s" @click="submit" :disabled="!fv">Create</button></div>
    </div>
    <add-stock-dialog :visible="showAdd" @close="showAdd=false" @added="onAdded" /></div>
  </div>
</template>
<script setup>
import { ref, computed, onMounted } from 'vue'; import { useAlertStore } from '../stores/alerts'; import { useWatchlistStore } from '../stores/watchlist'; import { useMarketStore } from '../stores/market'; import AddStockDialog from '../components/AddStockDialog.vue'
const store = useAlertStore(); const watchlist = useWatchlistStore(); const market = useMarketStore(); const showForm = ref(false); const showAdd = ref(false); const fe = ref(null)
const f = ref({stockId:0,condition:'above',price:0})
const wl = computed(() => watchlist.stocks); const fv = computed(() => f.value.stockId>0&&f.value.price>0&&wl.value.length>0)
store.initEvents()
onMounted(async () => { await market.loadExchanges(); if (market.activeExchangeId) await watchlist.fetchWatchlist(market.activeExchangeId); await store.loadAlerts() })
async function submit() { if (!fv.value) return; fe.value=null; try { await store.createAlert(f.value.stockId,f.value.condition,f.value.price); showForm.value=false; f.value={stockId:0,condition:'above',price:0} } catch(e) { fe.value=e.message||'Failed' } }
async function onAdded() { await watchlist.fetchWatchlist(market.activeExchangeId) }
</script>
<style scoped>
.av { display:flex;flex-direction:column;height:100%;overflow-y:auto;padding:16px;grid-column:1/-1; }
.ah { display:flex;align-items:center;justify-content:space-between;margin-bottom:16px; }
h2 { font-size:16px;margin:0; }
.btn { padding:6px 16px;background:#1f6feb22;border:1px solid #58a6ff;color:#58a6ff;font-family:inherit;font-size:12px;border-radius:4px;cursor:pointer; }
.notifs { margin-bottom:16px;padding:8px;background:var(--surface);border:1px solid #d29922;border-radius:6px; }
.ni { display:flex;align-items:flex-start;gap:8px;padding:4px 0;font-size:12px; }
.ni-icon { color:#d29922;font-weight:700;font-size:14px; }
.nic { color:var(--text-secondary); } .nit { color:var(--text-secondary);margin-left:8px;font-size:11px; }
.btn-clr { margin-top:4px;padding:2px 8px;background:none;border:1px solid var(--border);color:var(--text-secondary);font-family:inherit;font-size:11px;border-radius:3px;cursor:pointer; }
.al { margin-top:8px; }
.ar { display:grid;grid-template-columns:80px 1fr 70px 80px 90px 120px;align-items:center;padding:8px 0;border-bottom:1px solid var(--border);font-size:12px; }
.ar.tr { background:#d2992211; } .ar.dis { opacity:0.5; }
.cs { font-weight:600; } .cn { color:var(--text-secondary);overflow:hidden;text-overflow:ellipsis;white-space:nowrap; }
.st { font-size:11px; } .ar.tr .st { color:#d29922;font-weight:600; }
.ca { display:flex;gap:4px; }
.bs { padding:2px 8px;background:var(--bg);border:1px solid var(--border);color:var(--text);font-family:inherit;font-size:11px;border-radius:3px;cursor:pointer; }
.bs:hover { background:var(--border); } .bd { border-color:#f8514933;color:var(--red); }
.e { padding:32px;text-align:center;color:var(--text-secondary); }
.overlay { position:fixed;inset:0;background:#0d1117cc;display:flex;align-items:center;justify-content:center;z-index:100; }
.fc { background:var(--surface);border:1px solid var(--border);border-radius:8px;padding:20px;width:360px; }
h3 { margin:0 0 12px;font-size:16px; }
label { display:block;font-size:11px;color:var(--text-secondary);margin:8px 0 2px;text-transform:uppercase; }
.inp { width:100%;padding:6px 10px;background:var(--bg);border:1px solid var(--border);color:var(--text);font-family:inherit;font-size:13px;border-radius:4px;box-sizing:border-box;color-scheme:dark; }
.inp option { background:#161b22;color:#e6edf3; }
.acts { display:flex;gap:8px;margin-top:16px;justify-content:flex-end; }
.btn-c { padding:6px 16px;background:var(--bg);border:1px solid var(--border);color:var(--text);font-family:inherit;font-size:12px;border-radius:4px;cursor:pointer; }
.btn-s { padding:6px 16px;background:#1f6feb22;border:1px solid #58a6ff;color:#58a6ff;font-family:inherit;font-size:12px;border-radius:4px;cursor:pointer; }
.btn-s:disabled { opacity:0.4;cursor:default; }
.ns { padding:8px;text-align:center;font-size:12px;color:var(--text-secondary); }
.btn-add { background:#1f6feb22;border:1px solid #58a6ff;color:#58a6ff;font-family:inherit;font-size:12px;padding:4px 8px;border-radius:4px;cursor:pointer; }
.err { color:var(--red);font-size:12px;margin-top:8px; }
</style>
