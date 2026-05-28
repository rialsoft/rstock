package models

type Exchange struct {
	ID           int64  `json:"id"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	Country      string `json:"country"`
	Currency     string `json:"currency"`
	Timezone     string `json:"timezone"`
	YahooSuffix  string `json:"yahooSuffix"`
	MarketOpen   string `json:"marketOpen"`
	MarketClose  string `json:"marketClose"`
	DataProvider string `json:"dataProvider"`
}

type Stock struct {
	ID         int64 `json:"id"`
	Symbol     string `json:"symbol"`
	Name       string `json:"name"`
	ExchangeID int64  `json:"exchangeId"`
	Watchlist  int    `json:"watchlist"`
	CreatedAt  string `json:"createdAt"`
}

type Quote struct {
	ID        int64   `json:"id"`
	Symbol    string  `json:"symbol"`
	StockID   int64   `json:"stockId"`
	Price     float64 `json:"price"`
	PrevClose float64 `json:"prevClose"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Volume    int64   `json:"volume"`
	ChangePct float64 `json:"changePct"`
	FetchedAt string  `json:"fetchedAt"`
}

type OHLCV struct {
	Date   string  `json:"date"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume int64   `json:"volume"`
}

type StockWithQuote struct {
	Stock Stock  `json:"stock"`
	Quote *Quote `json:"quote,omitempty"`
}

type ChartData struct {
	Period string  `json:"period"`
	Data   []OHLCV `json:"data"`
}

type SymbolSuggestion struct {
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Exchange string `json:"exchange"`
	Type     string `json:"type"`
}

type Portfolio struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	BaseCurrency string  `json:"baseCurrency"`
	InitialCash  float64 `json:"initialCash"`
	CreatedAt    string  `json:"createdAt"`
}

type Transaction struct {
	ID          int64   `json:"id"`
	PortfolioID int64   `json:"portfolioId"`
	StockID     int64   `json:"stockId"`
	Type        string  `json:"type"`
	Quantity    int64   `json:"quantity"`
	Price       float64 `json:"price"`
	Fees        float64 `json:"fees"`
	Date        string  `json:"date"`
	Notes       string  `json:"notes"`
	CreatedAt   string  `json:"createdAt"`
}

type Holding struct {
	StockID     int64   `json:"stockId"`
	Symbol      string  `json:"symbol"`
	Name        string  `json:"name"`
	TotalQty    int64   `json:"totalQty"`
	AvgCost     float64 `json:"avgCost"`
	CostBasis   float64 `json:"costBasis"`
	MarketPrice float64 `json:"marketPrice"`
	MarketValue float64 `json:"marketValue"`
	GainLoss    float64 `json:"gainLoss"`
	GainLossPct float64 `json:"gainLossPct"`
}

type PortfolioSummary struct {
	Portfolio   Portfolio `json:"portfolio"`
	TotalCost   float64   `json:"totalCost"`
	MarketValue float64   `json:"marketValue"`
	TotalGL     float64   `json:"totalGL"`
	TotalGLPct  float64   `json:"totalGLPct"`
	Holdings    []Holding `json:"holdings"`
	RealizedGL  float64   `json:"realizedGL"`
}

type Alert struct {
	ID        int64   `json:"id"`
	StockID   int64   `json:"stockId"`
	Symbol    string  `json:"symbol"`
	Name      string  `json:"name"`
	Condition string  `json:"condition"`
	Price     float64 `json:"price"`
	Enabled   bool    `json:"enabled"`
	Triggered bool    `json:"triggered"`
	CreatedAt string  `json:"createdAt"`
}
