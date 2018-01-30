package kline

import (
	"encoding/json"
	"errors"
	"fmt"
	hxUtils "hx98/base/utils"
	"strconv"
	sxApi "sx98/quote2/api"
	"time"

	"github.com/robertkrimen/otto"
)

var vm *otto.Otto

func fetchSinaKlineFutureDo(symbol string, period string) (retList []K_MACDV, retE error) {

	validPeriods := sxApi.GetValidPeriods()

	s := ""
	found := false
	for _, v := range validPeriods {
		if v == period {
			found = true
			break
		}
	}
	if !found {
		retE = fmt.Errorf("周期只能为: %v", validPeriods)
		return
	}

	// http://stock.finance.sina.com.cn/futures/api/jsonp.php/var%20kke_future_nf_IF0=/InterfaceInfoService.getMarket?category=nf&symbol=IF
	// http://stock2.finance.sina.com.cn/futures/api/jsonp.php/ret=/InnerFuturesNewService.getFewMinLine?symbol=RB0&type=5
	// http://stock2.finance.sina.com.cn/futures/api/jsonp.php/var%20_IF02017_12_18=/InnerFuturesNewService.getDailyKLine?symbol=IF0&_=2017_12_18

	timeFormat := "2006-01-02 15:04:05"

	var url string
	if period == "1d" {
		url = fmt.Sprintf("http://stock2.finance.sina.com.cn/futures/api/jsonp.php/%%20/InnerFuturesNewService.getDailyKLine?symbol=%s&_=%s", symbol, time.Now().Format("2016-01-02"))
		timeFormat = "2006-01-02"
	} else {
		var min int
		cnt, err := fmt.Sscanf(period, "%dm", &min)
		if err != nil {
			retE = err
			return
		}
		if cnt < 1 {
			retE = fmt.Errorf("解析分钟时错误 返回个数[%d] < 1", cnt)
			return
		}
		url = fmt.Sprintf("http://stock2.finance.sina.com.cn/futures/api/jsonp.php/%%20/InnerFuturesNewService.getFewMinLine?symbol=%s&type=%d", symbol, min)
	}

	resp, err := hxUtils.HttpGet(url)

	if err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, errors.New("fetchSinaKline server response is empty")
	}

	if vm == nil {
		vm = otto.New()
	}

	jsValue, err := vm.Run(string(resp))
	if err != nil {
		retE = err
		return
	}

	if !jsValue.IsObject() {
		retE = fmt.Errorf("kline response root is not object, is %s", jsValue.Class())
		return
	}

	jsObj := jsValue.Object()
	jsObjKeys1 := jsObj.Keys()
	if len(jsObjKeys1) < 1 {
		retE = fmt.Errorf("parse kline response root keys error, found:%v", jsObjKeys1)
		return
	}

	retList = make([]K_MACDV, len(jsObjKeys1))

	for idx, jsObjKey1 := range jsObjKeys1 {
		jsValue, retE = jsObj.Get(jsObjKey1)
		if retE != nil {
			return
		}

		if !jsValue.IsObject() {
			retE = errors.New("kline response level2 is not object, is %s" + jsValue.Class())
			return
		}

		jsObj2 := jsValue.Object()
		jsObjKeys2 := jsObj2.Keys()
		if len(jsObjKeys2) < 6 {
			retE = fmt.Errorf("parse kline response level2 not 6 cnt, found:%v", jsObjKeys2)
			return
		}

		var item K_MACDV
		item.Index = int64(idx)

		for _, jsObjKey2 := range jsObjKeys2 {
			jsValue, retE = jsObj2.Get(jsObjKey2)
			if retE != nil {
				return
			}
			switch jsObjKey2 {
			case "d":
				s, retE = jsValue.ToString()
				if retE != nil {
					return
				}
				t, err := time.ParseInLocation(timeFormat, s, time.Local)
				if err != nil {
					retE = err
					return
				}
				item.Time = t.Unix() * 1000
			case "o":
				item.Open, retE = jsValue.ToFloat()

			case "h":
				item.High, retE = jsValue.ToFloat()
			case "l":
				item.Low, retE = jsValue.ToFloat()
			case "c":
				item.Close, retE = jsValue.ToFloat()
			case "v":
				item.Vol, retE = jsValue.ToInteger()
			}
			if retE != nil {
				return
			}
		}
		retList[idx] = item
	}

	return
}

// 旧的取K线期货数据，不支持中国金融交易所的股指期货
func fetchSinaKlineFutureDo_bak(symbol string, period string) ([]K_MACDV, error) {

	validPeriods := sxApi.GetValidPeriods()

	found := false
	for _, v := range validPeriods {
		if v == period {
			found = true
			break
		}
	}
	if !found {
		s := fmt.Sprintf("周期只能为: %v", validPeriods)
		return nil, errors.New(s)
	}

	var sinaServiceName string
	timeFormat := "2006-01-02 15:04:05"
	if period == "1d" {
		sinaServiceName = "IndexService.getInnerFuturesDailyKLine"
		timeFormat = "2006-01-02"
	} else {
		sinaServiceName = fmt.Sprintf("IndexService.getInnerFuturesMiniKLine%s", period)
	}

	url := fmt.Sprintf("http://stock2.finance.sina.com.cn/futures/api/json.php/%s?_=%d000/&symbol=%s", sinaServiceName, time.Now().Unix(), symbol)

	resp, err := hxUtils.HttpGet(url)

	if err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, errors.New("fetchSinaKline server response is empty")
	}

	var tmp [][6]string
	err = json.Unmarshal(resp, &tmp)
	if err != nil {
		return nil, err
	}
	l := len(tmp)
	ret := make([]K_MACDV, l)

	for i, r := range tmp {
		if len(r) < 6 {
			return nil, errors.New("数据长度不足6条")
		}

		t, err := time.ParseInLocation(timeFormat, r[0], time.Local)
		if err != nil {
			return nil, err
		}

		var item K_MACDV
		item.Time = t.Unix() * 1000

		item.Index = int64(i)

		item.Open, err = strconv.ParseFloat(r[1], 32)
		if err != nil {
			return nil, err
		}
		item.High, err = strconv.ParseFloat(r[2], 32)
		if err != nil {
			return nil, err
		}
		item.Low, err = strconv.ParseFloat(r[3], 32)
		if err != nil {
			return nil, err
		}
		item.Close, err = strconv.ParseFloat(r[4], 32)
		if err != nil {
			return nil, err
		}
		item.Vol, err = strconv.ParseInt(r[5], 10, 64)
		if err != nil {
			return nil, err
		}
		ret[i] = item
	}

	return ret, nil
}
