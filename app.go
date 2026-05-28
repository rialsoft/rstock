package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"rstock/models"
	"rstock/services"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx           context.Context
	db            *services.DatabaseService
	yahoo         *services.YahooFinanceService
	market        *services.MarketService
	watchlist     *services.WatchlistService
	chart         *services.ChartService
	settings      *services.SettingsService
	portfolio     *services.PortfolioService
	alert         *services.AlertService
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	configDir, _ := os.UserConfigDir()
	dbPath := filepath.Join(configDir, "rstock", "rstock.db")

	db, err := services.NewDatabaseService(dbPath)
	if err != nil { panic(err) }
	db.Migrate()
	db.CleanupOldQuotes()
	cleanupSymbols(db)
	a.db = db

	a.yahoo = services.NewYahooFinanceService()
	a.market = services.NewMarketService(db)
	a.watchlist = services.NewWatchlistService(db, a.yahoo)
	a.chart = services.NewChartService(db, a.yahoo)
	a.settings = services.NewSettingsService(db)
	a.portfolio = services.NewPortfolioService(db, a.yahoo)
	a.alert = services.NewAlertService(db, a.yahoo)

	a.market.SeedDefaultExchanges()

	a.watchlist.SetUpdateCallback(func(data []models.StockWithQuote) {
		runtime.EventsEmit(a.ctx, "stock:updated", data)
	})

	a.alert.SetAlertCallback(func(alert models.Alert, currentPrice float64) {
		runtime.EventsEmit(a.ctx, "alert:triggered", map[string]interface{}{
			"id": alert.ID, "symbol": alert.Symbol, "condition": alert.Condition,
			"targetPrice": alert.Price, "currentPrice": currentPrice,
		})
	})
}

func (a *App) shutdown(ctx context.Context) {
	a.watchlist.StopPolling()
	a.alert.StopChecking()
	if a.db != nil { a.db.Close() }
}

// Market
func (a *App) GetExchanges() []models.Exchange {
	result := a.market.GetExchanges()
	if result == nil { return []models.Exchange{} }
	return result
}
func (a *App) GetWatchlist(exchangeID int64) []models.StockWithQuote {
	result := a.watchlist.GetWatchlist(exchangeID)
	if result == nil { return []models.StockWithQuote{} }
	return result
}
func (a *App) AddSymbol(symbol, countryCode string) (*models.Stock, error) { return a.watchlist.AddSymbol(symbol, countryCode) }
func (a *App) RemoveSymbol(stockID int64) error { return a.watchlist.RemoveSymbol(stockID) }
func (a *App) RefreshAll(exchangeID int64) []models.StockWithQuote {
	result := a.watchlist.RefreshAll(exchangeID)
	if result == nil { return []models.StockWithQuote{} }
	return result
}
func (a *App) StartPolling(intervalSec int, exchangeID int64) { a.watchlist.StartPolling(intervalSec, exchangeID) }
func (a *App) StopPolling() { a.watchlist.StopPolling() }

// Chart
func (a *App) GetChartData(stockID int64, period string) *models.ChartData {
	result := a.chart.GetChartData(stockID, period)
	if result == nil { return &models.ChartData{Period: period, Data: []models.OHLCV{}} }
	if result.Data == nil { result.Data = []models.OHLCV{} }
	return result
}
func (a *App) Ping() string { return "pong" }

func cleanupSymbols(db *services.DatabaseService) {
	exchanges, _ := db.GetExchanges()
	for _, e := range exchanges {
		if e.YahooSuffix == "" { continue }
		// Find stocks with double suffix and fix them
		// Simple approach: trim suffix from symbol column
		db.DB.Exec(`UPDATE stocks SET symbol = REPLACE(symbol, ?, '') WHERE symbol LIKE ? AND exchange_id = ?`,
			e.YahooSuffix, "%"+e.YahooSuffix, e.ID)
	}
}
func (a *App) GetOHLCV(stockID int64, period string) string {
	result := a.chart.GetChartData(stockID, period)
	if result == nil { return `{"error":"nil result","data":[]}` }
	if result.Data == nil { return `{"error":"nil Data","data":[]}` }
	if len(result.Data) == 0 {
		from := "2020-01-01"
		to := "2030-12-31"
		h, _ := a.db.GetPriceHistory(stockID, from, to)
		s, _ := a.db.GetStockByID(stockID)
		var sym, suf string
		if s != nil {
			sym = s.Symbol
			exs, _ := a.db.GetExchanges()
			for _, e := range exs { if e.ID == s.ExchangeID { suf = e.YahooSuffix } }
		}
		return fmt.Sprintf(`{"error":"empty","stockID":%d,"symbol":"%s","suffix":"%s","dbEntries":%d}`, stockID, sym, suf, len(h))
	}
	b, err := json.Marshal(result.Data)
	if err != nil { return `{"error":"marshal err","data":[]}` }
	return string(b)
}

// Yahoo
func (a *App) SearchSymbol(query string) []models.SymbolSuggestion {
	result, _ := a.yahoo.SearchSymbol(query)
	if result == nil { return []models.SymbolSuggestion{} }
	return result
}

// Settings
func (a *App) GetPollingInterval() int { return a.settings.GetPollingInterval() }
func (a *App) SetPollingInterval(sec int) error { return a.settings.SetPollingInterval(sec) }
func (a *App) GetTheme() string { return a.settings.GetTheme() }
func (a *App) SetTheme(t string) error { return a.settings.SetTheme(t) }
func (a *App) GetWindowState() (int, int) { return a.settings.GetWindowState() }
func (a *App) SaveWindowState(w, h int) { a.settings.SaveWindowState(w, h) }

// Portfolio
func (a *App) CreatePortfolio(name, currency string, initialCash float64) (*models.Portfolio, error) { return a.portfolio.CreatePortfolio(name, currency, initialCash) }
func (a *App) GetPortfolios() []models.Portfolio {
	result := a.portfolio.GetPortfolios()
	if result == nil { return []models.Portfolio{} }
	return result
}
func (a *App) DeletePortfolio(id int64) error { return a.portfolio.DeletePortfolio(id) }
func (a *App) AddTransaction(t models.Transaction) (int64, error) { return a.portfolio.AddTransaction(t) }
func (a *App) GetTransactions(portfolioID int64) []models.Transaction {
	result := a.portfolio.GetTransactions(portfolioID)
	if result == nil { return []models.Transaction{} }
	return result
}
func (a *App) GetPortfolioSummary(portfolioID int64) *models.PortfolioSummary { return a.portfolio.GetPortfolioSummary(portfolioID) }

// Alerts
func (a *App) CreateAlert(stockID int64, condition string, price float64) (int64, error) { return a.alert.CreateAlert(stockID, condition, price) }
func (a *App) GetAlerts() []models.Alert {
	result := a.alert.GetAlerts()
	if result == nil { return []models.Alert{} }
	return result
}
func (a *App) ToggleAlert(id int64, enabled bool) error { return a.alert.ToggleAlert(id, enabled) }
func (a *App) DeleteAlert(id int64) error { return a.alert.DeleteAlert(id) }
