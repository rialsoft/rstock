<template>
  <div class="bar">
    <template v-if="store.portfolios.length">
      <div class="tabs">
        <button v-for="p in store.portfolios" :key="p.id" class="tab" :class="{a:store.activePortfolioId===p.id}" @click="store.setActivePortfolio(p.id)">
          {{ p.name }}
          <span class="tab-del" @click.stop="confirmRemove(p)">×</span>
        </button>
      </div>
      <button class="btn" @click="show = true">+ New</button>
    </template>

    <div v-if="show" class="overlay" @click.self="show=false">
      <div class="form">
        <h4>New Portfolio</h4>
        <label>Name</label>
        <input v-model="name" placeholder="e.g. IDX Trading" class="inp" />
        <label>Currency</label>
        <input v-model="currency" placeholder="IDR, USD" class="inp" />
        <label>Initial Cash</label>
        <input v-model.number="cash" type="number" placeholder="0" class="inp" />
        <div class="acts"><button class="btn-c" @click="show=false">Cancel</button><button class="btn-s" @click="create" :disabled="!name">Create</button></div>
      </div>
    </div>

    <div v-if="confirmTarget" class="overlay" @click.self="confirmTarget=null">
      <div class="confirm-box">
        <p>Delete portfolio "{{ confirmTarget.name }}"?</p>
        <div class="acts">
          <button class="btn-c" @click="confirmTarget=null">Cancel</button>
          <button class="btn-s btn-del" @click="doDelete">Delete</button>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup>
import { ref } from 'vue'; import { usePortfolioStore } from '../stores/portfolio'
const store = usePortfolioStore(); const show = ref(false); const name = ref(''); const currency = ref('USD'); const cash = ref(0)
const confirmTarget = ref(null)
async function create() { if (!name.value) return; await store.createPortfolio(name.value, currency.value, cash.value||0); show.value=false; name.value='' }
function confirmRemove(p) { confirmTarget.value = p }
async function doDelete() { if (!confirmTarget.value) return; await store.deletePortfolio(confirmTarget.value.id); confirmTarget.value = null }
</script>
<style scoped>
.bar { display:flex;align-items:center;justify-content:space-between;padding:0 12px;border-bottom:1px solid var(--border);min-height:44px; }
.empty-portfolio { display:flex;flex-direction:column;align-items:center;justify-content:center;gap:12px;padding:24px 0; }
.empty-portfolio p { color:var(--text-secondary);font-size:14px;margin:0; }
.btn-empty { padding:6px 20px;background:#1f6feb22;border:1px solid #58a6ff;color:#58a6ff;font-size:13px; }
.tabs { display:flex;gap:2px; }
.tab { padding:8px 12px 8px 16px;background:none;border:none;border-bottom:2px solid transparent;color:var(--text-secondary);font-family:inherit;font-size:13px;cursor:pointer;display:flex;align-items:center;gap:6px; }
.tab:hover { color:var(--text); } .tab.a { color:var(--text);border-bottom-color:#58a6ff; }
.tab-del { font-size:14px;color:var(--text-secondary);line-height:1;padding:0 2px; }
.tab-del:hover { color:var(--red); }
.btn { padding:4px 12px;background:var(--bg);border:1px solid var(--border);color:var(--text);font-family:inherit;font-size:12px;border-radius:4px;cursor:pointer; }
.overlay { position:fixed;inset:0;background:#0d1117cc;display:flex;align-items:center;justify-content:center;z-index:100; }
.form { background:var(--surface);border:1px solid var(--border);border-radius:8px;padding:20px;width:320px; }
.confirm-box { background:var(--surface);border:1px solid var(--border);border-radius:8px;padding:20px;width:360px;text-align:center; }
.confirm-box p { color:var(--text);font-size:14px;margin:0 0 16px; }
h4 { margin:0 0 12px;font-size:14px; }
.inp { width:100%;padding:6px 10px;margin-bottom:8px;background:var(--bg);border:1px solid var(--border);color:var(--text);font-family:inherit;font-size:13px;border-radius:4px;box-sizing:border-box; }
.acts { display:flex;gap:8px;justify-content:center; }
.btn-c { padding:6px 16px;background:var(--bg);border:1px solid var(--border);color:var(--text);font-family:inherit;font-size:12px;border-radius:4px;cursor:pointer; }
.btn-s { padding:6px 16px;background:#1f6feb22;border:1px solid #58a6ff;color:#58a6ff;font-family:inherit;font-size:12px;border-radius:4px;cursor:pointer; }
.btn-s:disabled { opacity:0.4;cursor:default; }
.btn-del { background:#f8514922 !important;border-color:#f85149 !important;color:#f85149 !important; }
.btn-del:hover { background:#f8514933 !important; }
</style>
