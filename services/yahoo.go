package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"rstock/models"
	"time"
)

type YahooFinanceService struct {
	client *http.Client
}

func NewYahooFinanceService() *YahooFinanceService {
	return &YahooFinanceService{client: &http.Client{Timeout: 10 * time.Second}}
}

func (y *YahooFinanceService) doGet(rawURL string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", rawURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	return y.client.Do(req)
}

func (y *YahooFinanceService) FetchQuote(symbol, suffix string) (*models.Quote, error) {
	resp, err := y.doGet(fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s%s?interval=1d&range=1d", symbol, suffix))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 429 {
		return nil, fmt.Errorf("rate limited")
	}
	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("symbol not found: %s", symbol)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("yahoo returned %d", resp.StatusCode)
	}
	var cr struct {
		Chart struct {
			Result []struct {
				Meta struct {
					RegularMarketPrice float64 `json:"regularMarketPrice"`
					PreviousClose      float64 `json:"previousClose"`
					Symbol             string  `json:"symbol"`
				} `json:"meta"`
				Indicators struct {
					Quote []struct {
						Open   []float64 `json:"open"`
						High   []float64 `json:"high"`
						Low    []float64 `json:"low"`
						Volume []int64   `json:"volume"`
					} `json:"quote"`
				} `json:"indicators"`
			} `json:"result"`
			Error interface{} `json:"error"`
		} `json:"chart"`
	}
	json.NewDecoder(resp.Body).Decode(&cr)
	if cr.Chart.Error != nil || len(cr.Chart.Result) == 0 {
		return nil, fmt.Errorf("no data for: %s", symbol)
	}
	r := cr.Chart.Result[0]
	q := &models.Quote{Symbol: symbol, Price: r.Meta.RegularMarketPrice, PrevClose: r.Meta.PreviousClose}
	if r.Meta.PreviousClose > 0 && r.Meta.RegularMarketPrice > 0 {
		q.ChangePct = ((r.Meta.RegularMarketPrice - r.Meta.PreviousClose) / r.Meta.PreviousClose) * 100
	}
	if len(r.Indicators.Quote) > 0 {
		ind := r.Indicators.Quote[0]
		if len(ind.Open) > 0 {
			q.Open = ind.Open[0]
		}
		if len(ind.High) > 0 {
			q.High = ind.High[0]
		}
		if len(ind.Low) > 0 {
			q.Low = ind.Low[0]
		}
		if len(ind.Volume) > 0 {
			q.Volume = ind.Volume[0]
		}
	}
	return q, nil
}

func (y *YahooFinanceService) FetchHistory(symbol, suffix, period string) ([]models.OHLCV, error) {
	resp, err := y.doGet(fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s%s?interval=1d&range=%s", symbol, suffix, period))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("yahoo returned %d", resp.StatusCode)
	}
	var cr struct {
		Chart struct {
			Result []struct {
				Timestamp  []int64 `json:"timestamp"`
				Indicators struct {
					Quote []struct {
						Open   []float64 `json:"open"`
						High   []float64 `json:"high"`
						Low    []float64 `json:"low"`
						Close  []float64 `json:"close"`
						Volume []int64   `json:"volume"`
					} `json:"quote"`
				} `json:"indicators"`
			} `json:"result"`
			Error interface{} `json:"error"`
		} `json:"chart"`
	}
	json.NewDecoder(resp.Body).Decode(&cr)
	if cr.Chart.Error != nil || len(cr.Chart.Result) == 0 {
		return nil, fmt.Errorf("no history data")
	}
	r := cr.Chart.Result[0]
	if len(r.Timestamp) == 0 || len(r.Indicators.Quote) == 0 {
		return nil, fmt.Errorf("empty data")
	}
	ind := r.Indicators.Quote[0]
	var history []models.OHLCV
	for i := 0; i < len(r.Timestamp); i++ {
		h := models.OHLCV{Date: time.Unix(r.Timestamp[i], 0).Format("2006-01-02")}
		if len(ind.Open) > i {
			h.Open = ind.Open[i]
		}
		if len(ind.High) > i {
			h.High = ind.High[i]
		}
		if len(ind.Low) > i {
			h.Low = ind.Low[i]
		}
		if len(ind.Close) > i {
			h.Close = ind.Close[i]
		}
		if len(ind.Volume) > i {
			h.Volume = ind.Volume[i]
		}
		// Skip incomplete current-day data where close is zero/null
		if h.Close <= 0 {
			continue
		}
		history = append(history, h)
	}
	return history, nil
}

func (y *YahooFinanceService) SearchSymbol(query string) ([]models.SymbolSuggestion, error) {
	resp, err := y.doGet(fmt.Sprintf("https://query2.finance.yahoo.com/v1/finance/search?q=%s&quotesCount=10&newsCount=0", url.QueryEscape(query)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result struct {
		Quotes []struct {
			Symbol    string `json:"symbol"`
			Shortname string `json:"shortname"`
			Exchange  string `json:"exchange"`
			QuoteType string `json:"quoteType"`
		} `json:"quotes"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	var suggestions []models.SymbolSuggestion
	for _, q := range result.Quotes {
		suggestions = append(suggestions, models.SymbolSuggestion{Symbol: q.Symbol, Name: q.Shortname, Exchange: q.Exchange, Type: q.QuoteType})
	}
	return suggestions, nil
}
