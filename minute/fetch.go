package minute

//

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/pkrss/go-sina/quote"

	pkNet "github.com/pkrss/go-utils/net"

	"github.com/robertkrimen/otto"
)

var vm *otto.Otto

func fetchSinaMinuteFutureDo(symbol string) (retList []MINUTEDATA, retE error) {

	url := fmt.Sprintf("http://stock2.finance.sina.com.cn/futures/api/jsonp.php/%20%20/InnerFuturesNewService.getMinLine?symbol=%s", symbol)

	resp, err := pkNet.HttpGet(url)

	if err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, errors.New("fetchSinaMinuteFutureDo server response is empty")
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
		retE = fmt.Errorf("future minute response root is not object, is %s", jsValue.Class())
		return
	}

	jsObj := jsValue.Object()
	jsObjKeys1 := jsObj.Keys()
	if len(jsObjKeys1) < 1 {
		retE = fmt.Errorf("parse future minute response root keys error, found:%v", jsObjKeys1)
		return
	}

	retList = make([]MINUTEDATA, len(jsObjKeys1))

	var s, sDate string
	var prevSettlement float64
	var now = time.Now().Unix()

	for idx, jsObjKey1 := range jsObjKeys1 {
		jsValue, retE = jsObj.Get(jsObjKey1)
		if retE != nil {
			return
		}

		if !jsValue.IsObject() {
			retE = errors.New("future minute response level2 is not object, is %s" + jsValue.Class())
			return
		}

		jsObj2 := jsValue.Object()
		jsObjKeys2 := jsObj2.Keys()
		if len(jsObjKeys2) < 5 {
			retE = fmt.Errorf("parse future minute response level2 not 5 cnt, found:%v", jsObjKeys2)
			return
		}

		var item MINUTEDATA
		var prevItem *MINUTEDATA
		item.Index = int64(idx)

		if len(jsObjKeys2) >= 7 {
			var sTime string
			jsValue, retE = jsObj2.Get("0")
			if retE != nil {
				return
			}
			sTime, retE = jsValue.ToString()
			if retE != nil {
				return
			}

			jsValue, retE = jsObj2.Get("5")
			if retE != nil {
				return
			}
			prevSettlement, retE = jsValue.ToFloat()
			if retE != nil {
				return
			}

			jsValue, retE = jsObj2.Get("6")
			if retE != nil {
				return
			}
			sDate, retE = jsValue.ToString()
			if retE != nil {
				return
			}

			t, err := time.ParseInLocation("2006-01-02 15:04", sDate+" "+sTime, time.Local)
			if err != nil {
				retE = err
				return
			}
			if t.Unix() > now {
				u := t.Unix()
				u = u - 24*60*60
				t = time.Unix(u, 0)
			}
			sDate = t.Format("2006-01-02")
		}

		for _, jsObjKey2 := range jsObjKeys2 {
			jsValue, retE = jsObj2.Get(jsObjKey2)
			if retE != nil {
				return
			}
			switch jsObjKey2 {
			case "0":
				s, retE = jsValue.ToString()
				if retE != nil {
					return
				}
				t, err := time.ParseInLocation("2006-01-02 15:04", sDate+" "+s, time.Local)
				if err != nil {
					retE = err
					return
				}
				item.Time = t.Unix() * 1000
				if prevItem != nil && item.Time < prevItem.Time {
					if item.Time+24*60*60 <= now {
						item.Time += 24 * 60 * 60

						if len(jsObjKeys2) < 7 {
							var ok bool
							prevSettlement, ok = quote.GetHqFloatField(symbol, "1d", "prevSettlement")
							if !ok {
								retE = fmt.Errorf("prevSettlement fetch error: %s", symbol)
								return
							}
						}
					}
				}
			case "1":
				item.Price, retE = jsValue.ToFloat()

				if prevSettlement > 0 {
					item.UpDn = math.Trunc((item.Price-prevSettlement)*1e3+0.5) * 1e-3
					item.UpDnPer = math.Trunc(((item.Price-prevSettlement)*100/prevSettlement)*1e3+0.5) * 1e-3
				}
			case "2":
				item.Avg, retE = jsValue.ToFloat()
			case "3":
				item.Vol, retE = jsValue.ToInteger()
			case "4":
				item.Positions, retE = jsValue.ToInteger()
			}
			if retE != nil {
				return
			}
		}
		retList[idx] = item
		prevItem = &retList[idx]
	}

	return
}
