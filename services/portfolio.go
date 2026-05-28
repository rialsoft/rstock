package services

import (
	"fmt"
	"rstock/models"
)

type PortfolioService struct {
	db    *DatabaseService
	yahoo *YahooFinanceService
}

func NewPortfolioService(db *DatabaseService, yahoo *YahooFinanceService) *PortfolioService {
	return &PortfolioService{db: db, yahoo: yahoo}
}

func (p *PortfolioService) CreatePortfolio(name, currency string, initialCash float64) (*models.Portfolio, error) {
	id, err := p.db.InsertPortfolio(models.Portfolio{Name: name, BaseCurrency: currency, InitialCash: initialCash})
	if err != nil {
		return nil, err
	}
	return p.db.GetPortfolioByID(id)
}

func (p *PortfolioService) GetPortfolios() []models.Portfolio {
	result, _ := p.db.GetPortfolios()
	return result
}

func (p *PortfolioService) DeletePortfolio(id int64) error {
	return p.db.DeletePortfolio(id)
}

func (p *PortfolioService) AddTransaction(t models.Transaction) (int64, error) {
	if t.Type != "buy" && t.Type != "sell" {
		return 0, fmt.Errorf("invalid type")
	}
	if t.Quantity <= 0 || t.Price <= 0 {
		return 0, fmt.Errorf("invalid quantity/price")
	}
	stock, _ := p.db.GetStockByID(t.StockID)
	if stock == nil {
		return 0, fmt.Errorf("stock not found")
	}
	if t.Type == "sell" {
		holdings, _ := p.db.GetHoldings(t.PortfolioID)
		var available int64
		for _, h := range holdings {
			if h.StockID == t.StockID {
				available = h.TotalQty
			}
		}
		if t.Quantity > available {
			return 0, fmt.Errorf("insufficient shares: have %d", available)
		}
	}
	return p.db.InsertTransaction(t)
}

func (p *PortfolioService) GetTransactions(portfolioID int64) []models.Transaction {
	result, _ := p.db.GetTransactions(portfolioID)
	return result
}

func (p *PortfolioService) GetPortfolioSummary(portfolioID int64) *models.PortfolioSummary {
	portfolio, err := p.db.GetPortfolioByID(portfolioID)
	if err != nil {
		return nil
	}
	holdings, _ := p.db.GetHoldings(portfolioID)
	var totalCost, totalValue float64
	for i := range holdings {
		q, err := p.db.GetLatestQuote(holdings[i].StockID)
		if err == nil {
			holdings[i].MarketPrice = q.Price
			holdings[i].MarketValue = q.Price * float64(holdings[i].TotalQty)
			holdings[i].GainLoss = holdings[i].MarketValue - holdings[i].CostBasis
			if holdings[i].CostBasis > 0 {
				holdings[i].GainLossPct = (holdings[i].GainLoss / holdings[i].CostBasis) * 100
			}
		}
		totalCost += holdings[i].CostBasis
		totalValue += holdings[i].MarketValue
	}
	realizedGL, _ := p.db.GetRealizedGL(portfolioID)
	totalGL := totalValue - totalCost + realizedGL
	totalGLPct := 0.0
	if totalCost > 0 {
		totalGLPct = (totalGL / totalCost) * 100
	}
	return &models.PortfolioSummary{
		Portfolio:   *portfolio,
		TotalCost:   totalCost,
		MarketValue: totalValue,
		TotalGL:     totalGL,
		TotalGLPct:  totalGLPct,
		Holdings:    holdings,
		RealizedGL:  realizedGL,
	}
}
