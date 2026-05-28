import { ref, onUnmounted } from 'vue'
import { StartPolling, StopPolling } from '../../wailsjs/go/main/App'
import { useMarketStore } from '../stores/market'

export function usePolling() {
  const market = useMarketStore()
  const isPolling = ref(false)
  async function start() {
    if (isPolling.value) return
    try { await StartPolling(market.pollingInterval, market.activeExchangeId); isPolling.value = true; market.polling = true } catch(e) { console.error(e) }
  }
  async function stop() {
    if (!isPolling.value) return
    try { await StopPolling(); isPolling.value = false; market.polling = false } catch(e) { console.error(e) }
  }
  function toggle() { isPolling.value ? stop() : start() }
  onUnmounted(() => { if (isPolling.value) stop() })
  return { isPolling, start, stop, toggle }
}
