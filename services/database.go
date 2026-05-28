package services

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"rstock/models"

	_ "github.com/mattn/go-sqlite3"
)

type DatabaseService struct {
	DB *sql.DB
}

func NewDatabaseService(dbPath string) (*DatabaseService, error) {
	dir := filepath.Dir(dbPath)
	os.MkdirAll(dir, 0755)
	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_foreign_keys=on")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	return &DatabaseService{DB: db}, nil
}

func (d *DatabaseService) Close() error { return d.DB.Close() }

func (d *DatabaseService) Migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS exchanges (id INTEGER PRIMARY KEY AUTOINCREMENT, code TEXT NOT NULL UNIQUE, name TEXT NOT NULL, country TEXT NOT NULL, currency TEXT NOT NULL, timezone TEXT NOT NULL, yahoo_suffix TEXT NOT NULL, market_open TEXT, market_close TEXT, data_provider TEXT DEFAULT 'yahoo');
	CREATE TABLE IF NOT EXISTS stocks (id INTEGER PRIMARY KEY AUTOINCREMENT, symbol TEXT NOT NULL, name TEXT NOT NULL, exchange_id INTEGER NOT NULL REFERENCES exchanges(id), watchlist INTEGER DEFAULT 1, created_at TEXT DEFAULT (datetime('now')), UNIQUE(symbol, exchange_id));
	CREATE TABLE IF NOT EXISTS quotes (id INTEGER PRIMARY KEY AUTOINCREMENT, stock_id INTEGER NOT NULL REFERENCES stocks(id), price REAL NOT NULL, prev_close REAL, open REAL, high REAL, low REAL, volume INTEGER, change_pct REAL, fetched_at TEXT DEFAULT (datetime('now')));
	CREATE INDEX IF NOT EXISTS idx_quotes_stock_time ON quotes(stock_id, fetched_at);
	CREATE TABLE IF NOT EXISTS price_history (id INTEGER PRIMARY KEY AUTOINCREMENT, stock_id INTEGER NOT NULL REFERENCES stocks(id), date TEXT NOT NULL, open REAL NOT NULL, high REAL NOT NULL, low REAL NOT NULL, close REAL NOT NULL, volume INTEGER NOT NULL, UNIQUE(stock_id, date));
	CREATE INDEX IF NOT EXISTS idx_history_stock_date ON price_history(stock_id, date);
	CREATE TABLE IF NOT EXISTS app_settings (key TEXT PRIMARY KEY, value TEXT NOT NULL);
	CREATE TABLE IF NOT EXISTS portfolios (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL, base_currency TEXT NOT NULL DEFAULT 'IDR', initial_cash REAL NOT NULL DEFAULT 0, created_at TEXT DEFAULT (datetime('now')));
	CREATE TABLE IF NOT EXISTS transactions (id INTEGER PRIMARY KEY AUTOINCREMENT, portfolio_id INTEGER NOT NULL REFERENCES portfolios(id), stock_id INTEGER NOT NULL REFERENCES stocks(id), type TEXT NOT NULL CHECK(type IN ('buy','sell')), quantity INTEGER NOT NULL, price REAL NOT NULL, fees REAL DEFAULT 0, date TEXT NOT NULL, notes TEXT, created_at TEXT DEFAULT (datetime('now')));
	CREATE INDEX IF NOT EXISTS idx_tx_portfolio ON transactions(portfolio_id, stock_id);
	CREATE TABLE IF NOT EXISTS alerts (id INTEGER PRIMARY KEY AUTOINCREMENT, stock_id INTEGER NOT NULL REFERENCES stocks(id), condition TEXT NOT NULL CHECK(condition IN ('above','below')), price REAL NOT NULL, enabled INTEGER DEFAULT 1, triggered INTEGER DEFAULT 0, created_at TEXT DEFAULT (datetime('now')));
	`
	_, err := d.DB.Exec(schema)
	return err
}

func (d *DatabaseService) InsertExchange(e models.Exchange) (int64, error) {
	res, err := d.DB.Exec(`INSERT INTO exchanges (code,name,country,currency,timezone,yahoo_suffix,market_open,market_close,data_provider) VALUES (?,?,?,?,?,?,?,?,?)`, e.Code, e.Name, e.Country, e.Currency, e.Timezone, e.YahooSuffix, e.MarketOpen, e.MarketClose, e.DataProvider)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (d *DatabaseService) GetExchanges() ([]models.Exchange, error) {
	rows, err := d.DB.Query(`SELECT id,code,name,country,currency,timezone,yahoo_suffix,market_open,market_close,data_provider FROM exchanges ORDER BY code`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.Exchange
	for rows.Next() {
		var e models.Exchange
		rows.Scan(&e.ID, &e.Code, &e.Name, &e.Country, &e.Currency, &e.Timezone, &e.YahooSuffix, &e.MarketOpen, &e.MarketClose, &e.DataProvider)
		result = append(result, e)
	}
	return result, nil
}

func (d *DatabaseService) InsertStock(s models.Stock) (int64, error) {
	res, err := d.DB.Exec(`INSERT INTO stocks (symbol,name,exchange_id,watchlist) VALUES (?,?,?,1) ON CONFLICT(symbol, exchange_id) DO UPDATE SET watchlist=1, name=excluded.name`, s.Symbol, s.Name, s.ExchangeID)
	if err != nil { return 0, err }
	id, _ := res.LastInsertId()
	// ON CONFLICT on existing row returns 0 for LastInsertId; fetch the existing id
	if id == 0 {
		var existingID int64
		d.DB.QueryRow(`SELECT id FROM stocks WHERE symbol=? AND exchange_id=?`, s.Symbol, s.ExchangeID).Scan(&existingID)
		return existingID, nil
	}
	return id, nil
}

func (d *DatabaseService) GetWatchlist(exchangeID int64) ([]models.StockWithQuote, error) {
	rows, err := d.DB.Query(`SELECT s.id,s.symbol,s.name,s.exchange_id,s.watchlist,s.created_at,q.id,q.stock_id,q.price,q.prev_close,q.open,q.high,q.low,q.volume,q.change_pct,q.fetched_at FROM stocks s LEFT JOIN quotes q ON q.id=(SELECT id FROM quotes WHERE stock_id=s.id ORDER BY fetched_at DESC LIMIT 1) WHERE s.exchange_id=? AND s.watchlist=1 ORDER BY s.symbol`, exchangeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []models.StockWithQuote
	for rows.Next() {
		var swq models.StockWithQuote
		var qID, qSID sql.NullInt64
		var price, prevClose, open, high, low sql.NullFloat64
		var volume sql.NullInt64
		var changePct sql.NullFloat64
		var fetchedAt sql.NullString
		rows.Scan(&swq.Stock.ID, &swq.Stock.Symbol, &swq.Stock.Name, &swq.Stock.ExchangeID, &swq.Stock.Watchlist, &swq.Stock.CreatedAt, &qID, &qSID, &price, &prevClose, &open, &high, &low, &volume, &changePct, &fetchedAt)
		if qID.Valid {
			swq.Quote = &models.Quote{ID: qID.Int64, StockID: qSID.Int64, Price: price.Float64, PrevClose: prevClose.Float64, Open: open.Float64, High: high.Float64, Low: low.Float64, Volume: volume.Int64, ChangePct: changePct.Float64, FetchedAt: fetchedAt.String}
		}
		results = append(results, swq)
	}
	return results, nil
}

func (d *DatabaseService) UpdateWatchlistStatus(stockID int64, active bool) error {
	var v int
	if active {
		v = 1
	}
	_, err := d.DB.Exec(`UPDATE stocks SET watchlist=? WHERE id=?`, v, stockID)
	return err
}

func (d *DatabaseService) InsertQuote(q models.Quote) error {
	_, err := d.DB.Exec(`INSERT INTO quotes (stock_id,price,prev_close,open,high,low,volume,change_pct) VALUES (?,?,?,?,?,?,?,?)`, q.StockID, q.Price, q.PrevClose, q.Open, q.High, q.Low, q.Volume, q.ChangePct)
	return err
}

func (d *DatabaseService) GetLatestQuote(stockID int64) (*models.Quote, error) {
	q := &models.Quote{}
	err := d.DB.QueryRow(`SELECT id,stock_id,price,prev_close,open,high,low,volume,change_pct,fetched_at FROM quotes WHERE stock_id=? ORDER BY fetched_at DESC LIMIT 1`, stockID).Scan(&q.ID, &q.StockID, &q.Price, &q.PrevClose, &q.Open, &q.High, &q.Low, &q.Volume, &q.ChangePct, &q.FetchedAt)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (d *DatabaseService) InsertPriceHistory(stockID int64, history []models.OHLCV) error {
	tx, _ := d.DB.Begin()
	defer tx.Rollback()
	stmt, _ := tx.Prepare(`INSERT OR REPLACE INTO price_history (stock_id,date,open,high,low,close,volume) VALUES (?,?,?,?,?,?,?)`)
	defer stmt.Close()
	for _, h := range history {
		stmt.Exec(stockID, h.Date, h.Open, h.High, h.Low, h.Close, h.Volume)
	}
	return tx.Commit()
}

func (d *DatabaseService) GetPriceHistory(stockID int64, from, to string) ([]models.OHLCV, error) {
	rows, err := d.DB.Query(`SELECT date,open,high,low,close,volume FROM price_history WHERE stock_id=? AND date>=? AND date<=? ORDER BY date ASC`, stockID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.OHLCV
	for rows.Next() {
		var h models.OHLCV
		rows.Scan(&h.Date, &h.Open, &h.High, &h.Low, &h.Close, &h.Volume)
		result = append(result, h)
	}
	return result, nil
}

func (d *DatabaseService) SetSetting(key, value string) error {
	_, err := d.DB.Exec(`INSERT OR REPLACE INTO app_settings (key,value) VALUES (?,?)`, key, value)
	return err
}

func (d *DatabaseService) GetSetting(key string) (string, error) {
	var v string
	err := d.DB.QueryRow(`SELECT value FROM app_settings WHERE key=?`, key).Scan(&v)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("not found: %s", key)
	}
	return v, err
}

func (d *DatabaseService) GetStockByID(id int64) (*models.Stock, error) {
	s := &models.Stock{}
	err := d.DB.QueryRow(`SELECT id,symbol,name,exchange_id,watchlist,created_at FROM stocks WHERE id=?`, id).Scan(&s.ID, &s.Symbol, &s.Name, &s.ExchangeID, &s.Watchlist, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (d *DatabaseService) CleanupOldQuotes() error {
	_, err := d.DB.Exec(`DELETE FROM quotes WHERE fetched_at < datetime('now','-30 days')`)
	return err
}

func (d *DatabaseService) InsertPortfolio(p models.Portfolio) (int64, error) {
	res, err := d.DB.Exec(`INSERT INTO portfolios (name,base_currency,initial_cash) VALUES (?,?,?)`, p.Name, p.BaseCurrency, p.InitialCash)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (d *DatabaseService) GetPortfolios() ([]models.Portfolio, error) {
	rows, err := d.DB.Query(`SELECT id,name,base_currency,initial_cash,created_at FROM portfolios ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.Portfolio
	for rows.Next() {
		var p models.Portfolio
		rows.Scan(&p.ID, &p.Name, &p.BaseCurrency, &p.InitialCash, &p.CreatedAt)
		result = append(result, p)
	}
	return result, nil
}

