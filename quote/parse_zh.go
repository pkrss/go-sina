package quote

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

func parseZh(quote map[string]interface{}, stk string, rspRowCols []string) (e error) {
	var open, price, f1 float64
	var v1 int64

	if len(rspRowCols) < 32 {
		return fmt.Errorf("fetchSinaHq row not have 32+ columns = %d", len(rspRowCols))
	}

	quote["id"] = stk
	quote["name"] = rspRowCols[0]
	f1, e = strconv.ParseFloat(rspRowCols[1], 64)
	if e == nil {
		open = f1
		quote["open"] = open
	}
	f1, e = strconv.ParseFloat(rspRowCols[2], 64)
	if e == nil {
		quote["prevClose"] = f1
	}
	f1, e = strconv.ParseFloat(rspRowCols[3], 64)
	if e == nil {
		price = f1
		quote["price"] = price
	}
	f1, e = strconv.ParseFloat(rspRowCols[4], 64)
	if e == nil {
		quote["high"] = f1
	}
	f1, e = strconv.ParseFloat(rspRowCols[5], 64)
	if e == nil {
		quote["low"] = f1
	}
	f1, e = strconv.ParseFloat(rspRowCols[6], 64)
	if e == nil {
		quote["buy1"] = f1
	}
	f1, e = strconv.ParseFloat(rspRowCols[7], 64)
	if e == nil {
		quote["sell1"] = f1
	}
	v1, e = strconv.ParseInt(rspRowCols[8], 10, 64)
	if e == nil {
		quote["vol"] = v1
	}
	f1, e = strconv.ParseFloat(rspRowCols[7], 64)
	if e == nil {
		quote["amt"] = f1
	}
	v1, e = strconv.ParseInt(rspRowCols[10], 10, 64)
	if e == nil {
		quote["buy1vol"] = v1
	}
	v1, e = strconv.ParseInt(rspRowCols[12], 10, 64)
	if e == nil {
		quote["buy2vol"] = v1
	}
	f1, e = strconv.ParseFloat(rspRowCols[13], 64)
	if e == nil {
		quote["buy2"] = f1
	}
	v1, e = strconv.ParseInt(rspRowCols[14], 10, 64)
	if e == nil {
		quote["buy3vol"] = v1
	}
	f1, e = strconv.ParseFloat(rspRowCols[15], 64)
	if e == nil {
		quote["buy3"] = f1
	}
	v1, e = strconv.ParseInt(rspRowCols[16], 10, 64)
	if e == nil {
		quote["buy4vol"] = v1
	}
	f1, e = strconv.ParseFloat(rspRowCols[17], 64)
	if e == nil {
		quote["buy4"] = f1
	}
	v1, e = strconv.ParseInt(rspRowCols[18], 10, 64)
	if e == nil {
		quote["buy5vol"] = v1
	}
	f1, e = strconv.ParseFloat(rspRowCols[19], 64)
	if e == nil {
		quote["buy5"] = f1
	}
	v1, e = strconv.ParseInt(rspRowCols[20], 10, 64)
	if e == nil {
		quote["sell1vol"] = v1
	}
	v1, e = strconv.ParseInt(rspRowCols[22], 10, 64)
	if e == nil {
		quote["sell2vol"] = v1
	}
	f1, e = strconv.ParseFloat(rspRowCols[23], 64)
	if e == nil {
		quote["sell2"] = f1
	}
	v1, e = strconv.ParseInt(rspRowCols[24], 10, 64)
	if e == nil {
		quote["sell3vol"] = v1
	}
	f1, e = strconv.ParseFloat(rspRowCols[25], 64)
	if e == nil {
		quote["sell3"] = f1
	}
	v1, e = strconv.ParseInt(rspRowCols[26], 10, 64)
	if e == nil {
		quote["sell4vol"] = v1
	}
	f1, e = strconv.ParseFloat(rspRowCols[27], 64)
	if e == nil {
		quote["sell4"] = f1
	}
	v1, e = strconv.ParseInt(rspRowCols[28], 10, 64)
	if e == nil {
		quote["sell5vol"] = v1
	}
	f1, e = strconv.ParseFloat(rspRowCols[29], 64)
	if e == nil {
		quote["sell5"] = f1
	}

	t, err4 := time.ParseInLocation("2006-01-02 15:04:05", rspRowCols[30]+" "+rspRowCols[31], time.Local)
	if err4 == nil {
		quote["time"] = t.Unix() // - 8*60*60
	}

	quote["updnPrice"] = math.Trunc((price-open)*1e3+0.5) * 1e-3
	if open > 0 {
		quote["updnPricePer"] = math.Trunc(((price-open)*100/open)*1e3+0.5) * 1e-3
	} else {
		quote["updnPricePer"] = nil
	}

	return
}
