<template>
  <div class="overlay" v-if="visible" @click.self="$emit('close')">
    <div class="dialog" role="dialog" aria-label="Add stock">
      <h3>Add Stock</h3>
      <input ref="inputRef" v-model="query" placeholder="Search symbol (e.g. BBRI, AAPL)..." class="input" @keydown.escape="$emit('close')" />
      <select v-model="country" class="input select">
        <option value="">— Select country —</option>
        <option v-for="ex in countries" :key="ex.country" :value="ex.country">{{ ex.country }} ({{ ex.code }})</option>
      </select>
      <ul class="suggestions" v-if="suggestions.length">
        <li v-for="s in suggestions" :key="s.symbol" class="item" @click="add(s)" tabindex="0" @keydown.enter="add(s)">
          <strong>{{ s.symbol }}</strong> <span>{{ s.name }}</span> <span class="exch">{{ s.exchange }}</span>
        </li>
      </ul>
      <div v-if="searching">Searching...</div>
      <div v-if="error" class="error">{{ error }}</div>
      <div class="actions"><button class="btn-cancel" @click="$emit('close')">Cancel</button></div>
    </div>
  </div>
</template>
<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { SearchSymbol } from '../../wailsjs/go/main/App'
import { useWatchlistStore } from '../stores/watchlist'; import { useMarketStore } from '../stores/market'
defineProps({ visible: Boolean }); const emit = defineEmits(['close', 'added'])
const market = useMarketStore(); const watchlist = useWatchlistStore()
const query = ref(''); const country = ref(''); const suggestions = ref([]); const searching = ref(false); const error = ref(null); const inputRef = ref(null)
const countries = computed(() => { const seen = new Set(); return market.exchanges.filter(e => { if (seen.has(e.country)) return false; seen.add(e.country); return true }) })
let timer = null
watch(query, () => { clearTimeout(timer); if (query.value.length < 2) { suggestions.value = [] } else { timer = setTimeout(async () => { searching.value = true; error.value = null; try { suggestions.value = await SearchSymbol(query.value) } catch(e) { error.value = 'Search failed' } finally { searching.value = false } }, 300) } })
async function add(s) { const c = country.value || market.activeExchange?.country; if (!c) { error.value = 'Select a country'; return }; try { await watchlist.addStock(s.symbol, c); emit('added'); emit('close') } catch(e) { error.value = e.message || 'Failed' } }
onMounted(() => inputRef.value?.focus())
</script>
<style scoped>
.overlay { position: fixed; inset: 0; background: #0d1117cc; display: flex; align-items: center; justify-content: center; z-index: 100; }
.dialog { background: var(--surface); border: 1px solid var(--border); border-radius: 8px; padding: 20px; width: 440px; max-height: 80vh; overflow-y: auto; }
h3 { margin: 0 0 12px; font-size: 16px; }
.input { width: 100%; padding: 8px 12px; background: var(--bg); border: 1px solid var(--border); color: var(--text); font-family: inherit; font-size: 13px; border-radius: 4px; box-sizing: border-box; margin-bottom: 8px; }
.input:focus { outline: 2px solid #58a6ff; outline-offset: -1px; }
.select option { background: #161b22; color: #e6edf3; }
.suggestions { list-style: none; margin: 0; padding: 0; max-height: 240px; overflow-y: auto; }
.item { display: flex; gap: 12px; padding: 8px 12px; cursor: pointer; border-bottom: 1px solid var(--border); font-size: 12px; }
.item:hover, .item:focus { background: #1f6feb22; outline: none; }
.item strong { min-width: 60px; }
.item span:last-child { color: var(--text-secondary); }
.exch { color: var(--text-secondary); margin-left: auto; }
.error { color: var(--red); font-size: 12px; padding: 8px; }
.actions { display: flex; justify-content: flex-end; margin-top: 8px; }
.btn-cancel { padding: 6px 16px; background: var(--bg); border: 1px solid var(--border); color: var(--text); font-family: inherit; font-size: 12px; border-radius: 4px; cursor: pointer; }
</style>
