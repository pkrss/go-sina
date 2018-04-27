package symbols

import (
	"encoding/json"
	"errors"

	"github.com/pkrss/go-sina/kline"
	pkNet "github.com/pkrss/go-utils/net/gbk"
)

func ctpFetchExchanges() (ret []Exchange2Items, retE error) {

	url := "http://" + openCtpMarketHost + "/quote2/exchanges"

	resp, err := pkNet.HttpGetEx(url)

	if err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, errors.New("fetchExchanges server response is empty")
	}

	ret = make([]Exchange2Items, 0)
	err = json.Unmarshal(resp, &ret)

	if err != nil {
		return nil, err
	}

	validPeriods := kline.GetValidPeriods()

	for idx := range ret {

		exchange2Items := &ret[idx]

		exchange2Items.KPeriods = validPeriods
		exchange2Items.ExchangeType = ExchangeType_Futures
		exchange2Items.Country = Country_CN
	}

	return

}

func ctpFetchExchangeSymbols(foundExchange *Exchange2Items) (retList []map[string]string, retE error) {

	switch foundExchange.ExchangeType {
	case ExchangeType_Futures:
		return ctpFetchExchangeSymbols_Futures(foundExchange)
	}
	retE = errors.New("Not implements:" + foundExchange.ExchangeType)
	return
}

func ctpFetchExchangeSymbols_Futures(foundExchange *Exchange2Items) (retList []map[string]string, retE error) {
	exchangeId := foundExchange.Id

	if Exchange_CustomFutures == exchangeId {
		return
	}

	retList = make([]map[string]string, 0)

	url := "http://" + openCtpMarketHost + "/quote2/exchanges/" + exchangeId + "/symbols"

	resp, err := pkNet.HttpGetEx(url)

	if err != nil {
		retE = err
		return
	}

	if len(resp) == 0 {
		retE = errors.New("fetchExchangeSymbols server response is empty")
		return
	}

	retList = make([]map[string]string, 0)
	retE = json.Unmarshal(resp, &retList)
	return

	// symbol:"TA0",market:"czce",contract:"PTA",name:"PTA连续",trade:"5394",settlement:"5408",prevsettlement:"5416",open:"5420",high:"5440",low:"5374",close:"5400",bid:"5392",ask:"5394",bidvol:"200",askvol:"36",volume:"581910",position:"981548",currentvol:"0",ticktime:"15:00:00",tradedate:"2017-12-11",changepercent:"-0.0040620"
}
