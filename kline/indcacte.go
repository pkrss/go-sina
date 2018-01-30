package kline

// import (
// 	"strconv"
// 	"github.com/thetruetrade/gotrade"
// 	"github.com/thetruetrade/gotrade/indicators"
// )

// func getIndicateParamInt(indicateParams map[string]string, param string, defValue int) (int){
// 	s, ok := indicateParams[param]
// 	if !ok {
// 		return defValue
// 	}
// 	v, e := strconv.Atoi(s)
// 	if e != nil {
// 		return defValue
// 	}
// 	return v
// }

// func handleKIndicate(macdvList []K_MACDV, indicate string, indicateParams map[string]float32) (retList interface{}, retE error) {

// 	if indicateParams == nil {
// 		indicateParams = make(map[string]float32, 0)
// 	}

// 	priceStream := gotrade.NewDailyDOHLCVStream()

// 	switch indicate {
// 	case "macd":
// 		p_short := getIndicateParamInt(indicateParams, "short", 12)
// 		p_long := getIndicateParamInt(indicateParams, "long", 26)
// 		p_mid := getIndicateParamInt(indicateParams, "mid", 9)
// 		indicators.NewMacdForStream(priceStream, p_long, p_short, p_mid, selectData gotrade.UseClosePrice) (indicator *Macd, err error)
// 	}

// 	sma, _ := indicators.NewSMAForStream(priceStream, 20, gotrade.UseClosePrice)

// 	dohlcv := gotrade.NewDOHLCVDataItem(date, open, high, low, close, volume)

// 	retList = macdvList
// 	retE = nil
// 	return
// }
