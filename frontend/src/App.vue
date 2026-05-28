<template>
  <app-shell>
    <router-view />
  </app-shell>
</template>
<script setup>
import { onMounted } from 'vue'
import AppShell from './components/AppShell.vue'
import { useMarketStore } from './stores/market'
const market = useMarketStore()
onMounted(() => { market.loadExchanges() })
</script>
