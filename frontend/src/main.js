import './assets/main.css'
import { createApp } from 'vue'
import { createRouter, createWebHashHistory } from 'vue-router'
import { createPinia } from 'pinia'
import App from './App.vue'
import WatchlistView from './views/WatchlistView.vue'
import ChartView from './views/ChartView.vue'
import PortfolioView from './views/PortfolioView.vue'
import AlertView from './views/AlertView.vue'

const routes = [
  { path: '/', name: 'watchlist', component: WatchlistView },
  { path: '/chart/:stockId', name: 'chart', component: ChartView },
  { path: '/portfolio', name: 'portfolio', component: PortfolioView },
  { path: '/alerts', name: 'alerts', component: AlertView },
]

const router = createRouter({ history: createWebHashHistory(), routes })
const pinia = createPinia()
const app = createApp(App)
app.use(router)
app.use(pinia)
app.mount('#app')
