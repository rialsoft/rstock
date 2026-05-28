package services

import "rstock/models"

type MarketService struct {
	db *DatabaseService
}

func NewMarketService(db *DatabaseService) *MarketService {
	return &MarketService{db: db}
}

func (m *MarketService) GetExchanges() []models.Exchange {
	exchanges, _ := m.db.GetExchanges()
	return exchanges
}

func (m *MarketService) SeedDefaultExchanges() {
	defaults := []models.Exchange{
		{Code: "NYSE", Name: "New York Stock Exchange", Country: "US", Currency: "USD", Timezone: "America/New_York", YahooSuffix: "", MarketOpen: "09:30", MarketClose: "16:00"},
		{Code: "NASDAQ", Name: "NASDAQ", Country: "US", Currency: "USD", Timezone: "America/New_York", YahooSuffix: "", MarketOpen: "09:30", MarketClose: "16:00"},
		{Code: "JKSE", Name: "Jakarta Stock Exchange", Country: "ID", Currency: "IDR", Timezone: "Asia/Jakarta", YahooSuffix: ".JK", MarketOpen: "09:00", MarketClose: "16:00"},
		{Code: "HKEX", Name: "Hong Kong Stock Exchange", Country: "HK", Currency: "HKD", Timezone: "Asia/Hong_Kong", YahooSuffix: ".HK", MarketOpen: "09:30", MarketClose: "16:00"},
		{Code: "TSE", Name: "Tokyo Stock Exchange", Country: "JP", Currency: "JPY", Timezone: "Asia/Tokyo", YahooSuffix: ".T", MarketOpen: "09:00", MarketClose: "15:00"},
		{Code: "LSE", Name: "London Stock Exchange", Country: "GB", Currency: "GBP", Timezone: "Europe/London", YahooSuffix: ".L", MarketOpen: "08:00", MarketClose: "16:30"},
		{Code: "Euronext", Name: "Euronext", Country: "NL", Currency: "EUR", Timezone: "Europe/Amsterdam", YahooSuffix: ".AS", MarketOpen: "09:00", MarketClose: "17:30"},
		{Code: "XETRA", Name: "Deutsche Borse XETRA", Country: "DE", Currency: "EUR", Timezone: "Europe/Berlin", YahooSuffix: ".DE", MarketOpen: "09:00", MarketClose: "17:30"},
		{Code: "SIX", Name: "SIX Swiss Exchange", Country: "CH", Currency: "CHF", Timezone: "Europe/Zurich", YahooSuffix: ".SW", MarketOpen: "09:00", MarketClose: "17:30"},
		{Code: "ASX", Name: "Australian Securities Exchange", Country: "AU", Currency: "AUD", Timezone: "Australia/Sydney", YahooSuffix: ".AX", MarketOpen: "10:00", MarketClose: "16:00"},
		{Code: "SGX", Name: "Singapore Exchange", Country: "SG", Currency: "SGD", Timezone: "Asia/Singapore", YahooSuffix: ".SI", MarketOpen: "09:00", MarketClose: "17:00"},
		{Code: "NSE", Name: "National Stock Exchange of India", Country: "IN", Currency: "INR", Timezone: "Asia/Kolkata", YahooSuffix: ".NS", MarketOpen: "09:15", MarketClose: "15:30"},
		{Code: "KRX", Name: "Korea Exchange", Country: "KR", Currency: "KRW", Timezone: "Asia/Seoul", YahooSuffix: ".KS", MarketOpen: "09:00", MarketClose: "15:30"},
		{Code: "TWSE", Name: "Taiwan Stock Exchange", Country: "TW", Currency: "TWD", Timezone: "Asia/Taipei", YahooSuffix: ".TW", MarketOpen: "09:00", MarketClose: "13:30"},
		{Code: "SSE", Name: "Shanghai Stock Exchange", Country: "CN", Currency: "CNY", Timezone: "Asia/Shanghai", YahooSuffix: ".SS", MarketOpen: "09:30", MarketClose: "15:00"},
		{Code: "SZSE", Name: "Shenzhen Stock Exchange", Country: "CN", Currency: "CNY", Timezone: "Asia/Shanghai", YahooSuffix: ".SZ", MarketOpen: "09:30", MarketClose: "15:00"},
		{Code: "TSX", Name: "Toronto Stock Exchange", Country: "CA", Currency: "CAD", Timezone: "America/Toronto", YahooSuffix: ".TO", MarketOpen: "09:30", MarketClose: "16:00"},
		{Code: "B3", Name: "Brasil Bolsa Balcao", Country: "BR", Currency: "BRL", Timezone: "America/Sao_Paulo", YahooSuffix: ".SA", MarketOpen: "10:00", MarketClose: "17:00"},
		{Code: "JSE", Name: "Johannesburg Stock Exchange", Country: "ZA", Currency: "ZAR", Timezone: "Africa/Johannesburg", YahooSuffix: ".JO", MarketOpen: "09:00", MarketClose: "17:00"},
		{Code: "OMX", Name: "Nasdaq Stockholm", Country: "SE", Currency: "SEK", Timezone: "Europe/Stockholm", YahooSuffix: ".ST", MarketOpen: "09:00", MarketClose: "17:30"},
	}
	for _, e := range defaults {
		m.db.InsertExchange(e)
	}
}
