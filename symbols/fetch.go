package symbols

type ExchangeItemsType struct {
	Id     string
	Name   string
	Number int
}

type Exchange2Items struct {
	Id           string              `json:"id"`
	Name         string              `json:"name"`
	ExchangeType string              `json:"exchangeType"`
	Country      string              `json:"country"`
	ItemTypes    []ExchangeItemsType `json:"-"`
	KPeriods     []string            `json:"kPeriods"`
}

func fetchExchanges() (ret []Exchange2Items, retE error) {
	if openCtpMarketHost != "" {
		return ctpFetchExchanges()
	}
	return sinaFetchExchanges()
}

func fetchExchangeSymbols(foundExchange *Exchange2Items) (retList []map[string]string, retE error) {
	if openCtpMarketHost != "" {
		return ctpFetchExchangeSymbols(foundExchange)
	}
	return sinaFetchExchangeSymbols(foundExchange)
}
