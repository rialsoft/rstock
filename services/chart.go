package services

import (
	"rstock/models"
	"strings"
	"time"
)

type ChartService struct {
	db    *DatabaseService
	yahoo *YahooFinanceService
}

func NewChartService(db *DatabaseService, yahoo *YahooFinanceService) *ChartService {
	return &ChartService{db: db, yahoo: yahoo}
}

func (c *ChartService) GetChartData(stockID int64, period string) *models.ChartData {
	from, to := dateRange(period)
	history, err := c.db.GetPriceHistory(stockID, from, to)
	if err == nil && len(history) > 0 {
		return &models.ChartData{Period: period, Data: history}
	}
	stock, _ := c.db.GetStockByID(stockID)
	if stock != nil && c.yahoo != nil {
		exchanges, _ := c.db.GetExchanges()
		for _, e := range exchanges {
			if e.ID == stock.ExchangeID {
				yahooRange := period
				if period == "1w" { yahooRange = "5d" }
				sym := strings.TrimSuffix(stock.Symbol, e.YahooSuffix)
				h, err := c.yahoo.FetchHistory(sym, e.YahooSuffix, yahooRange)
				if err == nil && len(h) > 0 {
					c.db.InsertPriceHistory(stockID, h)
					return &models.ChartData{Period: period, Data: h}
				}
				break
			}
		}
	}
	return &models.ChartData{Period: period, Data: []models.OHLCV{}}
}

func dateRange(period string) (string, string) {
	now := time.Now()
	to := now.Format("2006-01-02")
	switch strings.ToLower(period) {
	case "1d":
		return now.AddDate(0, 0, -1).Format("2006-01-02"), to
	case "1w":
		return now.AddDate(0, 0, -7).Format("2006-01-02"), to
	case "1mo":
		return now.AddDate(0, -1, 0).Format("2006-01-02"), to
	case "3mo":
		return now.AddDate(0, -3, 0).Format("2006-01-02"), to
	case "6mo":
		return now.AddDate(0, -6, 0).Format("2006-01-02"), to
	case "1y":
		return now.AddDate(-1, 0, 0).Format("2006-01-02"), to
	case "5y":
		return now.AddDate(-5, 0, 0).Format("2006-01-02"), to
	case "max":
		return "2000-01-01", to
	default:
		return now.AddDate(0, -1, 0).Format("2006-01-02"), to
	}
}