func (d *DatabaseService) GetPortfolioByID(id int64) (*models.Portfolio, error) {
	p := &models.Portfolio{}
	err := d.DB.QueryRow(`SELECT id,name,base_currency,initial_cash,created_at FROM portfolios WHERE id=?`, id).Scan(&p.ID, &p.Name, &p.BaseCurrency, &p.InitialCash, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (d *DatabaseService) DeletePortfolio(id int64) error {
	tx, _ := d.DB.Begin()
	defer tx.Rollback()
	tx.Exec(`DELETE FROM transactions WHERE portfolio_id=?`, id)
	tx.Exec(`DELETE FROM portfolios WHERE id=?`, id)
	return tx.Commit()
}

func (d *DatabaseService) InsertTransaction(t models.Transaction) (int64, error) {
	res, err := d.DB.Exec(`INSERT INTO transactions (portfolio_id,stock_id,type,quantity,price,fees,date,notes) VALUES (?,?,?,?,?,?,?,?)`, t.PortfolioID, t.StockID, t.Type, t.Quantity, t.Price, t.Fees, t.Date, t.Notes)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (d *DatabaseService) GetTransactions(portfolioID int64) ([]models.Transaction, error) {
	rows, err := d.DB.Query(`SELECT id,portfolio_id,stock_id,type,quantity,price,fees,date,notes,created_at FROM transactions WHERE portfolio_id=? ORDER BY date DESC, id DESC`, portfolioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.Transaction
	for rows.Next() {
		var t models.Transaction
		var notes sql.NullString
		rows.Scan(&t.ID, &t.PortfolioID, &t.StockID, &t.Type, &t.Quantity, &t.Price, &t.Fees, &t.Date, &notes, &t.CreatedAt)
		if notes.Valid {
			t.Notes = notes.String
		}
		result = append(result, t)
	}
	return result, nil
}

func (d *DatabaseService) GetHoldings(portfolioID int64) ([]models.Holding, error) {
	rows, err := d.DB.Query(`SELECT t.stock_id,s.symbol,s.name,SUM(CASE WHEN t.type='buy' THEN t.quantity ELSE -t.quantity END),SUM(CASE WHEN t.type='buy' THEN t.quantity*t.price+t.fees ELSE -t.quantity*t.price-t.fees END),CASE WHEN SUM(CASE WHEN t.type='buy' THEN t.quantity ELSE -t.quantity END)>0 THEN SUM(CASE WHEN t.type='buy' THEN t.quantity*t.price+t.fees ELSE -t.quantity*t.price-t.fees END)/SUM(CASE WHEN t.type='buy' THEN t.quantity ELSE -t.quantity END) ELSE 0 END FROM transactions t JOIN stocks s ON t.stock_id=s.id WHERE t.portfolio_id=? GROUP BY t.stock_id HAVING SUM(CASE WHEN t.type='buy' THEN t.quantity ELSE -t.quantity END)>0 ORDER BY s.symbol`, portfolioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.Holding
	for rows.Next() {
		var h models.Holding
		rows.Scan(&h.StockID, &h.Symbol, &h.Name, &h.TotalQty, &h.CostBasis, &h.AvgCost)
		result = append(result, h)
	}
	return result, nil
}

func (d *DatabaseService) GetRealizedGL(portfolioID int64) (float64, error) {
	rows, err := d.DB.Query(`SELECT t_sell.stock_id,t_sell.quantity,t_sell.price,t_sell.fees,COALESCE((SELECT SUM(CASE WHEN type='buy' THEN quantity*price+fees ELSE 0 END)/SUM(CASE WHEN type='buy' THEN quantity ELSE 0 END) FROM transactions WHERE portfolio_id=? AND stock_id=t_sell.stock_id AND type='buy' AND date<=t_sell.date),0) FROM transactions t_sell WHERE t_sell.portfolio_id=? AND t_sell.type='sell'`, portfolioID, portfolioID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var totalGL float64
	for rows.Next() {
		var stockID, qty int64
		var price, fees, avgBuy float64
		rows.Scan(&stockID, &qty, &price, &fees, &avgBuy)
		totalGL += (price-avgBuy)*float64(qty) - fees
	}
	return totalGL, nil
}

func (d *DatabaseService) InsertAlert(a models.Alert) (int64, error) {
	res, err := d.DB.Exec(`INSERT INTO alerts (stock_id,condition,price) VALUES (?,?,?)`, a.StockID, a.Condition, a.Price)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (d *DatabaseService) GetAlerts() ([]models.Alert, error) {
	rows, err := d.DB.Query(`SELECT a.id,a.stock_id,s.symbol,s.name,a.condition,a.price,a.enabled,a.triggered,a.created_at FROM alerts a JOIN stocks s ON a.stock_id=s.id ORDER BY s.symbol`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var alerts []models.Alert
	for rows.Next() {
		var a models.Alert
		var enabled, triggered int
		rows.Scan(&a.ID, &a.StockID, &a.Symbol, &a.Name, &a.Condition, &a.Price, &enabled, &triggered, &a.CreatedAt)
		a.Enabled = enabled != 0
		a.Triggered = triggered != 0
		alerts = append(alerts, a)
	}
	return alerts, nil
}

func (d *DatabaseService) GetEnabledAlerts() ([]models.Alert, error) {
	rows, err := d.DB.Query(`SELECT a.id,a.stock_id,s.symbol,s.name,a.condition,a.price,a.enabled,a.triggered,a.created_at FROM alerts a JOIN stocks s ON a.stock_id=s.id WHERE a.enabled=1 AND a.triggered=0 ORDER BY s.symbol`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var alerts []models.Alert
	for rows.Next() {
		var a models.Alert
		var enabled, triggered int
		rows.Scan(&a.ID, &a.StockID, &a.Symbol, &a.Name, &a.Condition, &a.Price, &enabled, &triggered, &a.CreatedAt)
		a.Enabled = enabled != 0
		a.Triggered = triggered != 0
		alerts = append(alerts, a)
	}
	return alerts, nil
}

func (d *DatabaseService) ToggleAlert(id int64, enabled bool) error {
	v := 0
	if enabled {
		v = 1
	}
	_, err := d.DB.Exec(`UPDATE alerts SET enabled=? WHERE id=?`, v, id)
	return err
}

func (d *DatabaseService) SetAlertTriggered(id int64, triggered bool) error {
	v := 0
	if triggered {
		v = 1
	}
	_, err := d.DB.Exec(`UPDATE alerts SET triggered=? WHERE id=?`, v, id)
	return err
}

func (d *DatabaseService) DeleteAlert(id int64) error {
	_, err := d.DB.Exec(`DELETE FROM alerts WHERE id=?`, id)
	return err
}
