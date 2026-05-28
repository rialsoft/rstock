export namespace models {
	
	export class Alert {
	    id: number;
	    stockId: number;
	    symbol: string;
	    name: string;
	    condition: string;
	    price: number;
	    enabled: boolean;
	    triggered: boolean;
	    createdAt: string;
	
	    static createFrom(source: any = {}) {
	        return new Alert(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.stockId = source["stockId"];
	        this.symbol = source["symbol"];
	        this.name = source["name"];
	        this.condition = source["condition"];
	        this.price = source["price"];
	        this.enabled = source["enabled"];
	        this.triggered = source["triggered"];
	        this.createdAt = source["createdAt"];
	    }
	}
	export class OHLCV {
	    date: string;
	    open: number;
	    high: number;
	    low: number;
	    close: number;
	    volume: number;
	
	    static createFrom(source: any = {}) {
	        return new OHLCV(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.date = source["date"];
	        this.open = source["open"];
	        this.high = source["high"];
	        this.low = source["low"];
	        this.close = source["close"];
	        this.volume = source["volume"];
	    }
	}
	export class ChartData {
	    period: string;
	    data: OHLCV[];
	
	    static createFrom(source: any = {}) {
	        return new ChartData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.period = source["period"];
	        this.data = this.convertValues(source["data"], OHLCV);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Exchange {
	    id: number;
	    code: string;
	    name: string;
	    country: string;
	    currency: string;
	    timezone: string;
	    yahooSuffix: string;
	    marketOpen: string;
	    marketClose: string;
	    dataProvider: string;
	
	    static createFrom(source: any = {}) {
	        return new Exchange(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.code = source["code"];
	        this.name = source["name"];
	        this.country = source["country"];
	        this.currency = source["currency"];
	        this.timezone = source["timezone"];
	        this.yahooSuffix = source["yahooSuffix"];
	        this.marketOpen = source["marketOpen"];
	        this.marketClose = source["marketClose"];
	        this.dataProvider = source["dataProvider"];
	    }
	}
	export class Holding {
	    stockId: number;
	    symbol: string;
	    name: string;
	    totalQty: number;
	    avgCost: number;
	    costBasis: number;
	    marketPrice: number;
	    marketValue: number;
	    gainLoss: number;
	    gainLossPct: number;
	
	    static createFrom(source: any = {}) {
	        return new Holding(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.stockId = source["stockId"];
	        this.symbol = source["symbol"];
	        this.name = source["name"];
	        this.totalQty = source["totalQty"];
	        this.avgCost = source["avgCost"];
	        this.costBasis = source["costBasis"];
	        this.marketPrice = source["marketPrice"];
	        this.marketValue = source["marketValue"];
	        this.gainLoss = source["gainLoss"];
	        this.gainLossPct = source["gainLossPct"];
	    }
	}
	
	export class Portfolio {
	    id: number;
	    name: string;
	    baseCurrency: string;
	    initialCash: number;
	    createdAt: string;
	
	    static createFrom(source: any = {}) {
	        return new Portfolio(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.baseCurrency = source["baseCurrency"];
	        this.initialCash = source["initialCash"];
	        this.createdAt = source["createdAt"];
	    }
	}
	export class PortfolioSummary {
	    portfolio: Portfolio;
	    totalCost: number;
	    marketValue: number;
	    totalGL: number;
	    totalGLPct: number;
	    holdings: Holding[];
	    realizedGL: number;
	
	    static createFrom(source: any = {}) {
	        return new PortfolioSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.portfolio = this.convertValues(source["portfolio"], Portfolio);
	        this.totalCost = source["totalCost"];
	        this.marketValue = source["marketValue"];
	        this.totalGL = source["totalGL"];
	        this.totalGLPct = source["totalGLPct"];
	        this.holdings = this.convertValues(source["holdings"], Holding);
	        this.realizedGL = source["realizedGL"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Quote {
	    id: number;
	    symbol: string;
	    stockId: number;
	    price: number;
	    prevClose: number;
	    open: number;
	    high: number;
	    low: number;
	    volume: number;
	    changePct: number;
	    fetchedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new Quote(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.symbol = source["symbol"];
	        this.stockId = source["stockId"];
	        this.price = source["price"];
	        this.prevClose = source["prevClose"];
	        this.open = source["open"];
	        this.high = source["high"];
	        this.low = source["low"];
	        this.volume = source["volume"];
	        this.changePct = source["changePct"];
	        this.fetchedAt = source["fetchedAt"];
	    }
	}
	export class Stock {
	    id: number;
	    symbol: string;
	    name: string;
	    exchangeId: number;
	    watchlist: number;
	    createdAt: string;
	
	    static createFrom(source: any = {}) {
	        return new Stock(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.symbol = source["symbol"];
	        this.name = source["name"];
	        this.exchangeId = source["exchangeId"];
	        this.watchlist = source["watchlist"];
	        this.createdAt = source["createdAt"];
	    }
	}
	export class StockWithQuote {
	    stock: Stock;
	    quote?: Quote;
	
	    static createFrom(source: any = {}) {
	        return new StockWithQuote(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.stock = this.convertValues(source["stock"], Stock);
	        this.quote = this.convertValues(source["quote"], Quote);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SymbolSuggestion {
	    symbol: string;
	    name: string;
	    exchange: string;
	    type: string;
	
	    static createFrom(source: any = {}) {
	        return new SymbolSuggestion(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.symbol = source["symbol"];
	        this.name = source["name"];
	        this.exchange = source["exchange"];
	        this.type = source["type"];
	    }
	}
	export class Transaction {
	    id: number;
	    portfolioId: number;
	    stockId: number;
	    type: string;
	    quantity: number;
	    price: number;
	    fees: number;
	    date: string;
	    notes: string;
	    createdAt: string;
	
	    static createFrom(source: any = {}) {
	        return new Transaction(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.portfolioId = source["portfolioId"];
	        this.stockId = source["stockId"];
	        this.type = source["type"];
	        this.quantity = source["quantity"];
	        this.price = source["price"];
	        this.fees = source["fees"];
	        this.date = source["date"];
	        this.notes = source["notes"];
	        this.createdAt = source["createdAt"];
	    }
	}

}

