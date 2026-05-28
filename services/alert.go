package services

import (
	"fmt"
	"rstock/models"
	"strings"
	"sync"
	"time"
)

type AlertService struct {
	db      *DatabaseService
	yahoo   *YahooFinanceService
	mu      sync.Mutex
	stopCh  chan struct{}
	onAlert func(models.Alert, float64)
}

func NewAlertService(db *DatabaseService, yahoo *YahooFinanceService) *AlertService {
	return &AlertService{db: db, yahoo: yahoo}
}

func (a *AlertService) SetAlertCallback(fn func(models.Alert, float64)) {
	a.onAlert = fn
}

func (a *AlertService) CreateAlert(stockID int64, condition string, price float64) (int64, error) {
	if condition != "above" && condition != "below" {
		return 0, fmt.Errorf("invalid condition")
	}
	if price <= 0 {
		return 0, fmt.Errorf("price must be positive")
	}
	stock, _ := a.db.GetStockByID(stockID)
	if stock == nil {
		return 0, fmt.Errorf("stock not found")
	}
	return a.db.InsertAlert(models.Alert{StockID: stockID, Condition: condition, Price: price})
}

func (a *AlertService) GetAlerts() []models.Alert {
	result, _ := a.db.GetAlerts()
	return result
}

func (a *AlertService) ToggleAlert(id int64, enabled bool) error {
	return a.db.ToggleAlert(id, enabled)
}

func (a *AlertService) DeleteAlert(id int64) error {
	return a.db.DeleteAlert(id)
}

func (a *AlertService) StartChecking(intervalSec int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.stopCh != nil {
		return
	}
	a.stopCh = make(chan struct{})
	ticker := time.NewTicker(time.Duration(intervalSec) * time.Second)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				a.checkAlerts()
			case <-a.stopCh:
				return
			}
		}
	}()
}

func (a *AlertService) StopChecking() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.stopCh != nil {
		close(a.stopCh)
		a.stopCh = nil
	}
}

func (a *AlertService) checkAlerts() {
	alerts, _ := a.db.GetEnabledAlerts()
	if len(alerts) == 0 {
		return
	}
	for _, alert := range alerts {
		stock, _ := a.db.GetStockByID(alert.StockID)
		if stock == nil {
			continue
		}
		exchanges, _ := a.db.GetExchanges()
		var suffix string
		for _, e := range exchanges {
			if e.ID == stock.ExchangeID {
				suffix = e.YahooSuffix
				break
			}
		}
		sym := strings.TrimSuffix(stock.Symbol, suffix)
		q, err := a.yahoo.FetchQuote(sym, suffix)
		if err != nil {
			continue
		}
		triggered := false
		if alert.Condition == "above" && q.Price >= alert.Price {
			triggered = true
		}
		if alert.Condition == "below" && q.Price <= alert.Price {
			triggered = true
		}
		if triggered {
			a.db.SetAlertTriggered(alert.ID, true)
			if a.onAlert != nil {
				a.onAlert(alert, q.Price)
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
}
