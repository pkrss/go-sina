package symbols

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/pkrss/go-sina/kline"
	pkJson "github.com/pkrss/go-utils/json"
	pkNet "github.com/pkrss/go-utils/net/gbk"

	"github.com/robertkrimen/otto"
)

var reg1 *regexp.Regexp
var vm *otto.Otto

func sinaFetchExchanges() (ret []Exchange2Items, retE error) {

	url := "http://vip.stock.finance.sina.com.cn/quotes_service/view/js/qihuohangqing.js"

	params := make(map[string]string, 0)
	params["charset"] = "gbk"

	resp, err := pkNet.HttpGetEx(url, params)

	if err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, errors.New("fetchExchanges server response is empty")
	}

	if reg1 == nil {
		reg1 = regexp.MustCompile(`(\{[^\}]*\})`)
		vm = otto.New()
	}

	regRst := reg1.FindSubmatch(resp)
	if len(regRst) < 2 {
		retE = fmt.Errorf("parse exchanges not found: ARRFUTURESNODES found:%v", regRst)
		return
	}
	resp = regRst[1]

	jsValue, err := vm.Run("(" + string(resp) + ")")
	if err != nil {
		retE = err
		return
	}

	if !jsValue.IsObject() {
		retE = errors.New("JSObject is not object, is:" + jsValue.Class())
		return
	}

	jsObj := jsValue.Object()
	exchangeIds := jsObj.Keys()
	if len(exchangeIds) < 1 {
		retE = fmt.Errorf("parse exchanges not found: exchanges keys, found:%v", exchangeIds)
		return
	}

	validPeriods := kline.GetValidPeriods()

	for _, exchangeId := range exchangeIds {

		jsValue, retE = jsObj.Get(exchangeId)
		if retE != nil {
			return
		}

		if !jsValue.IsObject() {
			retE = errors.New("JSObject2 is not object, is:" + jsValue.Class())
			return
		}

		jsObj2 := jsValue.Object()
		colsKeys := jsObj2.Keys()
		if len(colsKeys) < 1 {
			retE = fmt.Errorf("parse exchanges not found: colsKeys keys, found:%v", exchangeIds)
			return
		}

		var exchange2Items Exchange2Items
		exchange2Items.ItemTypes = make([]ExchangeItemsType, 0)
		exchange2Items.Id = exchangeId
		exchange2Items.KPeriods = validPeriods
		exchange2Items.ExchangeType = ExchangeType_Futures
		exchange2Items.Country = Country_CN

		for idx, colKey := range colsKeys {
			jsValue, retE = jsObj2.Get(colKey)
			if retE != nil {
				return
			}

			if idx == 0 {
				if !jsValue.IsString() {
					retE = errors.New("JSObject3 [0] is not string, is:" + jsValue.Class())
					return
				}
				exchange2Items.Name = jsValue.String()
			} else {
				if !jsValue.IsObject() {
					retE = errors.New("JSObject3 is not object, is:" + jsValue.Class())
					return
				}

				jsObj4 := jsValue.Object()
				typesKeys := jsObj4.Keys()
				if len(typesKeys) < 2 {
					retE = fmt.Errorf("parse exchanges not found: typesKeys keys, found:%v", typesKeys)
					return
				}

				var exchangeItemsType ExchangeItemsType

				jsValue, retE = jsObj4.Get("0")
				if retE != nil {
					return
				}

				if !jsValue.IsString() {
					retE = errors.New("typeId name is not string, is:" + jsValue.Class())
					return
				}
				exchangeItemsType.Name = jsValue.String()

				jsValue, retE = jsObj4.Get("1")
				if retE != nil {
					return
				}

				if !jsValue.IsString() {
					retE = errors.New("typeId id not string, is:" + jsValue.Class())
					return
				}
				exchangeItemsType.Id = jsValue.String()

				// exchange2Items.ItemTypes = append(exchange2Items.ItemTypes, exchangeItemsType)

			}
		}

		ret = append(ret, exchange2Items)

	}

	return

}

