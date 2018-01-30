package quote

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

/*
	var hq_str_CFF_RE_IF0="4044.0,4049.8,4022.2,4034.4,7067,8554560000,24339,0.0,,4442.8,3635.2,,,4044.0,4039.0,26389,0,0,--,--,--,--,--,--,--,--,0,0,--,--,--,--,--,--,--,--,2017-12-20,11:04:30,0,0,4061.600,3965.000,4072.000,3939.000,4271.000,3939.000,4271.000,3802.200,53.673";
	0: 4044.0,开
	1: 4049.8,高
	2: 4022.2,低
	3: 4034.4,最新价
	4: 7067,成交量
	5: 8554560000,
	6: 24339, 持仓量
	7: 0.0,
	8: ,
	9: 4442.8, 涨停价
	10: 3635.2, 跌停价
	11: ,
	12: ,
	13: 4044.0, 昨收
	14: 4039.0, 昨结算
	15: 26389, 昨持仓
	16: 0,
	17: 0,
	18: --,
	19: --,
	20: --,
	21: --,
	22: --,
	23: --,
	24: --,
	25: --,
	26: 0,
	27: 0,
	28: --,
	29: --,
	30: --,
	31: --,
	32: --,
	33: --,
	34: --,
	35: --,
	36: 2017-12-20,日期
	37: 11:04:30,时间
	38: 0,
	39: 0,
	40: 4061.600,
	41: 3965.000,
	42: 4072.000,
	43: 3939.000,
	44: 4271.000,
	45: 3939.000,
	46: 4271.000,
	47: 3802.200,
	48: 53.673
*/
//
//
func parseCtpCffre(quote map[string]interface{}, stk string, rspRowCols []string) (e error) {

	var open, price, prevSettlement, f1 float64
	var v1 int64

	if len(rspRowCols) < 17 {
		return fmt.Errorf("fetchSinaHq row not have 17+ columns = %d", len(rspRowCols))
	}

	quote["id"] = stk
	quote["name"] = stk
	f1, e = strconv.ParseFloat(rspRowCols[0], 64)
	if e == nil {
		open = f1
		quote["open"] = open
	}
	f1, e = strconv.ParseFloat(rspRowCols[1], 64)
	if e == nil {
		quote["high"] = f1
	}
	f1, e = strconv.ParseFloat(rspRowCols[2], 64)
	if e == nil {
		quote["low"] = f1
	}
	f1, e = strconv.ParseFloat(rspRowCols[3], 64)
	if e == nil {
		price = f1
		quote["price"] = price
	}
	v1, e = strconv.ParseInt(rspRowCols[4], 10, 64)
	if e == nil {
		quote["vol"] = v1
	}
	v1, e = strconv.ParseInt(rspRowCols[6], 10, 64)
	if e == nil {
		quote["positions"] = v1
	}
	f1, e = strconv.ParseFloat(rspRowCols[9], 64)
	if e == nil {
		quote["upStop"] = f1
	}
	f1, e = strconv.ParseFloat(rspRowCols[10], 64)
	if e == nil {
		quote["dnStop"] = f1
	}
	f1, e = strconv.ParseFloat(rspRowCols[13], 64)
	if e == nil {
		quote["prevClose"] = f1
	}
	f1, e = strconv.ParseFloat(rspRowCols[14], 64)
	if e == nil {
		prevSettlement = f1
		quote["prevSettlement"] = prevSettlement
	}
	v1, e = strconv.ParseInt(rspRowCols[15], 10, 64)
	if e == nil {
		quote["prevPositions"] = v1
	}

	t, err4 := time.ParseInLocation("2006-01-02 15:04:05", rspRowCols[36]+" "+rspRowCols[37], time.Local)
	if err4 == nil {
		quote["time"] = t.Unix() // - 8*60*60
	}

	if prevSettlement == 0 {
		prevSettlement = open
	}
	quote["updnPrice"] = price - prevSettlement
	if prevSettlement > 0 {
		quote["updnPricePer"] = math.Trunc((float64(price-prevSettlement)*100/float64(prevSettlement))*1e3+0.5) * 1e-3
	} else {
		quote["updnPricePer"] = nil
	}
	return nil
}

