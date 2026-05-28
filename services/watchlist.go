package services

import (
	"fmt"
	"rstock/models"
	"strings"
	"sync"
	"time"
)

type WatchlistService struct {
	db       *DatabaseService
	yahoo    *YahooFinanceService
	mu       sync.Mutex
	ticker   *time.Ticker
	stopCh   chan struct{}
	onUpdate func([]models.StockWithQuote)
}

func NewWatchlistService(db *DatabaseService, yahoo *YahooFinanceService) *WatchlistService {
	return &WatchlistService{db: db, yahoo: yahoo}
}

func (w *WatchlistService) SetUpdateCallback(fn func([]models.StockWithQuote)) {
	w.onUpdate = fn
}

func (w *WatchlistService) AddSymbol(symbol, countryCode string) (*models.Stock, error) {
	exchanges, _ := w.db.GetExchanges()
	var exchange *models.Exchange
	for _, e := range exchanges {
		if e.Country == countryCode {
			exchange = &e
			break
		}
	}
	if exchange == nil {
		return nil, fmt.Errorf("no exchange for country: %s", countryCode)
	}

	// Strip exchange suffix if symbol already includes it (Yahoo search returns "CPIN.JK")
	cleanSymbol := strings.TrimSuffix(symbol, exchange.YahooSuffix)
	if cleanSymbol == "" {
		cleanSymbol = symbol
	}

	name := cleanSymbol
	if w.yahoo != nil {
		q, err := w.yahoo.FetchQuote(cleanSymbol, exchange.YahooSuffix)
		if err == nil && q.Symbol != "" {
			name = q.Symbol
		}
	}

	id, err := w.db.InsertStock(models.Stock{Symbol: cleanSymbol, Name: name, ExchangeID: exchange.ID, Watchlist: 1})
	if err != nil {
		return nil, err
	}

	if w.yahoo != nil {
		q, err := w.yahoo.FetchQuote(cleanSymbol, exchange.YahooSuffix)
		if err == nil {
			q.StockID = id
			w.db.InsertQuote(*q)
		}
	}

	return &models.Stock{ID: id, Symbol: cleanSymbol, Name: name, ExchangeID: exchange.ID, Watchlist: 1}, nil
}

func (w *WatchlistService) RemoveSymbol(stockID int64) error {
	return w.db.UpdateWatchlistStatus(stockID, false)
}

func (w *WatchlistService) GetWatchlist(exchangeID int64) []models.StockWithQuote {
	result, _ := w.db.GetWatchlist(exchangeID)
	return result
}

func (w *WatchlistService) RefreshAll(exchangeID int64) []models.StockWithQuote {
	if w.yahoo == nil {
		r, _ := w.db.GetWatchlist(exchangeID)
		return r
	}
	exchanges, _ := w.db.GetExchanges()
	exchMap := make(map[int64]models.Exchange)
	for _, e := range exchanges {
		exchMap[e.ID] = e
	}
	stocks, _ := w.db.GetWatchlist(exchangeID)
	for _, swq := range stocks {
		exch, ok := exchMap[swq.Stock.ExchangeID]
		if !ok {
			continue
		}
		sym := strings.TrimSuffix(swq.Stock.Symbol, exch.YahooSuffix)
		q, err := w.yahoo.FetchQuote(sym, exch.YahooSuffix)
		if err == nil {
			q.StockID = swq.Stock.ID
			w.db.InsertQuote(*q)
		}
		time.Sleep(200 * time.Millisecond)
	}
	result, _ := w.db.GetWatchlist(exchangeID)
	return result
}

func (w *WatchlistService) StartPolling(intervalSec int, exchangeID int64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.ticker != nil {
		return
	}
	w.stopCh = make(chan struct{})
	w.ticker = time.NewTicker(time.Duration(intervalSec) * time.Second)
	go func() {
		result := w.RefreshAll(exchangeID)
		if w.onUpdate != nil {
			w.onUpdate(result)
		}
		for {
			select {
			case <-w.ticker.C:
				result := w.RefreshAll(exchangeID)
				if w.onUpdate != nil {
					w.onUpdate(result)
				}
			case <-w.stopCh:
				return
			}
		}
	}()
}

func (w *WatchlistService) StopPolling() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.ticker != nil {
		w.ticker.Stop()
		close(w.stopCh)
		w.ticker = nil
	}
}