func sinaFetchExchangeSymbols(foundExchange *Exchange2Items) (retList []map[string]string, retE error) {

	switch foundExchange.ExchangeType {
	case ExchangeType_Futures:
		return sinaFetchExchangeSymbols_Futures(foundExchange)
	}
	return
}

func sinaFetchExchangeSymbols_Futures(foundExchange *Exchange2Items) (retList []map[string]string, retE error) {
	exchangeId := foundExchange.Id

	if Exchange_CustomFutures == exchangeId {
		retList = []map[string]string{
			map[string]string{"id": "AU0", "name": "黄金连续"},
			map[string]string{"id": "AG0", "name": "白银连续"},
			map[string]string{"id": "CU0", "name": "沪铜连续"},
			map[string]string{"id": "NI0", "name": "沪镍连续"},
			map[string]string{"id": "P0", "name": "棕榈连续"},
			map[string]string{"id": "RU0", "name": "橡胶连续"},
			map[string]string{"id": "J0", "name": "焦炭连续"},
			map[string]string{"id": "M0", "name": "豆粕连续"},
			map[string]string{"id": "HC0", "name": "热轧卷板连续"},
			map[string]string{"id": "RB0", "name": "螺纹钢连续"},
			map[string]string{"id": "I0", "name": "铁矿石连续"},
			map[string]string{"id": "BU0", "name": "沥青连续"},
			map[string]string{"id": "Y0", "name": "豆油连续"},
		}
		return
	}
	if foundExchange.ItemTypes == nil {
		retE = errors.New("ItemTypes为空")
		return
	}

	retList = make([]map[string]string, 0)

	for _, itemTypes := range foundExchange.ItemTypes {
		l, e := sinaFetchExchangeSymbols_FuturesType(foundExchange, &itemTypes)
		if e != nil {
			retE = e
			return
		}

		retList = append(retList, l...)
	}

	return
}

func sinaFetchExchangeSymbols_FuturesType(foundExchange *Exchange2Items, exchangeItemsType *ExchangeItemsType, params ...string) (retList []map[string]string, retE error) {

	wsMethodName := "getHQFuturesData"
	if foundExchange.Id == "cffex" {
		wsMethodName = "getNameList"
	}

	url := fmt.Sprintf("http://vip.stock.finance.sina.com.cn/quotes_service/api/json_v2.php/Market_Center.%s?node=%s&base=futures", wsMethodName, exchangeItemsType.Id)

	resp, err := pkNet.HttpGetEx(url)

	if err != nil {
		retE = err
		return
	}

	if len(resp) == 0 {
		retE = errors.New("fetchExchangeSymbols server response is empty")
		return
	}

	resp = []byte(pkJson.FixJsonKey(string(resp)))

	var quotes []map[string]string
	retE = json.Unmarshal(resp, &quotes)
	if retE != nil {
		return
	}

	var lExchangeId, lExchangeId2 string
	if len(params) > 0 {
		lExchangeId = strings.ToLower(params[0])
	}

	retList = make([]map[string]string, 0)

	var ok bool
	for _, quote := range quotes {
		retItem := make(map[string]string, 0)

		if lExchangeId != "" {
			lExchangeId2, ok = quote["market"]
			if !ok {
				retE = errors.New("market not found")
				return
			}
			if strings.ToLower(lExchangeId2) != lExchangeId {
				continue
			}
		}

		retItem["id"], ok = quote["symbol"]
		if !ok {
			retE = errors.New("symbol not found")
			return
		}

		retItem["name"], ok = quote["name"]
		if !ok {
			// retE = errors.New("name not found")
			// return
			retItem["name"] = retItem["id"]
		}

		retList = append(retList, retItem)
	}

	return

	// symbol:"TA0",market:"czce",contract:"PTA",name:"PTA连续",trade:"5394",settlement:"5408",prevsettlement:"5416",open:"5420",high:"5440",low:"5374",close:"5400",bid:"5392",ask:"5394",bidvol:"200",askvol:"36",volume:"581910",position:"981548",currentvol:"0",ticktime:"15:00:00",tradedate:"2017-12-11",changepercent:"-0.0040620"
}
