# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [Unreleased]

### Planned
- Technical Indicator Editor (SMA, EMA, RSI, MACD)
- Indicator Scanner
- Stock Market News
- Cloud Sync (backup/restore)
- Currency Exchange
- Export/Import (CSV/JSON)
- i18n Localization
- Stock Comparison

## [0.1.0] — 2026-05-28

### Added
- 28+ world stock exchange support (NYSE, NASDAQ, JKSE, HKEX, TSE, LSE, ASX, SGX, NSE, KRX, TWSE, SSE, SZSE, TSX, B3, JSE, MOEX, and more)
- Real-time price polling with configurable interval
- Candlestick charts via TradingView Lightweight Charts (1D, 1W, 1M, 3M, 6M, 1Y, 5Y, MAX)
- Yahoo Finance API integration (quote, history, symbol search) with User-Agent header
- SQLite local storage (WAL mode, 7 tables)
- Dark theme — WAAG AAA compliant (14.3:1 contrast ratio)
- JetBrains Mono font
- Exchange selector with localStorage persistence
- Symbol search with Yahoo Finance autocomplete + country filter
- Price flash animation on update (green/red)

### Added — Portfolio Management
- Multi-portfolio support with named portfolios
- Buy/Sell transaction recording with quantity, price, fees, date
- Realized & unrealized gain/loss calculation
- Holdings view with cost basis, market value, per-position G/L
- Portfolio summary cards (market value, total cost, total G/L, G/L %, realized)
- Transaction history log
- Sell validation (prevents selling more shares than held)
- Portfolio delete with confirmation dialog
- Empty state with centered "No portfolios" + create button

### Added — Price Alerts
- Price trigger alerts (above/below threshold per stock)
- Background alert check loop (60s interval)
- In-app notification panel (last 5 alerts)
- Web Audio API sound notification (880Hz beep)
- Alert enable/disable toggle
- One-shot trigger (auto-disables after trigger)

### Fixed
- Exchange selector shows wrong exchange on first load → `loadExchanges()` reads localStorage first, `ExchangeSelector` follows `activeExchange` reactively
- Watchlist empty on startup → `watch()` on `activeExchangeId` auto-reloads stocks when exchange changes
- Exchange selector label too generic → show ISO Alpha-2 country code: "JKSE (ID)", "NYSE (US)"
- Yahoo Finance 429 rate limiting → added Mozilla User-Agent header
- Yahoo SearchSymbol deprecated endpoint → switched to query2 + URL encoding
- Yahoo price history `null` values for current day → filter `Close <= 0`
- Double suffix bug (CPIN.JK + .JK) → TrimSuffix in all Yahoo API callers
- Wails v2 model wrapper corrupting nested struct arrays → return raw JSON string
- Stock re-add after removal failing → `ON CONFLICT DO UPDATE SET watchlist=1`
- Go `nil` slices marshaling to JSON `null` → all slice returns default to empty `[]`
- Chart not rendering → `v-show` ref binding fix + lazy init
- Exchange selector not updating watchlist → reactive `watch` on activeExchange
- Empty state navigation on Portfolio/Alerts → inline AddStockDialog in TransactionForm