/*
0：豆粕连续，名字
1：145958，不明数字（难道是数据提供商代码？）
2：3170，开盘价
3：3190，最高价
4：3145，最低价
5：3178，昨日收盘价 （2013年6月27日）
6：3153，买价，即“买一”报价
7：3154，卖价，即“卖一”报价
8：3154，最新价，即收盘价
9：3162，结算价
10：3169，昨结算
11：1325，买  量
12：223，卖  量
13：1371608，持仓量
14：1611074，成交量
15：连，大连商品交易所简称
16：豆粕，品种名简称
17：2013-06-28，日期
*/
func parseCtp(quote map[string]interface{}, stk string, rspRowCols []string) (retStk string, e error) {
	var open, price, prevSettlement, f1 float64
	var v1 int64

	if strings.HasPrefix(stk, "CFF_RE_") {
		retStk = strings.Replace(stk, "CFF_RE_", "", 1)
		e = parseCtpCffre(quote, stk, rspRowCols)
		return
	}

	retStk = stk

	if len(rspRowCols) < 17 {
		return "", fmt.Errorf("fetchSinaHq row not have 17+ columns = %d", len(rspRowCols))
	}

	quote["id"] = stk
	quote["name"] = rspRowCols[0]
	f1, e = strconv.ParseFloat(rspRowCols[2], 64)
	if e == nil {
		open = f1
		quote["open"] = open
	}
	f1, e = strconv.ParseFloat(rspRowCols[3], 64)
	if e == nil {
		quote["high"] = f1
	}
	f1, e = strconv.ParseFloat(rspRowCols[4], 64)
	if e == nil {
		quote["low"] = f1
	}
	f1, e = strconv.ParseFloat(rspRowCols[5], 64)
	if e == nil {
		quote["prevClose"] = f1
	}
	f1, e = strconv.ParseFloat(rspRowCols[6], 64)
	if e == nil {
		quote["buy1"] = f1
	}
	f1, e = strconv.ParseFloat(rspRowCols[7], 64)
	if e == nil {
		quote["sell1"] = f1
	}

	f1, e = strconv.ParseFloat(rspRowCols[8], 64)
	if e == nil {
		price = f1
		quote["price"] = price
	}
	f1, e = strconv.ParseFloat(rspRowCols[9], 64)
	if e == nil {
		quote["settlement"] = f1
	}
	f1, e = strconv.ParseFloat(rspRowCols[10], 64)
	if e == nil {
		prevSettlement = f1
		quote["prevSettlement"] = prevSettlement
	}

	f1, e = strconv.ParseFloat(rspRowCols[11], 64)
	if e == nil {
		quote["buy1vol"] = f1
	}
	v1, e = strconv.ParseInt(rspRowCols[12], 10, 64)
	if e == nil {
		quote["sell1vol"] = v1
	}
	v1, e = strconv.ParseInt(rspRowCols[13], 10, 64)
	if e == nil {
		quote["positions"] = v1
	}
	v1, e = strconv.ParseInt(rspRowCols[14], 10, 64)
	if e == nil {
		quote["vol"] = v1
	}
	quote["exSName"] = rspRowCols[15]
	quote["sName"] = rspRowCols[16]

	t, err4 := time.ParseInLocation("2006-01-02", rspRowCols[17], time.Local)
	if err4 == nil {
		quote["time"] = t.Unix()
	}

	if prevSettlement == 0 {
		prevSettlement = open
	}
	quote["updnPrice"] = price - prevSettlement
	if prevSettlement > 0 {
		quote["updnPricePer"] = math.Trunc((float64(price-prevSettlement)*100/float64(prevSettlement))*1e3+0.5) * 1e-3
	} else {
		quote["updnPricePer"] = nil
	}

	return
}
