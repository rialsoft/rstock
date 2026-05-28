<template>
  <div class="panel">
    <div class="header"><span class="col-sym">Symbol</span><span class="col-name">Name</span><span class="col-price num">Price</span><span class="col-chg num">Chg%</span></div>
    <div class="body">
      <div v-if="loading && !ws.stocks.length" class="empty">Loading...</div>
      <div v-else-if="!ws.stocks.length" class="empty">
        <span>No stocks in watchlist</span>
        <button class="btn-empty" @click="$emit('add')">+ Add Stock</button>
      </div>
      <div v-for="swq in sorted" :key="swq.stock.id" class="row" :class="{ sel: selectedId === swq.stock.id, up: flashUp.has(swq.stock.id), down: flashDown.has(swq.stock.id) }" @click="$emit('select', swq)" tabindex="0" role="row" @keydown.enter="$emit('select', swq)">
        <span class="col-sym">{{ swq.stock.symbol }}</span><span class="col-name">{{ swq.stock.name }}</span><span class="col-price num">{{ fmtPrice(swq.quote) }}</span><span class="col-chg num" :class="chgClass(swq.quote)"><span aria-hidden="true">{{ chgIcon(swq.quote) }}</span> {{ fmtChg(swq.quote) }}</span>
      </div>
    </div>
    <div class="footer"><button class="btn-foot" @click="$emit('add')">+ Add Stock</button><button class="btn-foot btn-rem" @click="$emit('remove', selectedId)" :disabled="!selectedId">Remove</button></div>
  </div>
</template>
<script setup>
import { computed, ref, watch } from 'vue'; import { useWatchlistStore } from '../stores/watchlist'
defineProps({ selectedId: { type: Number, default: null }, loading: Boolean }); defineEmits(['select', 'add', 'remove'])
const ws = useWatchlistStore(); const flashUp = ref(new Set()); const flashDown = ref(new Set())
const sorted = computed(() => [...ws.stocks].sort((a,b) => a.stock.symbol.localeCompare(b.stock.symbol)))
watch(() => ws.stocks, (n, o) => { if (!o) return; n.forEach(s => { const old = o.find(x => x.stock.id === s.stock.id); if (old?.quote && s.quote && old.quote.price !== s.quote.price) { const set = s.quote.price > old.quote.price ? flashUp : flashDown; set.value.add(s.stock.id); setTimeout(() => { set.value.delete(s.stock.id); set.value = new Set(set.value) }, 1500) } }) }, {deep:true})
const fmtPrice = q => q ? Intl.NumberFormat('en-US', {minimumFractionDigits:0,maximumFractionDigits:4}).format(q.price) : '—'
const fmtChg = q => q?.changePct != null ? q.changePct.toFixed(2)+'%' : '—'
const chgClass = q => q?.changePct >= 0 ? 'pos' : 'neg'
const chgIcon = q => q?.changePct >= 0 ? '\u25B2' : '\u25BC'
</script>
<style scoped>
.panel { display: flex; flex-direction: column; height: 100%; border-right: 1px solid var(--border); background: var(--surface); }
.header { display: grid; grid-template-columns: 80px 1fr 90px 80px; padding: 6px 12px; font-size: 11px; color: var(--text-secondary); text-transform: uppercase; letter-spacing: 0.5px; border-bottom: 1px solid var(--border); }
.body { flex: 1; overflow-y: auto; }
.row { display: grid; grid-template-columns: 80px 1fr 90px 80px; padding: 6px 12px; cursor: pointer; border-bottom: 1px solid var(--border); transition: background 0.1s; }
.row:hover { background: #1f6feb11; }
.row:focus { outline: 2px solid #58a6ff; outline-offset: -2px; }
.row.sel { background: #1f6feb22; border-left: 2px solid #58a6ff; }
.up { animation: fg 1.5s; } .down { animation: fr 1.5s; }
@keyframes fg { 0% {background:#3fb95033} 100% {background:transparent} }
@keyframes fr { 0% {background:#f8514933} 100% {background:transparent} }
.num { text-align: right; font-variant-numeric: tabular-nums; }
.col-sym { font-weight: 600; } .col-name { color: var(--text-secondary); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.pos { color: var(--green); } .neg { color: var(--red); }
.empty { padding: 24px; text-align: center; color: var(--text-secondary); display: flex; flex-direction: column; align-items: center; gap: 12px; }
.btn-empty { padding: 6px 20px; background: #1f6feb22; border: 1px solid #58a6ff; color: #58a6ff; font-family: inherit; font-size: 13px; border-radius: 4px; cursor: pointer; }
.footer { display: flex; gap: 8px; padding: 8px 12px; border-top: 1px solid var(--border); }
.btn-foot { padding: 4px 12px; border: 1px solid var(--border); border-radius: 4px; font-family: inherit; font-size: 12px; cursor: pointer; background: var(--bg); color: var(--text); }
.btn-foot:hover { background: #1f6feb22; border-color: #58a6ff; }
.btn-rem:hover { background: #f8514922; border-color: #f85149; }
.btn-rem:disabled { opacity: 0.4; cursor: default; }
</style>
